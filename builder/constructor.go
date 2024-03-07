package builder

import (
	"fmt"
	"strings"

	"github.com/pshvedko/dbx/filter"
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

type Option struct {
	Created string
	Updated string
	Deleted string
}

type Constructor struct {
	Filter
	Column
	Option
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
	_, err := c.WriteString("SELECT")
	if err != nil {
		return nil, "", nil, nil, err
	}

	a := filter.And{}
	if f != nil {
		a = append(a, f)
	}

	k := 0
	nn := j.Names()
	vv := j.Values()
	for i, n := range nn {
		switch n {
		case c.Deleted:
			a = append(a, filter.Eq{n: nil})
		}
		if !c.Used(n) {
			continue
		}
		if k > 0 {
			err = c.WriteByte(',')
			if err != nil {
				return nil, "", nil, nil, err
			}
		}
		_, err = fmt.Fprintf(c, " %q", n)
		if err != nil {
			return nil, "", nil, nil, err
		}
		vv[k] = vv[i]
		k++
	}

	n := c.Len()

	_, err = fmt.Fprintf(c, " FROM %q", j.Table())
	if err != nil {
		return nil, "", nil, nil, err
	}

	_, err = fmt.Fprintf(c, " WHERE ")
	if err != nil {
		return nil, "", nil, nil, err
	}

	w := c.Len()

	err = a.To(c, j)
	if err != nil {
		return nil, "", nil, nil, err
	}
	if w == c.Len() {
		_, err = fmt.Fprintf(c, "TRUE")
		if err != nil {
			return nil, "", nil, nil, err
		}
	}

	m := c.Len()

	if len(c.y) > 0 {
		_, err = fmt.Fprintf(c, " ORDER BY")
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
		return nil, q, c.Values(), vv[:k], nil
	}

	return &Counter{q: q[n:m], z: z}, q, c.Values(), vv[:k], nil
}

func (c *Constructor) Range(o, l *uint) *Constructor {
	c.p.o, c.p.l = o, l
	return c
}

func (c *Constructor) Sort(y Order) *Constructor {
	c.y = y
	return c
}
