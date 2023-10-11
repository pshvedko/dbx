package filter

import (
	"github.com/pshvedko/db/driver/postgres"
	"reflect"
	"testing"
)

func TestFilter_To(t *testing.T) {
	tests := []struct {
		name    string
		f       Filter
		b       Builder
		want    string
		want1   map[string]any
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:  "",
			f:     Eq{},
			b:     postgres.NewBuilder(),
			want:  "",
			want1: map[string]any{},
		},
		{
			name:  "",
			f:     Eq{"int": 0},
			b:     postgres.NewBuilder(),
			want:  "int = :int0",
			want1: map[string]any{"int0": 0},
		},
		{
			name:  "",
			f:     Eq{"int": 0, "float": 3.14},
			b:     postgres.NewBuilder(),
			want:  "( float = :float0 AND int = :int0 )",
			want1: map[string]any{"int0": 0, "float0": 3.14},
		},
		{
			name:  "",
			f:     Eq{"int": 0, "float": 3.14, "string": "one"},
			b:     postgres.NewBuilder(),
			want:  "( float = :float0 AND int = :int0 AND string = :string0 )",
			want1: map[string]any{"int0": 0, "float0": 3.14, "string0": "one"},
		},
		{
			name:  "",
			f:     Eq{"int": 0, "float": 3.14, "string": "one", "null": nil},
			b:     postgres.NewBuilder(),
			want:  "( float = :float0 AND int = :int0 AND null IS NULL AND string = :string0 )",
			want1: map[string]any{"int0": 0, "float0": 3.14, "string0": "one"},
		},
		{
			name:  "",
			f:     Eq{"int": 0, "float": 3.14, "string": "one", "null": nil, "bool": true},
			b:     postgres.NewBuilder(),
			want:  "( bool IS TRUE AND float = :float0 AND int = :int0 AND null IS NULL AND string = :string0 )",
			want1: map[string]any{"int0": 0, "float0": 3.14, "string0": "one"},
		},
		{
			name:  "",
			f:     And{},
			b:     postgres.NewBuilder(),
			want:  "",
			want1: map[string]any{},
		},
		{
			name:  "",
			f:     And{Eq{"int": 0}},
			b:     postgres.NewBuilder(),
			want:  "int = :int0",
			want1: map[string]any{"int0": 0},
		},
		{
			name:  "",
			f:     And{Eq{"int": 1}, Eq{"float": 3.14}},
			b:     postgres.NewBuilder(),
			want:  "( int = :int0 AND float = :float0 )",
			want1: map[string]any{"int0": 1, "float0": 3.14},
		},
		{
			name:  "",
			f:     Or{},
			b:     postgres.NewBuilder(),
			want:  "",
			want1: map[string]any{},
		},
		{
			name:  "",
			f:     Or{Eq{"int": 0}},
			b:     postgres.NewBuilder(),
			want:  "int = :int0",
			want1: map[string]any{"int0": 0},
		},
		{
			name:  "",
			f:     Or{Eq{"int": 0}, Eq{"int": 1}},
			b:     postgres.NewBuilder(),
			want:  "( int = :int0 OR int = :int1 )",
			want1: map[string]any{"int0": 0, "int1": 1},
		},
		{
			name:  "",
			f:     Or{Eq{"int": 0, "bool": false}, Eq{"int": 1}},
			b:     postgres.NewBuilder(),
			want:  "( ( bool IS FALSE AND int = :int0 ) OR int = :int1 )",
			want1: map[string]any{"int0": 0, "int1": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.To(tt.b); (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.b.String()
			if got != tt.want {
				t.Errorf("To() got = %v, want %v", got, tt.want)
			}
			got1 := tt.b.Values()
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("To() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
