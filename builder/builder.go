package builder

import (
	"fmt"

	"github.com/pshvedko/dbx/filter"
)

func Conjunct(b filter.Builder, j filter.Projector, o string, ff []filter.Filter) (err error) {
	if len(ff) > 1 {
		_, err = b.WriteString("( ")
		if err != nil {
			return
		}
		defer func() {
			if err == nil {
				_, err = b.WriteString(" )")
			}
		}()
	}
	var i int
	for _, f := range ff {
		if filter.IsEmpty(f) {
			continue
		}
		if i > 0 {
			err = b.WriteByte(' ')
			if err != nil {
				return
			}
			_, err = b.WriteString(o)
			if err != nil {
				return
			}
			err = b.WriteByte(' ')
			if err != nil {
				return
			}
		}
		err = f.To(b, j)
		if err != nil {
			return
		}
		i++
	}
	return
}

type ErrNoSuchField map[string]struct{}

func (e ErrNoSuchField) Error() string {
	ff := make([]string, 0, len(e))
	for f := range e {
		ff = append(ff, f)
	}
	return fmt.Sprintln("no such field:", ff)
}

type Field []string

func (e Field) Erase() Field {
	clear(e)
	return e[:0]
}

var (
	poolErrNoSuchField = filter.Pool[ErrNoSuchField]{New: func() any { return make(ErrNoSuchField, 32) }}
	poolField          = filter.Pool[Field]{New: func() any { return make(Field, 0, 32) }}
)

func Straight[T any, M interface {
	~map[string]T
	Type() filter.Type
}](b filter.Builder, j filter.Projector, o string, oo M) (err error) {
	nn := poolErrNoSuchField.Get()
	for k := range oo {
		nn[k] = struct{}{}
	}
	ff := poolField.Get()
	for _, k := range j.Names() {
		_, ok := oo[k]
		if ok {
			ff = append(ff, k)
			delete(nn, k)
		}
	}
	if len(nn) > 0 {
		return nn
	}
	poolErrNoSuchField.Put(nn)
	defer func() { poolField.Put(ff.Erase()) }()
	if len(oo) > 1 {
		_, err = b.WriteString("( ")
		if err != nil {
			return
		}
		defer func() {
			if err == nil {
				_, err = b.WriteString(" )")
			}
		}()
	}
	t := j.Table()
	for i, f := range ff {
		if i > 0 {
			err = b.WriteByte(' ')
			if err != nil {
				return
			}
			_, err = b.WriteString(o)
			if err != nil {
				return
			}
			err = b.WriteByte(' ')
			if err != nil {
				return
			}
		}
		_, err = b.AppendColumn(t, f)
		if err != nil {
			return
		}
		_, err = b.Append(oo.Type(), oo[f])
		if err != nil {
			return
		}
	}
	return
}
