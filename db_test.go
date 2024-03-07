package dbx_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pshvedko/dbx"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/request"
	"github.com/pshvedko/dbx/t"
)

type DB struct {
	*dbx.DB
}

func openDB(t *testing.T) (*DB, error) {
	t.Helper()
	bd := os.Getenv("TEST_POSTGRES")
	if bd == "" {
		t.Skip("env var TEST_POSTGRES is not set")
	}
	db, err := dbx.Open(bd)
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	db.EnableLogger(help.LogHandler(t))
	return &DB{DB: db}, nil
}

func TestDB(t *testing.T) {
	tt, err := openDB(t)

	require.NoError(t, err)
	require.NotNil(t, tt)

	t.Run("Get", tt.TestGet)
	t.Run("List", tt.TestList)
}

func (db DB) TestGet(t *testing.T) {
	type args struct {
		ctx context.Context
		o   dbx.Object
		f   filter.Filter
		oo  []request.Option
	}
	ctx := context.TODO()
	tests := []struct {
		name    string
		args    args
		want    dbx.Object
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
				String1: help.PtrString("red"),
				String2: "hello",
				String3: "",
				Uint64:  nil,
				UUID1:   uuid.UUID{},
				UUID2:   help.PtrUUID(uuid.UUID{}),
				UUID3:   nil,
				UUID4:   uuid.UUID{},
				Time1:   time.Unix(0, 0).UTC(),
				Time2:   help.PtrTime(time.Unix(0, 0).UTC()),
				Time3:   nil,
				Time4:   time.Time{},
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				ctx: ctx,
				o:   &help.Object{},
				f:   filter.Eq{"id": 1},
				oo:  []request.Option{request.WithField{"id", "o_string_1"}},
			},
			want: &help.Object{
				ID:      1,
				Bool:    nil,
				Float32: 0,
				Float64: nil,
				Int:     0,
				Int16:   nil,
				Null:    nil,
				String1: help.PtrString("red"),
				Uint64:  nil,
				UUID1:   uuid.UUID{},
				UUID2:   nil,
				UUID3:   nil,
				UUID4:   uuid.UUID{},
			},
			wantErr: nil,
		}, {
			name: "",
			args: args{
				ctx: ctx,
				o:   &help.Object{},
				f:   filter.Eq{"id": 1},
				oo:  []request.Option{request.WithoutField(help.Object{}.Names())},
			},
			want:    &help.Object{},
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
			err := db.Get(tt.args.ctx, tt.args.o, tt.args.f, tt.args.oo...)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, tt.args.o)
		})
	}
}

func (db *DB) TestList(t *testing.T) {
	type args struct {
		ctx context.Context
		i   filter.Injector
		f   filter.Filter
		o   *uint
		l   *uint
		y   []string
		oo  []request.Option
	}
	tests := []struct {
		name    string
		args    args
		want    uint
		want1   filter.Injector
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				i:   &help.ObjectList{},
				f:   filter.Eq{"o_time_1": time.Unix(0, 0).UTC(), "o_uint_64": nil},
				o:   help.PtrUint(1),
				l:   help.PtrUint(3),
				y:   []string{"-id"},
				oo:  []request.Option{request.WithField{"id", "o_string_1"}},
			},
			want: 4,
			want1: &help.ObjectList{
				{
					ID:      3,
					String1: help.PtrString("white"),
				}, {
					ID:      2,
					String1: help.PtrString("black"),
				}, {
					ID:      1,
					String1: help.PtrString("red"),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.List(tt.args.ctx, tt.args.i, tt.args.f, tt.args.o, tt.args.l, tt.args.y, tt.args.oo...)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want1, tt.args.i)
		})
	}
}
