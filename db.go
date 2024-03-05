package db

import (
	"context"

	"github.com/pshvedko/db/filter"
	"github.com/pshvedko/db/request"
)

type Object interface {
	filter.Fielder
}

func (db *DB) Get(ctx context.Context, o Object, f filter.Filter, oo ...request.Option) (err error) {
	r, err := request.New(ctx, db, append(oo, request.WithTx(db))...)
	if err != nil {
		return
	}
	defer r.End(ctx, &err)
	return r.SelectOne(ctx, o, f)
}
