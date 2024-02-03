package filter

import (
	"fmt"
	"io"
	"sort"
)

type Builder interface {
	io.Writer
	Eq(string, any) error
	Ne(string, any) error
	Ge(string, any) error
	Gt(string, any) error
	Le(string, any) error
	Lt(string, any) error
	As(string, any) error
	In(string, ...any) error
	Ni(string, ...any) error
}

type Filter interface {
	To(Builder) error
}

func conjunction(o string, b Builder, ff []Filter) (err error) {
	var q, p string
	if len(ff) > 1 {
		q, p = "( ", " )"
	}
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

func straight(o string, b Builder, f map[string]any, t any) (err error) {
	var q, p string
	if len(f) > 1 {
		q, p = "( ", " )"
	}
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
		switch t.(type) {
		case Eq:
			err = b.Eq(k, f[k])
		case Ne:
			err = b.Ne(k, f[k])
		case Ge:
			err = b.Ge(k, f[k])
		case Gt:
			err = b.Gt(k, f[k])
		case Le:
			err = b.Le(k, f[k])
		case Lt:
			err = b.Lt(k, f[k])
		case As:
			err = b.As(k, f[k])
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
	return conjunction(" AND ", b, f)
}

type Or []Filter

func (f Or) To(b Builder) error {
	return conjunction(" OR ", b, f)
}

type Eq map[string]any

func (f Eq) To(b Builder) error {
	return straight(" AND ", b, f, f)
}

type Ne map[string]any

func (f Ne) To(b Builder) error {
	return straight(" AND ", b, f, f)
}

type Ge map[string]any

func (f Ge) To(b Builder) error {
	return straight(" AND ", b, f, f)
}

type Gt map[string]any

func (f Gt) To(b Builder) error {
	return straight(" AND ", b, f, f)
}

type Le map[string]any

func (f Le) To(b Builder) error {
	return straight(" AND ", b, f, f)
}

type Lt map[string]any

func (f Lt) To(b Builder) error {
	return straight(" AND ", b, f, f)
}

type As map[string]any

func (f As) To(b Builder) error {
	return straight(" AND ", b, f, f)
}

type In map[string][]any

func (f In) To(b Builder) error {
	return io.EOF
}

type Ni map[string][]any

func (f Ni) To(b Builder) error {
	return io.EOF
}
