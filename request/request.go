package request

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
)

type Request struct {
	c Connection
	e bool
	t bool
	b bool
	f map[string]struct{}
	o *sql.TxOptions
	u string
	g string
	m ReadDeleted
	x struct {
		d string
		u string
		c string
	}
}

func (r *Request) closer() io.Closer {
	if r.e {
		return r.c
	}
	return r
}

func (r *Request) Close() error {
	return nil
}

func (r *Request) makeConn(ctx context.Context, b Connector) error {
	if r.c == nil && ctx != nil && b != nil {
		c, err := b.Connx(ctx)
		if err != nil {
			return err
		}
		r.c = Conn{Conn: c, Logger: slog.New(b.Handler())}
		r.e = true
	}
	return nil
}

func (r *Request) makeTx(ctx context.Context) error {
	if r.o != nil && !r.t && ctx != nil {
		t, err := r.c.BeginTxx(ctx, r.o)
		if err != nil {
			return err
		}
		r.c = Tx{Tx: t, Logger: slog.New(r.c.Handler()), Closer: r.closer()}
		r.e = true
		r.t = true
	}
	return nil
}

func New(ctx context.Context, db Connector, oo ...Option) (*Request, error) {
	var r Request
	for _, o := range append(db.Option(), append(oo, makeConnect(ctx, db))...) {
		err := o.Apply(&r)
		if err != nil {
			return nil, err
		}
	}
	return &r, nil
}

func (r *Request) Apply(a *Request) error {
	switch a.c {
	case nil:
	default:
		a.c = r.c
		a.t = r.t
	}
	return nil
}

func (r *Request) End(err *error) {
	if r.e {
		err1 := r.c.End(*err)
		err2 := r.c.Close()
		if err1 == nil {
			*err = err2
		} else if err2 == nil {
			*err = err1
		} else {
			*err = errors.Join(err1, err2)
		}
	}
	r.c = nil
	r.o = nil
	r.f = nil
	r.b = false
	r.t = false
	r.e = false
}

func (r *Request) Get(ctx context.Context, j filter.Projector, f filter.Filter) error {
	_, q, aa, vv, err := r.Constructor().Select(j, f)
	if err != nil {
		return err
	}
	return r.c.QueryRowxContext(ctx, q, aa...).Scan(vv...)
}

func (r *Request) List(ctx context.Context, i filter.Injector, f filter.Filter, o, l *uint, y builder.Order) (uint, error) {
	j := i.Get()
	c := r.Constructor()
	z, q, aa, vv, err := c.Range(o, l).Sort(y).Select(j, f)
	if err != nil {
		return 0, err
	}
	rows, err := r.c.QueryxContext(ctx, q, aa...)
	if err != nil {
		return 0, err
	}
	var t uint
	for rows.Next() {
		err = rows.Scan(vv...)
		if err != nil {
			break
		}
		i.Put(j)
		t++
	}
	err2 := rows.Close()
	if err2 != nil {
		return 0, err2
	}
	if err != nil {
		return 0, err
	}
	err = rows.Err()
	if err != nil {
		return 0, err
	} else if z == nil {
		return t, nil
	}
	p, n, err := z.Count()
	if err != nil {
		return 0, err
	}
	err = r.c.QueryRowxContext(ctx, p, aa[:n]...).Scan(&t)
	if err != nil {
		return 0, err
	}
	return t, nil
}

func (r *Request) Constructor() *builder.Constructor {
	return &builder.Constructor{
		Column: func() builder.Column {
			if r.b || len(r.f) == 0 {
				return builder.ExcludedColumn(r.f)
			}
			return builder.AllowedColumn(r.f)
		}(),
		Access: builder.Access{
			Owner: r.u,
			Group: r.g,
		},
		Modify: builder.Modify{
			Created: r.x.c,
			Updated: r.x.u,
			Deleted: func() builder.Deleted {
				if r.m == DeletedFree {
					return builder.DeletedFree(r.x.d)
				} else if r.m == DeletedOnly {
					return builder.DeletedOnly(r.x.d)
				}
				return builder.DeletedNone(r.x.d)
			}(),
		},
	}
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
