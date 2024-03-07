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
	if len(bd) == 0 {
		t.Skip("env var TEST_POSTGRES is not set")
	}
	db, err := dbx.Open(bd)
	if err != nil {
		return nil, err
	}
	db.SetLogger(help.LogHandler(t))
	db.SetOption(request.Deleted("o_time_4"))
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return &DB{DB: db}, nil
}

func TestDB(t *testing.T) {
	db, err := openDB(t)
	require.NoError(t, err)
	require.NotNil(t, db)
	t.Run("Get", db.TestGet)
	t.Run("List", db.TestList)
	t.Run("ListIn", db.TestListIn)
	t.Run("ListAny", db.TestListAny)
}

func (db DB) TestListAny(t *testing.T) {
	var oo []help.Object
	err := db.Select(&oo, `SELECT "id" FROM "objects" WHERE "id" =ANY($1) AND "o_string_1" =ANY($2) AND "o_bool" =ANY($3) AND "o_float_64" =ANY($4) AND "o_uuid_2" =ANY($5) AND "o_time_1" <>ALL($6)`,
		filter.Array{1, 2, 3, 100},
		filter.Array{"red", "black", "white", "green", "yellow"},
		filter.Array{false, true},
		filter.Array{0, 3.14},
		filter.Array{uuid.UUID{}},
		filter.Array{time.Time{}},
	)
	require.NoError(t, err)
	require.ElementsMatch(t, []help.Object{{ID: 1}, {ID: 2}, {ID: 3}}, oo)
}

func (db DB) TestListIn(t *testing.T) {
	var oo help.ObjectList
	total, err := db.List(context.TODO(), &oo, filter.And{
		filter.In{"id": {1, 2, 3, 4, 5}, "o_string_1": {"red", "black", "white", "green", "yellow"}},
		filter.Ni{"o_time_1": {time.Time{}}}}, nil, nil, nil, request.WithField{"id"})
	require.NoError(t, err)
	require.EqualValues(t, 5, total)
	require.ElementsMatch(t, help.ObjectList{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, oo)
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
				f:   nil,
				o:   nil,
				l:   nil,
				y:   []string{"id"},
				oo:  []request.Option{request.WithField{"id"}},
			},
			want:    5,
			want1:   &help.ObjectList{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				i:   &help.ObjectList{},
				f:   nil,
				o:   nil,
				l:   nil,
				y:   []string{"id"},
				oo:  []request.Option{request.WithField{"id"}, request.DeletedOnly},
			},
			want:    1,
			want1:   &help.ObjectList{{ID: 6}},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				i:   &help.ObjectList{},
				f:   nil,
				o:   nil,
				l:   nil,
				y:   []string{"id"},
				oo:  []request.Option{request.WithField{"id"}, request.DeletedFree},
			},
			want:    6,
			want1:   &help.ObjectList{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}},
			wantErr: nil,
		},
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
