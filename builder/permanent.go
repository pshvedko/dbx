package builder

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pshvedko/dbx/filter"
)

type Permanent struct {
	I string
	Builder
}

func (p *Permanent) To(b filter.Builder, j filter.Projector) error {
	if p.I != j.Name() {
		return fmt.Errorf("permanent filter mismatch: %q <> %q", p.I, j.Table())
	}
	n := b.Size()
	for _, v := range p.Values() {
		b.Value(v)
	}
	if n == 0 {
		_, err := b.WriteString(p.String())
		return err
	}
	w := strings.Fields(p.String())
	if len(w) > 0 {
		_, err := b.WriteString(w[0])
		if err != nil {
			return err
		}
		for _, s := range w[1:] {
			if len(s) > 0 && s[0] == '$' {
				i, err := strconv.Atoi(s[1:])
				if err == nil && i > 0 {
					_, err = b.WriteString(" $")
					if err != nil {
						return err
					}
					_, err = b.AppendInt(i + n)
					if err != nil {
						return err
					}
					continue
				}
			}
			err = b.WriteByte(' ')
			if err != nil {
				return err
			}
			_, err = b.WriteString(s)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewPermanent(j filter.Projector, f filter.Filter) (filter.Filter, error) {
	p := Permanent{I: j.Name()}
	err := f.To(&p, j)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func NewBuilder() filter.Builder {
	return &Builder{}
}
