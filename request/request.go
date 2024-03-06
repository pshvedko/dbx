package request

import (
	"context"
	"log/slog"

	"github.com/pshvedko/db/builder"
	"github.com/pshvedko/db/filter"
)

type Request struct {
	c Connection
	t bool
	f map[string]struct{}
	b bool
}

func (r *Request) makeConn(ctx context.Context, db Connector) error {
	if r.c == nil && ctx != nil && db != nil {
		c, err := db.Connx(ctx)
		if err != nil {
			return err
		}
		r.c = Conn{Conn: c, Logger: slog.New(db.Handler())}
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
		r.c = Tx{Tx: c, Logger: slog.New(r.c.Handler())}
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

func (r *Request) List(ctx context.Context, j filter.Projector, f filter.Filter, o, l *uint, y builder.Order) (int, error) {
	q, aa, vv, err := r.Constructor().Range(o, l).Sort(y).Select(j, f)
	if err != nil {
		return 0, err
	}
	rows, err := r.c.QueryxContext(ctx, q, aa...)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err = rows.Scan(vv...)
		if err != nil {
			break
		}
	}
	err2 := rows.Close()
	if err2 != nil {
		return 0, err2
	}
	if err != nil {
		return 0, err
	}
	return 0, rows.Err()
}

func (r *Request) Constructor() builder.Constructor {
	return builder.Constructor{
		Column: r.fields(),
	}
}

func (r *Request) fields() builder.Column {
	if r.b || len(r.f) == 0 {
		return builder.ExcludedColumn(r.f)
	}
	return builder.AllowedColumn(r.f)
}

func (r *Request) withField(b bool, kk ...string) error {
	switch {
	case r.b == b:
		r.b = !r.b
		fallthrough
	case r.f == nil:
		r.f = map[string]struct{}{}
	}
	for _, k := range kk {
		r.f[k] = struct{}{}
	}
	return nil
}
