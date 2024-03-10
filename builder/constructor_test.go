package builder_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/request"
	"github.com/pshvedko/dbx/t"
)

type Object struct {
	help.Object
}

func (o Object) Names() []string { return o.Object.Names()[:8] }
func (o Object) Values() []any   { return o.Object.Values()[:8] }

func TestConstructor_Select(t *testing.T) {
	type args struct {
		j filter.Projector
		f filter.Filter
		o []request.Option
	}
	o := Object{}
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
					request.WithField{"o_bool", "o_float_32", "o_int", "o_null", "o_string_1"},
				},
			},
			want:    `SELECT "o_bool", "o_float_32", "o_int", "o_null", "o_string_1" FROM "objects" WHERE ( "objects"."o_bool" IS TRUE AND "objects"."o_int" = $1 )`,
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
					request.WithoutField{"o_bool", "o_float_32", "o_int", "o_null", "o_string_1"},
				},
			},
			want:    `SELECT "id", "o_float_64", "o_int_16" FROM "objects" WHERE ( "objects"."o_bool" IS TRUE AND "objects"."o_int" = $1 )`,
			want1:   []any{1},
			want2:   []any{&o.ID, &o.Float64, &o.Int16},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx context.Context
			r, err := request.New(ctx, help.DB{}, tt.args.o...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			_, got, got1, got2, err := r.Constructor().Select(tt.args.j, tt.args.f)
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
