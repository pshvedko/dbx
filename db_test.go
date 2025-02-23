package dbx_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/pshvedko/dbx"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test"
	"github.com/pshvedko/dbx/request"
	"github.com/pshvedko/dbx/util"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	*dbx.DB
}

func Open(t *testing.T) (*DB, error) {
	t.Helper()
	bd := os.Getenv("TEST_POSTGRES")
	if len(bd) == 0 {
		t.Skip("env var TEST_POSTGRES is not set")
	}
	db, err := sqlx.Open("pgx", bd)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	return &DB{
		DB: dbx.New(db).
			WithLogger(test.LogHandler(t)).
			WithOption(
				request.WithCreated("o_time_0"),
				request.WithUpdated("o_time_1"),
				request.WithDeleted("o_time_4"),
				request.WithTx{})}, nil
}

func TestDB(t *testing.T) {
	db, err := Open(t)
	require.NoError(t, err)
	require.NotNil(t, db)
	t.Run("Connect", db.TestConn)
	t.Run("TestScan", db.TestScan)
	t.Run("Get", db.TestGet)
	t.Run("List", db.TestList)
	t.Run("ListIn", db.TestListIn)
	t.Run("ListAny", db.TestListAny)
	t.Run("ListLike", db.TestListLike)
	t.Run("Put", db.TestPut)

}

const selectPidStmt = `SELECT pg_backend_pid()`

func (db DB) TestConn(t *testing.T) {
	ctx := context.TODO()
	conn, err := db.Connx(ctx)
	require.NoError(t, err)
	require.NotZero(t, t, conn)
	var pid, pid2 int
	err = conn.GetContext(ctx, &pid, selectPidStmt)
	require.NoError(t, err)
	require.NotZero(t, pid)
	err = conn.GetContext(ctx, &pid2, selectPidStmt)
	require.NoError(t, err)
	require.Equal(t, pid, pid2)
	err = conn.GetContext(ctx, &pid2, selectPidStmt)
	require.NoError(t, err)
	require.Equal(t, pid, pid2)

	tx, err := conn.BeginTxx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	require.NotZero(t, tx)
	err = tx.GetContext(ctx, &pid2, selectPidStmt)
	require.NoError(t, err)
	require.Equal(t, pid, pid2)
	err = tx.Rollback()
	require.NoError(t, err)
	err = tx.GetContext(ctx, &pid2, selectPidStmt)
	require.Error(t, err)

	tx, err = conn.BeginTxx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	require.NotZero(t, tx)
	err = tx.GetContext(ctx, &pid2, selectPidStmt)
	require.NoError(t, err)
	require.Equal(t, pid, pid2)
	err = tx.Commit()
	require.NoError(t, err)
	err = tx.GetContext(ctx, &pid2, selectPidStmt)
	require.Error(t, err)

	err = conn.GetContext(ctx, &pid2, selectPidStmt)
	require.NoError(t, err)
	require.Equal(t, pid, pid2)

	err = conn.Close()
	require.NoError(t, err)

	err = conn.GetContext(ctx, &pid2, selectPidStmt)
	require.Error(t, err)
}

func (db DB) TestListAny(t *testing.T) {
	var oo []test.Object
	err := db.SelectContext(context.TODO(), &oo, `--
			SELECT "id" 
			FROM "objects" 
			WHERE "id" =ANY($1) 
			  AND "o_string_1" =ANY($2) 
			  AND "o_bool" =ANY($3)
			  AND "o_float_64" =ANY($4)
			  AND "o_uuid_2" =ANY($5)
			  AND "o_time_1" <>ALL($6)`,
		filter.Array{1, 2, 3, 100},
		filter.Array{"red", "black", "white", "green", "yellow"},
		filter.Array{false, true},
		filter.Array{0, 3.14},
		filter.Array{uuid.UUID{}},
		filter.Array{time.Time{}},
	)
	require.NoError(t, err)
	require.ElementsMatch(t, []test.Object{{ID: 1}, {ID: 2}, {ID: 3}}, oo)
}

