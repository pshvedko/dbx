package request

import (
	"context"
	"github.com/pshvedko/db/driver/postgres"
	"github.com/pshvedko/db/filter"
	"github.com/pshvedko/db/model"
)

type Request struct {
	filter.Builder
}

func (r Request) Get(u *model.User) error {
	return nil
}

func New(ctx context.Context) Request {
	return Request{
		Builder: postgres.NewHolder(),
	}
}
