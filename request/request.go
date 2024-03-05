package request

import (
	"context"
	"database/sql"
	"io"

	"github.com/jmoiron/sqlx"

	"github.com/pshvedko/db/filter"
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

type Request struct {
	c Connection
}

func New(ctx context.Context, db Connector, oo ...Option) (*Request, error) {
	var r Request
	for _, o := range append(oo, WithConnect(db)) {
		err := o.Apply(ctx, &r)
		if err != nil {
			return nil, err
		}
	}
	return &r, nil
}

func (r *Request) End(ctx context.Context, err *error) {

}

func (r *Request) SelectOne(ctx context.Context, o filter.Fielder, f filter.Filter) error {
	return io.EOF
}
