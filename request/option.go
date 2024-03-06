package request

import (
	"context"
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
