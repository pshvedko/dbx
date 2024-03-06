package request

import (
	"context"
	"log/slog"
)

type Logger interface {
	Handler() slog.Handler
	With(...any) *slog.Logger
	WithGroup(string) *slog.Logger
	Enabled(context.Context, slog.Level) bool
	Log(context.Context, slog.Level, string, ...any)
	LogAttrs(context.Context, slog.Level, string, ...slog.Attr)
	Debug(string, ...any)
	DebugContext(context.Context, string, ...any)
	Info(string, ...any)
	InfoContext(context.Context, string, ...any)
	Warn(string, ...any)
	WarnContext(context.Context, string, ...any)
	Error(string, ...any)
	ErrorContext(context.Context, string, ...any)
}
