package dbx

import (
	"context"
	"log/slog"
)

type logHandler struct{}

func (h logHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (h logHandler) Handle(context.Context, slog.Record) error { return nil }
func (h logHandler) WithAttrs([]slog.Attr) slog.Handler        { return h }
func (h logHandler) WithGroup(string) slog.Handler             { return h }

func (db *DB) SetLogger(h slog.Handler) {
	db.Logger = slog.New(h)
}
