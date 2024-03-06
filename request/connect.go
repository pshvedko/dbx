package request

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Connector interface {
	Logger
	Connx(context.Context) (*sqlx.Conn, error)
}

type Connection interface {
	Logger
	sqlx.ExecerContext
	sqlx.QueryerContext
	BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error)
	End(error) error
}

type Conn struct {
	*sqlx.Conn
	*slog.Logger
}

func (c Conn) QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	c.Logger.WarnContext(ctx, query, placeholder(args)...)
	return c.Conn.QueryxContext(ctx, query, args...)
}

func (c Conn) QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	c.Logger.WarnContext(ctx, query, placeholder(args)...)
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
	err2 := c.Close()
	if err2 != nil {
		return err2
	}
	return err1
}

type Tx struct {
	*sqlx.Tx
	*slog.Logger
}

func (c Tx) QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	c.Logger.WarnContext(ctx, query, placeholder(args)...)
	return c.Tx.QueryxContext(ctx, query, args...)
}

func (c Tx) QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	c.Logger.WarnContext(ctx, query, placeholder(args)...)
	return c.Tx.QueryRowxContext(ctx, query, args...)
}

func (c Tx) End(err1 error) error {
	if err1 == nil {
		return c.Commit()
	}
	err2 := c.Rollback()
	if err2 != nil {
		return err2
	}
	return err1
}

func (c Tx) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	panic("dont implement me")
}
