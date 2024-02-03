package db

import (
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Builder struct {
	b strings.Builder
	v []any
}

func (b *Builder) Write(p []byte) (int, error) {
	return b.b.Write(p)
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
		_, err = fmt.Fprint(&b.b, k, o[0], "NULL")
	case bool:
		_, err = fmt.Fprint(&b.b, k, o[0], boolean(x))
	default:
		b.v = append(b.v, v)
		_, err = fmt.Fprint(&b.b, k, o[len(o)-1], "$", len(b.v))
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

func (b *Builder) As(s string, a any) error {
	//TODO implement me
	panic("implement me")
}
