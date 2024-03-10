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

var operation = [...]string{
	Eq: "= %v",
	Is: "IS %v",
	Ne: "<> %v",
	Si: "IS NOT %v",
	Ge: ">= %v",
	Gt: "> %v",
	Le: "<= %v",
	Lt: "< %v",
	In: "= ANY(%v)",
	Ni: "<> ALL(%v)",
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

func (f *Filter) Value(v any) fmt.Formatter {
	f.v = append(f.v, v)
	return Holder(len(f.v))
}

func (f *Filter) Op(k fmt.Formatter, o int, v any) error {
	var p fmt.Formatter
	switch x := v.(type) {
	case nil:
		p = NULL
		o++
	case bool:
		if x {
			p = TRUE
		} else {
			p = FALSE
		}
		o++
	default:
		p = f.Value(v)
	}
	_, err := fmt.Fprintf(f, "%v ", k)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(f, operation[o], p)
	return err
}

func (f *Filter) Eq(k fmt.Formatter, v any) error {
	return f.Op(k, Eq, v)
}

func (f *Filter) Ne(k fmt.Formatter, v any) error {
	return f.Op(k, Ne, v)
}

func (f *Filter) Ge(k fmt.Formatter, v any) error {
	return f.Op(k, Ge, v)
}

func (f *Filter) Gt(k fmt.Formatter, v any) error {
	return f.Op(k, Gt, v)
}

func (f *Filter) Le(k fmt.Formatter, v any) error {
	return f.Op(k, Le, v)
}

func (f *Filter) Lt(k fmt.Formatter, v any) error {
	return f.Op(k, Lt, v)
}

func (f *Filter) In(k fmt.Formatter, v any) error {
	return f.Op(k, In, v)
}

func (f *Filter) Ni(k fmt.Formatter, v any) error {
	return f.Op(k, Ni, v)
}

func (f *Filter) As(fmt.Formatter, any) error {
	//TODO implement me
	panic("implement me")
}
