package request

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
)

type Request struct {
	c Connection
	e bool
	t bool
	x [2]bool
	f [2]map[string]int
	o *sql.TxOptions
	a struct {
		u string
		g string
	}
	m ReadDeleted
	p PerformPut
	s struct {
		d string
		u string
		c string
	}
	z bool
	h struct {
		o builder.Origin
		k builder.Keeper
	}
	w func()
	b bool
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
		c, err := b.Connect(ctx)
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
			return r.End(err)
		}
		r.c = Tx{Tx: t, Logger: slog.New(r.c.Handler()), Closer: r.closer()}
		r.e = true
		r.t = true
	}
	return nil
}

func (r *Request) apply(oo ...Option) error {
	for _, o := range oo {
		err := o.Apply(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewWithOption(options ...[]Option) (*Request, error) {
	var r Request
	for _, oo := range options {
		err := r.apply(oo...)
		if err != nil {
			return nil, r.End(err)
		}
	}
	return &r, nil
}

func New(ctx context.Context, db Connector, oo ...Option) (*Request, error) {
	return NewWithOption([]Option{makeDefer(func() {})}, db.Option(), oo, []Option{makeConnect(ctx, db)})
}

func makeDefer(f func()) OptionFunc {
	return func(r *Request) error {
		w := r.w
		r.w = func() {
			if w != nil {
				w()
			}
			f()
		}
		return nil
	}
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

func (r *Request) End(err error) error {
	if r.e {
		err1 := r.c.End(err)
		err2 := r.c.Close()
		if err1 == nil {
			err = err2
		} else if err2 == nil {
			err = err1
		} else {
			err = errors.Join(err1, err2)
		}
	}
	r.w()
	r.w = nil
	r.c = nil
	r.o = nil
	r.f = [2]map[string]int{nil, nil}
	r.x = [2]bool{false, false}
	r.t = false
	r.e = false
	r.z = false
	r.h.k = nil
	return err
}

func (r *Request) withField(i int, b bool, kk ...string) error {
	switch {
	case r.x[i] == b:
		r.x[i] = !r.x[i]
		fallthrough
	case r.f[i] == nil:
		r.f[i] = map[string]int{}
	}
	for _, k := range kk {
		_, ok := r.f[i][k]
		if ok {
			return fmt.Errorf("repeated column: %s", k)
		}
		r.f[i][k] = 0
	}
	return nil
}

func (r *Request) Constructor() *builder.Constructor {
	return &builder.Constructor{
		Fielder: func() builder.Fielder {
			if r.x[0] || len(r.f[0]) == 0 {
				return builder.ExcludedColumn(r.f)
			}
			return builder.IncludedColumn(r.f)
		}(),
		Access: builder.Access{
			Owner: r.a.u,
			Group: r.a.g,
		},
		Aliases: make(builder.Aliases, 2),
		Donner:  func() {},
		Modify: builder.Modify{
			Created: r.s.c,
			Updated: r.s.u,
			Deleted: func() builder.Deleted {
				if r.m == DeletedFree {
					return builder.DeletedFree(r.s.d)
				} else if r.m == DeletedOnly {
					return builder.DeletedOnly(r.s.d)
				}
				return builder.DeletedNone(r.s.d)
			}(),
		},
		Cache: builder.Cache{Origin: r.h.o, Keeper: r.h.k},
		Mode:  r.p.Mode(),
		Z:     r.z,
		T:     r.b,
	}
}

func (r *Request) Get(ctx context.Context, j filter.Projector, f filter.Filter) error {
	_, q, aa, vv, err := r.Constructor().Select(j, f)
	if err != nil {
		return err
	}
	return r.c.QueryRow(ctx, q, aa...).Scan(vv...)
}

func (r *Request) List(ctx context.Context, i filter.Injector, f filter.Filter, o, l *uint, y builder.Order) (uint, error) {
	j := i.Element()
	z, q, aa, vv, err := r.Constructor().Range(o, l).Sort(y).Select(j, f)
	if err != nil {
		return 0, err
	}
	rows, err := r.c.Query(ctx, q, aa...)
	if err != nil {
		return 0, err
	}
	var t uint
	for rows.Next() {
		err = rows.Scan(vv...)
		if err != nil {
			break
		}
		i.Inject(j)
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
	}
	if z == nil {
		return t, nil
	}
	// TODO - FOR UPDATE / COUNT(*) OVER() / CTE
	// SELECT
	//    *, COUNT(*) OVER() AS total
	// FROM table
	// WHERE ...
	// ORDER BY ...
	// LIMIT 10 OFFSET 20
	//
	// WITH result AS (
	//    SELECT * FROM table WHERE ... LIMIT 10 OFFSET 20
	// )
	// SELECT
	//    *, (SELECT COUNT(*) FROM table WHERE ...) AS total
	// FROM result
	p, n, err := z.Count()
	if err != nil {
		return 0, err
	}
	var x uint
	err = r.c.QueryRow(ctx, p, aa[:n]...).Scan(&x)
	if err != nil {
		return 0, err
	}
	if x > t {
		return x, nil
	}
	return t, nil
}

func (r *Request) Put(ctx context.Context, j filter.Projector) error {
	q, aa, vv, err := r.Constructor().Insert(j)
	if err != nil {
		return err
	}
	return r.c.QueryRow(ctx, q, aa...).Scan(vv...)
}

func (r *Request) Delete(ctx context.Context, i filter.Injector, f filter.Filter) error {
	j := i.Element()
	q, aa, vv, err := r.Constructor().Delete(j, f)
	if err != nil {
		return err
	}
	rows, err := r.c.Query(ctx, q, aa...)
	if err != nil {
		return err
	}
	var t uint
	for rows.Next() {
		err = rows.Scan(vv...)
		if err != nil {
			break
		}
		i.Inject(j)
		t++
	}
	err2 := rows.Close()
	if err2 != nil {
		return err2
	}
	if err != nil {
		return err
	}
	return rows.Err()
}

func (r *Request) WithOption(oo ...Option) error {
	err := r.apply(oo...)
	if err != nil {
		return err
	}
	return nil
}
