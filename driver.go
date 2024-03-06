package db

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type logHandler struct{}

func (h logHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (h logHandler) Handle(context.Context, slog.Record) error { return nil }
func (h logHandler) WithAttrs([]slog.Attr) slog.Handler        { return h }
func (h logHandler) WithGroup(string) slog.Handler             { return h }

type DB struct {
	*sqlx.DB
	*slog.Logger
}

func Open(name string) (*DB, error) {
	db, err := sqlx.Open("pgx", name)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db, Logger: slog.New(logHandler{})}, nil
}

func (db *DB) EnableLogger(h slog.Handler) {
	db.Logger = slog.New(h)
}
