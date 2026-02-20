package builder

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/pshvedko/dbx/filter"
)

type Comma int

var (
	comma = []byte{','}
	space = []byte{' '}
	dummy = []byte{' ', '1'}
	among = []byte{'"', '.', '"'}
	quote = []byte{'"'}
	place = []byte{'$'}
)

func (c Comma) Format(f fmt.State, _ rune) {
	if c > 0 {
		_, _ = f.Write(comma)
	}
}

type Holder int

func (h Holder) Format(f fmt.State, _ rune) {
	_, _ = h.AppendTo(f)
}

func (h Holder) AppendTo(w io.Writer) (int, error) {
	n1, err := w.Write(place)
	if err != nil {
		return n1, err
	}
	n2, err := Integer(h).AppendTo(w)
	return n1 + n2, err
}

type Integer int

func (i Integer) Format(f fmt.State, _ rune) {
	_, _ = i.AppendTo(f)
}

func (i Integer) AppendTo(w io.Writer) (int, error) {
	if int(i) < len(integers) {
		return w.Write(integers[i])
	}
	var b [20]byte
	return w.Write(strconv.AppendInt(b[:0], int64(i), 10))
}

type Column [2]string

func (c Column) Format(f fmt.State, _ rune) {
	_, _ = c.AppendTo(f)
}

func (c Column) AppendTo(w io.Writer) (int, error) {
	u1, err := w.Write(quote)
	if err != nil {
		return u1, err
	}
	u2, err := io.WriteString(w, c[0])
	if err != nil {
		return u1 + u2, err
	}
	u3, err := w.Write(among)
	if err != nil {
		return u1 + u2 + u3, err
	}
	u4, err := io.WriteString(w, c[1])
	if err != nil {
		return u1 + u2 + u3 + u4, err
	}
	u5, err := w.Write(quote)
	return u1 + u2 + u3 + u4 + u5, err
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

type AppenderTo interface {
	AppendTo(io.Writer) (int, error)
}

type By [2]AppenderTo

func (b By) Format(f fmt.State, _ rune) {
	_, _ = b.AppendTo(f)
}

func (b By) AppendTo(w io.Writer) (int, error) {
	n1, err := b[0].AppendTo(w)
	if err != nil || b[1] == nil {
		return n1, err
	}
	_, err = w.Write(space)
	if err != nil {
		return n1, err
	}
	n2, err := b[1].AppendTo(w)
	return n1 + 1 + n2, err
}

type Builder struct {
	strings.Builder
	v []any
	f []string
	m map[string]int
}

func (f *Builder) Value(v any) fmt.Formatter {
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

func (f *Builder) Add(v any) fmt.Formatter {
	switch x := v.(type) {
	case fmt.Formatter:
		return x
	}
	f.v = append(f.v, v)
	return Holder(len(f.v))
}

func (f *Builder) Size() int {
	return len(f.v)
}

func (f *Builder) Values() []any {
	return f.v
}

func (f *Builder) Alloc(n int) {
	f.v = make([]any, 0, 8+n)
	f.f = make([]string, 0, 8+n>>1)
}

var (
	poolErrNoSuchField = filter.Pool[map[string]int]{New: func() any { return make(map[string]int, 32) }}
)

func (f *Builder) Supply() (map[string]int, []string) {
	m := poolErrNoSuchField.Get()
	return m, f.f[:0]
}

func (f *Builder) Reuse(m map[string]int, v []string) {
	poolErrNoSuchField.Put(m)
	f.f = v
}

func (f *Builder) Width() (int, bool) {
	return 0, false
}

func (f *Builder) Precision() (int, bool) {
	return 0, false
}

func (f *Builder) Flag(int) bool {
	return false
}

func (f *Builder) AppendInt(i int) (int, error) {
	return Integer(i).AppendTo(f)
}

func (f *Builder) AppendColumn(t string, c string, m map[string]int) (int, error) {
	return Column{t, c}.AppendTo(f)
}

func (f *Builder) AppendVia(b bool) (int, error) {
	if b {
		return f.WriteString(" AND ")
	}
	return f.WriteString(" OR ")
}

func (f *Builder) AppendParenthesis(b bool) (int, error) {
	if b {
		return f.WriteString("( ")
	}
	return f.WriteString(" )")
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

func (f *Builder) EQ(v fmt.Formatter) (int, error) {
	n1, err := f.WriteString(" = ")
	if err != nil {
		return n1, err
	}
	n2, err := f.AppendFormat(v)
	return n1 + n2, err
}

func (f *Builder) IN(v fmt.Formatter) (int, error) {
	n1, err := f.WriteString(" = ANY(")
	if err != nil {
		return n1, err
	}
	n2, err := f.AppendFormat(v)
	if err != nil {
		return n1 + n2, err
	}
	err = f.WriteByte(')')
	return n1 + n2 + 1, err
}

func (f *Builder) AppendValue(t filter.Type, v any) (int, error) {
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
		return f.AppendFormat(f.Value(v))
	case filter.AND, filter.OR:
		fallthrough
	default:
		panic(t)
	}
}

func (f *Builder) Conjunct(j filter.Projector, u bool, ff []filter.Filter) error {
	return Conjunct(f, j, u, ff)
}

func (f *Builder) Straight(j filter.Projector, u bool, x filter.Filter) error {
	return Straight(f, j, u, x)
}

func (f *Builder) AppendFormat(x fmt.Formatter) (int, error) {
	switch x := x.(type) {
	case AppenderTo:
		return x.AppendTo(f)
	}
	return fmt.Fprint(f, x)
}
