package builder

import (
	"fmt"

	"github.com/pshvedko/db/filter"
)

type Column interface {
	Used(string) bool
}

type AllowedColumn map[string]struct{}

func (f AllowedColumn) Used(k string) bool {
	_, ok := f[k]
	return ok
}

type ExcludedColumn map[string]struct{}

func (f ExcludedColumn) Used(k string) bool {
	_, ok := f[k]
	return !ok
}

type Order []string

type Ranger struct {
	o *uint
	l *uint
}

type Constructor struct {
	Filter
	Column
	p Ranger
	y Order
}

func (c Constructor) Select(j filter.Projector, f filter.Filter) (string, []any, []any, error) {
	_, err := c.WriteString("SELECT")
	if err != nil {
		return "", nil, nil, err
	}
	k := 0
	nn := j.Names()
	vv := j.Values()
	for i, n := range nn {
		if !c.Used(n) {
			continue
		}
		if k > 0 {
			err = c.WriteByte(',')
			if err != nil {
				return "", nil, nil, err
			}
		}
		_, err = fmt.Fprintf(&c, " %q", n)
		if err != nil {
			return "", nil, nil, err
		}
		vv[k] = vv[i]
		k++
	}
	{
		_, err = fmt.Fprintf(&c, " FROM %q", j.Table())
		if err != nil {
			return "", nil, nil, err
		}
	}
	{
		_, err = fmt.Fprintf(&c, " WHERE ")
		if err != nil {
			return "", nil, nil, err
		}
		n := c.Len()
		if f != nil {
			err = f.To(&c, j)
			if err != nil {
				return "", nil, nil, err
			}
		}
		if n == c.Len() {
			_, err = fmt.Fprintf(&c, "TRUE")
			if err != nil {
				return "", nil, nil, err
			}
		}
	}
	{
		if c.p.o != nil {
			_, err = fmt.Fprintf(&c, " OFFSET %d", *c.p.o)
			if err != nil {
				return "", nil, nil, err
			}
		}
		if c.p.l != nil {
			_, err = fmt.Fprintf(&c, " LIMIT %d", *c.p.l)
			if err != nil {
				return "", nil, nil, err
			}
		}
	}
	if len(c.y) > 0 {
		_, err = fmt.Fprintf(&c, " ORDER BY")
		if err != nil {
			return "", nil, nil, err
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
					return "", nil, nil, err
				}
			}
			_, err = fmt.Fprintf(&c, " %q%s", y, o)
			if err != nil {
				return "", nil, nil, err
			}
		}
	}
	return c.String(), c.Values(), vv[:k], nil
}

func (c Constructor) Range(o, l *uint) Constructor {
	c.p.o, c.p.l = o, l
	return c
}

func (c Constructor) Sort(y Order) Constructor {
	c.y = y
	return c
}
