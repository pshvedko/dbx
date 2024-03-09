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
	w int
	m int
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
		case c.HasDeleted(n):
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
		_, err = fmt.Fprintf(c, " OFFSET %s", c.Hold(*c.p.o))
		if err != nil {
			return nil, "", nil, nil, err
		}
	}
	if c.p.l != nil {
		_, err = fmt.Fprintf(c, " LIMIT %s", c.Hold(*c.p.l))
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

func (c *Constructor) Insert(j filter.Projector) (string, []any, []any, error) {
	c.Grow(256)
	_, err := c.WriteString("INSERT INTO")
	if err != nil {
		return "", nil, nil, err
	}
	_, err = fmt.Fprintf(c, " %q (", j.Table())
	if err != nil {
		return "", nil, nil, err
	}
	v, nn, vv := 0, j.Names(), j.Values()
	a, aa := 0, make([]any, len(vv))
	for i, n := range nn {
		vv[v] = vv[i]
		nn[v] = nn[i]
		v++
		o, none, auto := j.Value(i)
		if none && auto {
			continue
		}
		_, err = c.WithComma(a).Printf(" %q", n)
		if err != nil {
			return "", nil, nil, err
		}
		aa[a] = o
		a++
	}
	_, err = c.WriteString(" ) VALUES (")
	if err != nil {
		return "", nil, nil, err
	}
	for i := range aa[:a] {
		_, err = c.WithComma(i).Printf(" $%d", i+1)
		if err != nil {
			return "", nil, nil, err
		}
	}
	_, err = c.WriteString(" )")
	if err != nil {
		return "", nil, nil, err
	}
	pk := j.PK()
	if len(pk) > 0 {
		_, err = c.Printf(" ON CONFLICT ( %v ) DO UPDATE SET", pk)
		if err != nil {
			return "", nil, nil, err
		}
		var i int
		for _, n := range nn[:v] {
			if pk.Have(n) || !c.Used(n) {
				continue
			}
			_, err = c.WithComma(i).Printf(" %q = EXCLUDED.%q", n, n)
			if err != nil {
				return "", nil, nil, err
			}
			i++
		}
	}
	_, err = c.WriteString(" RETURNING")
	if err != nil {
		return "", nil, nil, err
	}
	for i, n := range nn[:v] {
		_, err = c.WithComma(i).Printf(" %q", n)
		if err != nil {
			return "", nil, nil, err
		}
	}
	return c.String(), aa[:a], vv[:v], nil
}

type Printer interface {
	Printf(string, ...any) (int, error)
}

type ErrPrint struct {
	error
}

func (e ErrPrint) Printf(string, ...any) (int, error) { return 0, e }

func (c *Constructor) WithComma(i int) Printer {
	if i > 0 {
		err := c.WriteByte(',')
		if err != nil {
			return ErrPrint{error: err}
		}
	}
	return c
}

func (c *Constructor) Printf(format string, a ...any) (int, error) {
	return fmt.Fprintf(c, format, a...)
}
