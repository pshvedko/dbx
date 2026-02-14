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

func New(db *sqlx.DB) *DB {
	return &DB{
		DB:     db,
		Logger: slog.New(logHandler{}),
	}
}

func (db *DB) Option() []request.Option {
	return db.oo
}

func (db *DB) WithOption(oo ...request.Option) *DB {
	db.oo = append(db.oo, oo...)
	return db
}

func (db *DB) Connect(ctx context.Context) (*sqlx.Conn, error) {
	return db.Connx(ctx)
}

type logHandler struct{}

func (h logHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (h logHandler) Handle(context.Context, slog.Record) error { return nil }
func (h logHandler) WithAttrs([]slog.Attr) slog.Handler        { return h }
func (h logHandler) WithGroup(string) slog.Handler             { return h }

func (db *DB) WithLogger(h slog.Handler) *DB {
	db.Logger = slog.New(h)
	return db
}

func Get(ctx context.Context, db request.Connector, j filter.Projector, f filter.Filter, oo ...request.Option) error {
	r, err := request.New(ctx, db, oo...)
	if err != nil {
		return err
	}
	err = r.Get(ctx, j, f)
	return r.End(err)
}

func List(ctx context.Context, db request.Connector, i filter.Injector, f filter.Filter, o, l *uint, y []string, oo ...request.Option) (uint, error) {
	r, err := request.New(ctx, db, oo...)
	if err != nil {
		return 0, err
	}
	total, err := r.List(ctx, i, f, o, l, y)
	return total, r.End(err)
}

func Put(ctx context.Context, db request.Connector, j filter.Projector, oo ...request.Option) error {
	r, err := request.New(ctx, db, oo...)
	if err != nil {
		return err
	}
	err = r.Put(ctx, j)
	return r.End(err)
}

func Delete(ctx context.Context, db request.Connector, i filter.Injector, f filter.Filter, oo ...request.Option) error {
	r, err := request.New(ctx, db, oo...)
	if err != nil {
		return err
	}
	err = r.Delete(ctx, i, f)
	return r.End(err)
}
