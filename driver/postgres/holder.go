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

func (b Builder) Values() map[string]any {
	return b.v
}

func (b *Builder) Eq(k string, v any) (err error) {
	switch x := v.(type) {
	case nil:
		_, err = fmt.Fprint(b, k, " IS NULL")
	case bool:
		_, err = fmt.Fprint(b, k, func() string {
			if x {
				return " IS TRUE"
			}
			return " IS FALSE"
		}())
	default:
		p := fmt.Sprint(k, b.k[k])
		b.k[k]++
		b.v[p] = v
		_, err = fmt.Fprint(b, k, " = :", p)
	}
	return
}

func NewBuilder() *Builder {
	return &Builder{k: map[string]int{}, v: map[string]any{}}
}
