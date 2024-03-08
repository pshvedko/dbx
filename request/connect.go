package request

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Beginner interface {
	Logger
	io.Closer
	BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error)
}

type Connector interface {
	Beginner
	Connect(context.Context) (*sqlx.Conn, error)
	Option() []Option
}

type Connection interface {
	Beginner
	Query(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) *sqlx.Row
	End(error) error
}

type Conn struct {
	*sqlx.Conn
	*slog.Logger
}

func (c Conn) Query(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	c.DebugContext(ctx, query, placeholder(args)...)
	return c.QueryxContext(ctx, query, args...)
}

func (c Conn) QueryRow(ctx context.Context, query string, args ...any) *sqlx.Row {
	c.DebugContext(ctx, query, placeholder(args)...)
	return c.QueryRowxContext(ctx, query, args...)
}

func placeholder(vv []any) []any {
	aa := make([]any, 0, 2*len(vv))
	for i, v := range vv {
		aa = append(aa, fmt.Sprint("$", i+1), v)
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

func (c Tx) Query(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	c.DebugContext(ctx, query, placeholder(args)...)
	return c.QueryxContext(ctx, query, args...)
}

func (c Tx) QueryRow(ctx context.Context, query string, args ...any) *sqlx.Row {
	c.DebugContext(ctx, query, placeholder(args)...)
	return c.QueryRowxContext(ctx, query, args...)
}

func (c Tx) End(err1 error) error {
	if err1 == nil {
		return c.Commit()
	}
	err2 := c.Rollback()
	if err2 == nil {
		return err1
	}
	return errors.Join(err1, err2)
}

func (c Tx) BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error) {
	panic("don't implement me")
}
