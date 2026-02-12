package builder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test/model"
	"github.com/pshvedko/dbx/request"
)

func TestConstructor_Select(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := request.NewWithOption(tt.args.o)
			require.ErrorIs(t, err, tt.wantErr)
			_, got, got1, got2, err := r.Constructor().Select(tt.args.j, tt.args.f)
			t.Log(got)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want1, got1)
			require.Equal(t, tt.want2, got2)
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
			want:    `INSERT INTO "objects" ( "uuid_2", "uuid_4", "bool_1", "bool_2", "bool_3", "bool_4", "float_32", "float_64", "int_16", "int_32", "int_64", "string_1", "string_2", "string_4" ) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14 ) ON CONFLICT ( "id" ) DO UPDATE SET "uuid_2" = EXCLUDED."uuid_2", "uuid_4" = EXCLUDED."uuid_4", "bool_1" = EXCLUDED."bool_1", "bool_2" = EXCLUDED."bool_2", "bool_3" = EXCLUDED."bool_3", "bool_4" = EXCLUDED."bool_4", "float_32" = EXCLUDED."float_32", "float_64" = EXCLUDED."float_64", "int_16" = EXCLUDED."int_16", "int_32" = EXCLUDED."int_32", "int_64" = EXCLUDED."int_64", "string_1" = EXCLUDED."string_1", "string_2" = EXCLUDED."string_2", "string_4" = EXCLUDED."string_4", "time_2" = DEFAULT WHERE "objects"."time_4" IS NULL RETURNING "id"`,
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
			want:    `INSERT INTO "objects" ( "uuid_2", "uuid_4", "bool_1", "bool_2", "bool_3", "bool_4", "float_32", "float_64", "int_16", "int_32", "int_64", "string_1", "string_2", "string_4" ) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14 ) ON CONFLICT ( "id" ) DO UPDATE SET "uuid_2" = EXCLUDED."uuid_2", "uuid_4" = EXCLUDED."uuid_4", "bool_1" = EXCLUDED."bool_1", "bool_2" = EXCLUDED."bool_2", "bool_3" = EXCLUDED."bool_3", "bool_4" = EXCLUDED."bool_4", "float_32" = EXCLUDED."float_32", "float_64" = EXCLUDED."float_64", "int_16" = EXCLUDED."int_16", "int_32" = EXCLUDED."int_32", "int_64" = EXCLUDED."int_64", "string_1" = EXCLUDED."string_1", "string_2" = EXCLUDED."string_2", "string_4" = EXCLUDED."string_4", "time_2" = DEFAULT WHERE "objects"."time_4" IS NOT NULL RETURNING "id"`,
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
			want:    `INSERT INTO "objects" ( "uuid_2", "uuid_4", "bool_1", "bool_2", "bool_3", "bool_4", "float_32", "float_64", "int_16", "int_32", "int_64", "string_1", "string_2", "string_4" ) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14 ) ON CONFLICT ( "id" ) DO UPDATE SET "uuid_2" = EXCLUDED."uuid_2", "uuid_4" = EXCLUDED."uuid_4", "bool_1" = EXCLUDED."bool_1", "bool_2" = EXCLUDED."bool_2", "bool_3" = EXCLUDED."bool_3", "bool_4" = EXCLUDED."bool_4", "float_32" = EXCLUDED."float_32", "float_64" = EXCLUDED."float_64", "int_16" = EXCLUDED."int_16", "int_32" = EXCLUDED."int_32", "int_64" = EXCLUDED."int_64", "string_1" = EXCLUDED."string_1", "string_2" = EXCLUDED."string_2", "string_4" = EXCLUDED."string_4", "time_2" = DEFAULT RETURNING "id"`,
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
