package builder

import (
	"github.com/pshvedko/dbx/filter"
)

type Column interface {
	Used(string) bool
	Names() map[string]struct{}
	Allowed() bool
}

type AllowedColumn map[string]struct{}

func (c AllowedColumn) Used(k string) bool {
	_, ok := c[k]
	return ok
}

func (c AllowedColumn) Names() map[string]struct{} {
	return c
}

func (c AllowedColumn) Allowed() bool {
	return true
}

type ExcludedColumn map[string]struct{}

func (c ExcludedColumn) Used(k string) bool {
	_, ok := c[k]
	return !ok
}

func (c ExcludedColumn) Names() map[string]struct{} {
	return c
}

func (c ExcludedColumn) Allowed() bool {
	return false
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
	DeletionClause(filter.And) filter.And
	IsDeleted(string) bool
}

type DeletedOnly string

func (o DeletedOnly) IsDeleted(n string) bool {
	return n == string(o)
}

func (o DeletedOnly) DeletionClause(a filter.And) filter.And {
	return append(a, filter.Ne{string(o): nil})
}

type DeletedNone string

func (o DeletedNone) IsDeleted(n string) bool {
	return n == string(o)
}

func (o DeletedNone) DeletionClause(a filter.And) filter.And {
	return append(a, filter.Eq{string(o): nil})
}

type DeletedFree string

func (o DeletedFree) IsDeleted(n string) bool {
	return n == string(o)
}

func (DeletedFree) DeletionClause(a filter.And) filter.And {
	return a
}
