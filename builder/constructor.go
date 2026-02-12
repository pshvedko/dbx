package builder

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pshvedko/dbx/filter"
)

type Order []string

type Ranger struct {
	o *uint
	l *uint
}

type Access struct {
	Group string
	Owner string
}

type Aliases map[rune]int

func (a Aliases) Alias(t string) string {
	var i int
	var r, k rune
	for i, r = range t {
		if i > 0 {
			break
		}
		k = r
	}
	n := a[k]
	a[k]++
	if n == 0 {
		return t[:i]
	}
	return t[:i] + strconv.Itoa(n)
}

type Mode int

func (m Mode) IsModify() bool { return m == 0 }

func (m Mode) IsCreate() bool { return m == 1 }

func (m Mode) IsUpdate() bool { return m == 2 }

type Constructor struct {
	Filter
	Column
	Modify
	Access
	Aliases
	Mode
	p Ranger
	y Order
}

type Counter struct {
	strings.Builder
	q string
	z int
}

func (c *Counter) Count() (string, int, error) {
	c.Grow(256)
	_, err := c.WriteString("SELECT COUNT(*)")
	if err != nil {
		return "", 0, err
	}
	_, err = c.WriteString(c.q)
	if err != nil {
		return "", 0, err
	}
	return c.String(), c.z, nil
}

func (c *Constructor) Select(j filter.Projector, f filter.Filter) (*Counter, string, []any, []any, error) {
	c.Grow(256)
	err := c.Validate(j)
	if err != nil {
		return nil, "", nil, nil, err
	}
	_, err = c.WriteString("SELECT")
	if err != nil {
		return nil, "", nil, nil, err
	}
	var a filter.And
	if f != nil {
		a = append(a, f)
	}
	v, nn, vv, t := 0, j.Names(), j.Values(), c.Alias(j.Table())
	for i, n := range nn {
		if c.IsDeleted(n) {
			a = c.Visibility(a)
		}
		if !c.Used(n) {
			continue
		}
		if v > 0 {
			err = c.WriteByte(',')
			if err != nil {
				return nil, "", nil, nil, err
			}
		}
		_, err = fmt.Fprintf(c, " %q", filter.Column{t, n})
		if err != nil {
			return nil, "", nil, nil, err
		}
		vv[v] = vv[i]
		v++
	}
	n := c.Len()
	_, err = fmt.Fprintf(c, " FROM %q AS %q", j.Table(), t)
	if err != nil {
		return nil, "", nil, nil, err
	}
	_, err = c.WriteString(" WHERE ")
	if err != nil {
		return nil, "", nil, nil, err
	}
	w := c.Len()
	err = a.To(c, filter.Table{Projector: j, Alias: t})
	if err != nil {
		return nil, "", nil, nil, err
	}
	if w == c.Len() {
		_, err = c.WriteString("TRUE")
		if err != nil {
			return nil, "", nil, nil, err
		}
	}
	m := c.Len()
	if len(c.y) > 0 {
		_, err = c.WriteString(" ORDER BY")
		if err != nil {
			return nil, "", nil, nil, err
		}
		for i, y := range c.y {
			if len(y) == 0 {
				continue
			}
			var o string
			switch y[0] {
			case '-':
				o = " DESC"
				fallthrough
			case '+':
				y = y[1:]
			}
			if len(y) == 0 {
				continue
			}
			if i > 0 {
				err = c.WriteByte(',')
				if err != nil {
					return nil, "", nil, nil, err
				}
			}
			_, err = fmt.Fprintf(c, " %q%s", filter.Column{t, y}, o)
			if err != nil {
				return nil, "", nil, nil, err
			}
		}
	}
	z := c.Size()
	if c.p.o != nil {
		_, err = fmt.Fprintf(c, " OFFSET %v", c.Add(*c.p.o))
		if err != nil {
			return nil, "", nil, nil, err
		}
	}
	if c.p.l != nil {
		_, err = fmt.Fprintf(c, " LIMIT %v", c.Add(*c.p.l))
		if err != nil {
			return nil, "", nil, nil, err
		}
	}
	q := c.String()
	if z == c.Size() {
		return nil, q, c.Values(), vv[:v], nil
	}
	return &Counter{q: q[n:m], z: z}, q, c.Values(), vv[:v], nil
}

func (c *Constructor) Range(o, l *uint) *Constructor {
	c.p.o, c.p.l = o, l
	return c
}

func (c *Constructor) Sort(y Order) *Constructor {
	c.y = y
	return c
}

