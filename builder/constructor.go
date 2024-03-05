package builder

import (
	"fmt"

	"github.com/pshvedko/db/filter"
)

type Constructor struct {
	Filter
}

func (c Constructor) Select(p filter.Projector, f filter.Filter) (string, []any, []any, error) {
	_, err := c.WriteString("SELECT")
	if err != nil {
		return "", nil, nil, err
	}
	j := 0
	nn := p.Names()
	vv := p.Values()
	for i, n := range nn {
		if i > 0 {
			err = c.WriteByte(',')
			if err != nil {
				return "", nil, nil, err
			}
		}
		_, err = fmt.Fprintf(&c, " %q", n)
		if err != nil {
			return "", nil, nil, err
		}
		vv[j] = vv[i]
		j++
	}
	{
		_, err = fmt.Fprintf(&c, " FROM %q", p.Table())
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
			err = f.To(&c, p)
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
	return c.String(), c.Values(), vv[:j], nil
}
