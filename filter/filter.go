package filter

import (
	"fmt"
	"sort"
	"strings"
)

const (
	AND = " AND "
	OR  = " OR "
	SP  = " "
	LB  = "( "
	RB  = " )"
)

type Builder interface {
	Eq(string, any) string
	Values() map[string]any
}

func conjunction(o string, h Builder, a ...Filter) string {
	var s []string
	for _, f := range a {
		s = append(s, f.SQL(h))
	}
	return concat(o, s...)
}

func concat(o string, s ...string) string {
	sort.Strings(s)
	p := strings.Join(s, o)
	if len(s) > 1 {
		return fmt.Sprint(LB, p, RB)
	}
	return p
}

type Filter interface {
	SQL(Builder) string
}

type And []Filter

func (f And) SQL(h Builder) string {
	return conjunction(AND, h, f...)
}

type Or []Filter

func (f Or) SQL(h Builder) string {
	return conjunction(OR, h, f...)
}

type Eq map[string]any

func (f Eq) SQL(h Builder) string {
	var s []string
	for k, v := range f {
		s = append(s, h.Eq(k, v))
	}
	return concat(AND, s...)
}

type Ge map[string]any

type Gt map[string]any

type Le map[string]any

type Lt map[string]any
