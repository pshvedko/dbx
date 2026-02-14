package builder

import (
	"github.com/pshvedko/dbx/filter"
)

type Column interface {
	Used(string) bool
	Names() [2]map[string]struct{}
	Returned(string) bool
}

type IncludedColumn [2]map[string]struct{}

func (c IncludedColumn) Returned(k string) bool {
	if c[1] != nil {
		_, ok := c[1][k]
		return ok
	}
	_, ok := c[0][k]
	return ok
}

func (c IncludedColumn) Used(k string) bool {
	_, ok := c[0][k]
	return ok
}

func (c IncludedColumn) Names() [2]map[string]struct{} {
	return c
}

type ExcludedColumn [2]map[string]struct{}

func (c ExcludedColumn) Returned(k string) bool {
	if c[1] != nil {
		_, ok := c[1][k]
		return ok
	}
	_, ok := c[0][k]
	return !ok
}

func (c ExcludedColumn) Used(k string) bool {
	_, ok := c[0][k]
	return !ok
}

func (c ExcludedColumn) Names() [2]map[string]struct{} {
	return c
}

type Modify struct {
	Created string
	Updated string
	Deleted
}

func (m Modify) IsCreated(n string) bool { return m.Created == n }

func (m Modify) IsUpdated(n string) bool { return m.Updated == n }

type Deleted interface {
	DeletionClause(filter.And) filter.And
	IsDeleted(string) bool
	Name() string
}

type DeletedOnly string

func (o DeletedOnly) Name() string { return string(o) }

func (o DeletedOnly) IsDeleted(n string) bool { return n == string(o) }

func (o DeletedOnly) DeletionClause(a filter.And) filter.And {
	return append(a, filter.Ne{string(o): nil})
}

type DeletedNone string

func (o DeletedNone) Name() string { return string(o) }

func (o DeletedNone) IsDeleted(n string) bool { return n == string(o) }

func (o DeletedNone) DeletionClause(a filter.And) filter.And {
	return append(a, filter.Eq{string(o): nil})
}

type DeletedFree string

func (o DeletedFree) Name() string { return string(o) }

func (o DeletedFree) IsDeleted(n string) bool { return n == string(o) }

func (DeletedFree) DeletionClause(a filter.And) filter.And {
	return a
}
