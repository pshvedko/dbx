package dbx

import (
	"context"

	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/request"
)

type Object interface {
	filter.Projector
}

func (db *DB) Apply(ctx context.Context, r *request.Request) error {
	return r.Tx(ctx, db)
}

func (db *DB) Get(ctx context.Context, o Object, f filter.Filter, oo ...request.Option) (err error) {
	r, err := request.New(ctx, db, oo...)
	if err != nil {
		return
	}
	defer r.End(&err)
	return r.Get(ctx, o, f)
}

func (db *DB) List(ctx context.Context, i filter.Injector, f filter.Filter, o, l *uint, y []string, oo ...request.Option) (total uint, err error) {
	r, err := request.New(ctx, db, oo...)
	if err != nil {
		return
	}
	defer r.End(&err)
	return r.List(ctx, i, f, o, l, y)
}
