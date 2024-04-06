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

func (f *Filter) Eq(k, v any) (int, error) {
	switch v.(type) {
	case nil, bool:
		return fmt.Fprintf(f, Is, k, f.Value(v))
	}
	return fmt.Fprintf(f, Eq, k, f.Value(v))
}

func (f *Filter) Ne(k, v any) (int, error) {
	switch v.(type) {
	case nil, bool:
		return fmt.Fprintf(f, Si, k, f.Value(v))
	}
	return fmt.Fprintf(f, Ne, k, f.Value(v))
}

func (f *Filter) Ge(k, v any) (int, error) {
	return fmt.Fprintf(f, Ge, k, f.Value(v))
}

func (f *Filter) Gt(k, v any) (int, error) {
	return fmt.Fprintf(f, Gt, k, f.Value(v))
}

func (f *Filter) Le(k, v any) (int, error) {
	return fmt.Fprintf(f, Le, k, f.Value(v))
}

func (f *Filter) Lt(k, v any) (int, error) {
	return fmt.Fprintf(f, Lt, k, f.Value(v))
}

func (f *Filter) In(k, v any) (int, error) {
	return fmt.Fprintf(f, In, k, f.Value(v))
}

func (f *Filter) Ni(k, v any) (int, error) {
	return fmt.Fprintf(f, Ni, k, f.Value(v))
}

func (f *Filter) As(k, v any) (int, error) {
	return fmt.Fprintf(f, As, k, f.Value(v))
}

func (f *Filter) Na(k, v any) (int, error) {
	return fmt.Fprintf(f, Na, k, f.Value(v))
}
