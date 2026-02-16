package builder

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pshvedko/dbx/filter"
)

type Comma int

var (
	place = []byte{'$'}
	comma = []byte{','}
	space = []byte{' '}
)

func (c Comma) Format(f fmt.State, _ rune) {
	if c > 0 {
		_, _ = f.Write(comma)
	}
}

type Holder int

func (h Holder) Format(f fmt.State, _ rune) {
	_, _ = f.Write(place)
	var buf [20]byte
	b := strconv.AppendInt(buf[:0], int64(h), 10)
	_, _ = f.Write(b)
}

type Keyword = filter.Special

const (
	ASC     Keyword = "ASC"
	DESC    Keyword = "DESC"
	NULL    Keyword = "NULL"
	TRUE    Keyword = "TRUE"
	FALSE   Keyword = "FALSE"
	DEFAULT Keyword = "DEFAULT"
)

type By [2]fmt.Formatter

func (b By) Format(f fmt.State, r rune) {
	b[0].Format(f, r)
	if b[1] != nil {
		_, _ = f.Write(space)
		b[1].Format(f, r)
	}
}

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

func (f *Filter) Width() (int, bool) {
	return 0, false
}

func (f *Filter) Precision() (int, bool) {
	return 0, false
}

func (f *Filter) Flag(int) bool {
	return false
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

var (
	equal = []byte{' ', '=', ' '}
	inner = []byte{' ', '=', ' ', 'A', 'N', 'Y', '('}
	end   = []byte{')'}
)

func (f *Filter) EQ(k, v fmt.Formatter) (int, error) {
	_, _ = fmt.Fprint(f, k)
	_, _ = f.Write(equal)
	_, _ = fmt.Fprint(f, v)
	return 0, nil
}

func (f *Filter) IN(k, v fmt.Formatter) (int, error) {
	_, _ = fmt.Fprint(f, k)
	_, _ = f.Write(inner)
	_, _ = fmt.Fprint(f, v)
	_, _ = f.Write(end)
	return 0, nil
}

func (f *Filter) Collation(t filter.Type, k fmt.Formatter, v any) (int, error) {
	switch t {
	case filter.EQ:
		switch v.(type) {
		case nil, bool:
			if k == nil {
				return fmt.Fprint(f, f.Value(v))
			}
			return fmt.Fprintf(f, Is, k, f.Value(v))
		}
		return f.EQ(k, f.Value(v))
	case filter.NE:
		switch v.(type) {
		case nil, bool:
			return fmt.Fprintf(f, Si, k, f.Value(v))
		}
		return fmt.Fprintf(f, Ne, k, f.Value(v))
	case filter.GE:
		return fmt.Fprintf(f, Ge, k, f.Value(v))
	case filter.GT:
		return fmt.Fprintf(f, Gt, k, f.Value(v))
	case filter.LE:
		return fmt.Fprintf(f, Le, k, f.Value(v))
	case filter.LT:
		return fmt.Fprintf(f, Lt, k, f.Value(v))
	case filter.AS:
		return fmt.Fprintf(f, As, k, f.Value(v))
	case filter.NA:
		return fmt.Fprintf(f, Na, k, f.Value(v))
	case filter.IN:
		return f.IN(k, f.Value(v))
	case filter.NI:
		return fmt.Fprintf(f, Ni, k, f.Value(v))
	case filter.FALSE, filter.TRUE:
		return fmt.Fprint(f, f.Value(v))
	case filter.AND, filter.OR:
		fallthrough
	default:
		panic(t)
	}
}

func (f *Filter) Output(n int) error {
	f.v = make([]any, 0, n)
	return nil
}
