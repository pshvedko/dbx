package request

import (
	"context"
	"database/sql"

	"github.com/pshvedko/dbx/builder"
)

type Option interface {
	Apply(r *Request) error
}

type OptionFunc func(r *Request) error

func (f OptionFunc) Apply(r *Request) error {
	return f(r)
}

func makeConnect(ctx context.Context, db Connector) OptionFunc {
	return func(r *Request) error {
		err := r.makeConn(ctx, db)
		if err != nil {
			return err
		}
		return r.makeTx(ctx)
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

type BeginTx sql.TxOptions

func (o BeginTx) Apply(r *Request) error {
	r.o = (*sql.TxOptions)(&o)
	return nil
}

type Owner string

func (o Owner) Apply(r *Request) error {
	r.owner = string(o)
	return nil
}

type Group string

func (o Group) Apply(r *Request) error {
	r.group = string(o)
	return nil
}

type Deleted string

func (o Deleted) Apply(r *Request) error {
	r.deleted = string(o)
	return nil
}

type ReadDeleted int

const (
	DeletedOnly ReadDeleted = iota
	DeletedNone
	DeletedFree
)

func (o ReadDeleted) Apply(r *Request) error {
	switch o {
	case DeletedFree:
		r.m = builder.DeletedFree{Ability: r.m}
	case DeletedNone:
		r.m = builder.DeletedNone{Ability: r.m}
	case DeletedOnly:
		r.m = builder.DeletedOnly{Ability: r.m}
	}
	return nil
}
