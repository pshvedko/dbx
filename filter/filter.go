package filter

import (
	"fmt"
	"io"
	"sort"
)

type Builder interface {
	io.Writer
	fmt.Stringer
	Eq(string, any) error
	Values() map[string]any
}

type Filter interface {
	To(Builder) error
}

func conjunction(q, o, p string, b Builder, ff []Filter) (err error) {
	_, err = fmt.Fprint(b, q)
	if err != nil {
		return
	}
	for i, f := range ff {
		if i > 0 {
			_, err = fmt.Fprint(b, o)
			if err != nil {
				return
			}
		}
		err = f.To(b)
		if err != nil {
			return
		}
	}
	_, err = fmt.Fprint(b, p)
	return
}

func keys(f map[string]any) (key []string) {
	for k := range f {
		key = append(key, k)
	}
	sort.Strings(key)
	return
}

const (
	eq = iota
)

func straight(q, o, p string, t int, b Builder, f map[string]any) (err error) {
	_, err = fmt.Fprint(b, q)
	if err != nil {
		return
	}
	for i, k := range keys(f) {
		if i > 0 {
			_, err = fmt.Fprint(b, o)
			if err != nil {
				return
			}
		}
		switch t {
		case eq:
			err = b.Eq(k, f[k])
		default:
			return io.EOF
		}
		if err != nil {
			return
		}
	}
	_, err = fmt.Fprint(b, p)
	return
}

type And []Filter

func (f And) To(b Builder) error {
	if len(f) > 1 {
		return conjunction("( ", " AND ", " )", b, f)
	}
	return conjunction("", " AND ", "", b, f)
}

type Or []Filter

func (f Or) To(b Builder) error {
	if len(f) > 1 {
		return conjunction("( ", " OR ", " )", b, f)
	}
	return conjunction("", " OR ", "", b, f)
}

type Eq map[string]any

func (f Eq) To(b Builder) error {
	if len(f) > 1 {
		return straight("( ", " AND ", " )", eq, b, f)
	}
	return straight("", " AND ", "", eq, b, f)
}

type Ge map[string]any

type Gt map[string]any

type Le map[string]any

type Lt map[string]any
