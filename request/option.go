package request

import (
	"context"
	"database/sql"
)

type Option interface {
	Apply(r *Request) error
}

type OptionFunc func(r *Request) error

func (f OptionFunc) Apply(r *Request) error {
	return f(r)
}

func makeConnect(ctx context.Context, c Connector) OptionFunc {
	return func(r *Request) error {
		err := r.makeConn(ctx, c)
		if err != nil {
			return err
		}
		return r.makeTx(ctx, c)
	}
}

type WithField []string

func (o WithField) Apply(r *Request) error {
	return r.withField(true, o...)
}

type WithoutField []string

func (o WithoutField) Apply(r *Request) error {
	return r.withField(false, o...)
}

type WithTx sql.TxOptions

func (o WithTx) Apply(r *Request) error {
	r.o = (*sql.TxOptions)(&o)
	return nil
}

type WithOwner string

func (o WithOwner) Apply(r *Request) error {
	r.owner = string(o)
	return nil
}

type WithGroup string

func (o WithGroup) Apply(r *Request) error {
	r.group = string(o)
	return nil
}

type WithDeleted string

func (o WithDeleted) Apply(r *Request) error {
	r.x.d = string(o)
	return nil
}

type ReadDeleted int

const (
	DeletedNone ReadDeleted = iota
	DeletedOnly
	DeletedFree
)

func (o ReadDeleted) Apply(r *Request) error {
	r.m = o
	return nil
}
