package builder

import (
	"fmt"
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

type Constructor struct {
	Filter
	Column
	Modify
	Access
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
	_, err := c.WriteString("SELECT")
	if err != nil {
		return nil, "", nil, nil, err
	}
	a := filter.And{}
	if f != nil {
		a = append(a, f)
	}
	v, nn, vv := 0, j.Names(), j.Values()
	for i, n := range nn {
		switch {
		case c.IsDeleted(n):
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
		_, err = fmt.Fprintf(c, " %q", n)
		if err != nil {
			return nil, "", nil, nil, err
		}
		vv[v] = vv[i]
		v++
	}
	n := c.Len()
	_, err = fmt.Fprintf(c, " FROM %q", j.Table())
	if err != nil {
		return nil, "", nil, nil, err
	}
	_, err = c.WriteString(" WHERE ")
	if err != nil {
		return nil, "", nil, nil, err
	}
	w := c.Len()
	err = a.To(c, j)
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
			_, err = fmt.Fprintf(c, " %q%s", y, o)
			if err != nil {
				return nil, "", nil, nil, err
			}
		}
	}
	z := c.Size()
	if c.p.o != nil {
		_, err = fmt.Fprintf(c, " OFFSET %v", c.Value(*c.p.o))
		if err != nil {
			return nil, "", nil, nil, err
		}
	}
	if c.p.l != nil {
		_, err = fmt.Fprintf(c, " LIMIT %v", c.Value(*c.p.l))
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

func (c *Constructor) Update(j filter.Projector, f filter.Filter) (string, []any, []any, error) {
	c.Grow(256)
	_, err := c.Printf("UPDATE %q SET", j.Table())
	if err != nil {
		return "", nil, nil, err
	}
	u, nn, vv, pk := 0, j.Names(), j.Values(), j.PK()
	if len(pk) == 0 {
		return "", nil, nil, fmt.Errorf("unknown primary key")
	}
	w := filter.Eq{}
	switch f {
	case nil:
		f = w
	default:
		f = filter.And{f, w}
	}
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
			v = c.Value(o)
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
	err = f.To(c, j)
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

func (c *Constructor) Insert(j filter.Projector, m int) (string, []any, []any, error) {
	c.Grow(256)
	if m == 2 {
		return c.Update(j, nil)
	}
	_, err := c.WriteString("INSERT INTO")
	if err != nil {
		return "", nil, nil, err
	}
	_, err = c.Printf(" %q (", j.Table())
	if err != nil {
		return "", nil, nil, err
	}
	a, nn, vv, pk := 0, j.Names(), j.Values(), j.PK()
	uu := make([]string, 0, len(vv)-len(pk))
	for i, n := range nn {
		switch {
		case c.IsUpdated(n):
			uu = append(uu, n)
			fallthrough
		case c.IsCreated(n) || c.IsDeleted(n):
			continue
		case pk.Have(n) || c.Unused(n):
		default:
			uu = append(uu, n)
		}
		o, none, auto := j.Value(i)
		if none && auto {
			continue
		}
		_, err = c.Printf("%v %q", Comma(a), n)
		if err != nil {
			return "", nil, nil, err
		}
		c.Value(o)
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
	if m == 0 && len(pk) > 0 {
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
		if n, ok := c.HaveDeleted(); ok {
			_, err = c.Printf(" WHERE %q.%q IS NULL", j.Table(), n)
			if err != nil {
				return "", nil, nil, err
			}
		}
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

func (c *Constructor) Printf(format string, a ...any) (int, error) {
	return fmt.Fprintf(c, format, a...)
}

func (c *Constructor) Unused(n string) bool {
	return !c.Used(n)
}

func (c *Constructor) HaveDeleted() (string, bool) {
	return c.Deleted.Name()
}
