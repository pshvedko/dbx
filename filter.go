package db

import (
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Numbered struct {
	v []any
}

func (n *Numbered) Append(v any) string {
	n.v = append(n.v, v)
	return fmt.Sprintf("$%d", len(n.v))
}

func (n *Numbered) Values() any {
	return n.v
}

type Holder interface {
	Append(v any) string
	Values() any
}

type Builder struct {
	b strings.Builder
	Holder
}

func (b *Builder) Write(p []byte) (int, error) {
	return b.b.Write(p)
}

func (b *Builder) String() string {
	return b.b.String()
}

const (
	Eq = iota
	Ei
	Ne
	Ni
	Ge
	Gt
	Le
	Lt
)

var q = [...]string{
	Eq: "=",
	Ei: "IS",
	Ne: "<>",
	Ni: "IS NOT",
	Ge: ">=",
	Gt: ">",
	Le: "<=",
	Lt: "<",
}

func (b *Builder) Op(k string, o int, v any) error {
	var p string
	switch x := v.(type) {
	case nil:
		p = "NULL"
		o++
	case bool:
		if x {
			p = "TRUE"
		} else {
			p = "FALSE"
		}
		o++
	default:
		p = b.Append(v)
	}
	_, err := fmt.Fprint(&b.b, k, " ", q[o], " ", p)
	return err
}

func (b *Builder) Eq(k string, v any) error {
	return b.Op(k, Eq, v)
}

func (b *Builder) Ne(k string, v any) error {
	return b.Op(k, Ne, v)
}

func (b *Builder) Ge(k string, v any) error {
	return b.Op(k, Ge, v)
}

func (b *Builder) Gt(k string, v any) error {
	return b.Op(k, Gt, v)
}

func (b *Builder) Le(k string, v any) error {
	return b.Op(k, Le, v)
}

func (b *Builder) Lt(k string, v any) error {
	return b.Op(k, Lt, v)
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
