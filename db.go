package dbx

import (
	"context"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/request"
)

type Object interface {
	filter.Projector
}

func (db *DB) Get(ctx context.Context, o Object, f filter.Filter, oo ...request.Option) (err error) {
	r, err := request.New(ctx, db, oo...)
	if err != nil {
		return
	}
	defer r.End(&err)
	return r.Get(ctx, o, f)
}
