package dbx_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pshvedko/dbx"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test"
	"github.com/pshvedko/dbx/internal/test/model"
	"github.com/pshvedko/dbx/request"
	"github.com/pshvedko/dbx/util"
	"github.com/stretchr/testify/require"
	"os"
	"slices"
	"testing"
	"time"
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
				request.WithCreated("time_1"),
				request.WithUpdated("time_2"),
				request.WithDeleted("time_4"),
				request.WithTx{},
			)}, nil
}

func TestDB(t *testing.T) {
	db, err := Open(t)
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

var (
	ID1  = uuid.MustParse(`26dd056a-fddf-4992-9e7f-1af09d610c5c`)
	ID1a = uuid.MustParse(`167b738c-d795-44d2-a8d6-2724e9fda43b`)
	ID2  = uuid.MustParse(`4fc2afae-80b1-4045-845b-c373d180beda`)
	ID3  = uuid.MustParse(`6dc2ed41-03ce-47e3-937e-6b4a7b159ec2`)
	ID4  = uuid.MustParse(`e8d3d9a4-9aba-4303-95bf-4dec730131e6`)
	ID5  = uuid.MustParse(`62549d95-ba11-450b-a4a2-0911afd3a73f`)
	ID6  = uuid.MustParse(`9c0c9146-cb6d-4184-80e4-acd26c2c1881`)
	ID7  = uuid.MustParse(`1ad724f7-7d0e-4ac2-8723-7ec2ed36ad2f`)
	ID8  = uuid.MustParse(`22ea94c5-a7f3-44a3-a108-31d58d94f88e`)
	ID9  = uuid.MustParse(`e1ba1c92-1487-40be-ac87-708d805d6fd7`)
)

func (db DB) TestListAny(t *testing.T) {
	var oo []model.Object
	err := db.SelectContext(context.TODO(), &oo, `
			SELECT "id"
			FROM "objects"
			WHERE "id" =ANY($1)
 			  AND "string_2" =ANY($2)
 			  AND "bool_1" =ANY($3)
			  AND "float_32" =ANY($4)
 			  AND "uuid_2" <>ANY($5)
 			  AND "time_1" <>ALL($6)
			  `,
		filter.Array{ID6, ID1, ID2},
		filter.Array{"amber", "black", "white", "green", "yellow"},
		filter.Array{true},
		filter.Array{0, 1e0, 1e1, 1e2},
		filter.Array{uuid.UUID{}},
		filter.Array{time.Time{}},
	)
	require.NoError(t, err)
	require.ElementsMatch(t, []model.Object{{ID: ID1}, {ID: ID2}, {ID: ID6}}, oo)
}

func (db DB) TestListIn(t *testing.T) {
	var oo model.ObjectList
	total, err := dbx.List(context.TODO(), db, &oo,
		filter.And{
			filter.In{"id": {ID1, ID2, ID3, ID4, ID5}, "string_2": {"amber", "black", "white", "green", "yellow"}},
			filter.Ni{"time_1": {time.Time{}}},
		}, nil, nil, nil, request.WithField{"id"})
	require.NoError(t, err)
	require.EqualValues(t, 5, total)
	require.ElementsMatch(t, model.ObjectList{{ID: ID1}, {ID: ID2}, {ID: ID3}, {ID: ID4}, {ID: ID5}}, oo)
	oo = nil
	total, err = dbx.List(context.TODO(), db, &oo,
		filter.In{"time_3": {"1970-01-01T00:00:00Z"}}, nil, util.PtrUint(5), nil, request.WithField{"id"})
	require.NoError(t, err)
	require.EqualValues(t, 5, total)
	require.ElementsMatch(t, model.ObjectList{{ID: ID1}, {ID: ID2}, {ID: ID3}, {ID: ID4}, {ID: ID5}}, oo)
	oo = nil
	total, err = dbx.List(context.TODO(), db, &oo,
		filter.In{"time_3": {"YESTERDAY", filter.Now(), time.Now(), time.UnixMicro(0)}}, nil, nil, []string{"id"}, request.WithField{"id"})
	require.NoError(t, err)
	require.EqualValues(t, 5, total)
	require.ElementsMatch(t, model.ObjectList{{ID: ID1}, {ID: ID2}, {ID: ID3}, {ID: ID4}, {ID: ID5}}, oo)
}

func (db DB) TestListLike(t *testing.T) {
	var oo model.ObjectList
	total, err := dbx.List(context.TODO(), db, &oo,
		filter.And{
			filter.As{"string_2": "%a%"},
			filter.Lt{"time_1": filter.Now()},
		}, nil, nil, nil, request.WithField{"id", "string_1", "string_2", "string_3", "string_4"}, request.DeletedOnly)
	require.NoError(t, err)
	require.EqualValues(t, 2, total)
	require.ElementsMatch(t, model.ObjectList{
		{ID: ID6, String1: "green", String2: "black", String3: util.PtrString("red")},
		{ID: ID7, String1: "green", String2: "orange", String3: util.PtrString("red")},
	}, oo)
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
				o:  &model.Object{},
				f:  filter.Eq{"id": ID1},
				oo: nil,
			},
			want: &model.Object{
				ID:      ID1,
				UUID2:   ID1a,
				UUID3:   nil,
				UUID4:   nil,
				Bool1:   true,
				Bool2:   true,
				Bool3:   util.Ptr(false),
				Bool4:   nil,
				Float32: 1e2,
				Float64: util.Ptr(3.14),
				Int8:    pgtype.Bits{Bytes: []byte{7}, Len: 8, Valid: true},
				Int16:   16,
				Int32:   util.Ptr(int32(0)),
				Int64:   nil,
				String1: "green",
				String2: "yellow",
				String3: util.PtrString("red"),
				String4: nil,
				Time1:   time.Date(1970, 01, 01, 03, 0, 0, 0, time.Local),
				Time2:   time.Date(2025, 02, 24, 16, 8, 6, 482425000, time.Local),
				Time3:   util.Ptr(time.Unix(0, 0)),
				Time4:   nil,
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o:  &model.Object{},
				f:  filter.Eq{"id": ID1},
				oo: []request.Option{request.WithField{"id", "string_1"}},
			},
			want: &model.Object{
				ID:      ID1,
				String1: "green",
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o:  &model.Object{},
				f:  filter.And{filter.Eq{"id": ID1}, filter.Le{"time_2": "YESTERDAY"}},
				oo: []request.Option{request.WithoutField(model.Object{}.Names())},
			},
			want:    &model.Object{},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o:  &model.Object{},
				f:  filter.Eq{"id": nil},
				oo: nil,
			},
			want:    &model.Object{},
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
				i:  &model.ObjectList{},
				f:  nil,
				o:  nil,
				l:  nil,
				y:  []string{"id"},
				oo: []request.Option{request.WithField{"id"}},
			},
			want:    5,
			want1:   &model.ObjectList{{ID: ID1}, {ID: ID2}, {ID: ID5}, {ID: ID3}, {ID: ID4}},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				i:  &model.ObjectList{},
				f:  nil,
				o:  nil,
				l:  nil,
				y:  []string{"id"},
				oo: []request.Option{request.WithField{"id"}, request.DeletedOnly},
			},
			want:    2,
			want1:   &model.ObjectList{{ID: ID7}, {ID: ID6}},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				i:  &model.ObjectList{},
				f:  nil,
				o:  nil,
				l:  nil,
				y:  []string{"id"},
				oo: []request.Option{request.WithField{"id"}, request.DeletedFree},
			},
			want:    7,
			want1:   &model.ObjectList{{ID: ID7}, {ID: ID1}, {ID: ID2}, {ID: ID5}, {ID: ID3}, {ID: ID6}, {ID: ID4}},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				i:  &model.ObjectList{},
				f:  filter.Eq{"time_1": time.Unix(0, 0), "int_64": nil},
				o:  util.PtrUint(1),
				l:  util.PtrUint(3),
				y:  []string{"-id"},
				oo: []request.Option{request.WithField{"id", "string_2"}},
			},
			want: 5,
			want1: &model.ObjectList{
				{
					ID:      ID3,
					String2: "white",
				}, {
					ID:      ID5,
					String2: "black",
				}, {
					ID:      ID2,
					String2: "amber",
				},
			},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				i:  &model.ObjectList{},
				f:  filter.Eq{"float_32": "100", "float_64": "3.14", "int_16": "16", "string_1": "green", "string_2": "yellow", "time_1": "1970-01-01T00:00:00Z", "bool_1": "true", "int_64": nil},
				o:  nil,
				l:  nil,
				y:  nil,
				oo: nil,
			},
			want: 1,
			want1: &model.ObjectList{{
				ID:      ID1,
				UUID2:   ID1a,
				UUID3:   nil,
				UUID4:   nil,
				Bool1:   true,
				Bool2:   true,
				Bool3:   util.Ptr(false),
				Bool4:   nil,
				Float32: 1e2,
				Float64: util.Ptr(3.14),
				Int8:    pgtype.Bits{Bytes: []byte{7}, Len: 8, Valid: true},
				Int16:   16,
				Int32:   util.Ptr(int32(0)),
				Int64:   nil,
				String1: "green",
				String2: "yellow",
				String3: util.Ptr("red"),
				String4: nil,
				Time1:   time.Date(1970, 01, 01, 03, 0, 0, 0, time.Local),
				Time2:   time.Date(2025, 02, 24, 16, 8, 6, 482425000, time.Local),
				Time3:   util.Ptr(time.Unix(0, 0)),
				Time4:   nil,
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
		wantEx  []int
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				o: &model.Object{
					ID:    ID8,
					UUID2: ID9,
					UUID4: util.PtrUUID(ID5),
					Int8:  pgtype.Bits{Bytes: []byte{7}, Len: 8, Valid: true},
					Time2: time.Unix(1, 1),
				},
				oo: []request.Option{request.PutCreate},
			},
			want: &model.Object{
				ID:      ID8,
				UUID2:   ID9,
				UUID4:   util.PtrUUID(ID5),
				Int8:    pgtype.Bits{Bytes: []byte{7}, Len: 8, Valid: true},
				String3: util.PtrString("red"),
			},
			wantEx:  []int{2, 18, 19, 20},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o: &model.Object{
					ID:    ID8,
					UUID2: ID9,
					UUID4: util.PtrUUID(ID5),
					Int8:  pgtype.Bits{Bytes: []byte{1}, Len: 8, Valid: true},
					Time1: time.Unix(2, 2),
					Time2: time.Unix(2, 2),
				},
				oo: []request.Option{
					request.WithField{"id", "uuid_2", "uuid_3", "uuid_4", "int_8", "string_3", "time_1", "time_2", "time_3"},
				},
			},
			want: &model.Object{
				ID:      ID8,
				UUID2:   ID9,
				UUID4:   util.PtrUUID(ID5),
				Int8:    pgtype.Bits{Bytes: []byte{1}, Len: 8, Valid: true},
				String3: util.PtrString("red"),
			},
			wantEx:  []int{2, 18, 19, 20},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o: &model.Object{
					ID:    ID8,
					UUID2: ID9,
					UUID4: util.PtrUUID(ID5),
					Int8:  pgtype.Bits{Bytes: []byte{1}, Len: 8, Valid: true},
					Time1: time.Unix(2, 2),
					Time2: time.Unix(2, 2),
				},
				oo: []request.Option{
					request.WithoutField{"id", "uuid_2", "uuid_3", "uuid_4", "time_1", "time_2", "time_3", "time_4"},
					request.WithoutField{"bool_1", "bool_2", "bool_3", "bool_4", "float_32", "float_64"},
					request.WithoutField{"int_8", "int_16", "int_32", "int_64"},
					request.WithoutField{"string_1", "string_2", "string_3", "string_4"},
				},
			},
			want: &model.Object{
				ID:      ID8,
				UUID2:   ID9,
				UUID4:   util.PtrUUID(ID5),
				Int8:    pgtype.Bits{Bytes: []byte{1}, Len: 8, Valid: true},
				String3: util.PtrString("red"),
			},
			wantEx:  []int{2, 16, 18, 19, 20},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o: &model.Object{
					ID:      ID9,
					UUID2:   ID9,
					UUID4:   util.PtrUUID(ID5),
					Bool3:   util.PtrBool(false),
					Float64: util.PtrFloat64(1e4),
					Int16:   4,
					String2: "orange",
				},
				oo: []request.Option{request.PutUpdate},
			},
			want:    nil,
			wantEx:  []int{2, 16, 18, 19, 20},
			wantErr: sql.ErrNoRows,
		},
		{
			name: "",
			args: args{
				o: &model.Object{
					ID:      ID8,
					UUID2:   ID9,
					UUID4:   util.PtrUUID(ID5),
					Bool3:   util.PtrBool(true),
					Float64: util.PtrFloat64(1e4),
					Int16:   4,
					String2: "orange",
				},
				oo: []request.Option{request.PutUpdate},
			},
			want: &model.Object{
				ID:      ID8,
				UUID2:   ID9,
				UUID4:   util.PtrUUID(ID5),
				Bool3:   util.PtrBool(true),
				Int8:    pgtype.Bits{Bytes: []byte{1}, Len: 8, Valid: true},
				Float64: util.PtrFloat64(1e4),
				Int16:   4,
				String2: "orange",
			},
			wantEx:  []int{2, 16, 18, 19, 20},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o: &model.Object{
					ID:      ID8,
					UUID2:   ID9,
					UUID4:   util.PtrUUID(ID5),
					Bool3:   util.PtrBool(false),
					Float64: util.PtrFloat64(111.111),
					Int16:   11,
					String2: "magenta",
				},
				oo: []request.Option{
					request.PutUpdate,
					request.WithField{"id"},
					request.WithReturnField{"id", "bool_3", "float_64", "int_8", "int_16", "string_2"}},
			},
			want: &model.Object{
				ID:      ID8,
				UUID2:   ID9,
				UUID4:   util.PtrUUID(ID5),
				Bool3:   util.PtrBool(true),
				Int8:    pgtype.Bits{Bytes: []byte{1}, Len: 8, Valid: true},
				Float64: util.PtrFloat64(1e4),
				Int16:   4,
				String2: "orange",
			},
			wantEx:  []int{2, 16, 18, 19, 20},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				o: &model.Object{
					ID:      ID8,
					UUID2:   ID9,
					UUID4:   util.PtrUUID(ID5),
					Bool3:   util.PtrBool(false),
					Float64: util.PtrFloat64(111.111),
					Int16:   11,
					String2: "magenta",
				},
				oo: []request.Option{
					request.PutUpdate,
					request.WithField{"bool_3"},
					request.WithReturnField{"id", "bool_3", "float_64", "int_8", "int_16", "string_2"}},
			},
			want: &model.Object{
				ID:      ID8,
				UUID2:   ID9,
				UUID4:   util.PtrUUID(ID5),
				Bool3:   util.PtrBool(false),
				Int8:    pgtype.Bits{Bytes: []byte{1}, Len: 8, Valid: true},
				Float64: util.PtrFloat64(1e4),
				Int16:   4,
				String2: "orange",
			},
			wantEx:  []int{2, 16, 18, 19, 20},
			wantErr: nil,
		},
	}
	var ids test.Map
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dbx.Put(ctx, db, tt.args.o, tt.args.oo...)
			require.ErrorIs(t, tt.wantErr, err)
			require.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil {
				id, _ := tt.args.o.Get(0)
				t.Log(ids.Add(tt.args.o.Table(), id))
				for i, n := range tt.args.o.Names() {
					if slices.Contains(tt.wantEx, i) {
						t.Log(tt.args.o.Get(i))
						continue
					}
					v1, n1 := tt.want.Get(i)
					v2, n2 := tt.args.o.Get(i)
					require.Equal(t, v1, v2, "[%v] %v", i, n)
					require.Equal(t, n1, n2, "[%v] %v", i, n)
				}
			}
		})
	}
	t.Cleanup(func() {
		for k, v := range ids.M {
			s := fmt.Sprintf(`DELETE FROM`+` %q WHERE "id" =ANY($1)`, k)
			_, _ = db.Exec(s, v)
		}
	})
}

func ExampleUnmarshal() {
	var a, b pgtype.FlatArray[string]
	err := json.Unmarshal([]byte(`["a", "b"]`), &a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a, a == nil)
	err = json.Unmarshal([]byte(`null`), &a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b, b == nil)
	// Output:
	// [a b] false
	// [] true
}
