package builder

import (
	"errors"
	"github.com/pshvedko/dbx/filter"
)

type Key struct {
	*Constructor
}

func (k Key) AppendValue(t filter.Type, v any) (n int, err error) {
	var u int
	switch t {
	case filter.TRUE, filter.FALSE:
	default:
		n, err = k.WriteString(t.String())
		if err != nil {
			return
		}
	}
	z := k.Size()
	f := k.Value(v)
	switch x := f.(type) {
	case nil:
	case Keyword:
		err = k.WriteByte(x[0])
		n++
	default:
		if z == k.Size() {
			u, err = k.AppendFormat(x)
			n += u
		}
	}
	return
}

func (k Key) AppendColumn(t string, n string, j filter.Projector) (int, error) {
	return k.AppendInt(j.Index(t, n))
}

func (k Key) AppendVia(b bool) (int, error) {
	if b {
		return 1, k.WriteByte('&')
	}
	return 1, k.WriteByte('|')
}

func (k Key) AppendParenthesis(b bool) (int, error) {
	if b {
		return 1, k.WriteByte('[')
	}
	return 1, k.WriteByte(']')
}

func (k Key) Straight(j filter.Projector, u bool, x filter.Filter) (err error) {
	return Straight(k, j, u, x)
}

func (k Key) Conjunct(j filter.Projector, u bool, ff []filter.Filter) (err error) {
	_, err = k.AppendParenthesis(true)
	if err != nil {
		return
	}
	for i, f := range ff {
		if i > 0 {
			_, err = k.AppendVia(u)
			if err != nil {
				return
			}
		}
		err = f.To(k, j)
		if err != nil {
			return
		}
	}
	_, err = k.AppendParenthesis(false)
	return
}

func (k Key) WriteRange(i, z int) (n int, err error) {
	b := byte('-')
	switch i - z {
	case 0:
	case 1:
		b = ','
		fallthrough
	default:
		err = k.WriteByte(b)
		if err != nil {
			return
		}
		n, err = k.AppendInt(i - 1)
		n++
	}
	return
}

func (k Key) WriteIndices(j filter.Projector) (err error) {
	y := true
	z := 0
	t := k.Alias(0, j.Table())
	for i, n := range j.Names() {
		if k.Unreturned(t, n) {
			if !y {
				_, err = k.WriteRange(i, z)
				if err != nil {
					return
				}
				y = true
			}
			continue
		}
		if y {
			if z > 0 {
				err = k.WriteByte(',')
				if err != nil {
					return
				}
			}
			_, err = k.AppendInt(i)
			if err != nil {
				return
			}
			z = i + 1
			y = false
		}
	}
	if !y {
		_, err = k.WriteRange(j.Len(), z)
	}
	return
}

func (k Key) WriteDeleted(j filter.Projector) (err error) {
	b := byte('+')
	switch k.Deleted.(type) {
	case DeletedFree, nil:
		return
	case DeletedNone:
	case DeletedOnly:
		b = '-'
	}
	err = k.WriteByte(b)
	if err != nil {
		return
	}
	_, err = k.AppendInt(j.Index(j.Table(), k.AsDeleted()))
	return
}

func (k Key) WriteHash(j filter.Projector, f filter.Filter) (int, error) {
	err := k.WriteIndices(j)
	if err != nil {
		return 0, err
	}
	n := k.Len()
	if f != nil {
		err = f.To(k, j)
		if err != nil {
			return 0, err
		}
	}
	err = k.WriteDeleted(j)
	if err != nil {
		return 0, err
	}
	return n, k.WriteOrderOffsetLimit(j)
}

func (k Key) WriteOrderOffsetLimit(j filter.Projector) (err error) {
	var z byte
	for _, o := range k.O {
		switch x := o.(type) {
		case int:
			if x < 0 {
				x = -x
				z = '-'
			} else {
				z = '+'
			}
			err = k.WriteByte(z)
			if err != nil {
				return err
			}
			_, err = k.AppendInt(x)
		case string:
			switch x[0] {
			case '+', '-':
				x = x[1:]
				z = x[0]
			default:
				z = '+'
			}
			err = k.WriteByte(z)
			if err != nil {
				return err
			}
			t, x, _ := k.Regular(x)
			_, err = k.AppendInt(j.Index(t, x))
		}
	}
	if err != nil {
		return err
	}
	if k.R.O != nil {
		err = k.WriteByte('O')
		if err != nil {
			return
		}
		k.Add(*k.R.O)
	}
	if k.R.L != nil {
		err = k.WriteByte('L')
		if err != nil {
			return
		}
		k.Add(*k.R.L)
	}
	return
}

type Origin struct {
	File string `json:"file,omitempty"`
	Line int    `json:"line,omitempty"`
}

const (
	SELECT byte = iota
	INSERT
	UPDATE
	DELETE
)

type Static struct {
	Type byte   `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	List string `json:"list,omitempty"`
	Sign string `json:"sign,omitempty"`
}

type Hash struct {
	Origin `json:"origin"`
	Static `json:"static"`
}

var (
	ErrInvalidRange = errors.New("invalid range")
)

func (s Static) Places(j filter.Placer) ([]any, error) {
	v, i, e, r, l, p := 0, 0, 0, -1, s.List, j.Places()
	for k := 0; k <= len(l); k++ {
		var b byte
		if k == len(l) {
			if k == 0 {
				continue
			}
			b = ','
		} else {
			b = l[k]
		}
		switch {
		case b >= '0' && b <= '9':
			e *= 10
			e += int(b - '0')
		case b == '-':
			if r != -1 {
				return nil, ErrInvalidRange
			}
			r = e
			e = 0
		case b == ',':
			if i > e {
				return nil, ErrInvalidRange
			} else if r == -1 {
				r = e
			} else if r >= e {
				return nil, ErrInvalidRange
			}
			for r <= e {
				p[v] = p[r]
				v++
				r++
			}
			i = r
			r = -1
			e = 0
		default:
			return nil, ErrInvalidRange
		}
	}
	return p[:v], nil
}

type Query struct {
	Data string `json:"data,omitempty"`
	Size [4]int `json:"size,omitempty"`
}

type Keeper interface {
	Get(Hash) (Query, bool)
	Put(Hash, Query)
}

type Cache struct {
	Origin
	Keeper
}

func (c Cache) IsEnabled() bool {
	return c.Keeper != nil && c.Origin.File != "" && c.Origin.Line != 0
}
