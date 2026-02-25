package filter

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
)

type Map[K comparable, V any] sync.Map

func (m *Map[K, V]) Load(k K) (v V, ok bool) {
	x, ok := (*sync.Map)(m).Load(k)
	if ok {
		v = x.(V)
	}
	return
}

func (m *Map[K, V]) Store(k K, v V) { (*sync.Map)(m).Store(k, v) }

type Pool[T any] sync.Pool

func (p *Pool[T]) Get() T { return (*sync.Pool)(p).Get().(T) }

func (p *Pool[T]) Put(x T) { (*sync.Pool)(p).Put(x) }

type Table struct {
	Projector
	Alias string
}

func (t Table) Table() string { return t.Alias }

type Special string

func (s Special) Format(f fmt.State, _ rune) {
	_, _ = s.AppendTo(f)
}

func (s Special) AppendTo(w io.Writer) (int, error) {
	return io.WriteString(w, string(s))
}

func Now() Special {
	return "NOW()"
}

func Random() Special {
	return "RANDOM()"
}

func NilIfZero[T comparable](v T) (any, bool) {
	var z T
	if v == z {
		return nil, true
	}
	return v, false
}

type Column string

// T1 { [INNER] | { LEFT | RIGHT | FULL } [OUTER] } JOIN T2 ON EXPRESSION
// T1 NATURAL { [INNER] | { LEFT | RIGHT | FULL } [OUTER] } JOIN T2
const (
	InnerJoin = iota
	LeftOuterJoin
	RightOuterJoin
	FullOuterJoin
	CrossJoin
	NaturalInnerJoin
	NaturalLeftOuterJoin
	NaturalRightOuterJoin
	NaturalFullOuterJoin
)

type Joiner interface {
	Projector
	Right() Projector
	Type() int
	On(string) And
}

type Join[A Projector, B Projector] struct {
	A A
	B B
	F Filter
}

func (o Join[A, B]) Type() int {
	return InnerJoin
}

func (o Join[A, B]) Right() Projector {
	return o.B
}

func (o Join[A, B]) On(string) And {
	return And{o.F}
}

func (o Join[A, B]) PK() PK {
	return o.A.PK()
}

func (o Join[A, B]) Len() int {
	return o.A.Len() + o.B.Len()
}

func (o Join[A, B]) Name() string {
	return o.A.Name()
}

func (o Join[A, B]) Names() []string {
	return o.A.Names()
}

func (o Join[A, B]) Columns() map[string]int {
	//m := maps.Clone(o.A.Columns())
	//n := len(m)
	//for i, k := range o.B.Names() { FIXME
	//	m[k] = i + n
	//}
	return o.A.Columns()
}

func (o Join[A, B]) Exists(k string) bool {
	return o.A.Exists(k)
}

func (o Join[A, B]) Value(i int) (any, bool, bool) {
	if i < o.A.Len() {
		return o.A.Value(i)
	}
	return o.B.Value(i - o.A.Len())
}

func (o Join[A, B]) Field(i int) (any, bool) {
	if i < o.A.Len() {
		return o.A.Field(i)
	}
	return o.B.Field(i - o.A.Len())
}

func (o Join[A, B]) Self() Projector {
	return &o
}

func (o Join[A, B]) Copy() Copier {
	return o
}

func (o Join[A, B]) Places() []any {
	return append(append(make([]any, 0, o.A.Len()+o.B.Len()), o.A.Places()...), o.B.Places()...)
}

func (o Join[A, B]) Table() string {
	return o.A.Table()
}

func NewJoin[A Projector, B Projector](a A, b B, f Filter) Join[A, B] {
	return Join[A, B]{
		A: a,
		B: b,
		F: f,
	}
}

var _ Joiner = &Join[Projector, Projector]{}

type Injectable[T Fielder] []T

func (o Injectable[T]) Element() Projector {
	var x T
	return x.Self()
}

func (o *Injectable[T]) Inject(j Projector) {
	switch v := j.Copy().(type) {
	case T:
		*o = append(*o, v)
	default:
		panic(v)
	}
}

var _ Injector = &Injectable[Projector]{}

type Rectifier interface {
	Straight(Projector, bool, Filter) error
}

type Unifier interface {
	Conjunct(Projector, bool, []Filter) error
}

type Formatter interface {
	Size() int
	Value(any) fmt.Formatter
	Values() []any
	Len() int
	Rectifier
	Unifier
}

type Appender interface {
	AppendValue(Type, any) (int, error)
	AppendInt(int) (int, error)
	AppendColumn(string, string, Projector) (int, error)
	AppendParenthesis(bool) (int, error)
	AppendVia(bool) (int, error)
	AppendFormat(fmt.Formatter) (int, error)
}

