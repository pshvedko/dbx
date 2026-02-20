package request

import (
	"context"
	"database/sql"
	"runtime"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
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
		return r.makeTx(ctx)
	}
}

func WithCount() OptionFunc {
	return func(r *Request) error {
		r.z = false
		return nil
	}
}

func WithoutCount() OptionFunc {
	return func(r *Request) error {
		r.z = true
		return nil
	}
}

func WithoutReturning() OptionFunc {
	return func(r *Request) error {
		return r.withField(1, false)
	}
}

type WithReturnField []string

func (o WithReturnField) Apply(r *Request) error {
	return r.withField(1, true, o...)
}

type WithField []string

func (o WithField) Apply(r *Request) error {
	return r.withField(0, true, o...)
}

type WithoutField []string

func (o WithoutField) Apply(r *Request) error {
	return r.withField(0, false, o...)
}

type WithTx sql.TxOptions

func (o WithTx) Apply(r *Request) error {
	r.o = &sql.TxOptions{Isolation: o.Isolation, ReadOnly: o.ReadOnly}
	return nil
}

type WithOwner string

func (o WithOwner) Apply(r *Request) error {
	r.a.u = string(o)
	return nil
}

type WithGroup string

func (o WithGroup) Apply(r *Request) error {
	r.a.g = string(o)
	return nil
}

type WithDeleted string // filed must be defined as DEFAULT NULL and NOW compatible

func (o WithDeleted) Apply(r *Request) error {
	r.s.d = string(o)
	return nil
}

type WithUpdated string // filed must be defined as DEFAULT NOW

func (o WithUpdated) Apply(r *Request) error {
	r.s.u = string(o)
	return nil
}

type WithCreated string // filed must be defined as DEFAULT NOW

func (o WithCreated) Apply(r *Request) error {
	r.s.c = string(o)
	return nil
}

type ReadDeleted int

func (o ReadDeleted) Apply(r *Request) error {
	r.m = o
	return nil
}

const (
	DeletedNone ReadDeleted = iota // IS NULL
	DeletedOnly                    // IS NOT NULL
	DeletedFree                    // IS NULL AND IS NOT NULL
)

type PerformPut int

func (o PerformPut) Apply(r *Request) error {
	r.p = o
	return nil
}

func (o PerformPut) Mode() builder.Mode {
	return builder.Mode(o)
}

const (
	PutModify PerformPut = iota
	PutCreate
	PutUpdate
)

func WithCache(c builder.Keeper) OptionFunc {
	_, file, line, _ := runtime.Caller(1)
	return func(r *Request) error {
		r.h.o = builder.Origin{File: file, Line: line}
		r.h.c = c
		return nil
	}
}

type Cache struct {
	filter.Map[builder.Hash, string]
}
