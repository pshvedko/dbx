package request

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Option interface {
	Apply(ctx context.Context, r *Request) error
}

type OptionFunc func(ctx context.Context, r *Request) error

func (f OptionFunc) Apply(ctx context.Context, r *Request) error {
	return f(ctx, r)
}

type Conn struct {
	*sqlx.Conn
}

func (c Conn) CanTxx() bool {
	return true
}

func WithConnect(db Connector) OptionFunc {
	return func(ctx context.Context, r *Request) error {
		if r.c == nil {
			c, err := db.Connx(ctx)
			if err != nil {
				return err
			}
			r.c = Conn{Conn: c}
		}
		return nil
	}
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

func WithTx(db Connector) OptionFunc {
	return func(ctx context.Context, r *Request) error {
		if r.c == nil {
			err := WithConnect(db)(ctx, r)
			if err != nil {
				return err
			}
		}
		if r.c.CanTxx() {
			c, err := r.c.BeginTxx(ctx, nil)
			if err != nil {
				return err
			}
			r.c = Tx{Tx: c}
		}
		return nil
	}
}
