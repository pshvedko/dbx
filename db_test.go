package db_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pshvedko/db"
	"github.com/pshvedko/db/filter"
	"github.com/pshvedko/db/request"
	"github.com/pshvedko/db/t"
)

type DB struct {
	*db.DB
}

func openDB(t *testing.T) (*DB, error) {
	t.Helper()
	bd := os.Getenv("TEST_POSTGRES")
	if bd == "" {
		t.Skip("env var TEST_POSTGRES is not set")
	}
	c, err := db.Open(bd)
	if err != nil {
		return nil, err
	}
	return &DB{DB: c}, nil
}

func TestDB(t *testing.T) {
	tt, err := openDB(t)

	require.NoError(t, err)
	require.NotNil(t, tt)

	t.Run("Get", tt.TestGet)
}

func (bd DB) TestGet(t *testing.T) {
	type args struct {
		ctx context.Context
		o   db.Object
		f   filter.Filter
		oo  []request.Option
	}
	ctx := context.TODO()
	tests := []struct {
		name    string
		args    args
		want    db.Object
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				ctx: ctx,
				o:   &help.Object{},
				f:   filter.Eq{"id": 1},
				oo:  nil,
			},
			want: &help.Object{
				Id:      1,
				Bool:    true,
				Float32: 0,
				Float64: 3.14,
				Int:     0,
				Int16:   16,
				Null:    nil,
				String:  "red",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bd.Get(tt.args.ctx, tt.args.o, tt.args.f, tt.args.oo...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.Equal(t, tt.want, tt.args.o)
		})
	}
}
