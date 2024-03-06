package db_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
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
		wantErr error
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
				ID:      1,
				Bool:    help.PtrBool(true),
				Float32: 1e2,
				Float64: help.PtrFloat64(3.14),
				Int:     0,
				Int16:   help.PtrInt16(16),
				Null:    nil,
				String:  help.PtrString("red"),
				Uint64:  nil,
				UUID1:   uuid.UUID{},
				UUID2:   help.PtrUUID(uuid.UUID{}),
				UUID3:   nil,
				UUID4:   uuid.UUID{},
				Time1:   time.Unix(0, 0).UTC(),
				Time2:   help.PtrTime(time.Unix(0, 0).UTC()),
				Time3:   nil,
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				ctx: ctx,
				o:   &help.Object{},
				f:   filter.Eq{"id": 1},
				oo:  []request.Option{request.WithField{"id", "o_string"}},
			},
			want: &help.Object{
				ID:      1,
				Bool:    nil,
				Float32: 0,
				Float64: nil,
				Int:     0,
				Int16:   nil,
				Null:    nil,
				String:  help.PtrString("red"),
				Uint64:  nil,
				UUID1:   uuid.UUID{},
				UUID2:   nil,
				UUID3:   nil,
				UUID4:   uuid.UUID{},
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				ctx: ctx,
				o:   &help.Object{},
				f:   filter.Eq{"id": nil},
				oo:  nil,
			},
			want:    &help.Object{},
			wantErr: sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bd.Get(tt.args.ctx, tt.args.o, tt.args.f, tt.args.oo...)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, tt.args.o)
		})
	}
}
