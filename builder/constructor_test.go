package builder_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test/model"
	"github.com/pshvedko/dbx/request"
	"github.com/pshvedko/dbx/util"
)

func TestConstructor_Select(t *testing.T) {
	var o model.Object
	type args struct {
		j filter.Projector
		f filter.Filter
		o []request.Option
		b *uint
		l *uint
		y builder.Order
	}
	tests := []struct {
		name     string
		args     args
		want     string
		want1    []any
		want2    []any
		want3    string
		want4    int
		wantErr  error
		wantErr1 error
		wantErr2 error
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{"int_8": 'A', "bool_2": true},
				o: []request.Option{
					request.WithField{"bool_1", "float_32", "int_32", "uuid_4", "string_1"},
					request.WithDeleted("time_4"),
				},
			},
			want:  `SELECT "o"."uuid_4", "o"."bool_1", "o"."float_32", "o"."int_32", "o"."string_1" FROM "objects" AS "o" WHERE ( ( "o"."bool_2" IS TRUE AND "o"."int_8" = $1 ) AND "o"."time_4" IS NULL )`,
			want1: []any{'A'},
			want2: []any{&o.UUID4, &o.Bool1, &o.Float32, &o.Int32, &o.String1},
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{"int_32": 1, "bool_1": true},
				o: []request.Option{
					request.WithoutField{
						"uuid_2", "uuid_3", "uuid_4",
						"bool_1", "bool_2", "bool_3", "bool_4",
						"float_32", "int_8", "int_32", "int_64",
						"string_1", "string_2", "string_3", "string_4",
						"time_1", "time_2", "time_3", "time_4",
					},
					request.DeletedFree,
					request.WithDeleted("time_4"),
				},
			},
			want:  `SELECT "o"."id", "o"."float_64", "o"."int_16" FROM "objects" AS "o" WHERE ( "o"."bool_1" IS TRUE AND "o"."int_32" = $1 )`,
			want1: []any{1},
			want2: []any{&o.ID, &o.Float64, &o.Int16},
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{},
				o: []request.Option{request.WithField{"id"}, request.WithCount(), request.WithoutReturning()},
				b: util.PtrUint(10),
				l: util.PtrUint(20),
				y: nil,
			},
			want:  `SELECT FROM "objects" AS "o" WHERE TRUE OFFSET $1 LIMIT $2`,
			want1: []any{uint(10), uint(20)},
			want2: []any{},
			want3: `SELECT COUNT(*) FROM "objects" AS "o" WHERE TRUE`,
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{},
				o: []request.Option{request.WithField{"id"}, request.WithoutReturning(), request.WithoutCount()},
				b: util.PtrUint(10),
				l: util.PtrUint(20),
				y: nil,
			},
			want:  `SELECT FROM "objects" AS "o" WHERE TRUE OFFSET $1 LIMIT $2`,
			want1: []any{uint(10), uint(20)},
			want2: []any{},
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{},
				o: []request.Option{request.WithField{"id"}},
				y: []any{"+"},
			},
			wantErr1: builder.Order{"+"},
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{},
				o: []request.Option{request.WithField{"id"}},
				y: []any{"any"},
			},
			wantErr1: fmt.Errorf("unknown column: any"),
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{},
				o: []request.Option{request.WithField{"id"}},
				y: []any{"-id"},
			},
			want:  `SELECT "o"."id" FROM "objects" AS "o" WHERE TRUE ORDER BY "o"."id" DESC`,
			want1: []any{},
			want2: []any{&o.ID},
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{},
				o: []request.Option{request.WithField{"id"}},
				y: []any{"1"},
			},
			want:  `SELECT "o"."id" FROM "objects" AS "o" WHERE TRUE ORDER BY 1`,
			want1: []any{},
			want2: []any{&o.ID},
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{},
				o: []request.Option{request.WithField{"id"}},
				y: []any{-1},
			},
			want:  `SELECT "o"."id" FROM "objects" AS "o" WHERE TRUE ORDER BY 1 DESC`,
			want1: []any{},
			want2: []any{&o.ID},
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{},
				o: []request.Option{request.WithField{"id"}},
				y: []any{2},
			},
			wantErr1: fmt.Errorf("illegal position: 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := request.NewWithOption(tt.args.o)
			require.ErrorIs(t, err, tt.wantErr)
			z, got, got1, got2, err := r.Constructor().Range(tt.args.b, tt.args.l).Sort(tt.args.y).Select(tt.args.j, tt.args.f)
			t.Log(got)
			require.Equal(t, err, tt.wantErr1)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want1, got1)
			require.Equal(t, tt.want2, got2)
			require.Equal(t, z == nil, tt.want3 == "")
			if z == nil {
				return
			}
			got3, got4, err := z.Count()
			t.Log(got3)
			require.Equal(t, err, tt.wantErr2)
			require.Equal(t, tt.want3, got3)
			require.Equal(t, tt.want4, got4)
		})
	}
}

