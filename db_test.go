package dbx_test

import (
	"context"
	"database/sql"
	"fmt"
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
	context.Context
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
	db.SetOption(request.WithCreated("o_time_0"), request.WithUpdated("o_time_1"), request.WithDeleted("o_time_4"), request.WithTx{})
	t.Cleanup(func() {
		_ = db.Close()
	})
	return &DB{DB: db, Context: context.TODO()}, nil
}

func TestDB(t *testing.T) {
	db, err := openDB(t)
	require.NoError(t, err)
	require.NotNil(t, db)
	t.Run("Connect", db.TestConn)
	t.Run("Get", db.TestGet)
	t.Run("List", db.TestList)
	t.Run("ListIn", db.TestListIn)
	t.Run("ListAny", db.TestListAny)
	t.Run("ListLike", db.TestListLike)
	t.Run("Put", db.TestPut)
}

func (db DB) TestConn(t *testing.T) {
	conn, err := db.Connx(db)
	require.NoError(t, err)
	require.NotZero(t, t, conn)
	var pid, pid2 int
	err = conn.GetContext(db, &pid, `SELECT pg_backend_pid()`)
	require.NoError(t, err)
	require.NotZero(t, pid)
	err = conn.GetContext(db, &pid2, `SELECT pg_backend_pid()`)
	require.NoError(t, err)
	require.Equal(t, pid, pid2)
	err = conn.GetContext(db, &pid2, `SELECT pg_backend_pid()`)
	require.NoError(t, err)
	require.Equal(t, pid, pid2)

	tx, err := conn.BeginTxx(db, &sql.TxOptions{})
	require.NoError(t, err)
	require.NotZero(t, tx)
	err = tx.GetContext(db, &pid2, `SELECT pg_backend_pid()`)
	require.NoError(t, err)
	require.Equal(t, pid, pid2)
	err = tx.Rollback()
	require.NoError(t, err)
	err = tx.GetContext(db, &pid2, `SELECT pg_backend_pid()`)
	require.Error(t, err)

	tx, err = conn.BeginTxx(db, &sql.TxOptions{})
	require.NoError(t, err)
	require.NotZero(t, tx)
	err = tx.GetContext(db, &pid2, `SELECT pg_backend_pid()`)
	require.NoError(t, err)
	require.Equal(t, pid, pid2)
	err = tx.Commit()
	require.NoError(t, err)
	err = tx.GetContext(db, &pid2, `SELECT pg_backend_pid()`)
	require.Error(t, err)

	err = conn.GetContext(db, &pid2, `SELECT pg_backend_pid()`)
	require.NoError(t, err)
	require.Equal(t, pid, pid2)

	err = conn.Close()
	require.NoError(t, err)

	err = conn.GetContext(db, &pid2, `SELECT pg_backend_pid()`)
	require.Error(t, err)
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
	oo = nil
	total, err = db.List(context.TODO(), &oo, filter.In{"o_time_0": {"1970-01-01T00:00:00Z"}}, nil, nil, nil, request.WithField{"id"})
	require.NoError(t, err)
	require.EqualValues(t, 5, total)
	require.ElementsMatch(t, help.ObjectList{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, oo)
	oo = nil
	total, err = db.List(context.TODO(), &oo, filter.In{"o_time_0": {"YESTERDAY", filter.Now(), time.Now(), time.UnixMicro(0)}}, nil, nil, nil, request.WithField{"id"})
	require.NoError(t, err)
	require.EqualValues(t, 5, total)
	require.ElementsMatch(t, help.ObjectList{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, oo)
}

func (db DB) TestListLike(t *testing.T) {
	var oo help.ObjectList
	total, err := db.List(context.TODO(), &oo, filter.And{
		filter.As{"o_string_1": "%a%"},
		filter.Lt{"o_time_1": filter.Now()},
	}, nil, nil, nil, request.WithField{"id", "o_string_1"}, request.DeletedOnly)
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.ElementsMatch(t, help.ObjectList{{ID: 7, String1: help.PtrString("gray")}}, oo)
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
				Time0:   time.Unix(0, 0),
				Time1:   time.Unix(0, 0),
				Time2:   help.PtrTime(time.Unix(0, 0)),
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
				f:   filter.And{filter.Eq{"id": 1}, filter.Le{"o_time_0": "YESTERDAY"}},
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
			want:    2,
			want1:   &help.ObjectList{{ID: 6}, {ID: 7}},
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
			want:    7,
			want1:   &help.ObjectList{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				i:   &help.ObjectList{},
				f:   filter.Eq{"o_time_1": time.Unix(0, 0), "o_uint_64": nil},
				o:   help.PtrUint(1),
				l:   help.PtrUint(3),
				y:   []string{"-id"},
				oo:  []request.Option{request.WithField{"id", "o_absent_0", "o_string_1"}},
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
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				i:   &help.ObjectList{},
				f:   filter.Eq{"o_float_32": "100", "o_float_64": "3.14", "o_int_16": "16", "o_string_1": "red", "o_string_2": "hello", "o_time_0": "1970-01-01T00:00:00Z", "o_bool": "true", "o_null": nil},
				o:   nil,
				l:   nil,
				y:   nil,
				oo:  nil,
			},
			want: 1,
			want1: &help.ObjectList{{
				ID:      1,
				Bool:    help.PtrBool(true),
				Float32: 100,
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
				Time0:   time.Unix(0, 0),
				Time1:   time.Unix(0, 0),
				Time2:   help.PtrTime(time.Unix(0, 0)),
				Time3:   nil,
				Time4:   time.Time{},
			}},
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

func (db DB) TestPut(t *testing.T) {
	type args struct {
		ctx context.Context
		o   dbx.Object
		oo  []request.Option
	}
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
				ctx: context.TODO(),
				o: &help.Object{
					ID:      9,
					Bool:    help.PtrBool(false),
					Float64: help.PtrFloat64(0),
					Int16:   help.PtrInt16(0),
				},
				oo: []request.Option{request.PutCreate},
			},
			want: &help.Object{
				ID:      9,
				Bool:    help.PtrBool(false),
				Float64: help.PtrFloat64(0),
				Int16:   help.PtrInt16(0),
				String1: help.PtrString("green"),
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				o: &help.Object{
					ID:      9,
					Bool:    help.PtrBool(true),
					Float64: help.PtrFloat64(1e1),
					Int16:   help.PtrInt16(1),
					String3: "orange",
				},
				oo: []request.Option{request.WithField{"o_bool", "o_string_3", "o_float_64", "o_int_16"}},
			},
			want: &help.Object{
				ID:      9,
				Bool:    help.PtrBool(true),
				Float64: help.PtrFloat64(1e1),
				Int16:   help.PtrInt16(1),
				String1: help.PtrString("green"),
				String3: "orange",
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				o: &help.Object{
					ID:      9,
					Bool:    help.PtrBool(true),
					Float64: help.PtrFloat64(1e2),
					Int16:   help.PtrInt16(2),
					String1: help.PtrString("orange"),
				},
				oo: []request.Option{request.WithField{"o_bool", "o_string_1", "o_float_64", "o_int_16"}},
			},
			want: &help.Object{
				ID:      9,
				Bool:    help.PtrBool(true),
				Float64: help.PtrFloat64(1e2),
				Int16:   help.PtrInt16(2),
				String1: help.PtrString("orange"),
				String3: "orange",
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				o: &help.Object{
					ID:      9,
					Bool:    help.PtrBool(true),
					Float64: help.PtrFloat64(1e3),
					Int16:   help.PtrInt16(3),
					String3: "green",
				},
				oo: []request.Option{},
			},
			want: &help.Object{
				ID:      9,
				Bool:    help.PtrBool(true),
				Float64: help.PtrFloat64(1e3),
				Int16:   help.PtrInt16(3),
				String1: help.PtrString("green"),
				String3: "green",
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				o: &help.Object{
					ID:      9,
					Bool:    help.PtrBool(false),
					Float64: help.PtrFloat64(1e4),
					Int16:   help.PtrInt16(4),
					String3: "orange",
				},
				oo: []request.Option{request.PutUpdate, request.WithField{"o_bool", "o_string_3"}},
			},
			want: &help.Object{
				ID:      9,
				Bool:    help.PtrBool(false),
				Float64: help.PtrFloat64(1e3),
				Int16:   help.PtrInt16(3),
				String1: help.PtrString("green"),
				String3: "orange",
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				o: &help.Object{
					ID:      7,
					Bool:    help.PtrBool(true),
					Float64: help.PtrFloat64(1e3),
					Int16:   help.PtrInt16(3),
					String3: "green",
				},
				oo: []request.Option{request.WithField{"o_bool", "o_string_3"}},
			},
			want:    nil,
			wantErr: sql.ErrNoRows,
		},
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				o: &help.Object{
					ID:      7,
					Bool:    help.PtrBool(false),
					Float64: help.PtrFloat64(1e4),
					Int16:   help.PtrInt16(4),
					String3: "orange",
				},
				oo: []request.Option{request.PutUpdate, request.WithField{"o_bool", "o_string_3"}},
			},
			want:    nil,
			wantErr: sql.ErrNoRows,
		},
	}
	var ids help.Map
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Put(tt.args.ctx, tt.args.o, tt.args.oo...)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				t.Log(ids.Add(tt.args.o.Table(), tt.args.o.Get(0)))
				for i := 0; i < 15; i++ {
					require.Equal(t, tt.want.Get(i), tt.args.o.Get(i), i)
				}
				t.Log(tt.args.o.Get(15))
				t.Log(tt.args.o.Get(16))
			}
		})
	}
	t.Cleanup(func() {
		for k, v := range ids.M {
			_, _ = db.Exec(fmt.Sprintf(`DELETE FROM %q WHERE "id" =ANY($1)`, k), v)
		}
	})
}
