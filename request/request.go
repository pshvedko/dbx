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
	e       bool
	t       bool
	f       map[string]struct{}
	b       bool
	deleted string
	updated string
	created string
	m       builder.Ability
	o       *sql.TxOptions
	owner   string
	group   string
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

func (r *Request) makeTx(ctx context.Context, b Beginner) error {
	if r.o != nil && !r.t && ctx != nil {
		if r.e {
			b = r.c
		}
		t, err := b.BeginTxx(ctx, r.o)
		if err != nil {
			return err
		}
		r.c = Tx{Tx: t, Logger: slog.New(b.Handler())}
		r.e = true
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
	if r.e {
		*err = r.c.End(*err)
	}
	r.c = nil
	r.o = nil
	r.m = nil
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
		Column: r.fields(),
		Access: builder.Access{
			Owner: r.owner,
			Group: r.group,
		},
		Modify: builder.Modify{
			Created: r.created,
			Updated: r.updated,
			Deleted: r.deleted,
			Ability: r.m,
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
