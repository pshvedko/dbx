package filter

import (
	"fmt"
	"io"
)

func Conjunction(b Builder, j Projector, o string, ff []Filter) (err error) {
	if len(ff) > 1 {
		_, err = fmt.Fprint(b, "( ")
		if err != nil {
			return
		}
		defer func() {
			if err == nil {
				_, err = fmt.Fprint(b, " )")
			}
		}()
	}
	for i, f := range ff {
		if i > 0 {
			_, err = fmt.Fprint(b, " ", o, " ")
			if err != nil {
				return
			}
		}
		err = f.To(b, j)
		if err != nil {
			return
		}
	}
	return
}

func Straight[T any](b Builder, j Projector, o string, oo map[string]T, v any) (err error) {
	ff := make([]string, 0, len(oo))
	for _, k := range j.Names() {
		_, ok := oo[k]
		if ok {
			ff = append(ff, k)
		}
	}
	if len(ff) != cap(ff) {
		return io.EOF
	}
	if len(oo) > 1 {
		_, err = fmt.Fprint(b, "( ")
		if err != nil {
			return
		}
		defer func() {
			if err == nil {
				_, err = fmt.Fprint(b, " )")
			}
		}()
	}
	t := j.Table()
	for i, f := range ff {
		if i > 0 {
			_, err = fmt.Fprint(b, " ", o, " ")
			if err != nil {
				return
			}
		}
		switch v.(type) {
		case Eq:
			err = b.Eq(Column{t, f}, oo[f])
		case Ne:
			err = b.Ne(Column{t, f}, oo[f])
		case Ge:
			err = b.Ge(Column{t, f}, oo[f])
		case Gt:
			err = b.Gt(Column{t, f}, oo[f])
		case Le:
			err = b.Le(Column{t, f}, oo[f])
		case Lt:
			err = b.Lt(Column{t, f}, oo[f])
		case As:
			err = b.As(Column{t, f}, oo[f])
		case In:
			err = b.In(Column{t, f}, oo[f])
		case Ni:
			err = b.Ni(Column{t, f}, oo[f])
		default:
			return io.EOF
		}
		if err != nil {
			return
		}
	}
	return
}

func Nil[T comparable](v T) any {
	var z T
	if v != z {
		return v
	}
	return nil
}

type Injectable[T Projector] []T

func (o Injectable[T]) Get() Projector {
	var x T
	return x.Copy()
}

func (o *Injectable[T]) Put(j Projector) {
	switch v := j.(type) {
	case T:
		x := v.Copy()
		switch t := x.(type) {
		case T:
			*o = append(*o, t)
		default:
			panic("invalid copy")
		}
	default:
		panic("invalid injection")
	}
}
