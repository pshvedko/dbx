package builder

import (
	"fmt"
	"strings"
)

const (
	Eq = iota
	Is
	Ne
	Si
	Ge
	Gt
	Le
	Lt
	In
	Ni
)

var operation = [...][3]string{
	Eq: {"=", " "},
	Is: {"IS", " "},
	Ne: {"<>", " "},
	Si: {"IS NOT", " "},
	Ge: {">=", " "},
	Gt: {">", " "},
	Le: {"<=", " "},
	Lt: {"<", " "},
	In: {"= ANY", "(", ")"},
	Ni: {"<> ALL", "(", ")"},
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
	_, err := fmt.Fprintf(f, "%q %s%s%s%s", k, operation[o][0], operation[o][1], p, operation[o][2])
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

func (f *Filter) In(k string, v any) error {
	return f.Op(k, In, v)
}

func (f *Filter) Ni(k string, v any) error {
	return f.Op(k, Ni, v)
}

func (f *Filter) As(s string, a any) error {
	//TODO implement me
	panic("implement me")
}