func (db DB) TestListIn(t *testing.T) {
	var oo test.ObjectList
	total, err := dbx.List(context.TODO(), db, &oo,
		filter.And{
			filter.In{"id": {1, 2, 3, 4, 5}, "o_string_1": {"red", "black", "white", "green", "yellow"}},
			filter.Ni{"o_time_1": {time.Time{}}},
		}, nil, nil, nil, request.WithField{"id"})
	require.NoError(t, err)
	require.EqualValues(t, 5, total)
	require.ElementsMatch(t, test.ObjectList{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, oo)
	oo = nil
	total, err = dbx.List(context.TODO(), db, &oo,
		filter.In{"o_time_0": {"1970-01-01T00:00:00Z"}}, nil, nil, nil, request.WithField{"id"})
	require.NoError(t, err)
	require.EqualValues(t, 5, total)
	require.ElementsMatch(t, test.ObjectList{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, oo)
	oo = nil
	total, err = dbx.List(context.TODO(), db, &oo,
		filter.In{"o_time_0": {"YESTERDAY", filter.Now(), time.Now(), time.UnixMicro(0)}}, nil, nil, nil, request.WithField{"id"})
	require.NoError(t, err)
	require.EqualValues(t, 5, total)
	require.ElementsMatch(t, test.ObjectList{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, oo)
}

func (db DB) TestListLike(t *testing.T) {
	var oo test.ObjectList
	total, err := dbx.List(context.TODO(), db, &oo,
		filter.And{
			filter.As{"o_string_1": "%a%"},
			filter.Lt{"o_time_1": filter.Now()},
		}, nil, nil, nil, request.WithField{"id", "o_string_1"}, request.DeletedOnly)
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.ElementsMatch(t, test.ObjectList{{ID: 7, String1: util.PtrString("gray")}}, oo)
}

func (db DB) TestGet(t *testing.T) {
	ctx := context.TODO()
	type args struct {
		o  dbx.Object
		f  filter.Filter
		oo []request.Option
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
				o:  &test.Object{},
				f:  filter.Eq{"id": 1},
				oo: nil,
			},
			want: &test.Object{
				ID:      1,
				Bool:    util.PtrBool(true),
				Float32: 1e2,
				Float64: util.PtrFloat64(3.14),
				Int:     0,
				Int16:   util.PtrInt16(16),
				Null:    nil,
				String1: util.PtrString("red"),
				String2: "hello",
				String3: "",
				Uint64:  nil,
				UUID1:   uuid.UUID{},
				UUID2:   util.PtrUUID(uuid.UUID{}),
				UUID3:   nil,
				UUID4:   uuid.UUID{},
				Time0:   time.Unix(0, 0),
				Time1:   time.Unix(0, 0),
				Time2:   util.PtrTime(time.Unix(0, 0)),
				Time3:   nil,
				Time4:   time.Time{},
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o:  &test.Object{},
				f:  filter.Eq{"id": 1},
				oo: []request.Option{request.WithField{"id", "o_string_1"}},
			},
			want: &test.Object{
				ID:      1,
				Bool:    nil,
				Float32: 0,
				Float64: nil,
				Int:     0,
				Int16:   nil,
				Null:    nil,
				String1: util.PtrString("red"),
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
				o:  &test.Object{},
				f:  filter.And{filter.Eq{"id": 1}, filter.Le{"o_time_0": "YESTERDAY"}},
				oo: []request.Option{request.WithoutField(test.Object{}.Names())},
			},
			want:    &test.Object{},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o:  &test.Object{},
				f:  filter.Eq{"id": nil},
				oo: nil,
			},
			want:    &test.Object{},
			wantErr: sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dbx.Get(ctx, db, tt.args.o, tt.args.f, tt.args.oo...)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, tt.args.o)
		})
	}
}

