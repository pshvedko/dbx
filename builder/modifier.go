package builder

import "github.com/pshvedko/dbx/filter"

type Field struct {
	Table  string
	Column string
}

type Fields map[string]int

func (m Fields) Find(t, n string) bool {
	if _, ok := m[n]; ok {
		return true
	}
	if _, ok := m[t+"."+n]; ok {
		return true
	}
	return false
}

type Fielder interface {
	Used(string, string) bool
	Names() [2]Fields
	Returned(string, string) bool
}

type IncludedColumn [2]Fields

func (c IncludedColumn) Returned(t string, n string) bool {
	if c[1] != nil {
		return c[1].Find(t, n)
	}
	return c.Used(t, n)
}

func (c IncludedColumn) Used(t string, n string) bool {
	return c[0].Find(t, n)
}

func (c IncludedColumn) Names() [2]Fields {
	return c
}

type ExcludedColumn [2]Fields

func (c ExcludedColumn) Returned(t string, n string) bool {
	if c[1] != nil {
		return c[1].Find(t, n)
	}
	return c.Used(t, n)
}

func (c ExcludedColumn) Used(t string, n string) bool {
	return !c[0].Find(t, n)
}

func (c ExcludedColumn) Names() [2]Fields {
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
	WithDeleted(filter.And) filter.And
	IsDeleted(string) bool
	AsDeleted() string
}

type DeletedOnly string

func (o DeletedOnly) AsDeleted() string { return string(o) }

func (o DeletedOnly) IsDeleted(n string) bool { return n == string(o) }

func (o DeletedOnly) WithDeleted(a filter.And) filter.And {
	return append(a, filter.Ne{string(o): nil})
}

type DeletedNone string

func (o DeletedNone) AsDeleted() string { return string(o) }

func (o DeletedNone) IsDeleted(n string) bool { return n == string(o) }

func (o DeletedNone) WithDeleted(a filter.And) filter.And {
	return append(a, filter.Eq{string(o): nil})
}

type DeletedFree string

func (o DeletedFree) AsDeleted() string { return string(o) }

func (o DeletedFree) IsDeleted(n string) bool { return n == string(o) }

func (DeletedFree) WithDeleted(a filter.And) filter.And {
	return a
}