func TestConstructor_Insert(t *testing.T) {
	var o model.Object
	type args struct {
		j filter.Projector
		o []request.Option
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []any
		want2   []any
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				j: &o,
				o: []request.Option{
					request.WithCreated("time_1"),
					request.WithUpdated("time_2"),
					request.WithDeleted("time_4"),
					request.PutModify,
					request.WithField{"id"},
					request.DeletedNone,
				},
			},
			want:    `INSERT INTO "objects" AS "o" ( "uuid_2", "uuid_4", "bool_1", "bool_2", "bool_3", "bool_4", "float_32", "float_64", "int_16", "int_32", "int_64", "string_1", "string_2", "string_4" ) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14 ) ON CONFLICT ( "id" ) DO UPDATE SET "uuid_2" = EXCLUDED."uuid_2", "uuid_4" = EXCLUDED."uuid_4", "bool_1" = EXCLUDED."bool_1", "bool_2" = EXCLUDED."bool_2", "bool_3" = EXCLUDED."bool_3", "bool_4" = EXCLUDED."bool_4", "float_32" = EXCLUDED."float_32", "float_64" = EXCLUDED."float_64", "int_16" = EXCLUDED."int_16", "int_32" = EXCLUDED."int_32", "int_64" = EXCLUDED."int_64", "string_1" = EXCLUDED."string_1", "string_2" = EXCLUDED."string_2", "string_4" = EXCLUDED."string_4", "time_2" = DEFAULT WHERE "o"."time_4" IS NULL RETURNING "o"."id"`,
			want1:   []any{o.UUID2, nil, o.Bool1, o.Bool2, nil, nil, o.Float32, nil, o.Int16, nil, nil, o.String1, o.String2, nil},
			want2:   []any{&o.ID},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				j: &o,
				o: []request.Option{
					request.WithCreated("time_1"),
					request.WithUpdated("time_2"),
					request.WithDeleted("time_4"),
					request.PutModify,
					request.WithField{"id"},
					request.DeletedOnly,
				},
			},
			want:    `INSERT INTO "objects" AS "o" ( "uuid_2", "uuid_4", "bool_1", "bool_2", "bool_3", "bool_4", "float_32", "float_64", "int_16", "int_32", "int_64", "string_1", "string_2", "string_4" ) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14 ) ON CONFLICT ( "id" ) DO UPDATE SET "uuid_2" = EXCLUDED."uuid_2", "uuid_4" = EXCLUDED."uuid_4", "bool_1" = EXCLUDED."bool_1", "bool_2" = EXCLUDED."bool_2", "bool_3" = EXCLUDED."bool_3", "bool_4" = EXCLUDED."bool_4", "float_32" = EXCLUDED."float_32", "float_64" = EXCLUDED."float_64", "int_16" = EXCLUDED."int_16", "int_32" = EXCLUDED."int_32", "int_64" = EXCLUDED."int_64", "string_1" = EXCLUDED."string_1", "string_2" = EXCLUDED."string_2", "string_4" = EXCLUDED."string_4", "time_2" = DEFAULT WHERE "o"."time_4" IS NOT NULL RETURNING "o"."id"`,
			want1:   []any{o.UUID2, nil, o.Bool1, o.Bool2, nil, nil, o.Float32, nil, o.Int16, nil, nil, o.String1, o.String2, nil},
			want2:   []any{&o.ID},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				j: &o,
				o: []request.Option{
					request.WithCreated("time_1"),
					request.WithUpdated("time_2"),
					request.WithDeleted("time_4"),
					request.PutModify,
					request.WithField{"id"},
					request.DeletedFree,
				},
			},
			want:    `INSERT INTO "objects" AS "o" ( "uuid_2", "uuid_4", "bool_1", "bool_2", "bool_3", "bool_4", "float_32", "float_64", "int_16", "int_32", "int_64", "string_1", "string_2", "string_4" ) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14 ) ON CONFLICT ( "id" ) DO UPDATE SET "uuid_2" = EXCLUDED."uuid_2", "uuid_4" = EXCLUDED."uuid_4", "bool_1" = EXCLUDED."bool_1", "bool_2" = EXCLUDED."bool_2", "bool_3" = EXCLUDED."bool_3", "bool_4" = EXCLUDED."bool_4", "float_32" = EXCLUDED."float_32", "float_64" = EXCLUDED."float_64", "int_16" = EXCLUDED."int_16", "int_32" = EXCLUDED."int_32", "int_64" = EXCLUDED."int_64", "string_1" = EXCLUDED."string_1", "string_2" = EXCLUDED."string_2", "string_4" = EXCLUDED."string_4", "time_2" = DEFAULT RETURNING "o"."id"`,
			want1:   []any{o.UUID2, nil, o.Bool1, o.Bool2, nil, nil, o.Float32, nil, o.Int16, nil, nil, o.String1, o.String2, nil},
			want2:   []any{&o.ID},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := request.NewWithOption(tt.args.o)
			require.ErrorIs(t, err, tt.wantErr)
			got, got1, got2, err := r.Constructor().Insert(tt.args.j)
			t.Log(got)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want1, got1)
			require.Equal(t, tt.want2, got2)
		})
	}
}

