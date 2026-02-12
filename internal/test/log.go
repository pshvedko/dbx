package test

import (
	"context"
	"log/slog"
	"testing"
)

type logHandler struct {
	testing.TB
}

func (h logHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h logHandler) Handle(_ context.Context, r slog.Record) error {
	h.Helper()
	h.Log(r.Level, r.Message)
	r.Attrs(func(a slog.Attr) bool {
		h.Helper()
		h.Log(r.Level, a)
		return true
	})
	return nil
}

func (h logHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h logHandler) WithGroup(string) slog.Handler {
	return h
}

func LogHandler(t *testing.T) slog.Handler {
	return logHandler{TB: t}
}
