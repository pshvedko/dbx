package builder

import (
	"fmt"
	"github.com/pshvedko/dbx/filter"
)

var (
	andor = [2]byte{'|', '&'}
)

type Key struct {
	Filter
}

func (k *Key) Append(t filter.Type, a any) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (k *Key) AppendInt(i int) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (k *Key) AppendColumn(s string, s2 string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (k *Key) Value(a any) fmt.Formatter {
	//TODO implement me
	panic("implement me")
}

func (k *Key) Straight(j filter.Projector, u uint8, f filter.Filter) error {
	return nil
}

func (k *Key) Conjunct(j filter.Projector, u uint8, ff []filter.Filter) error {
	err := k.WriteByte('[')
	if err != nil {
		return err
	}
	for i, f := range ff {
		if i > 0 {
			err = k.WriteByte(andor[u])
			if err != nil {
				return err
			}
		}
		err = f.To(k, j)
		if err != nil {
			return err
		}
	}
	return k.WriteByte(']')
}

var _ filter.Builder = &Key{}
