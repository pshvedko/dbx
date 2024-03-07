package builder

import "github.com/pshvedko/dbx/filter"

type Modify interface {
	Visibility(a filter.And, n string) filter.And
}

type DefaultAvailability struct {
	DeletedNone
}

type DeletedOnly struct {
	Modify
}

func (DeletedOnly) Visibility(a filter.And, n string) filter.And {
	return append(a, filter.Ne{n: nil})
}

type DeletedNone struct {
	Modify
}

func (DeletedNone) Visibility(a filter.And, n string) filter.And {
	return append(a, filter.Eq{n: nil})
}

type DeletedFree struct {
	Modify
}

func (DeletedFree) Visibility(a filter.And, _ string) filter.And {
	return a
}
