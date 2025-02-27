package builder_test

import (
	"reflect"
	"testing"

	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test/model"
	"github.com/pshvedko/dbx/request"
)

func TestConstructor_Select(t *testing.T) {
	type args struct {
		j filter.Projector
		f filter.Filter
		o []request.Option
	}
	o := model.Object{}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []any
		want2   []any
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{"o_int_8": 'A', "o_bool_2": true},
				o: []request.Option{
					request.WithField{"o_bool_1", "o_float_32", "o_int_32", "o_uuid_4", "o_string_1"},
					request.WithDeleted("o_time_4"),
				},
			},
			want:    `SELECT "o"."o_bool_1", "o"."o_float_32", "o"."o_int_32", "o"."o_uuid_4", "o"."o_string_1" FROM "objects" AS "o" WHERE ( ( "o"."o_bool_2" IS TRUE AND "o"."o_int_8" = $1 ) AND "o"."o_time_4" IS NULL )`,
			want1:   []any{1},
			want2:   []any{&o.Bool, &o.Float32, &o.Int, &o.Null, &o.String1},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{"o_int": 1, "o_bool": true},
				o: []request.Option{
					request.WithoutField{
						"o_bool", "o_float_32", "o_int", "o_null", "o_uint_64",
						"o_uuid_1", "o_uuid_2", "o_uuid_3", "o_uuid_4",
						"o_time_0", "o_time_1", "o_time_2", "o_time_3", "o_time_4",
						"o_string_2", "o_string_3", "o_string_1"},
					request.DeletedFree,
					request.WithDeleted("o_time_4"),
				},
			},
			want:    `SELECT "o"."id", "o"."o_float_64", "o"."o_int_16" FROM "objects" AS "o" WHERE ( "o"."o_bool" IS TRUE AND "o"."o_int" = $1 )`,
			want1:   []any{1},
			want2:   []any{&o.ID, &o.Float64, &o.Int16},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := request.NewWithOption(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			_, got, got1, got2, err := r.Constructor().Select(tt.args.j, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Select() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("makeSelect() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("makeSelect() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("makeSelect() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
