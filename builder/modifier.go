package builder

import (
	"github.com/pshvedko/dbx/filter"
)

type Column interface {
	Used(string) bool
}

type AllowedColumn map[string]struct{}

func (f AllowedColumn) Used(k string) bool {
	_, ok := f[k]
	return ok
}

type ExcludedColumn map[string]struct{}

func (f ExcludedColumn) Used(k string) bool {
	_, ok := f[k]
	return !ok
}

type Modify struct {
	Created string
	Updated string
	Deleted
}

func (m Modify) IsCreated(n string) bool {
	return m.Created == n
}

func (m Modify) IsUpdated(n string) bool {
	return m.Updated == n
}

type Deleted interface {
	Visibility(filter.And) filter.And
	IsDeleted(string) bool
	Name() (string, bool)
}

type DeletedOnly string

func (o DeletedOnly) Name() (string, bool) {
	return string(o), len(o) > 0
}

func (o DeletedOnly) IsDeleted(n string) bool {
	return n == string(o)
}

func (o DeletedOnly) Visibility(a filter.And) filter.And {
	return append(a, filter.Ne{string(o): nil})
}

type DeletedNone string

func (o DeletedNone) Name() (string, bool) {
	return string(o), len(o) > 0
}

func (o DeletedNone) IsDeleted(n string) bool {
	return n == string(o)
}

func (o DeletedNone) Visibility(a filter.And) filter.And {
	return append(a, filter.Eq{string(o): nil})
}

type DeletedFree string

func (o DeletedFree) Name() (string, bool) {
	return string(o), len(o) > 0
}

func (o DeletedFree) IsDeleted(n string) bool {
	return n == string(o)
}

func (DeletedFree) Visibility(a filter.And) filter.And {
	return a
}
