package builder

import (
	"fmt"
	"strings"
)

const (
	Eq = iota
	Ei
	Ne
	Ni
	Ge
	Gt
	Le
	Lt
)

var operation = [...]string{
	Eq: "=",
	Ei: "IS",
	Ne: "<>",
	Ni: "IS NOT",
	Ge: ">=",
	Gt: ">",
	Le: "<=",
	Lt: "<",
}

type Filter struct {
	strings.Builder
	v []any
}

func (f *Filter) Size() int {
	return len(f.v)
}

func (f *Filter) Values() []any {
	return f.v
}

func (f *Filter) Hold(v any) string {
	f.v = append(f.v, v)
	return fmt.Sprintf("$%d", len(f.v))
}

func (f *Filter) Op(k string, o int, v any) error {
	var p string
	switch x := v.(type) {
	case nil:
		p = "NULL"
		o++
	case bool:
		if x {
			p = "TRUE"
		} else {
			p = "FALSE"
		}
		o++
	default:
		p = f.Hold(v)
	}
	_, err := fmt.Fprintf(f, "%q %s %s", k, operation[o], p)
	return err
}

func (f *Filter) Eq(k string, v any) error {
	return f.Op(k, Eq, v)
}

func (f *Filter) Ne(k string, v any) error {
	return f.Op(k, Ne, v)
}

func (f *Filter) Ge(k string, v any) error {
	return f.Op(k, Ge, v)
}

func (f *Filter) Gt(k string, v any) error {
	return f.Op(k, Gt, v)
}

func (f *Filter) Le(k string, v any) error {
	return f.Op(k, Le, v)
}

func (f *Filter) Lt(k string, v any) error {
	return f.Op(k, Lt, v)
}

func (f *Filter) In(s string, a ...any) error {
	//TODO implement me
	panic("implement me")
}

func (f *Filter) Ni(s string, a ...any) error {
	//TODO implement me
	panic("implement me")
}

func (f *Filter) As(s string, a any) error {
	//TODO implement me
	panic("implement me")
}
