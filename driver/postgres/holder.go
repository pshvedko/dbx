package postgres

import "fmt"

type Holder struct {
	k map[string]int
	v map[string]any
}

func (h Holder) Values() map[string]any {
	return h.v
}

func (h Holder) Eq(k string, v any) string {
	p := fmt.Sprint(k, h.k[k])
	h.k[k]++
	switch x := v.(type) {
	case nil:
		return fmt.Sprint(k, " IS NULL")
	case bool:
		return fmt.Sprint(k, " IS ", func(x bool) string {
			if x {
				return "TRUE"
			}
			return "FALSE"
		}(x))
	}
	h.v[p] = v
	return fmt.Sprint(k, " = :", p)
}

func NewHolder() Holder {
	return Holder{k: map[string]int{}, v: map[string]any{}}
}
