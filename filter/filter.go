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
	Value(int) any
	Auto(int) (any, bool)
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