func (c *Constructor) Update(j filter.Projector) (string, []any, []any, error) {
	c.Grow(256)
	err := c.Validate(j)
	if err != nil {
		return "", nil, nil, err
	}
	_, err = c.Printf("UPDATE %q SET", j.Table())
	if err != nil {
		return "", nil, nil, err
	}
	u, nn, vv, pk := 0, j.Names(), j.Values(), j.PK()
	if len(pk) == 0 {
		return "", nil, nil, fmt.Errorf("unknown primary key")
	}
	w := filter.Eq{}
	for i, n := range nn {
		var v fmt.Formatter
		switch {
		case c.IsUpdated(n):
			v = DEFAULT
		case pk.Have(n):
			o, none, auto := j.Value(i)
			if none && auto {
				return "", nil, nil, fmt.Errorf("invalid primary key")
			}
			w[n] = o
			continue
		case c.IsDeleted(n):
			w[n] = nil
			continue
		case c.Unused(n) || c.IsCreated(n):
			continue
		default:
			o, none, auto := j.Value(i)
			if none && auto {
				continue
			}
			v = c.Add(o)
		}
		_, err = c.Printf("%v %q = %v", Comma(u), n, v)
		if err != nil {
			return "", nil, nil, err
		}
		u++
	}
	_, err = c.WriteString(" WHERE ")
	if err != nil {
		return "", nil, nil, err
	}
	err = w.To(c, j)
	if err != nil {
		return "", nil, nil, err
	}
	_, err = c.WriteString(" RETURNING")
	if err != nil {
		return "", nil, nil, err
	}
	for i, n := range nn {
		_, err = c.Printf("%v %q", Comma(i), n)
		if err != nil {
			return "", nil, nil, err
		}
	}
	return c.String(), c.Values(), vv, nil
}

func (c *Constructor) Insert(j filter.Projector) (string, []any, []any, error) {
	c.Grow(256)
	if c.IsUpdate() {
		return c.Update(j)
	}
	err := c.Validate(j)
	if err != nil {
		return "", nil, nil, err
	}
	_, err = c.WriteString("INSERT INTO")
	if err != nil {
		return "", nil, nil, err
	}
	_, err = c.Printf(" %q (", j.Table())
	if err != nil {
		return "", nil, nil, err
	}
	a, nn, vv, pk := 0, j.Names(), j.Values(), j.PK()
	uu := make([]string, 0, len(vv)-len(pk))
	w := filter.And{}
	var up bool
	for i, n := range nn {
		o, none, auto := j.Value(i)
		switch {
		case c.IsUpdated(n):
			up = true
			continue
		case c.IsDeleted(n):
			w = c.Visibility(w)
			continue
		case c.IsCreated(n):
			continue
		case none && auto:
			continue
		case pk.Have(n):
		default:
			uu = append(uu, n)
		}
		_, err = c.Printf("%v %q", Comma(a), n)
		if err != nil {
			return "", nil, nil, err
		}
		c.Add(o)
		a++
	}
	_, err = c.WriteString(" ) VALUES (")
	if err != nil {
		return "", nil, nil, err
	}
	for i := 0; i < a; i++ {
		_, err = c.Printf("%v %v", Comma(i), Holder(i+1))
		if err != nil {
			return "", nil, nil, err
		}
	}
	_, err = c.WriteString(" )")
	if err != nil {
		return "", nil, nil, err
	}
	if c.IsModify() && len(pk) > 0 {
		_, err = c.Printf(" ON CONFLICT ( %v ) DO UPDATE SET", pk)
		if err != nil {
			return "", nil, nil, err
		}
		for i, u := range uu {
			_, err = c.Printf("%v %q = EXCLUDED.%q", Comma(i), u, u)
			if err != nil {
				return "", nil, nil, err
			}
		}
		if up {
			_, err = c.Printf("%v %q = DEFAULT", Comma(a), c.Updated)
			if err != nil {
				return "", nil, nil, err
			}
		}
		if len(w) > 0 {
			_, err = c.WriteString(" WHERE ")
			if err != nil {
				return "", nil, nil, err
			}
			err = w.To(c, j)
		}
	}
	_, err = c.WriteString(" RETURNING")
	if err != nil {
		return "", nil, nil, err
	}
	var v int
	for i, n := range nn {
		if c.Unused(n) {
			continue
		}
		_, err = c.Printf("%v %q", Comma(v), n)
		if err != nil {
			return "", nil, nil, err
		}
		vv[v] = vv[i]
		v++
	}
	if v == 0 {
		_, err = c.Write([]byte{' ', '1'})
		if err != nil {
			return "", nil, nil, err
		}
		vv[0] = new(int64)
		v++
	}
	return c.String(), c.Values(), vv[:v], nil
}

func (c *Constructor) Printf(format string, a ...any) (int, error) {
	return fmt.Fprintf(c, format, a...)
}

func (c *Constructor) Unused(n string) bool {
	return !c.Used(n)
}

func (c *Constructor) Validate(f filter.Fielder) error {
	names := f.Names()
	for name := range c.Column.Names() {
		if !slices.Contains(names, name) {
			return fmt.Errorf("unknown column: %s", name)
		}
	}
	return nil
}
