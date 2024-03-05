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

func WithConnect(db Connector) OptionFunc {
	return func(ctx context.Context, r *Request) error {
		return r.makeConn(ctx, db)
	}
}

func WithTx(db Connector) OptionFunc {
	return func(ctx context.Context, r *Request) error {
		err := r.makeConn(ctx, db)
		if err != nil {
			return err
		}
		return r.makeTx(ctx)
	}
}