type Builder interface {
	io.Writer
	io.ByteWriter
	io.StringWriter
	fmt.Stringer
	Formatter
	Appender
	IsTrusted() bool
}

type PK []string

func (pk PK) Contains(n string) bool {
	return slices.Contains(pk, n)
}

type Placer interface {
	Places() []any
}

type Fielder interface {
	PK() PK
	Len() int
	Name() string
	Names() []string
	Columns() map[string]int
	Exists(string) bool
	Value(int) (any, bool, bool) // value, none, auto
	Field(int) (any, bool)       // value, none
	Copier
}

type Injector interface {
	Element() Projector
	Inject(Projector)
}

type Copier interface {
	Self() Projector
	Copy() Copier
}

type Projector interface {
	Fielder
	Placer
	Table() string
}

type Filter interface {
	To(Builder, Projector) error
}

type And []Filter

func (f And) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f And) To(b Builder, j Projector) error { return b.Conjunct(j, true, f) }

func (f And) Type() Type { return AND }

type Or []Filter

func (f Or) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Or) To(b Builder, j Projector) error { return b.Conjunct(j, false, f) }

func (f Or) Type() Type { return OR }

type Eq map[string]any

func (f Eq) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Eq) To(b Builder, j Projector) error { return b.Straight(j, true, f) }

func (f Eq) Type() Type { return EQ }

type Ne map[string]any

func (f Ne) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Ne) To(b Builder, j Projector) error { return b.Straight(j, true, f) }

func (f Ne) Type() Type { return NE }

type Ge map[string]any

func (f Ge) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Ge) To(b Builder, j Projector) error { return b.Straight(j, true, f) }

func (f Ge) Type() Type { return GE }

type Gt map[string]any

func (f Gt) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Gt) To(b Builder, j Projector) error { return b.Straight(j, true, f) }

func (f Gt) Type() Type { return GT }

type Le map[string]any

func (f Le) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Le) To(b Builder, j Projector) error { return b.Straight(j, true, f) }

func (f Le) Type() Type { return LE }

type Lt map[string]any

func (f Lt) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Lt) To(b Builder, j Projector) error { return b.Straight(j, true, f) }

func (f Lt) Type() Type { return LT }

type As map[string]string

func (f As) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f As) To(b Builder, j Projector) error { return b.Straight(j, true, f) }

func (f As) Type() Type { return AS }

type Na map[string]string

func (f Na) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Na) To(b Builder, j Projector) error { return b.Straight(j, true, f) }

func (f Na) Type() Type { return NA }

type In map[string]Array

func (f In) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f In) To(b Builder, j Projector) error { return b.Straight(j, true, f) }

func (f In) Type() Type { return IN }

type Ni map[string]Array

func (f Ni) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Ni) To(b Builder, j Projector) error { return b.Straight(j, true, f) }

func (f Ni) Type() Type { return NI }

const RFC3339MICRO = "2006-01-02T15:04:05.999999Z07:00"

type Time interface {
	AppendFormat([]byte, string) []byte
}

type ByteArray []byte

func (a ByteArray) MarshalJSON() ([]byte, error) {
	if a == nil {
		return nil, nil
	}
	h := make([]byte, len(a)<<1+2+3)
	n := copy(h, "\"\\\\x")
	n += hex.Encode(h[n:], a)
	n += copy(h[n:], "\"")
	return h[:n], nil
}

type Array []any

func (a Array) Value() (driver.Value, error) {
	var b strings.Builder
	b.Grow(32)
	err := b.WriteByte('{')
	if err != nil {
		return nil, err
	}
	for i, v := range a {
		if i > 0 {
			err = b.WriteByte(',')
			if err != nil {
				return nil, err
			}
		}
		switch x := v.(type) {
		case nil:
			_, err = b.WriteString("NULL")
		case bool:
			switch x {
			case true:
				_, err = b.WriteString("TRUE")
			case false:
				_, err = b.WriteString("FALSE")
			}
		case Time:
			var t [len(RFC3339MICRO)]byte
			_, err = b.Write(x.AppendFormat(t[:0], RFC3339MICRO))
		case fmt.Stringer:
			err = b.WriteByte('"')
			if err != nil {
				return nil, err
			}
			_, err = b.WriteString(x.String())
			if err != nil {
				return nil, err
			}
			err = b.WriteByte('"')
		case string:
			err = b.WriteByte('"')
			if err != nil {
				return nil, err
			}
			_, err = b.WriteString(x)
			if err != nil {
				return nil, err
			}
			err = b.WriteByte('"')
		default:
			_, err = fmt.Fprint(&b, v)
		}
		if err != nil {
			return nil, err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return nil, err
	}
	return b.String(), nil
}
