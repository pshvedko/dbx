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

func placeholder(args []any) []any {
	i := len(args)
	args = append(args, make([]any, i)...)
	i <<= 1
	for i > 0 {
		args[i-1] = args[i>>1-1]
		args[i-2] = fmt.Sprint("$", i>>1)
		i -= 2
	}
	return args
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
