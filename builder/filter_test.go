package builder_test

import (
	"reflect"
	"testing"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/t"
)

func TestFilter_To(t *testing.T) {
	tests := []struct {
		name    string
		f       filter.Filter
		want    string
		want1   []any
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:  "",
			f:     filter.Eq{},
			want:  ``,
			want1: nil,
		},
		{
			name:  "",
			f:     filter.Eq{"o_int": 0},
			want:  `"objects"."o_int" = $1`,
			want1: []any{0},
		},
		{
			name:  "",
			f:     filter.Eq{"o_int": 0, "o_float_64": 3.14},
			want:  `( "objects"."o_float_64" = $1 AND "objects"."o_int" = $2 )`,
			want1: []any{3.14, 0},
		},
		{
			name:  "",
			f:     filter.Eq{"o_int": 0, "o_float_64": 3.14, "o_string_1": "one"},
			want:  `( "objects"."o_float_64" = $1 AND "objects"."o_int" = $2 AND "objects"."o_string_1" = $3 )`,
			want1: []any{3.14, 0, "one"},
		},
		{
			name:  "",
			f:     filter.Eq{"o_int": 0, "o_float_64": 3.14, "o_string_1": "one", "o_null": nil},
			want:  `( "objects"."o_float_64" = $1 AND "objects"."o_int" = $2 AND "objects"."o_null" IS NULL AND "objects"."o_string_1" = $3 )`,
			want1: []any{3.14, 0, "one"},
		},
		{
			name:  "",
			f:     filter.Eq{"o_int": 0, "o_float_64": 3.14, "o_string_1": "one", "o_null": nil, "o_bool": true},
			want:  `( "objects"."o_bool" IS TRUE AND "objects"."o_float_64" = $1 AND "objects"."o_int" = $2 AND "objects"."o_null" IS NULL AND "objects"."o_string_1" = $3 )`,
			want1: []any{3.14, 0, "one"},
		},
		{
			name:  "",
			f:     filter.And{},
			want:  ``,
			want1: nil,
		},
		{
			name:  "",
			f:     filter.And{filter.Eq{"o_int": 0}},
			want:  `"objects"."o_int" = $1`,
			want1: []any{0},
		},
		{
			name:  "",
			f:     filter.And{filter.Eq{"o_int": 1}, filter.Eq{"o_float_64": 3.14}},
			want:  `( "objects"."o_int" = $1 AND "objects"."o_float_64" = $2 )`,
			want1: []any{1, 3.14},
		},
		{
			name:  "",
			f:     filter.Or{},
			want:  ``,
			want1: nil,
		},
		{
			name:  "",
			f:     filter.Or{filter.Eq{"o_int": 0}},
			want:  `"objects"."o_int" = $1`,
			want1: []any{0},
		},
		{
			name:  "",
			f:     filter.Or{filter.Eq{"o_int": 0}, filter.Eq{"o_int": 1}},
			want:  `( "objects"."o_int" = $1 OR "objects"."o_int" = $2 )`,
			want1: []any{0, 1},
		},
		{
			name:  "",
			f:     filter.Or{filter.Eq{"o_int": 0, "o_bool": false}, filter.Eq{"o_int": 1}},
			want:  `( ( "objects"."o_bool" IS FALSE AND "objects"."o_int" = $1 ) OR "objects"."o_int" = $2 )`,
			want1: []any{0, 1},
		},
		{
			name:  "",
			f:     filter.Ne{"o_int": 0, "o_float_64": 3.14, "o_string_1": "one", "o_null": nil, "o_bool": true},
			want:  `( "objects"."o_bool" IS NOT TRUE AND "objects"."o_float_64" <> $1 AND "objects"."o_int" <> $2 AND "objects"."o_null" IS NOT NULL AND "objects"."o_string_1" <> $3 )`,
			want1: []any{3.14, 0, "one"},
		},
		{
			name:    "",
			f:       filter.Ge{"o_time_1": filter.Now()},
			want:    `"objects"."o_time_1" >= NOW()`,
			want1:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := builder.Filter{}
			if err := tt.f.To(&b, &help.Object{}); (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, got1 := b.String(), b.Values()
			if got != tt.want {
				t.Errorf("To() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("To() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
