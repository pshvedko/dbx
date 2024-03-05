package request

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Connector interface {
	Connx(ctx context.Context) (*sqlx.Conn, error)
}

type Connection interface {
	CanTxx() bool
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

type Conn struct {
	*sqlx.Conn
}

func (c Conn) CanTxx() bool {
	return true
}

type Tx struct {
	*sqlx.Tx
}

func (t Tx) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	panic("dont implement me")
}

func (t Tx) CanTxx() bool {
	return false
}
