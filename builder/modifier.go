package builder

import "github.com/pshvedko/dbx/filter"

type Modify struct {
	Created string
	Updated string
	Deleted string
	Ability
}

func (m Modify) Visibility(a filter.And) filter.And {
	return m.Ability.Visibility(a, m.Deleted)
}

type Ability interface {
	Visibility(a filter.And, n string) filter.And
}

type DefaultAvailability struct {
	DeletedNone
}

type DeletedOnly struct {
	Ability
}

func (DeletedOnly) Visibility(a filter.And, n string) filter.And {
	return append(a, filter.Ne{n: nil})
}

type DeletedNone struct {
	Ability
}

func (DeletedNone) Visibility(a filter.And, n string) filter.And {
	return append(a, filter.Eq{n: nil})
}

type DeletedFree struct {
	Ability
}

func (DeletedFree) Visibility(a filter.And, _ string) filter.And {
	return a
}
