package request

import (
	"context"
	"github.com/pshvedko/dbx/builder"
)

type Option interface {
	Apply(ctx context.Context, r *Request) error
}

type OptionFunc func(ctx context.Context, r *Request) error

func (f OptionFunc) Apply(ctx context.Context, r *Request) error {
	return f(ctx, r)
}

func makeConnect(db Connector) OptionFunc {
	return func(ctx context.Context, r *Request) error {
		return r.makeConn(ctx, db)
	}
}

type WithField []string

func (o WithField) Apply(_ context.Context, r *Request) error {
	return r.withField(true, o...)
}

type WithoutField []string

func (o WithoutField) Apply(_ context.Context, r *Request) error {
	return r.withField(false, o...)
}

type Deleted string

func (w Deleted) Apply(_ context.Context, r *Request) error {
	r.deleted = string(w)
	return nil
}

type ReadDeleted int

const (
	DeletedOnly ReadDeleted = iota
	DeletedNone
	DeletedFree
)

func (w ReadDeleted) Apply(_ context.Context, r *Request) error {
	switch w {
	case DeletedFree:
		r.m = builder.DeletedFree{Modify: r.m}
	case DeletedNone:
		r.m = builder.DeletedNone{Modify: r.m}
	case DeletedOnly:
		r.m = builder.DeletedOnly{Modify: r.m}
	}
	return nil
}
