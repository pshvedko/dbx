package builder

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/pshvedko/dbx/filter"
)

type FormatFunc func(fmt.State, rune)

func (f FormatFunc) Format(w fmt.State, r rune) {
	f(w, r)
}

type Comma int

var (
	place = []byte{'$'}
	comma = []byte{','}
	space = []byte{' '}
	dummy = []byte{' ', '1'}
)

func (c Comma) Format(f fmt.State, _ rune) {
	if c > 0 {
		_, _ = f.Write(comma)
	}
}

type Holder int

func (h Holder) Format(f fmt.State, _ rune) {
	_, _ = h.WriteTo(f)
}

func (h Holder) WriteTo(w io.Writer) (int64, error) {
	u1, err := w.Write(place)
	if err != nil {
		return int64(u1), err
	}
	var b [20]byte
	u2, err := w.Write(strconv.AppendInt(b[:0], int64(h), 10))
	return int64(u1 + u2), err
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

type By [2]io.WriterTo

func (b By) Format(f fmt.State, _ rune) {
	_, _ = b.WriteTo(f)
}

func (b By) WriteTo(w io.Writer) (int64, error) {
	u1, err := b[0].WriteTo(w)
	if err != nil || b[1] == nil {
		return u1, err
	}
	_, err = w.Write(space)
	if err != nil {
		return u1, err
	}
	u2, err := b[1].WriteTo(w)
	return u1 + 1 + u2, err
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

const DefaultHolderCacheSize = 32

var holders = make([]fmt.Formatter, 0, DefaultHolderCacheSize)

func init() {
	for i := 0; i < cap(holders); i++ {
		h := [21]byte{'$'}
		b := strconv.AppendInt(h[:1], int64(i), 10)
		holders = append(holders, func() FormatFunc {
			return func(f fmt.State, r rune) {
				_, _ = f.Write(b)
			}
		}())
	}
}

func (f *Filter) Add(v any) fmt.Formatter {
	switch x := v.(type) {
	case fmt.Formatter:
		return x
	}
	f.v = append(f.v, v)
	if len(f.v) < len(holders) {
		return holders[len(f.v)]
	}
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

func (f *Filter) Alloc(n int) {
	f.v = make([]any, 0, n)
}
