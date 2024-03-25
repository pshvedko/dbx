package builder

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pshvedko/dbx/filter"
)

type Permanent struct {
	Filter
}

func (p *Permanent) To(b filter.Builder, j filter.Projector) error {
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
				if i, err := strconv.Atoi(s[1:]); err == nil && i > 0 {
					s = fmt.Sprint('$', i+n)
				}
			}
			_, err = b.Write([]byte{' '})
			_, err = b.WriteString(s)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewPermanent(f filter.Filter, j filter.Projector) (filter.Filter, error) {
	var p Permanent
	err := f.To(&p, j)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