func TestConstructor_Update(t *testing.T) {
	o := model.Object{ID: uuid.New()}
	type args struct {
		j filter.Projector
		o []request.Option
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []any
		want2   []any
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				j: &o,
				o: []request.Option{
					request.WithCreated("time_1"),
					request.WithUpdated("time_2"),
					request.WithDeleted("time_4"),
					request.PutUpdate,
					request.WithField{"id", "uuid_3"},
					request.DeletedNone,
				},
			},
			want:    `UPDATE "objects" AS "o" SET "time_2" = DEFAULT WHERE ( "o"."id" = $1 AND "o"."time_4" IS NULL ) RETURNING "o"."id", "o"."uuid_3"`,
			want1:   []any{o.ID},
			want2:   []any{&o.ID, &o.UUID3},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				j: &o,
				o: []request.Option{
					request.WithCreated("time_1"),
					request.WithUpdated("time_2"),
					request.WithDeleted("time_4"),
					request.PutUpdate,
					request.WithField{"id", "uuid_3"},
					request.DeletedOnly,
				},
			},
			want:    `UPDATE "objects" AS "o" SET "time_2" = DEFAULT WHERE ( "o"."id" = $1 AND "o"."time_4" IS NOT NULL ) RETURNING "o"."id", "o"."uuid_3"`,
			want1:   []any{o.ID},
			want2:   []any{&o.ID, &o.UUID3},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				j: &o,
				o: []request.Option{
					request.WithCreated("time_1"),
					request.WithUpdated("time_2"),
					request.WithDeleted("time_4"),
					request.PutUpdate,
					request.WithField{"id", "uuid_3"},
					request.WithReturnField{"id", "uuid_3", "time_1"},
					request.DeletedFree,
				},
			},
			want:    `UPDATE "objects" AS "o" SET "time_2" = DEFAULT WHERE "o"."id" = $1 RETURNING "o"."id", "o"."uuid_3", "o"."time_1"`,
			want1:   []any{o.ID},
			want2:   []any{&o.ID, &o.UUID3, &o.Time1},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := request.NewWithOption(tt.args.o)
			require.ErrorIs(t, err, tt.wantErr)
			got, got1, got2, err := r.Constructor().Update(tt.args.j)
			t.Log(got)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want1, got1)
			require.Equal(t, tt.want2, got2)
		})
	}
}

