package builder_test

import (
	"bytes"
	"github.com/pshvedko/dbx/internal/test/model"
	"reflect"
	"testing"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
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
			want1: []any{},
		},
		{
			name:    "",
			f:       filter.Eq{"not_exists": 0},
			want:    ``,
			want1:   []any{},
			wantErr: true,
		},
		{
			name:  "",
			f:     filter.Eq{"int_32": 0},
			want:  `"objects"."int_32" = $1`,
			want1: []any{0},
		},
		{
			name:  "",
			f:     filter.Eq{"int_32": 0, "float_64": 3.14},
			want:  `( "objects"."float_64" = $1 AND "objects"."int_32" = $2 )`,
			want1: []any{3.14, 0},
		},
		{
			name:  "",
			f:     filter.Eq{"int_32": 0, "float_64": 3.14, "string_1": "one"},
			want:  `( "objects"."float_64" = $1 AND "objects"."int_32" = $2 AND "objects"."string_1" = $3 )`,
			want1: []any{3.14, 0, "one"},
		},
		{
			name:  "",
			f:     filter.Eq{"int_32": 0, "float_64": 3.14, "string_1": "one", "bool_4": nil},
			want:  `( "objects"."bool_4" IS NULL AND "objects"."float_64" = $1 AND "objects"."int_32" = $2 AND "objects"."string_1" = $3 )`,
			want1: []any{3.14, 0, "one"},
		},
		{
			name:  "",
			f:     filter.Eq{"int_32": 0, "float_64": 3.14, "string_1": "one", "bool_4": nil, "bool_1": true},
			want:  `( "objects"."bool_1" IS TRUE AND "objects"."bool_4" IS NULL AND "objects"."float_64" = $1 AND "objects"."int_32" = $2 AND "objects"."string_1" = $3 )`,
			want1: []any{3.14, 0, "one"},
		},
		{
			name:  "",
			f:     filter.And{},
			want:  ``,
			want1: []any{},
		},
		{
			name:  "",
			f:     filter.And{filter.Eq{"int_32": 0}},
			want:  `"objects"."int_32" = $1`,
			want1: []any{0},
		},
		{
			name:  "",
			f:     filter.And{filter.Eq{"int_32": 1}, filter.Eq{"float_64": 3.14}},
			want:  `( "objects"."int_32" = $1 AND "objects"."float_64" = $2 )`,
			want1: []any{1, 3.14},
		},
		{
			name:  "",
			f:     filter.Or{},
			want:  ``,
			want1: []any{},
		},
		{
			name:  "",
			f:     filter.Or{filter.Eq{"int_32": 0}, filter.Eq{}, filter.And{filter.Eq{}, filter.Eq{}}},
			want:  `( "objects"."int_32" = $1 )`,
			want1: []any{0},
		},
		{
			name:  "",
			f:     filter.Or{filter.Eq{"int_32": 0}, filter.Eq{"int_32": 1}},
			want:  `( "objects"."int_32" = $1 OR "objects"."int_32" = $2 )`,
			want1: []any{0, 1},
		},
		{
			name:  "",
			f:     filter.Or{filter.Eq{"int_32": 0, "bool_1": false}, filter.Eq{"int_32": 1}},
			want:  `( ( "objects"."bool_1" IS FALSE AND "objects"."int_32" = $1 ) OR "objects"."int_32" = $2 )`,
			want1: []any{0, 1},
		},
		{
			name:  "",
			f:     filter.Ne{"int_32": 0, "float_64": 3.14, "string_1": "one", "bool_4": nil, "bool_1": true},
			want:  `( "objects"."bool_1" IS NOT TRUE AND "objects"."bool_4" IS NOT NULL AND "objects"."float_64" <> $1 AND "objects"."int_32" <> $2 AND "objects"."string_1" <> $3 )`,
			want1: []any{3.14, 0, "one"},
		},
		{
			name:  "",
			f:     filter.Ge{"time_1": filter.Now()},
			want:  `"objects"."time_1" >= NOW()`,
			want1: []any{},
		},
		{
			name:  "",
			f:     filter.True{},
			want:  `TRUE`,
			want1: []any{},
		},
		{
			name:  "",
			f:     filter.False{},
			want:  `FALSE`,
			want1: []any{},
		},
		{
			name:  "",
			f:     filter.Or{filter.And{}},
			want:  ``,
			want1: []any{},
		},
		{
			name:  "",
			f:     filter.Or{filter.And{}, filter.And{}},
			want:  `( TRUE )`,
			want1: []any{},
		},
		{
			name:  "",
			f:     filter.Or{filter.And{}, filter.And{filter.Eq{}, filter.Ne{}}},
			want:  `( TRUE )`,
			want1: []any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := builder.Builder{}
			b.Alloc(32)
			j := model.Object{}
			if err := tt.f.To(&b, &j); (err != nil) != tt.wantErr {
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

func TestHolder_AppendTo(t *testing.T) {
	tests := []struct {
		name    string
		h       builder.Holder
		wantW   string
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "",
			h:       0,
			wantW:   "$0",
			want:    2,
			wantErr: false,
		},
		{
			name:    "",
			h:       10,
			wantW:   "$10",
			want:    3,
			wantErr: false,
		},
		{
			name:    "",
			h:       100,
			wantW:   "$100",
			want:    4,
			wantErr: false,
		},
		{
			name:    "",
			h:       1000,
			wantW:   "$1000",
			want:    5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := tt.h.AppendTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppendTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("AppendTo() gotW = %v, want %v", gotW, tt.wantW)
			}
			if got != tt.want {
				t.Errorf("AppendTo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
