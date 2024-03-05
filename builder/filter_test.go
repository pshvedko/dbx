package builder_test

import (
	help "github.com/pshvedko/db/t"
	"reflect"
	"testing"

	"github.com/pshvedko/db/builder"
	"github.com/pshvedko/db/filter"
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
			want:  "",
			want1: nil,
		},
		{
			name:  "",
			f:     filter.Eq{"o_int": 0},
			want:  "o_int = $1",
			want1: []any{0},
		},
		{
			name:  "",
			f:     filter.Eq{"o_int": 0, "o_float_64": 3.14},
			want:  "( o_float_64 = $1 AND o_int = $2 )",
			want1: []any{3.14, 0},
		},
		{
			name:  "",
			f:     filter.Eq{"o_int": 0, "o_float_64": 3.14, "o_string": "one"},
			want:  "( o_float_64 = $1 AND o_int = $2 AND o_string = $3 )",
			want1: []any{3.14, 0, "one"},
		},
		{
			name:  "",
			f:     filter.Eq{"o_int": 0, "o_float_64": 3.14, "o_string": "one", "o_null": nil},
			want:  "( o_float_64 = $1 AND o_int = $2 AND o_null IS NULL AND o_string = $3 )",
			want1: []any{3.14, 0, "one"},
		},
		{
			name:  "",
			f:     filter.Eq{"o_int": 0, "o_float_64": 3.14, "o_string": "one", "o_null": nil, "o_bool": true},
			want:  "( o_bool IS TRUE AND o_float_64 = $1 AND o_int = $2 AND o_null IS NULL AND o_string = $3 )",
			want1: []any{3.14, 0, "one"},
		},
		{
			name:  "",
			f:     filter.And{},
			want:  "",
			want1: nil,
		},
		{
			name:  "",
			f:     filter.And{filter.Eq{"o_int": 0}},
			want:  "o_int = $1",
			want1: []any{0},
		},
		{
			name:  "",
			f:     filter.And{filter.Eq{"o_int": 1}, filter.Eq{"o_float_64": 3.14}},
			want:  "( o_int = $1 AND o_float_64 = $2 )",
			want1: []any{1, 3.14},
		},
		{
			name:  "",
			f:     filter.Or{},
			want:  "",
			want1: nil,
		},
		{
			name:  "",
			f:     filter.Or{filter.Eq{"o_int": 0}},
			want:  "o_int = $1",
			want1: []any{0},
		},
		{
			name:  "",
			f:     filter.Or{filter.Eq{"o_int": 0}, filter.Eq{"o_int": 1}},
			want:  "( o_int = $1 OR o_int = $2 )",
			want1: []any{0, 1},
		},
		{
			name:  "",
			f:     filter.Or{filter.Eq{"o_int": 0, "o_bool": false}, filter.Eq{"o_int": 1}},
			want:  "( ( o_bool IS FALSE AND o_int = $1 ) OR o_int = $2 )",
			want1: []any{0, 1},
		},
		{
			name:  "",
			f:     filter.Ne{"o_int": 0, "o_float_64": 3.14, "o_string": "one", "o_null": nil, "o_bool": true},
			want:  "( o_bool IS NOT TRUE AND o_float_64 <> $1 AND o_int <> $2 AND o_null IS NOT NULL AND o_string <> $3 )",
			want1: []any{3.14, 0, "one"},
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
