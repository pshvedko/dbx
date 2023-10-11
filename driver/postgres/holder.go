package postgres

import (
	"fmt"
	"strings"
)

type Builder struct {
	strings.Builder
	k map[string]int
	v map[string]any
}

func NewBuilder() *Builder {
	return &Builder{k: map[string]int{}, v: map[string]any{}}
}

func (b Builder) Values() map[string]any {
	return b.v
}

func boolean(x bool) string {
	if x {
		return "TRUE"
	}
	return "FALSE"
}

func (b *Builder) Operation(k string, v any, o ...string) (err error) {
	switch x := v.(type) {
	case nil:
		_, err = fmt.Fprint(b, k, o[0], "NULL")
	case bool:
		_, err = fmt.Fprint(b, k, o[0], boolean(x))
	default:
		p := fmt.Sprint(k, b.k[k])
		b.k[k]++
		b.v[p] = v
		_, err = fmt.Fprint(b, k, o[len(o)-1], ":", p)
	}
	return
}

func (b *Builder) Eq(k string, v any) (err error) {
	return b.Operation(k, v, " IS ", " = ")
}

func (b *Builder) Ne(k string, v any) (err error) {
	return b.Operation(k, v, " IS NOT ", " <> ")
}

func (b *Builder) Ge(k string, v any) (err error) {
	return b.Operation(k, v, " >= ")
}

func (b *Builder) Gt(k string, v any) (err error) {
	return b.Operation(k, v, " > ")
}

func (b *Builder) Le(k string, v any) (err error) {
	return b.Operation(k, v, " <= ")
}

func (b *Builder) Lt(k string, v any) (err error) {
	return b.Operation(k, v, " < ")
}

func (b *Builder) In(s string, a ...any) error {
	//TODO implement me
	panic("implement me")
}

func (b *Builder) Ni(s string, a ...any) error {
	//TODO implement me
	panic("implement me")
}

func (b Builder) As(s string, a any) error {
	//TODO implement me
	panic("implement me")
}