func TestConstructor_Delete(t *testing.T) {
	var o model.Object
	type args struct {
		j filter.Projector
		f filter.Filter
		o []request.Option
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []any
		want2   []any
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				j: &o,
				f: nil,
				o: nil,
			},
			want:    `DELETE FROM "objects" AS "o" RETURNING "o"."id", "o"."uuid_2", "o"."uuid_3", "o"."uuid_4", "o"."bool_1", "o"."bool_2", "o"."bool_3", "o"."bool_4", "o"."float_32", "o"."float_64", "o"."int_8", "o"."int_16", "o"."int_32", "o"."int_64", "o"."string_1", "o"."string_2", "o"."string_3", "o"."string_4", "o"."time_1", "o"."time_2", "o"."time_3", "o"."time_4"`,
			want1:   []any{},
			want2:   o.Places(),
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				j: &o,
				f: nil,
				o: []request.Option{request.WithField{"id", "time_1"}},
			},
			want:    `DELETE FROM "objects" AS "o" RETURNING "o"."id", "o"."time_1"`,
			want1:   []any{},
			want2:   []any{&o.ID, &o.Time1},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				j: &o,
				f: nil,
				o: []request.Option{request.WithField{"id", "time_1"}, request.WithReturnField{"id"}},
			},
			want:    `DELETE FROM "objects" AS "o" RETURNING "o"."id"`,
			want1:   []any{},
			want2:   []any{&o.ID},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				j: &o,
				f: nil,
				o: []request.Option{request.WithField{"id", "time_1"}, request.WithReturnField{"id"}, request.WithoutReturning()},
			},
			want:    `DELETE FROM "objects" AS "o" RETURNING 1`,
			want1:   []any{},
			want2:   []any{new(int64)},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{"id": uuid.UUID{}},
				o: []request.Option{request.WithField{"id"}},
			},
			want:    `DELETE FROM "objects" AS "o" WHERE "o"."id" = $1 RETURNING "o"."id"`,
			want1:   []any{uuid.UUID{}},
			want2:   []any{&o.ID},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{"id": uuid.UUID{0x48, 0xf9, 0xf2, 0x4e, 0x12, 0xbe, 0x48, 0xe8, 0xa1, 0x61, 0xe0, 0x9, 0xc4, 0xb7, 0x28, 0xd3}},
				o: []request.Option{
					request.WithDeleted("time_4"),
					request.WithField{"id"},
				},
			},
			want:    `UPDATE "objects" AS "o" SET "time_4" = NOW() WHERE ( "o"."id" = $1 AND "o"."time_4" IS NULL ) RETURNING "o"."id"`,
			want1:   []any{uuid.UUID{0x48, 0xf9, 0xf2, 0x4e, 0x12, 0xbe, 0x48, 0xe8, 0xa1, 0x61, 0xe0, 0x9, 0xc4, 0xb7, 0x28, 0xd3}},
			want2:   []any{&o.ID},
			wantErr: nil,
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{"id": uuid.UUID{0x48, 0xf9, 0xf2, 0x4e, 0x12, 0xbe, 0x48, 0xe8, 0xa1, 0x61, 0xe0, 0x9, 0xc4, 0xb7, 0x28, 0xd3}},
				o: []request.Option{
					request.WithCreated("time_1"),
					request.WithUpdated("time_2"),
					request.WithDeleted("time_4"),
					request.WithField{"id"},
				},
			},
			want:    `UPDATE "objects" AS "o" SET "time_2" = DEFAULT, "time_4" = NOW() WHERE ( "o"."id" = $1 AND "o"."time_4" IS NULL ) RETURNING "o"."id"`,
			want1:   []any{uuid.UUID{0x48, 0xf9, 0xf2, 0x4e, 0x12, 0xbe, 0x48, 0xe8, 0xa1, 0x61, 0xe0, 0x9, 0xc4, 0xb7, 0x28, 0xd3}},
			want2:   []any{&o.ID},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := request.NewWithOption(tt.args.o)
			require.ErrorIs(t, err, tt.wantErr)
			got, got1, got2, err := r.Constructor().Delete(tt.args.j, tt.args.f)
			t.Log(got)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want1, got1)
			require.Equal(t, tt.want2, got2)
		})
	}
}

func BenchmarkConstructor_Select(b *testing.B) {
	r, err := request.NewWithOption([]request.Option{
		request.WithCreated("time_1"),
		request.WithUpdated("time_2"),
		request.WithDeleted("time_4"),
	})
	if err != nil {
		b.Fatal(err)
	}
	var j model.Object
	var o, l uint = 100, 200
	var f filter.Filter = filter.And{
		filter.Eq{"int_8": `111`, "bool_2": true},
		filter.In{"string_1": []any{"yellow", "green"}, "int_16": []any{1, 2, 3}, "float_32": []any{0.0}},
		filter.Or{filter.Gt{"int_64": 0}, filter.Na{"string_2": `%ing`}},
	}
	for i := 0; i < b.N; i++ {
		_, _, _, _, err = r.Constructor().Range(&o, &l).Sort([]any{-1, "time_1"}).Select(&j, f)
		if err != nil {
			b.Fatal(err)
		}
	}
	// SELECT "o"."id", "o"."uuid_2", "o"."uuid_3", "o"."uuid_4", "o"."bool_1", "o"."bool_2", "o"."bool_3", "o"."bool_4", "o"."float_32","o"."float_64", "o"."int_8", "o"."int_16", "o"."int_32", "o"."int_64", "o"."string_1", "o"."string_2", "o"."string_3", "o"."string_4", "o"."time_1", "o"."time_2", "o"."time_3", "o"."time_4" FROM "objects" AS "o" WHERE ( ( ( "o"."bool_2" IS TRUE AND "o"."int_8" = $1 ) AND ( "o"."float_32" = ANY($2) AND "o"."int_16" = ANY($3) AND "o"."string_1" = ANY($4) ) AND ( "o"."int_64" > NULL OR "o"."string_2" NOT LIKE $5 ) ) AND "o"."time_4" IS NULL ) ORDER BY 1 DESC, "o"."time_1" OFFSET $6 LIMIT $7
}
