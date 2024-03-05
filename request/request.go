package request

import (
	"context"

	"github.com/pshvedko/db/builder"
	"github.com/pshvedko/db/filter"
)

type Request struct {
	c Connection
	t bool
}

func (r *Request) makeConn(ctx context.Context, db Connector) error {
	if r.c == nil && ctx != nil && db != nil {
		c, err := db.Connx(ctx)
		if err != nil {
			return err
		}
		r.c = Conn{Conn: c}
		r.t = false
	}
	return nil
}

func (r *Request) makeTx(ctx context.Context) error {
	if r.t == false && ctx != nil {
		c, err := r.c.BeginTxx(ctx, nil)
		if err != nil {
			return err
		}
		r.c = Tx{Tx: c}
		r.t = true
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

func (r *Request) End(err *error) {
	*err = r.c.End(*err)
	r.c = nil
	r.t = false
}

func (r *Request) Get(ctx context.Context, j filter.Projector, f filter.Filter) error {
	q, aa, vv, err := r.Constructor().Select(j, f)
	if err != nil {
		return err
	}
	return r.c.QueryRowxContext(ctx, q, aa...).Scan(vv...)
}

func (r *Request) Constructor() builder.Constructor {
	return builder.Constructor{}
}
