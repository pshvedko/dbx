package filter

import (
	"database/sql/driver"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"
	"sync"
)

type Table struct {
	Projector
	Alias string
}

func (t Table) Table() string {
	return t.Alias
}

func Conjunction(b Builder, j Projector, o string, ff []Filter) (err error) {
	if len(ff) > 1 {
		_, err = b.WriteString("( ")
		if err != nil {
			return
		}
		defer func() {
			if err == nil {
				_, err = b.WriteString(" )")
			}
		}()
	}
	var i int
	for _, f := range ff {
		if IsEmpty(f) {
			continue
		}
		if i > 0 {
			err = b.WriteByte(' ')
			if err != nil {
				return
			}
			_, err = b.WriteString(o)
			if err != nil {
				return
			}
			err = b.WriteByte(' ')
			if err != nil {
				return
			}
		}
		err = f.To(b, j)
		if err != nil {
			return err
		}
		i++
	}
	return
}

type ErrNoSuchField map[string]struct{}

func (e ErrNoSuchField) Error() string {
	ff := make([]string, 0, len(e))
	for f := range e {
		ff = append(ff, f)
	}
	return fmt.Sprintln("no field:", ff)
}

type Pool[T any] sync.Pool

func (p *Pool[T]) Get() T { return (*sync.Pool)(p).Get().(T) }

func (p *Pool[T]) Put(x T) { (*sync.Pool)(p).Put(x) }

var (
	poolErrNoSuchField = Pool[ErrNoSuchField]{New: func() any { return make(ErrNoSuchField, 32) }}
	poolNameField      = Pool[[]string]{New: func() any { return make([]string, 0, 32) }}
)

func Straight[T any, M interface {
	~map[string]T
	Type() Type
}](b Builder, j Projector, o string, oo M) (err error) {
	nn := poolErrNoSuchField.Get()
	for k := range oo {
		nn[k] = struct{}{}
	}
	ff := poolNameField.Get()
	for _, k := range j.Names() {
		_, ok := oo[k]
		if ok {
			ff = append(ff, k)
			delete(nn, k)
		}
	}
	if len(nn) > 0 {
		return nn
	}
	poolErrNoSuchField.Put(nn)
	defer func() { poolNameField.Put(ff[:0]) }()
	if len(oo) > 1 {
		_, err = b.WriteString("( ")
		if err != nil {
			return
		}
		defer func() {
			if err == nil {
				_, err = b.WriteString(" )")
			}
		}()
	}
	t := j.Table()
	for i, f := range ff {
		if i > 0 {
			err = b.WriteByte(' ')
			if err != nil {
				return
			}
			_, err = b.WriteString(o)
			if err != nil {
				return
			}
			err = b.WriteByte(' ')
			if err != nil {
				return
			}
		}
		_, err = b.AppendColumn(t, f)
		if err != nil {
			return
		}
		_, err = b.Append(oo.Type(), oo[f])
		if err != nil {
			return
		}
	}
	return
}

type Special string

func (s Special) Format(f fmt.State, _ rune) {
	_, _ = s.WriteTo(f)
}

func (s Special) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(s))
	return int64(n), err
}

func Now() Special {
	return "NOW()"
}

func Random() Special {
	return "RANDOM()"
}

func NilIfZero[T any](v T) (any, bool) {
	if reflect.ValueOf(v).IsZero() {
		return nil, true
	}
	return v, false
}

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

type Formatter interface {
	Size() int
	Value(any) fmt.Formatter
	Values() []any
	Len() int
}

type Builder interface {
	io.Writer
	io.ByteWriter
	io.StringWriter
	fmt.Stringer
	Append(Type, any) (int, error)
	AppendInt(int) (int64, error)
	AppendColumn(string, string) (int64, error)
	Formatter
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
	Names() []string
	Columns() map[string]int
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

func (f And) To(b Builder, j Projector) error { return Conjunction(b, j, "AND", f) }

func (f And) Type() Type { return AND }

type Or []Filter

func (f Or) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Or) To(b Builder, j Projector) error { return Conjunction(b, j, "OR", f) }

func (f Or) Type() Type { return OR }

type Eq map[string]any

func (f Eq) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Eq) To(b Builder, j Projector) error { return Straight(b, j, "AND", f) }

func (f Eq) Type() Type { return EQ }

type Ne map[string]any

func (f Ne) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Ne) To(b Builder, j Projector) error { return Straight(b, j, "AND", f) }

func (f Ne) Type() Type { return NE }

type Ge map[string]any

func (f Ge) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Ge) To(b Builder, j Projector) error { return Straight(b, j, "AND", f) }

func (f Ge) Type() Type { return GE }

type Gt map[string]any

func (f Gt) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Gt) To(b Builder, j Projector) error { return Straight(b, j, "AND", f) }

func (f Gt) Type() Type { return GT }

type Le map[string]any

func (f Le) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Le) To(b Builder, j Projector) error { return Straight(b, j, "AND", f) }

func (f Le) Type() Type { return LE }

type Lt map[string]any

func (f Lt) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Lt) To(b Builder, j Projector) error { return Straight(b, j, "AND", f) }

func (f Lt) Type() Type { return LT }

type As map[string]string

func (f As) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f As) To(b Builder, j Projector) error { return Straight(b, j, "AND", f) }

func (f As) Type() Type { return AS }

type Na map[string]string

func (f Na) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Na) To(b Builder, j Projector) error { return Straight(b, j, "AND", f) }

func (f Na) Type() Type { return NA }

type In map[string]Array

func (f In) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f In) To(b Builder, j Projector) error { return Straight(b, j, "AND", f) }

func (f In) Type() Type { return IN }

type Ni map[string]Array

func (f Ni) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Ni) To(b Builder, j Projector) error { return Straight(b, j, "AND", f) }

func (f Ni) Type() Type { return NI }

const RFC3339MICRO = "2006-01-02T15:04:05.999999Z07:00"

type Time interface {
	AppendFormat([]byte, string) []byte
}

type Array []any

func (a Array) Value() (driver.Value, error) {
	var b strings.Builder
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
