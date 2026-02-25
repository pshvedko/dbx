package builder

import "github.com/pshvedko/dbx/filter"

type Joiner []func() error

func (j Joiner) Join() error {
	i := len(j)
	for i > 0 {
		i--
		err := j[i]()
		if err != nil {
			return err
		}
	}
	return nil
}

type Join struct {
	*Constructor
	Alias string
}

func (b Join) Conjunct(j filter.Projector, u bool, ff []filter.Filter) error {
	return Conjunct(b, j, u, ff)
}

func (b Join) Straight(j filter.Projector, u bool, x filter.Filter) error {
	return Straight(b, j, u, x)
}

func (b Join) AppendValue(t filter.Type, v any) (int, error) {
	switch v := v.(type) {
	case filter.Column:
		return b.Constructor.Builder.AppendValue(t, Column{b.Alias, string(v)})
	default:
		return b.Constructor.Builder.AppendValue(t, v)
	}
}
