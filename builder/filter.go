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
	Eq: " = %v",
	Is: " IS %v",
	Ne: " <> %v",
	Si: " IS NOT %v",
	Ge: " >= %v",
	Gt: " > %v",
	Le: " <= %v",
	Lt: " < %v",
	In: " = ANY(%v)",
	Ni: " <> ALL(%v)",
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
	switch p.(type) {
	case Keyword:
		switch o {
		case Si, Is:
		default:
			return fmt.Errorf("illegal operation")
		}
	}
	_, err := fmt.Fprint(f, k)
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
