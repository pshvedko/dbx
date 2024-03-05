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
	NoTxx() bool
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	End(error) error
}

type Conn struct {
	*sqlx.Conn
}

func (c Conn) End(err1 error) error {
	err2 := c.Close()
	if err2 != nil {
		return err2
	}
	return err1
}

func (c Conn) NoTxx() bool {
	return true
}

type Tx struct {
	*sqlx.Tx
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

func (c Tx) NoTxx() bool {
	return false
}