func (db DB) TestList(t *testing.T) {
	ctx := context.TODO()
	type args struct {
		i  filter.Injector
		f  filter.Filter
		o  *uint
		l  *uint
		y  []string
		oo []request.Option
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
				i:  &test.ObjectList{},
				f:  nil,
				o:  nil,
				l:  nil,
				y:  []string{"id"},
				oo: []request.Option{request.WithField{"id"}},
			},
			want:    5,
			want1:   &test.ObjectList{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				i:  &test.ObjectList{},
				f:  nil,
				o:  nil,
				l:  nil,
				y:  []string{"id"},
				oo: []request.Option{request.WithField{"id"}, request.DeletedOnly},
			},
			want:    2,
			want1:   &test.ObjectList{{ID: 6}, {ID: 7}},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				i:  &test.ObjectList{},
				f:  nil,
				o:  nil,
				l:  nil,
				y:  []string{"id"},
				oo: []request.Option{request.WithField{"id"}, request.DeletedFree},
			},
			want:    7,
			want1:   &test.ObjectList{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				i:  &test.ObjectList{},
				f:  filter.Eq{"o_time_1": time.Unix(0, 0), "o_uint_64": nil},
				o:  util.PtrUint(1),
				l:  util.PtrUint(3),
				y:  []string{"-id"},
				oo: []request.Option{request.WithField{"id", "o_absent_0", "o_string_1"}},
			},
			want: 4,
			want1: &test.ObjectList{
				{
					ID:      3,
					String1: util.PtrString("white"),
				}, {
					ID:      2,
					String1: util.PtrString("black"),
				}, {
					ID:      1,
					String1: util.PtrString("red"),
				},
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				i:  &test.ObjectList{},
				f:  filter.Eq{"o_float_32": "100", "o_float_64": "3.14", "o_int_16": "16", "o_string_1": "red", "o_string_2": "hello", "o_time_0": "1970-01-01T00:00:00Z", "o_bool": "true", "o_null": nil},
				o:  nil,
				l:  nil,
				y:  nil,
				oo: nil,
			},
			want: 1,
			want1: &test.ObjectList{{
				ID:      1,
				Bool:    util.PtrBool(true),
				Float32: 100,
				Float64: util.PtrFloat64(3.14),
				Int:     0,
				Int16:   util.PtrInt16(16),
				Null:    nil,
				String1: util.PtrString("red"),
				String2: "hello",
				String3: "",
				Uint64:  nil,
				UUID1:   uuid.UUID{},
				UUID2:   util.PtrUUID(uuid.UUID{}),
				UUID3:   nil,
				UUID4:   uuid.UUID{},
				Time0:   time.Unix(0, 0),
				Time1:   time.Unix(0, 0),
				Time2:   util.PtrTime(time.Unix(0, 0)),
				Time3:   nil,
				Time4:   time.Time{},
			}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dbx.List(ctx, db, tt.args.i, tt.args.f, tt.args.o, tt.args.l, tt.args.y, tt.args.oo...)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want1, tt.args.i)
		})
	}
}

