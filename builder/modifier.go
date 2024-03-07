package builder

import "github.com/pshvedko/dbx/filter"

type Column interface {
	HasColumn(string) bool
}

type AllowedColumn map[string]struct{}

func (f AllowedColumn) HasColumn(k string) bool {
	_, ok := f[k]
	return ok
}

type ExcludedColumn map[string]struct{}

func (f ExcludedColumn) HasColumn(k string) bool {
	_, ok := f[k]
	return !ok
}

type Modify struct {
	Created string
	Updated string
	Deleted
}

type Deleted interface {
	Visibility(filter.And) filter.And
	HasDeleted(string) bool
}

type DeletedOnly string

func (o DeletedOnly) HasDeleted(n string) bool {
	return n == string(o)
}

func (o DeletedOnly) Visibility(a filter.And) filter.And {
	return append(a, filter.Ne{string(o): nil})
}

type DeletedNone string

func (o DeletedNone) HasDeleted(n string) bool {
	return n == string(o)
}

func (o DeletedNone) Visibility(a filter.And) filter.And {
	return append(a, filter.Eq{string(o): nil})
}

type DeletedFree string

func (o DeletedFree) HasDeleted(n string) bool {
	return n == string(o)
}

func (DeletedFree) Visibility(a filter.And) filter.And {
	return a
}
