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
	comma = []byte{','}
	space = []byte{' '}
	dummy = []byte{' ', '1'}
	among = []byte{'"', '.', '"'}
	quote = []byte{'"'}
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
	b := [21]byte{'$'}
	n, err := w.Write(strconv.AppendInt(b[:1], int64(h), 10))
	return int64(n), err
}

type Int int

func (i Int) Format(f fmt.State, _ rune) {
	_, _ = i.WriteTo(f)
}

func (i Int) WriteTo(w io.Writer) (int64, error) {
	var b [20]byte
	n, err := w.Write(strconv.AppendInt(b[:0], int64(i), 10))
	return int64(n), err
}

type Column [2]string

func (c Column) Format(f fmt.State, _ rune) {
	_, _ = c.WriteTo(f)
}

func (c Column) WriteTo(w io.Writer) (int64, error) {
	u1, err := w.Write(quote)
	if err != nil {
		return int64(u1), err
	}
	u2, err := io.WriteString(w, c[0])
	if err != nil {
		return int64(u1 + u2), err
	}
	u3, err := w.Write(among)
	if err != nil {
		return int64(u1 + u2 + u3), err
	}
	u4, err := io.WriteString(w, c[1])
	if err != nil {
		return int64(u1 + u2 + u3 + u4), err
	}
	u5, err := w.Write(quote)
	return int64(u1 + u2 + u3 + u4 + u5), err
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
	n1, err := b[0].WriteTo(w)
	if err != nil || b[1] == nil {
		return n1, err
	}
	_, err = w.Write(space)
	if err != nil {
		return n1, err
	}
	n2, err := b[1].WriteTo(w)
	return n1 + 1 + n2, err
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
		holders = append(holders, func(b []byte) FormatFunc {
			return func(f fmt.State, r rune) {
				_, _ = f.Write(b)
			}
		}(strconv.AppendInt(h[:1], int64(i), 10)))
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

func (f *Filter) AppendInt(i int) (int64, error) {
	return Int(i).WriteTo(f)
}

func (f *Filter) AppendColumn(t string, c string) (int64, error) {
	return Column{t, c}.WriteTo(f)
}

const (
	Eq = " = %v"
	Is = " IS %v"
	Ne = " <> %v"
	Si = " IS NOT %v"
	Ge = " >= %v"
	Gt = " > %v"
	Le = " <= %v"
	Lt = " < %v"
	In = " = ANY(%v)"
	Ni = " <> ALL(%v)"
	As = " LIKE %v"
	Na = " NOT LIKE %v"
)

var (
	equal = []byte{' ', '=', ' '}
	inner = []byte{' ', '=', ' ', 'A', 'N', 'Y', '('}
	end   = []byte{')'}
)

func (f *Filter) EQ(v fmt.Formatter) (int, error) {
	n1, err := f.Write(equal)
	if err != nil {
		return n1, err
	}
	n2, err := fmt.Fprint(f, v)
	return n1 + n2, err
}

func (f *Filter) IN(v fmt.Formatter) (int, error) {
	n1, err := f.Write(inner)
	if err != nil {
		return n1, err
	}
	n2, err := fmt.Fprint(f, v)
	if err != nil {
		return n1 + n2, err
	}
	n3, err := f.Write(end)
	return n1 + n2 + n3, err
}

func (f *Filter) Append(t filter.Type, v any) (int, error) {
	switch t {
	case filter.EQ:
		switch v.(type) {
		case nil, bool:
			return fmt.Fprintf(f, Is, f.Value(v))
		}
		return f.EQ(f.Value(v))
	case filter.NE:
		switch v.(type) {
		case nil, bool:
			return fmt.Fprintf(f, Si, f.Value(v))
		}
		return fmt.Fprintf(f, Ne, f.Value(v))
	case filter.GE:
		return fmt.Fprintf(f, Ge, f.Value(v))
	case filter.GT:
		return fmt.Fprintf(f, Gt, f.Value(v))
	case filter.LE:
		return fmt.Fprintf(f, Le, f.Value(v))
	case filter.LT:
		return fmt.Fprintf(f, Lt, f.Value(v))
	case filter.AS:
		return fmt.Fprintf(f, As, f.Value(v))
	case filter.NA:
		return fmt.Fprintf(f, Na, f.Value(v))
	case filter.IN:
		return f.IN(f.Value(v))
	case filter.NI:
		return fmt.Fprintf(f, Ni, f.Value(v))
	case filter.FALSE, filter.TRUE:
		return fmt.Fprint(f, f.Value(v))
	case filter.AND, filter.OR:
		fallthrough
	default:
		panic(t)
	}
}

func (f *Filter) Conjunct(j filter.Projector, o string, ff []filter.Filter) error {
	return Conjunction(f, j, o, ff)
}

func (f *Filter) Straight(j filter.Projector, o string, x filter.Filter) error {
	switch x := x.(type) {
	case filter.Eq:
		return Straight(f, j, o, x)
	case filter.Ne:
		return Straight(f, j, o, x)
	case filter.Ge:
		return Straight(f, j, o, x)
	case filter.Gt:
		return Straight(f, j, o, x)
	case filter.Le:
		return Straight(f, j, o, x)
	case filter.Lt:
		return Straight(f, j, o, x)
	case filter.As:
		return Straight(f, j, o, x)
	case filter.Na:
		return Straight(f, j, o, x)
	case filter.In:
		return Straight(f, j, o, x)
	case filter.Ni:
		return Straight(f, j, o, x)
	default:
		panic(x)
	}
}

func (f *Filter) Alloc(n int) {
	f.v = make([]any, 0, n)
}
