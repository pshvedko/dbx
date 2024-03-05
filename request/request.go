package request

import (
	"context"
	"io"

	"github.com/pshvedko/db/filter"
)

type Request struct {
	c Connection
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

func (r *Request) End(ctx context.Context, err *error) {

}

func (r *Request) SelectOne(ctx context.Context, o filter.Fielder, f filter.Filter) error {
	return io.EOF
}
