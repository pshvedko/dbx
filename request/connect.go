package request

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Beginner interface {
	Logger
	BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error)
}

type Connector interface {
	Beginner
	Connx(context.Context) (*sqlx.Conn, error)
	Option() []Option
}

type Connection interface {
	Beginner
	sqlx.ExecerContext
	sqlx.QueryerContext
	End(error) error
	Close() error
}

type Conn struct {
	*sqlx.Conn
	*slog.Logger
}

func (c Conn) QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	c.Logger.DebugContext(ctx, query, placeholder(args)...)
	return c.Conn.QueryxContext(ctx, query, args...)
}

func (c Conn) QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	c.Logger.DebugContext(ctx, query, placeholder(args)...)
	return c.Conn.QueryRowxContext(ctx, query, args...)
}

func placeholder(vv []any) []any {
	aa := make([]any, 0, 2*len(vv))
	for i, v := range vv {
		aa = append(aa, fmt.Sprint("$", i), v)
	}
	return aa
}

func (c Conn) End(err1 error) error {
	return err1
}

type Tx struct {
	*sqlx.Tx
	*slog.Logger
	io.Closer
}

func (c Tx) QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	c.Logger.DebugContext(ctx, query, placeholder(args)...)
	return c.Tx.QueryxContext(ctx, query, args...)
}

func (c Tx) QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	c.Logger.DebugContext(ctx, query, placeholder(args)...)
	return c.Tx.QueryRowxContext(ctx, query, args...)
}

func (c Tx) End(err1 error) error {
	if err1 == nil {
		return c.Commit()
	}
	err := c.Rollback()
	if err == nil {
		return err1
	}
	return err
}

func (c Tx) BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error) {
	panic("don't implement me")
}
