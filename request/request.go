package request

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
)

type Request struct {
	c       Connection
	t       bool
	f       map[string]struct{}
	b       bool
	deleted string
	updated string
	created string
	m       builder.Modify
	o       *sql.TxOptions
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
	if r.o != nil && !r.t && ctx != nil {
		c, err := r.c.BeginTxx(ctx, r.o)
		if err != nil {
			return err
		}
		r.c = Tx{Tx: c, Logger: slog.New(r.c.Handler())}
		r.t = true
	}
	return nil
}

func New(ctx context.Context, db Connector, oo ...Option) (*Request, error) {
	r := Request{
		m: builder.DefaultAvailability{},
	}
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
	*err = r.c.End(*err)
	r.c = nil
	r.t = false
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
		Column: r.fields(),
		Modify: r.m,
		Option: builder.Option{
			Created: r.created,
			Updated: r.updated,
			Deleted: r.deleted,
		},
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
