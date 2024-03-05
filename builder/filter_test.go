package builder_test

import (
	"reflect"
	"testing"

	"github.com/pshvedko/db/builder"
	"github.com/pshvedko/db/filter"
)

type Object1 struct {
	Bool   bool    `json:"bool,omitempty"`
	Float  float64 `json:"float,omitempty"`
	Int    int     `json:"int,omitempty"`
	Null   any     `json:"null,omitempty"`
	String string  `json:"string,omitempty"`
}

func (o Object1) Table() string {
	return "objects"
}

func (o Object1) Names() []string {
	return []string{"bool", "float", "int", "null", "string"}
}

func (o *Object1) Values() []any {
	return []any{&o.Bool, &o.Float, &o.Int, &o.Null, &o.String}
}

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
			f:     filter.Eq{"int": 0},
			want:  "int = $1",
			want1: []any{0},
		},
		{
			name:  "",
			f:     filter.Eq{"int": 0, "float": 3.14},
			want:  "( float = $1 AND int = $2 )",
			want1: []any{3.14, 0},
		},
		{
			name:  "",
			f:     filter.Eq{"int": 0, "float": 3.14, "string": "one"},
			want:  "( float = $1 AND int = $2 AND string = $3 )",
			want1: []any{3.14, 0, "one"},
		},
		{
			name:  "",
			f:     filter.Eq{"int": 0, "float": 3.14, "string": "one", "null": nil},
			want:  "( float = $1 AND int = $2 AND null IS NULL AND string = $3 )",
			want1: []any{3.14, 0, "one"},
		},
		{
			name:  "",
			f:     filter.Eq{"int": 0, "float": 3.14, "string": "one", "null": nil, "bool": true},
			want:  "( bool IS TRUE AND float = $1 AND int = $2 AND null IS NULL AND string = $3 )",
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
			f:     filter.And{filter.Eq{"int": 0}},
			want:  "int = $1",
			want1: []any{0},
		},
		{
			name:  "",
			f:     filter.And{filter.Eq{"int": 1}, filter.Eq{"float": 3.14}},
			want:  "( int = $1 AND float = $2 )",
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
			f:     filter.Or{filter.Eq{"int": 0}},
			want:  "int = $1",
			want1: []any{0},
		},
		{
			name:  "",
			f:     filter.Or{filter.Eq{"int": 0}, filter.Eq{"int": 1}},
			want:  "( int = $1 OR int = $2 )",
			want1: []any{0, 1},
		},
		{
			name:  "",
			f:     filter.Or{filter.Eq{"int": 0, "bool": false}, filter.Eq{"int": 1}},
			want:  "( ( bool IS FALSE AND int = $1 ) OR int = $2 )",
			want1: []any{0, 1},
		},
		{
			name:  "",
			f:     filter.Ne{"int": 0, "float": 3.14, "string": "one", "null": nil, "bool": true},
			want:  "( bool IS NOT TRUE AND float <> $1 AND int <> $2 AND null IS NOT NULL AND string <> $3 )",
			want1: []any{3.14, 0, "one"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := builder.Filter{}
			if err := tt.f.To(&b, &Object1{}); (err != nil) != tt.wantErr {
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
