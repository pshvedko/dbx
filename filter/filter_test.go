package filter

import (
	"github.com/pshvedko/db/driver/postgres"
	"reflect"
	"testing"
)

func TestEq_SQL(t *testing.T) {
	tests := []struct {
		name   string
		f      Filter
		h      Builder
		want   string
		values map[string]interface{}
	}{
		// TODO: Add test cases.
		{
			name:   "",
			f:      Eq{"int": 0},
			h:      postgres.NewHolder(),
			want:   "int = :int0",
			values: map[string]interface{}{"int0": 0},
		},
		{
			name:   "",
			f:      Eq{"int": 0, "float": 3.14},
			h:      postgres.NewHolder(),
			want:   "float = :float0 AND int = :int0",
			values: map[string]interface{}{"int0": 0, "float0": 3.14},
		},
		{
			name:   "",
			f:      Eq{"int": 0, "float": 3.14, "string": "one"},
			h:      postgres.NewHolder(),
			want:   "float = :float0 AND int = :int0 AND string = :string0",
			values: map[string]interface{}{"int0": 0, "float0": 3.14, "string0": "one"},
		},
		{
			name:   "",
			f:      Eq{"int": 0, "float": 3.14, "string": "one", "null": nil},
			h:      postgres.NewHolder(),
			want:   "float = :float0 AND int = :int0 AND null IS NULL AND string = :string0",
			values: map[string]interface{}{"int0": 0, "float0": 3.14, "string0": "one"},
		},
		{
			name:   "",
			f:      Eq{"int": 0, "float": 3.14, "string": "one", "null": nil, "bool": true},
			h:      postgres.NewHolder(),
			want:   "bool IS TRUE AND float = :float0 AND int = :int0 AND null IS NULL AND string = :string0",
			values: map[string]interface{}{"int0": 0, "float0": 3.14, "string0": "one"},
		},
		{
			name:   "",
			f:      And{},
			h:      postgres.NewHolder(),
			want:   "",
			values: map[string]interface{}{},
		},
		{
			name:   "",
			f:      And{Eq{"int": 0}},
			h:      postgres.NewHolder(),
			want:   "int = :int0",
			values: map[string]interface{}{"int0": 0},
		},
		{
			name:   "",
			f:      And{Eq{"int": 1}, Eq{"float": 3.14}},
			h:      postgres.NewHolder(),
			want:   "float = :float0 AND int = :int0",
			values: map[string]interface{}{"int0": 1, "float0": 3.14},
		},
		{
			name:   "",
			f:      Or{Eq{"int": 0}, Eq{"int": 1}},
			h:      postgres.NewHolder(),
			want:   "int = :int0 OR int = :int1",
			values: map[string]interface{}{"int0": 0, "int1": 1},
		},
		{
			name:   "",
			f:      Or{Eq{"int": 0, "bool": false}, Eq{"int": 1}},
			h:      postgres.NewHolder(),
			want:   "bool IS FALSE AND int = :int0 OR int = :int1",
			values: map[string]interface{}{"int0": 0, "int1": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f.SQL(tt.h)
			if got != tt.want {
				t.Errorf("SQL() got = %v, want %v", got, tt.want)
			}
			got1 := tt.h.Values()
			if !reflect.DeepEqual(got1, tt.values) {
				t.Errorf("SQL() got1 = %v, want %v", got1, tt.values)
			}
		})
	}
}
