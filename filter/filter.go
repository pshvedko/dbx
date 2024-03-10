package filter

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Column [2]string

func (c Column) Format(f fmt.State, _ rune) {
	_, _ = fmt.Fprintf(f, "%q.%q", c[0], c[1])
}

func Conjunction(b Builder, j Projector, o string, ff []Filter) (err error) {
	if len(ff) > 1 {
		_, err = fmt.Fprint(b, "( ")
		if err != nil {
			return
		}
		defer func() {
			if err == nil {
				_, err = fmt.Fprint(b, " )")
			}
		}()
	}
	for i, f := range ff {
		if i > 0 {
			_, err = fmt.Fprint(b, " ", o, " ")
			if err != nil {
				return
			}
		}
		err = f.To(b, j)
		if err != nil {
			return
		}
	}
	return
}

func Straight[T any](b Builder, j Projector, o string, oo map[string]T, v any) (err error) {
	ff := make([]string, 0, len(oo))
	for _, k := range j.Names() {
		_, ok := oo[k]
		if ok {
			ff = append(ff, k)
		}
	}
	if len(ff) != cap(ff) {
		return io.EOF
	}
	if len(oo) > 1 {
		_, err = fmt.Fprint(b, "( ")
		if err != nil {
			return
		}
		defer func() {
			if err == nil {
				_, err = fmt.Fprint(b, " )")
			}
		}()
	}
	t := j.Table()
	for i, f := range ff {
		if i > 0 {
			_, err = fmt.Fprint(b, " ", o, " ")
			if err != nil {
				return
			}
		}
		switch v.(type) {
		case Eq:
			err = b.Eq(Column{t, f}, oo[f])
		case Ne:
			err = b.Ne(Column{t, f}, oo[f])
		case Ge:
			err = b.Ge(Column{t, f}, oo[f])
		case Gt:
			err = b.Gt(Column{t, f}, oo[f])
		case Le:
			err = b.Le(Column{t, f}, oo[f])
		case Lt:
			err = b.Lt(Column{t, f}, oo[f])
		case As:
			err = b.As(Column{t, f}, oo[f])
		case In:
			err = b.In(Column{t, f}, oo[f])
		case Ni:
			err = b.Ni(Column{t, f}, oo[f])
		default:
			return io.EOF
		}
		if err != nil {
			return
		}
	}
	return
}

func Nil[T comparable](v T) any {
	var z T
	if v != z {
		return v
	}
	return nil
}

type Injectable[T Projector] []T

func (o Injectable[T]) Get() Projector {
	var x T
	return x.Copy()
}

func (o *Injectable[T]) Put(j Projector) {
	switch v := j.(type) {
	case T:
		x := v.Copy()
		switch t := x.(type) {
		case T:
			*o = append(*o, t)
		default:
			panic("invalid copy")
		}
	default:
		panic("invalid injection")
	}
}

type Builder interface {
	io.Writer
	fmt.Stringer
	Eq(fmt.Formatter, any) error
	Ne(fmt.Formatter, any) error
	Ge(fmt.Formatter, any) error
	Gt(fmt.Formatter, any) error
	Le(fmt.Formatter, any) error
	Lt(fmt.Formatter, any) error
	As(fmt.Formatter, any) error
	In(fmt.Formatter, any) error
	Ni(fmt.Formatter, any) error
}

type PK []string

func (pk PK) Format(f fmt.State, _ rune) {
	if len(pk) > 0 {
		_, _ = fmt.Fprintf(f, "%q", pk[0])
		for _, k := range pk[1:] {
			_, _ = fmt.Fprintf(f, ", %q", k)
		}
	}
}

func (pk PK) Have(n string) bool {
	for _, k := range pk {
		if k == n {
			return true
		}
	}
	return false
}

type Fielder interface {
	PK() PK
	Names() []string
	Values() []any
	Value(int) (any, bool, bool)
	Get(int) any
}

type Injector interface {
	Get() Projector
	Put(Projector)
}

type Copier interface {
	Copy() Projector
}

type Projector interface {
	Fielder
	Copier
	Table() string
}

type Filter interface {
	To(Builder, Projector) error
}

type And []Filter

func (f And) To(b Builder, j Projector) error {
	return Conjunction(b, j, "AND", f)
}

type Or []Filter

func (f Or) To(b Builder, j Projector) error {
	return Conjunction(b, j, "OR", f)
}

type Eq map[string]any

func (f Eq) To(b Builder, j Projector) error {
	return Straight(b, j, "AND", f, f)
}

type Ne map[string]any

func (f Ne) To(b Builder, j Projector) error {
	return Straight(b, j, "AND", f, f)
}

type Ge map[string]any

func (f Ge) To(b Builder, j Projector) error {
	return Straight(b, j, "AND", f, f)
}

type Gt map[string]any

func (f Gt) To(b Builder, j Projector) error {
	return Straight(b, j, "AND", f, f)
}

type Le map[string]any

func (f Le) To(b Builder, j Projector) error {
	return Straight(b, j, "AND", f, f)
}

type Lt map[string]any

func (f Lt) To(b Builder, j Projector) error {
	return Straight(b, j, "AND", f, f)
}

type As map[string]any

func (f As) To(b Builder, j Projector) error {
	return Straight(b, j, "AND", f, f)
}

type In map[string]Array

func (f In) To(b Builder, j Projector) error {
	return Straight(b, j, "AND", f, f)
}

type Ni map[string]Array

func (f Ni) To(b Builder, j Projector) error {
	return Straight(b, j, "AND", f, f)
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
			b.WriteByte(',')
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
		case time.Time:
			_, err = b.Write(pq.FormatTimestamp(x))
		case string, fmt.Stringer:
			_, err = fmt.Fprintf(&b, "%q", v)
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
