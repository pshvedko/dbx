package dbx

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/request"
)

type Object interface {
	filter.Projector
}

type DB struct {
	*sqlx.DB
	*slog.Logger
	oo []request.Option
}

func Open(name string) (*DB, error) {
	db, err := sqlx.Open("pgx", name)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB:     db,
		Logger: slog.New(logHandler{}),
	}, nil
}

func (db *DB) Option() []request.Option {
	return db.oo
}

func (db *DB) SetOption(oo ...request.Option) {
	db.oo = append(db.oo, oo...)
}

func (db *DB) Get(ctx context.Context, o Object, f filter.Filter, oo ...request.Option) (err error) {
	r, err := request.New(ctx, db, append(oo, request.WithTx{})...)
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
