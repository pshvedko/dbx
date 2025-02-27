package filter

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strings"
)

type Table struct {
	Projector
	Alias string
}

func (t Table) Table() string {
	return t.Alias
}

type Column [2]string

func (c Column) Format(f fmt.State, _ rune) {
	_, _ = fmt.Fprintf(f, "%q.%q", c[0], c[1])
}

func Conjunction(b Builder, j Projector, o string, ff []Filter) (err error) {
	if len(ff) > 1 {
		_, err = fmt.Fprint(b, "( ")
		if err != nil {
			return err
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
				return err
			}
		}
		err = f.To(b, j)
		if err != nil {
			return err
		}
	}
	return err
}

type ErrNoSuchField map[string]struct{}

func (e ErrNoSuchField) Error() string {
	ff := make([]string, 0, len(e))
	for f := range e {
		ff = append(ff, f)
	}
	return fmt.Sprintln("no field:", ff)
}

func Straight[T any](b Builder, j Projector, o string, oo map[string]T, t Type) (err error) {
	nn := make(ErrNoSuchField, len(oo))
	for k := range oo {
		nn[k] = struct{}{}
	}
	ff := make([]string, 0, len(oo))
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
	for i, f := range ff {
		if i > 0 {
			_, err = fmt.Fprint(b, " ", o, " ")
			if err != nil {
				return
			}
		}
		_, err = b.Print(t, Column{j.Table(), f}, oo[f])
		if err != nil {
			return
		}
	}
	return
}

type Special string

func (s Special) Format(f fmt.State, _ rune) {
	_, _ = fmt.Fprint(f, string(s))
}

func Now() Special {
	return "NOW()"
}

//type Joiner interface {
//	To() Projector
//}
//
//type To[T Projector] struct {
//	v T
//	k []string
//}
//
//func (j To[T]) To() Projector {
//	var t T
//	return t
//}
//
//func JoinTo[T Projector](v T, k ...string) Joiner {
//	return To[T]{v: v, k: k}
//}
//
//type Ref[T Projector] struct {
//	v *T
//	k []string
//}
//
//func (j Ref[T]) To() Projector {
//	var t T
//	return t
//}
//
//func JoinAny[T Projector](v *T, k ...string) Joiner {
//	return Ref[T]{v: v, k: k}
//}
//
//type Many[T Projector] struct {
//	v *[]T
//	k []string
//}
//
//func (j Many[T]) To() Projector {
//	var t T
//	return t
//}
//
//func JoinMany[T Projector](v *[]T, k ...string) Joiner {
//	return Many[T]{v: v, k: k}
//}

func NilIfZero[T comparable](v T) (any, bool) {
	var z T
	if v != z {
		return v, false
	}
	return nil, true
}

type Nullable[T comparable] struct {
	v *T
}

func (n Nullable[T]) Scan(v any) error {
	switch x := v.(type) {
	case nil:
	case T:
		*n.v = x
	}
	return nil
}

func Nil[T comparable](v *T) Nullable[T] {
	return Nullable[T]{v: v}
}

type Injectable[T Projector] []T

func (o Injectable[T]) Get() Projector {
	var x T
	return x.Copy()
}

func (o *Injectable[T]) Put(j Projector) {
	switch v := j.Copy().(type) {
	case T:
		*o = append(*o, v)
	default:
		panic("invalid injection")
	}
}

type Valuer interface {
	Size() int
	Value(any) fmt.Formatter
	Values() []any
}

type Builder interface {
	io.Writer
	io.StringWriter
	fmt.Stringer
	Print(Type, any, any) (int, error)
	Valuer
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
	Value(int) (any, bool, bool) // value, none, auto
	Get(int) (any, bool)         // value, none
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

func (f And) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f And) To(b Builder, j Projector) error { return Conjunction(b, j, "AND", f) }

type Or []Filter

func (f Or) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Or) To(b Builder, j Projector) error { return Conjunction(b, j, "OR", f) }

type Eq map[string]any

func (f Eq) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Eq) To(b Builder, j Projector) error { return Straight(b, j, "AND", f, EQ) }

type Ne map[string]any

func (f Ne) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Ne) To(b Builder, j Projector) error { return Straight(b, j, "AND", f, NE) }

type Ge map[string]any

func (f Ge) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Ge) To(b Builder, j Projector) error { return Straight(b, j, "AND", f, GE) }

type Gt map[string]any

func (f Gt) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Gt) To(b Builder, j Projector) error { return Straight(b, j, "AND", f, GT) }

type Le map[string]any

func (f Le) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Le) To(b Builder, j Projector) error { return Straight(b, j, "AND", f, LE) }

type Lt map[string]any

func (f Lt) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Lt) To(b Builder, j Projector) error { return Straight(b, j, "AND", f, LT) }

type As map[string]string

func (f As) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f As) To(b Builder, j Projector) error { return Straight(b, j, "AND", f, AS) }

type Na map[string]string

func (f Na) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Na) To(b Builder, j Projector) error { return Straight(b, j, "AND", f, NA) }

type In map[string]Array

func (f In) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f In) To(b Builder, j Projector) error { return Straight(b, j, "AND", f, IN) }

type Ni map[string]Array

func (f Ni) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f Ni) To(b Builder, j Projector) error { return Straight(b, j, "AND", f, NI) }

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
		case Time:
			_, err = b.Write(x.AppendFormat(make([]byte, 0, len(RFC3339MICRO)), RFC3339MICRO))
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