func (db DB) TestPut(t *testing.T) {
	ctx := context.TODO()
	type args struct {
		o  dbx.Object
		oo []request.Option
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
				o: &test.Object{
					ID:      9,
					Bool:    util.PtrBool(false),
					Float64: util.PtrFloat64(0),
					Int16:   util.PtrInt16(0),
				},
				oo: []request.Option{request.PutCreate},
			},
			want: &test.Object{
				ID:      9,
				Bool:    util.PtrBool(false),
				Float64: util.PtrFloat64(0),
				Int16:   util.PtrInt16(0),
				String1: util.PtrString("green"),
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o: &test.Object{
					ID:      9,
					Bool:    util.PtrBool(true),
					Float64: util.PtrFloat64(1e1),
					Int16:   util.PtrInt16(1),
					String3: "orange",
				},
				oo: []request.Option{request.WithField{"o_bool", "o_string_3", "o_float_64", "o_int_16"}},
			},
			want: &test.Object{
				ID:      9,
				Bool:    util.PtrBool(true),
				Float64: util.PtrFloat64(1e1),
				Int16:   util.PtrInt16(1),
				String1: util.PtrString("green"),
				String3: "orange",
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o: &test.Object{
					ID:      9,
					Bool:    util.PtrBool(true),
					Float64: util.PtrFloat64(1e2),
					Int16:   util.PtrInt16(2),
					String1: util.PtrString("orange"),
				},
				oo: []request.Option{request.WithField{"o_bool", "o_string_1", "o_float_64", "o_int_16"}},
			},
			want: &test.Object{
				ID:      9,
				Bool:    util.PtrBool(true),
				Float64: util.PtrFloat64(1e2),
				Int16:   util.PtrInt16(2),
				String1: util.PtrString("orange"),
				String3: "orange",
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o: &test.Object{
					ID:      9,
					Bool:    util.PtrBool(true),
					Float64: util.PtrFloat64(1e3),
					Int16:   util.PtrInt16(3),
					String3: "green",
				},
				oo: []request.Option{},
			},
			want: &test.Object{
				ID:      9,
				Bool:    util.PtrBool(true),
				Float64: util.PtrFloat64(1e3),
				Int16:   util.PtrInt16(3),
				String1: util.PtrString("green"),
				String3: "green",
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o: &test.Object{
					ID:      9,
					Bool:    util.PtrBool(false),
					Float64: util.PtrFloat64(1e4),
					Int16:   util.PtrInt16(4),
					String3: "orange",
				},
				oo: []request.Option{request.PutUpdate, request.WithField{"o_bool", "o_string_3"}},
			},
			want: &test.Object{
				ID:      9,
				Bool:    util.PtrBool(false),
				Float64: util.PtrFloat64(1e3),
				Int16:   util.PtrInt16(3),
				String1: util.PtrString("green"),
				String3: "orange",
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o: &test.Object{
					ID:      7,
					Bool:    util.PtrBool(true),
					Float64: util.PtrFloat64(1e3),
					Int16:   util.PtrInt16(3),
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
				o: &test.Object{
					ID:      7,
					Bool:    util.PtrBool(false),
					Float64: util.PtrFloat64(1e4),
					Int16:   util.PtrInt16(4),
					String3: "orange",
				},
				oo: []request.Option{request.PutUpdate, request.WithField{"o_bool", "o_string_3"}},
			},
			want:    nil,
			wantErr: sql.ErrNoRows,
		},
	}
	var ids test.Map
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dbx.Put(ctx, db, tt.args.o, tt.args.oo...)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				id, _ := tt.args.o.Get(0)
				t.Log(ids.Add(tt.args.o.Table(), id))
				for i, n := range tt.args.o.Names() {
					v1, n1 := tt.want.Get(i)
					v2, n2 := tt.args.o.Get(i)
					require.Equal(t, v1, v2, n)
					require.Equal(t, n1, n2, n)
				}
				t.Log(tt.args.o.Get(15))
				t.Log(tt.args.o.Get(16))
			}
		})
	}
	t.Cleanup(func() {
		for k, v := range ids.M {
			s := fmt.Sprintf(`DELETE FROM`+` %q`+` WHERE "id" =ANY($1)`, k)
			_, _ = db.Exec(s, v)
		}
	})
}

type Scan struct {
	bytes.Buffer
}

func (s *Scan) Scan(v any) error {
	_, _ = fmt.Fprintf(&s.Buffer, "%T::[%v],", v, v)
	return nil
}

var _ sql.Scanner = &Scan{}

func (db DB) TestScan(t *testing.T) {
	var x Scan
	row := db.QueryRow(`SELECT 0::int2, '11d7c4b2-9546-4296-b97b-f93a068f3d72'::uuid, '2025-02-23 05:32:48.899094'::timestamp, '123'::bytea`)
	err := row.Scan(&x, &x, &x, &x)
	require.NoError(t, err)
	require.Equal(t, "int64::[0],string::[11d7c4b2-9546-4296-b97b-f93a068f3d72],time.Time::[2025-02-23 05:32:48.899094 +0000 UTC],[]uint8::[[49 50 51]],", x.String())

	x.Reset()

	var o test.Object
	err = json.NewEncoder(&x).Encode(o)
	require.NoError(t, err)
	require.Equal(t, "", x.String())
}
