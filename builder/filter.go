package builder

import (
	"fmt"
	"strings"
)

type Comma int

func (c Comma) Format(f fmt.State, _ rune) {
	if c > 0 {
		_, _ = f.Write([]byte{','})
	}
}

type Holder int

func (h Holder) Format(f fmt.State, _ rune) {
	_, _ = fmt.Fprint(f, "$", int(h))
}

type Keyword string

func (k Keyword) Format(f fmt.State, _ rune) {
	_, _ = fmt.Fprint(f, string(k))
}

const (
	NULL    Keyword = "NULL"
	TRUE    Keyword = "TRUE"
	FALSE   Keyword = "FALSE"
	DEFAULT Keyword = "DEFAULT"
)

type Filter struct {
	strings.Builder
	v []any
}

func (f *Filter) Value(v any) fmt.Formatter {
	switch x := v.(type) {
	case nil:
		return NULL
	case bool:
		if x {
			return TRUE
		}
		return FALSE
	default:
		return f.Add(v)
	}
}

func (f *Filter) Size() int {
	return len(f.v)
}

func (f *Filter) Values() []any {
	return f.v
}

func (f *Filter) Add(v any) fmt.Formatter {
	switch x := v.(type) {
	case fmt.Formatter:
		return x
	}
	f.v = append(f.v, v)
	return Holder(len(f.v))
}

const (
	Eq = "%v = %v"
	Is = "%v IS %v"
	Ne = "%v <> %v"
	Si = "%v IS NOT %v"
	Ge = "%v >= %v"
	Gt = "%v > %v"
	Le = "%v <= %v"
	Lt = "%v < %v"
	In = "%v = ANY(%v)"
	Ni = "%v <> ALL(%v)"
	As = "%v LIKE %v"
	Na = "%v NOT LIKE %v"
)

func (f *Filter) Eq(k fmt.Formatter, v any) error {
	switch v.(type) {
	case nil, bool:
		return f.Op(k, Is, v)
	}
	return f.Op(k, Eq, v)
}

func (f *Filter) Ne(k fmt.Formatter, v any) error {
	switch v.(type) {
	case nil, bool:
		return f.Op(k, Si, v)
	}
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

func (f *Filter) As(k fmt.Formatter, v any) error {
	return f.Op(k, As, v)
}

func (f *Filter) Na(k fmt.Formatter, v any) error {
	return f.Op(k, Na, v)
}

func (f *Filter) Op(k fmt.Formatter, o string, v any) error {
	_, err := fmt.Fprintf(f, o, k, f.Value(v))
	return err
}
