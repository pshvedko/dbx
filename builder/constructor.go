package builder

import (
	"fmt"
	"github.com/pshvedko/dbx/filter"
	"strconv"
	"strings"
)

type Table = filter.Table

type Field = filter.Column

type Order []any

func (o Order) Error() string {
	return fmt.Sprintf("invalid order: %s", []any(o))
}

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
	R Ranger
	O Order
	Z bool
}

func (c *Constructor) Printf(format string, a ...any) (int, error) {
	return fmt.Fprintf(c, format, a...)
}

func (c *Constructor) Unused(n string) bool {
	return !c.Used(n)
}

func (c *Constructor) Unreturned(n string) bool {
	return !c.Returned(n)
}

func (c *Constructor) Validate(f filter.Fielder) error {
	columns := f.Columns()
	for _, fields := range c.Column.Names() {
		for name := range fields {
			_, ok := columns[name]
			if !ok {
				return fmt.Errorf("unknown column: %s", name)
			}
		}
	}
	return nil
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
	v, nn, vv, t := 0, j.Names(), j.Places(), c.Alias(j.Table())
	for i, n := range nn {
		if c.IsDeleted(n) {
			a = c.DeleteClause(a)
		}
		if c.Unreturned(n) || c.Unused(n) {
			continue
		}
		if v > 0 {
			err = c.WriteByte(',')
			if err != nil {
				return nil, "", nil, nil, err
			}
		}
		_, err = fmt.Fprintf(c, " %q", Field{t, n})
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
	err = a.To(c, Table{Projector: j, Alias: t})
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
	err = c.WriteOrder(j, t, v)
	if err != nil {
		return nil, "", nil, nil, err
	}

	z := c.Size()
	if c.R.o != nil {
		_, err = fmt.Fprintf(c, " OFFSET %v", c.Add(*c.R.o))
		if err != nil {
			return nil, "", nil, nil, err
		}
	}
	if c.R.l != nil {
		_, err = fmt.Fprintf(c, " LIMIT %v", c.Add(*c.R.l))
		if err != nil {
			return nil, "", nil, nil, err
		}
	}
	q := c.String()
	if c.Z || z == c.Size() {
		return nil, q, c.Places(), vv[:v], nil
	}
	return &Counter{q: q[n:m], z: z}, q, c.Places(), vv[:v], nil
}

func (c *Constructor) WriteOrder(j filter.Projector, t string, v int) error {
	if len(c.O) == 0 {
		return nil
	}
	_, err := c.WriteString(" ORDER BY")
	if err != nil {
		return err
	}
	for i, y := range c.O {
		var f fmt.Formatter
		switch y := y.(type) {
		case int:
			if y == 0 {
				f = filter.Random()
			} else if y > v || y < -v {
				return fmt.Errorf("illegal position: %d", y)
			} else if y < 0 {
				f = By{filter.Special(strconv.Itoa(-y)), DESC}
			} else {
				f = By{filter.Special(strconv.Itoa(y)), ASC}
			}
		case string:
			if len(y) == 0 {
				continue
			}
			var o fmt.Formatter
			switch y[0] {
			case '-':
				o = DESC
				fallthrough
			case '+':
				y = y[1:]
				if len(y) == 0 {
					return c.O
				}
			}
			_, ok := j.Columns()[y]
			if !ok {
				if len(y) == 0 || strings.ContainsFunc(y, func(r rune) bool {
					return r < '0' || r > '9'
				}) {
					return fmt.Errorf("unknown column: %s", y)
				}
				f = By{filter.Special(y), o}
			} else {
				f = By{Field{t, y}, o}
			}
		default:
			return fmt.Errorf("unknown column: %v", y)
		}
		if i > 0 {
			err = c.WriteByte(',')
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintf(c, " %v", f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Constructor) Range(o, l *uint) *Constructor {
	c.R.o, c.R.l = o, l
	return c
}

func (c *Constructor) Sort(y Order) *Constructor {
	c.O = y
	return c
}

func (c *Constructor) Update(j filter.Projector, ff ...filter.Filter) (string, []any, []any, error) {
	c.Grow(256)
	err := c.Validate(j)
	if err != nil {
		return "", nil, nil, err
	}
	t := c.Alias(j.Table())
	_, err = c.Printf("UPDATE %q AS %q SET", j.Table(), t)
	if err != nil {
		return "", nil, nil, err
	}
	u, nn, vv, pk := 0, j.Names(), j.Places(), j.PK()
	if len(ff) == 0 && len(pk) == 0 {
		return "", nil, nil, fmt.Errorf("unknown primary key")
	}
	k := filter.Eq{}
	w := filter.And{k}
	w = append(w, ff...)
	for i, n := range nn {
		var v fmt.Formatter
		o, none, auto := j.Value(i)
		switch {
		case c.IsUpdated(n):
			v = DEFAULT
		case pk.Have(n):
			if none {
				return "", nil, nil, fmt.Errorf("invalid primary key")
			}
			k[n] = o
			continue
		case c.IsDeleted(n):
			w = c.DeleteClause(w)
			if none {
				continue
			}
			v = c.Add(o)
		case c.Unused(n) || c.IsCreated(n):
			continue
		case none && auto:
			continue
		default:
			v = c.Add(o)
		}
		_, err = c.Printf("%v %q = %v", Comma(u), n, v)
		if err != nil {
			return "", nil, nil, err
		}
		u++
	}
	err = c.WriteWhere(j, t, w)
	if err != nil {
		return "", nil, nil, err
	}
	return c.WriteReturning(t, nn, vv)
}

func (c *Constructor) WriteWhere(j filter.Projector, t string, f filter.Filter) error {
	if filter.IsEmpty(f) {
		return nil
	}
	_, err := c.WriteString(" WHERE ")
	if err != nil {
		return err
	}
	err = f.To(c, Table{Projector: j, Alias: t})
	if err != nil {
		return err
	}
	return nil
}

func (c *Constructor) WriteReturning(t string, nn []string, vv []any) (string, []any, []any, error) {
	_, err := c.WriteString(" RETURNING")
	if err != nil {
		return "", nil, nil, err
	}
	var v int
	for i, n := range nn {
		if c.Unreturned(n) {
			continue
		}
		_, err = c.Printf("%v %q", Comma(v), Field{t, n})
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
		vv[v] = new(int64)
		v++
	}
	return c.String(), c.Places(), vv[:v], nil
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
	t := c.Alias(j.Table())
	_, err = c.Printf(" %q AS %q (", j.Table(), t)
	if err != nil {
		return "", nil, nil, err
	}
	a, nn, vv, pk := 0, j.Names(), j.Places(), j.PK()
	uu := make([]string, 0, len(vv)-len(pk))
	w := filter.And{}
	var up string
	for i, n := range nn {
		o, none, auto := j.Value(i)
		switch {
		case c.IsUpdated(n):
			up = n
			continue
		case c.IsDeleted(n):
			w = c.DeleteClause(w)
			if none {
				continue
			}
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
		if len(up) > 0 {
			_, err = c.Printf("%v %q = DEFAULT", Comma(a), up)
			if err != nil {
				return "", nil, nil, err
			}
		}
		err = c.WriteWhere(j, t, w)
		if err != nil {
			return "", nil, nil, err
		}
	}
	return c.WriteReturning(t, nn, vv)
}

type UnusedColumn struct {
	Column
}

func (UnusedColumn) Used(string) bool {
	return false
}

func (c *Constructor) SoftDelete() *Constructor {
	c.Column = UnusedColumn{Column: c.Column}
	return c
}

func (c *Constructor) Delete(j filter.Projector, f filter.Filter) (string, []any, []any, error) {
	c.Grow(256)
	if !c.IsDeleted("") {
		return c.SoftDelete().Update(filter.NewProjector(j).WithPK().WithValue(c.AsDeleted(), filter.Now()), f)
	}
	err := c.Validate(j)
	if err != nil {
		return "", nil, nil, err
	}
	nn, vv, t := j.Names(), j.Places(), c.Alias(j.Table())
	_, err = c.Printf(`DELETE FROM`+` %q AS %q`, j.Table(), t)
	if err != nil {
		return "", nil, nil, err
	}
	w := filter.And{f}
	err = c.WriteWhere(j, t, w)
	if err != nil {
		return "", nil, nil, err
	}
	return c.WriteReturning(t, nn, vv)
}
