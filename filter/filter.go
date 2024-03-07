package filter

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Builder interface {
	io.Writer
	fmt.Stringer
	Eq(string, any) error
	Ne(string, any) error
	Ge(string, any) error
	Gt(string, any) error
	Le(string, any) error
	Lt(string, any) error
	As(string, any) error
	In(string, any) error
	Ni(string, any) error
}

type Fielder interface {
	Names() []string
	Values() []any
}
type Injector interface {
	Get() Projector
	Put(Projector)
}

type Projector interface {
	Fielder
	Table() string
}

type Filter interface {
	To(Builder, Projector) error
}

func conjunction(b Builder, j Projector, o string, ff []Filter) (err error) {
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

func straight(b Builder, j Projector, o string, oo map[string]any, t any) (err error) {
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
	for i, k := range ff {
		if i > 0 {
			_, err = fmt.Fprint(b, " ", o, " ")
			if err != nil {
				return
			}
		}
		switch t.(type) {
		case Eq:
			err = b.Eq(k, oo[k])
		case Ne:
			err = b.Ne(k, oo[k])
		case Ge:
			err = b.Ge(k, oo[k])
		case Gt:
			err = b.Gt(k, oo[k])
		case Le:
			err = b.Le(k, oo[k])
		case Lt:
			err = b.Lt(k, oo[k])
		case As:
			err = b.As(k, oo[k])
		case In:
			err = b.In(k, oo[k])
		case Ni:
			err = b.Ni(k, oo[k])
		default:
			return io.EOF
		}
		if err != nil {
			return
		}
	}
	return
}

type And []Filter

func (f And) To(b Builder, j Projector) error {
	return conjunction(b, j, "AND", f)
}

type Or []Filter

func (f Or) To(b Builder, j Projector) error {
	return conjunction(b, j, "OR", f)
}

type Eq map[string]any

func (f Eq) To(b Builder, j Projector) error {
	return straight(b, j, "AND", f, f)
}

type Ne map[string]any

func (f Ne) To(b Builder, j Projector) error {
	return straight(b, j, "AND", f, f)
}

type Ge map[string]any

func (f Ge) To(b Builder, j Projector) error {
	return straight(b, j, "AND", f, f)
}

type Gt map[string]any

func (f Gt) To(b Builder, j Projector) error {
	return straight(b, j, "AND", f, f)
}

type Le map[string]any

func (f Le) To(b Builder, j Projector) error {
	return straight(b, j, "AND", f, f)
}

type Lt map[string]any

func (f Lt) To(b Builder, j Projector) error {
	return straight(b, j, "AND", f, f)
}

type As map[string]any

func (f As) To(b Builder, j Projector) error {
	return straight(b, j, "AND", f, f)
}

type In map[string]Array

func (f In) To(b Builder, j Projector) error {
	m := make(map[string]any, len(f))
	for k, v := range f {
		m[k] = v
	}
	return straight(b, j, "AND", m, f)
}

type Ni map[string]Array

func (f Ni) To(b Builder, j Projector) error {
	m := make(map[string]any, len(f))
	for k, v := range f {
		m[k] = v
	}
	return straight(b, j, "AND", m, f)
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
