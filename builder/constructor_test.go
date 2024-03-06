package builder_test

import (
	help "github.com/pshvedko/db/t"
	"reflect"
	"testing"

	"github.com/pshvedko/db/filter"
	"github.com/pshvedko/db/request"
)

func TestConstructor_Select(t *testing.T) {
	type args struct {
		j filter.Projector
		f filter.Filter
		o []request.Option
	}
	o := help.Object{}
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
				f: filter.Eq{"o_int": 1, "o_bool": true},
				o: []request.Option{
					request.WithField{"o_bool", "o_float_32", "o_int", "o_null", "o_string"},
				},
			},
			want:    `SELECT "o_bool", "o_float_32", "o_int", "o_null", "o_string" FROM "objects" WHERE ( "o_bool" IS TRUE AND "o_int" = $1 )`,
			want1:   []any{1},
			want2:   []any{&o.Bool, &o.Float32, &o.Int, &o.Null, &o.String},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				j: &o,
				f: filter.Eq{"o_int": 1, "o_bool": true},
				o: []request.Option{
					request.WithoutField{"id", "o_bool", "o_float_32", "o_float_64", "o_int", "o_int_16", "o_null", "o_string", "o_uint_64"},
				},
			},
			want:    `SELECT "o_uuid_1", "o_uuid_2", "o_uuid_3", "o_uuid_4" FROM "objects" WHERE ( "o_bool" IS TRUE AND "o_int" = $1 )`,
			want1:   []any{1},
			want2:   []any{&o.UUID1, &o.UUID2, &o.UUID3, &o.UUID4},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := request.New(nil, nil, tt.args.o...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, got1, got2, err := r.Constructor().Select(tt.args.j, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Fatalf("makeSelect() error = %v, wantErr %v", err, tt.wantErr)
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
