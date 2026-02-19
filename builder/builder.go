package builder

import (
	"fmt"
	"github.com/pshvedko/dbx/filter"
)

func Conjunct(b filter.Builder, j filter.Projector, u bool, ff []filter.Filter) (err error) {
	if len(ff) > 1 {
		_, err = b.AppendParenthesis(true)
		if err != nil {
			return
		}
	}
	var i int
	for _, f := range ff {
		if filter.IsEmpty(f) {
			continue
		}
		if i > 0 {
			_, err = b.AppendVia(u)
			if err != nil {
				return
			}
		}
		err = f.To(b, j)
		if err != nil {
			return
		}
		i++
	}
	if len(ff) > 1 {
		if i == 0 {
			_, err = b.WriteString("TRUE")
			if err != nil {
				return
			}
		}
		_, err = b.AppendParenthesis(false)
	}
	return
}

type ErrNoSuchField map[string]int

func (e ErrNoSuchField) Error() string {
	return fmt.Sprintf("no such field: %v", map[string]int(e))
}

var (
	poolErrNoSuchField = filter.Pool[ErrNoSuchField]{New: func() any { return make(ErrNoSuchField, 32) }}
)

func StraightTo[T any, M interface {
	~map[string]T
	Type() filter.Type
}](b filter.Builder, j filter.Projector, u bool, oo M) (err error) {
	nn, ff := b.Supply()
	for k := range oo {
		nn[k]++
	}
	for _, k := range j.Names() {
		_, ok := oo[k]
		if ok {
			ff = append(ff, k)
			delete(nn, k)
		}
	}
	if len(nn) > 0 {
		return ErrNoSuchField(nn)
	}
	b.Reuse(nn, ff)
	if len(oo) > 1 {
		_, err = b.AppendParenthesis(true)
		if err != nil {
			return
		}
	}
	t := j.Table()
	for i, f := range ff {
		if i > 0 {
			_, err = b.AppendVia(u)
			if err != nil {
				return
			}
		}
		_, err = b.AppendColumn(t, f, j.Columns())
		if err != nil {
			return
		}
		_, err = b.AppendValue(oo.Type(), oo[f])
		if err != nil {
			return
		}
	}
	if len(oo) > 1 {
		_, err = b.AppendParenthesis(false)
	}
	return
}

func Straight(b filter.Builder, j filter.Projector, u bool, f filter.Filter) error {
	switch x := f.(type) {
	case filter.Eq:
		return StraightTo(b, j, u, x)
	case filter.Ne:
		return StraightTo(b, j, u, x)
	case filter.Ge:
		return StraightTo(b, j, u, x)
	case filter.Gt:
		return StraightTo(b, j, u, x)
	case filter.Le:
		return StraightTo(b, j, u, x)
	case filter.Lt:
		return StraightTo(b, j, u, x)
	case filter.As:
		return StraightTo(b, j, u, x)
	case filter.Na:
		return StraightTo(b, j, u, x)
	case filter.In:
		return StraightTo(b, j, u, x)
	case filter.Ni:
		return StraightTo(b, j, u, x)
	default:
		panic(x)
	}
}
