package filter_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/t"
)

func TestMarshalJSON(t *testing.T) {
	type args struct {
		f filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name:    "",
			args:    args{f: filter.Eq{"f": "abc"}},
			want:    []byte(`[["f","EQ","abc"]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": 3.14}},
			want:    []byte(`[["f","EQ",3.14]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": 100}},
			want:    []byte(`[["f","EQ",100]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": true}},
			want:    []byte(`[["f","EQ",true]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": nil}},
			want:    []byte(`[["f","EQ",null]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": help.PtrTime(time.Unix(0, 0))}},
			want:    []byte(`[["f","EQ","1970-01-01T03:00:00+03:00"]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": time.Unix(0, 0).UTC()}},
			want:    []byte(`[["f","EQ","1970-01-01T00:00:00Z"]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": uuid.UUID{}}},
			want:    []byte(`[["f","EQ","00000000-0000-0000-0000-000000000000"]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.And{filter.Ge{"f": 0}, filter.Eq{"b": false}}},
			want:    []byte(`[[[["f","GE",0]],[["b","EQ",false]]]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Or{filter.Ge{"f": 0}, filter.Eq{"b": false}}},
			want:    []byte(`[[[["f","GE",0]]],[[["b","EQ",false]]]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Or{filter.And{filter.Ge{"f": 0}, filter.Eq{"b": false}}, filter.And{filter.Le{"f": 0}, filter.Ne{"b": false}}}},
			want:    []byte(`[[[[[["f","GE",0]],[["b","EQ",false]]]]],[[[[["f","LE",0]],[["b","NE",false]]]]]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.And{filter.Or{filter.Ge{"f": 0}, filter.Eq{"b": false}}, filter.Or{filter.Le{"f": 0}, filter.Ne{"b": false}}}},
			want:    []byte(`[[[[[["f","GE",0]]],[[["b","EQ",false]]]],[[[["f","LE",0]]],[[["b","NE",false]]]]]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Or{filter.And{filter.Ge{"f": 0}, filter.Eq{"b": false}}, filter.Le{"f": 0}}},
			want:    []byte(`[[[[[["f","GE",0]],[["b","EQ",false]]]]],[[["f","LE",0]]]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.And{filter.Or{filter.Ge{"f": 0}, filter.Eq{"b": false}}, filter.Le{"f": 0}}},
			want:    []byte(`[[[[[["f","GE",0]]],[[["b","EQ",false]]]],[["f","LE",0]]]]`),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filter.MarshalJSON(tt.args.f)
			require.ErrorIs(t, tt.wantErr, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    filter.Expression
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name:    "",
			args:    args{b: []byte(`[["f","EQ",3.14]]`)},
			want:    filter.Expression{filter.Operation{"f", "EQ", 3.14}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{b: []byte(`[[[["f","GE",0]],[["b","EQ",false]]]]`)},
			want:    filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", .0}}, filter.Expression{filter.Operation{"b", "EQ", false}}}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{b: []byte(`[[[["f","GE",0]]],[[["b","EQ",false]]]]`)},
			want:    filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", .0}}}, filter.Expression{filter.Expression{filter.Operation{"b", "EQ", false}}}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{b: []byte(`[[[[[["f","GE",0]],[["b","EQ",false]]]]],[[[[["f","LE",0]],[["b","NE",false]]]]]]`)},
			want:    filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", .0}}, filter.Expression{filter.Operation{"b", "EQ", false}}}}}, filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "LE", .0}}, filter.Expression{filter.Operation{"b", "NE", false}}}}}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{b: []byte(`[[[[[["f","GE",0]]],[[["b","EQ",false]]]],[[[["f","LE",0]]],[[["b","NE",false]]]]]]`)},
			want:    filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", .0}}}, filter.Expression{filter.Expression{filter.Operation{"b", "EQ", false}}}}, filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "LE", .0}}}, filter.Expression{filter.Expression{filter.Operation{"b", "NE", false}}}}}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f filter.Expression
			err := json.Unmarshal(tt.args.b, &f)
			require.ErrorIs(t, tt.wantErr, err)
			t.Logf("%v", f)
			require.EqualValues(t, tt.want, f)
		})
	}
}
