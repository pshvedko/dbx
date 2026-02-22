package builder

import (
	"github.com/pshvedko/dbx/filter"
)

type Key struct {
	*Constructor
}

func (k *Key) AppendValue(t filter.Type, v any) (n int, err error) {
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

func (k *Key) AppendColumn(t string, c string, m map[string]int) (int, error) {
	return k.AppendInt(m[c])
}

func (k *Key) AppendVia(b bool) (int, error) {
	if b {
		return 1, k.WriteByte('&')
	}
	return 1, k.WriteByte('|')
}

func (k *Key) AppendParenthesis(b bool) (int, error) {
	if b {
		return 1, k.WriteByte('[')
	}
	return 1, k.WriteByte(']')
}

func (k *Key) Straight(j filter.Projector, u bool, x filter.Filter) (err error) {
	return Straight(k, j, u, x)
}

func (k *Key) Conjunct(j filter.Projector, u bool, ff []filter.Filter) (err error) {
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

func (k *Key) WriteRange(i, z int) (n int, err error) {
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

func (k *Key) WriteIndices(j filter.Fielder) (err error) {
	y := true
	z := 0
	for i, n := range j.Names() {
		if !k.Returned(n) {
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

func (k *Key) WriteDeleted(j filter.Fielder) (err error) {
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
	_, err = k.AppendInt(j.Columns()[k.AsDeleted()])
	return
}

func (k *Key) WriteHash(j filter.Projector, f filter.Filter) (n int, err error) {
	err = k.WriteIndices(j)
	if err != nil {
		return
	}
	n = k.Len()
	if f != nil {
		err = f.To(k, j)
		if err != nil {
			return
		}
	}
	err = k.WriteDeleted(j)
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
	Origin
	Static
}

type Keeper interface {
	Get(Hash) (string, bool)
	Put(Hash, string)
}

type Cache struct {
	Origin
	Keeper
}

func (c Cache) IsEnabled() bool {
	return c.Keeper != nil && c.Origin.File != "" && c.Origin.Line != 0
}
