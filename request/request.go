package request

import (
	"context"
	"io"

	"github.com/pshvedko/db/filter"
)

type Request struct {
	c Connection
}

func (r *Request) makeConn(ctx context.Context, db Connector) error {
	if r.c != nil {
		c, err := db.Connx(ctx)
		if err != nil {
			return err
		}
		r.c = Conn{Conn: c}
	}
	return nil
}

func (r *Request) makeTx(ctx context.Context) error {
	if r.c.CanTxx() {
		c, err := r.c.BeginTxx(ctx, nil)
		if err != nil {
			return err
		}
		r.c = Tx{Tx: c}
	}
	return nil
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
