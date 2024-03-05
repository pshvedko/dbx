package request

import (
	"context"
	"fmt"
	"github.com/pshvedko/db/builder"
	"github.com/pshvedko/db/filter"
)

type Request struct {
	c Connection
}

func (r *Request) makeConn(ctx context.Context, db Connector) error {
	if r.c == nil && ctx != nil && db != nil {
		c, err := db.Connx(ctx)
		if err != nil {
			return err
		}
		r.c = Conn{Conn: c}
	}
	return nil
}

func (r *Request) makeTx(ctx context.Context) error {
	if r.c.NoTxx() && ctx != nil {
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

func (r *Request) End(err *error) {
	*err = r.c.End(*err)
	r.c = nil
}

func (r *Request) Get(ctx context.Context, j filter.Projector, f filter.Filter) error {
	q, aa, vv, err := r.makeSelect(j, f)
	if err != nil {
		return err
	}
	return r.c.QueryRowxContext(ctx, q, aa...).Scan(vv...)
}

func (r *Request) makeSelect(p filter.Projector, f filter.Filter) (string, []any, []any, error) {
	var b builder.Filter
	_, err := b.WriteString("SELECT")
	if err != nil {
		return "", nil, nil, err
	}
	j := 0
	nn := p.Names()
	vv := p.Values()
	for i, n := range nn {
		if i > 0 {
			err = b.WriteByte(',')
			if err != nil {
				return "", nil, nil, err
			}
		}
		_, err = fmt.Fprintf(&b, " %q", n)
		if err != nil {
			return "", nil, nil, err
		}
		vv[j] = vv[i]
		j++
	}
	{
		_, err = fmt.Fprintf(&b, " FROM %q", p.Table())
		if err != nil {
			return "", nil, nil, err
		}
	}
	{
		_, err = fmt.Fprintf(&b, " WHERE ")
		if err != nil {
			return "", nil, nil, err
		}
		n := b.Len()
		if f != nil {
			err = f.To(&b, p)
			if err != nil {
				return "", nil, nil, err
			}
		}
		if n == b.Len() {
			_, err = fmt.Fprintf(&b, "TRUE")
			if err != nil {
				return "", nil, nil, err
			}
		}
	}
	return b.String(), b.Values(), vv[:j], nil
}
