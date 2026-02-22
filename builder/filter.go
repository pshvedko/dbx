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
}

func (b *Builder) IsTrusted() bool { return false }

func (b *Builder) Value(v any) fmt.Formatter {
	switch x := v.(type) {
	case nil:
		return NULL
	case bool:
		if x {
			return TRUE
		}
		return FALSE
	default:
		return b.Add(v)
	}
}

func (b *Builder) Add(v any) fmt.Formatter {
	switch x := v.(type) {
	case fmt.Formatter:
		return x
	}
	b.v = append(b.v, v)
	return Holder(len(b.v))
}

func (b *Builder) Size() int {
	return len(b.v)
}

func (b *Builder) Values() []any {
	return b.v
}

func (b *Builder) Alloc(n int) {
	b.v = make([]any, 0, n)
}

func (b *Builder) Width() (int, bool) {
	return 0, false
}

func (b *Builder) Precision() (int, bool) {
	return 0, false
}

func (b *Builder) Flag(int) bool {
	return false
}

func (b *Builder) AppendInt(i int) (int, error) {
	return Integer(i).AppendTo(b)
}

func (b *Builder) AppendColumn(t string, c string, m map[string]int) (int, error) {
	return Column{t, c}.AppendTo(b)
}

func (b *Builder) AppendVia(u bool) (int, error) {
	if u {
		return b.WriteString(" AND ")
	}
	return b.WriteString(" OR ")
}

func (b *Builder) AppendParenthesis(u bool) (int, error) {
	if u {
		return b.WriteString("( ")
	}
	return b.WriteString(" )")
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

func (b *Builder) EQ(v fmt.Formatter) (int, error) {
	n1, err := b.WriteString(" = ")
	if err != nil {
		return n1, err
	}
	n2, err := b.AppendFormat(v)
	return n1 + n2, err
}

func (b *Builder) IN(v fmt.Formatter) (int, error) {
	n1, err := b.WriteString(" = ANY(")
	if err != nil {
		return n1, err
	}
	n2, err := b.AppendFormat(v)
	if err != nil {
		return n1 + n2, err
	}
	err = b.WriteByte(')')
	return n1 + n2 + 1, err
}

func (b *Builder) AppendValue(t filter.Type, v any) (int, error) {
	switch t {
	case filter.EQ:
		switch v.(type) {
		case nil, bool:
			return fmt.Fprintf(b, Is, b.Value(v))
		}
		return b.EQ(b.Value(v))
	case filter.NE:
		switch v.(type) {
		case nil, bool:
			return fmt.Fprintf(b, Si, b.Value(v))
		}
		return fmt.Fprintf(b, Ne, b.Value(v))
	case filter.GE:
		return fmt.Fprintf(b, Ge, b.Value(v))
	case filter.GT:
		return fmt.Fprintf(b, Gt, b.Value(v))
	case filter.LE:
		return fmt.Fprintf(b, Le, b.Value(v))
	case filter.LT:
		return fmt.Fprintf(b, Lt, b.Value(v))
	case filter.AS:
		return fmt.Fprintf(b, As, b.Value(v))
	case filter.NA:
		return fmt.Fprintf(b, Na, b.Value(v))
	case filter.IN:
		return b.IN(b.Value(v))
	case filter.NI:
		return fmt.Fprintf(b, Ni, b.Value(v))
	case filter.FALSE, filter.TRUE:
		return b.AppendFormat(b.Value(v))
	case filter.AND, filter.OR:
		fallthrough
	default:
		panic(t)
	}
}

func (b *Builder) Conjunct(j filter.Projector, u bool, ff []filter.Filter) error {
	return Conjunct(b, j, u, ff)
}

func (b *Builder) Straight(j filter.Projector, u bool, x filter.Filter) error {
	return Straight(b, j, u, x)
}

func (b *Builder) AppendFormat(x fmt.Formatter) (int, error) {
	switch x := x.(type) {
	case AppenderTo:
		return x.AppendTo(b)
	}
	return fmt.Fprint(b, x)
}
