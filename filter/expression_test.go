package filter_test

import (
	"encoding/json"
	"fmt"
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
		want1   filter.Expression
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name:    "",
			args:    args{f: filter.Eq{"f": "abc"}},
			want:    []byte(`[["f","EQ","abc"]]`),
			want1:   filter.Expression{filter.Operation{"f", "EQ", "abc"}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": 3.14}},
			want:    []byte(`[["f","EQ",3.14]]`),
			want1:   filter.Expression{filter.Operation{"f", "EQ", 3.14}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": 100}},
			want:    []byte(`[["f","EQ",100]]`),
			want1:   filter.Expression{filter.Operation{"f", "EQ", 1e2}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": true}},
			want:    []byte(`[["f","EQ",true]]`),
			want1:   filter.Expression{filter.Operation{"f", "EQ", true}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": nil}},
			want:    []byte(`[["f","EQ",null]]`),
			want1:   filter.Expression{filter.Operation{"f", "EQ", nil}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": help.PtrTime(time.Unix(0, 0))}},
			want:    []byte(`[["f","EQ","1970-01-01T03:00:00+03:00"]]`),
			want1:   filter.Expression{filter.Operation{"f", "EQ", "1970-01-01T03:00:00+03:00"}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": time.Unix(0, 0).UTC()}},
			want:    []byte(`[["f","EQ","1970-01-01T00:00:00Z"]]`),
			want1:   filter.Expression{filter.Operation{"f", "EQ", "1970-01-01T00:00:00Z"}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Eq{"f": uuid.UUID{}}},
			want:    []byte(`[["f","EQ","00000000-0000-0000-0000-000000000000"]]`),
			want1:   filter.Expression{filter.Operation{"f", "EQ", "00000000-0000-0000-0000-000000000000"}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.And{filter.Ge{"f": 0}, filter.Eq{"b": false}}},
			want:    []byte(`[[[["f","GE",0]],[["b","EQ",false]]]]`),
			want1:   filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", .0}}, filter.Expression{filter.Operation{"b", "EQ", false}}}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Or{filter.Ge{"f": 0}, filter.Eq{"b": false}}},
			want:    []byte(`[[[["f","GE",0]]],[[["b","EQ",false]]]]`),
			want1:   filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", .0}}}, filter.Expression{filter.Expression{filter.Operation{"b", "EQ", false}}}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Or{filter.And{filter.Ge{"f": 0}, filter.Eq{"b": false}}, filter.And{filter.Le{"f": 0}, filter.Ne{"b": false}}}},
			want:    []byte(`[[[[[["f","GE",0]],[["b","EQ",false]]]]],[[[[["f","LE",0]],[["b","NE",false]]]]]]`),
			want1:   filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", .0}}, filter.Expression{filter.Operation{"b", "EQ", false}}}}}, filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "LE", .0}}, filter.Expression{filter.Operation{"b", "NE", false}}}}}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.And{filter.Or{filter.Ge{"f": 0}, filter.Eq{"b": false}}, filter.Or{filter.Le{"f": 0}, filter.Ne{"b": false}}}},
			want:    []byte(`[[[[[["f","GE",0]]],[[["b","EQ",false]]]],[[[["f","LE",0]]],[[["b","NE",false]]]]]]`),
			want1:   filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", .0}}}, filter.Expression{filter.Expression{filter.Operation{"b", "EQ", false}}}}, filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "LE", .0}}}, filter.Expression{filter.Expression{filter.Operation{"b", "NE", false}}}}}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.Or{filter.And{filter.Ge{"f": 0}, filter.Eq{"b": false}}, filter.Le{"f": 0}}},
			want:    []byte(`[[[[[["f","GE",0]],[["b","EQ",false]]]]],[[["f","LE",0]]]]`),
			want1:   filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", .0}}, filter.Expression{filter.Operation{"b", "EQ", false}}}}}, filter.Expression{filter.Expression{filter.Operation{"f", "LE", .0}}}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.And{filter.Or{filter.Ge{"f": 0}, filter.Eq{"b": false}}, filter.Le{"f": 0}}},
			want:    []byte(`[[[[[["f","GE",0]]],[[["b","EQ",false]]]],[["f","LE",0]]]]`),
			want1:   filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", .0}}}, filter.Expression{filter.Expression{filter.Operation{"b", "EQ", false}}}}, filter.Expression{filter.Operation{"f", "LE", .0}}}},
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: filter.And{filter.In{"f": []any{1, 2}}, filter.Ni{"f": []any{"a", "b"}}}},
			want:    []byte(`[[[["f","IN",[1,2]]],[["f","NI",["a","b"]]]]]`),
			want1:   filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "IN", []any{1., 2.}}}, filter.Expression{filter.Operation{"f", "NI", []any{"a", "b"}}}}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filter.MarshalJSON(tt.args.f)
			require.ErrorIs(t, tt.wantErr, err)
			require.Equal(t, tt.want, got)
			var e filter.Expression
			err = json.Unmarshal(got, &e)
			require.NoError(t, err)
			require.Equal(t, tt.want1, e)
			_, err = e.Filter()
			require.NoError(t, err)
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
		{
			name:    "",
			args:    args{b: []byte(`[[[[[["f","GE",0]]],[[["b","EQ",false]]]],[["f","LE",0]]]]`)},
			want:    filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", .0}}}, filter.Expression{filter.Expression{filter.Operation{"b", "EQ", false}}}}, filter.Expression{filter.Operation{"f", "LE", .0}}}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f filter.Expression
			err := json.Unmarshal(tt.args.b, &f)
			require.ErrorIs(t, tt.wantErr, err)
			require.EqualValues(t, tt.want, f)
		})
	}
}

func TestOperation_Filter(t *testing.T) {
	tests := []struct {
		name    string
		op      filter.Operation
		want    filter.Filter
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name:    "",
			op:      filter.Operation{"f", "EQ", 1},
			want:    filter.Eq{"f": 1},
			wantErr: nil,
		},
		{
			name:    "",
			op:      filter.Operation{"f", "NE", 1e1},
			want:    filter.Ne{"f": 1e1},
			wantErr: nil,
		},
		{
			name:    "",
			op:      filter.Operation{"f", "GT", 0},
			want:    filter.Gt{"f": 0},
			wantErr: nil,
		},
		{
			name:    "",
			op:      filter.Operation{"f", "GE", 1},
			want:    filter.Ge{"f": 1},
			wantErr: nil,
		},
		{
			name:    "",
			op:      filter.Operation{"f", "LE", .0},
			want:    filter.Le{"f": .0},
			wantErr: nil,
		},
		{
			name:    "",
			op:      filter.Operation{"f", "LT", 1.},
			want:    filter.Lt{"f": 1.},
			wantErr: nil,
		},
		{
			name:    "",
			op:      filter.Operation{"f", "IN", []any{"green", "yellow"}},
			want:    filter.In{"f": {"green", "yellow"}},
			wantErr: nil,
		},
		{
			name:    "",
			op:      filter.Operation{"f", "NI", []any{"green", "yellow"}},
			want:    filter.Ni{"f": {"green", "yellow"}},
			wantErr: nil,
		},
		{
			name:    "",
			op:      filter.Operation{"f", "AS", "%123%"},
			want:    filter.As{"f": "%123%"},
			wantErr: nil,
		},
		{
			name:    "",
			op:      filter.Operation{"f", "NA", "%123%"},
			want:    filter.Na{"f": "%123%"},
			wantErr: nil,
		},
		{
			name:    "",
			op:      filter.Operation{},
			want:    nil,
			wantErr: filter.ErrMalformedOperation,
		},
		{
			name:    "",
			op:      filter.Operation{"f", 11},
			want:    nil,
			wantErr: filter.ErrIllegalOperation,
		},
		{
			name:    "",
			op:      filter.Operation{"f", "XX"},
			want:    nil,
			wantErr: filter.ErrUnknownOperation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.op.Filter()
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestExpression_Filter(t *testing.T) {
	tests := []struct {
		name    string
		ex      filter.Expression
		want    filter.Filter
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name:    "",
			ex:      filter.Expression{filter.Operation{"f", "EQ", 3.14}},
			want:    filter.Eq{"f": 3.14},
			wantErr: nil,
		},
		{
			name:    "",
			ex:      filter.Expression{filter.Operation{"f", "EQ", 3.14}, filter.Operation{"b", "EQ", true}, filter.Operation{"n", "EQ", nil}},
			want:    filter.Eq{"b": true, "f": 3.14, "n": nil},
			wantErr: nil,
		},
		{
			name:    "",
			ex:      filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", 3.14}}}, filter.Expression{filter.Expression{filter.Operation{"b", "EQ", true}}}},
			want:    filter.Or{filter.Ge{"f": 3.14}, filter.Eq{"b": true}},
			wantErr: nil,
		},
		{
			name:    "",
			ex:      filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", 3.14}}, filter.Expression{filter.Operation{"b", "EQ", false}}}},
			want:    filter.And{filter.Ge{"f": 3.14}, filter.Eq{"b": false}},
			wantErr: nil,
		},
		{
			name:    "",
			ex:      filter.Expression{filter.Operation{"f", "GE", 3.14}, filter.Operation{"f", "GE", 3.14}},
			want:    filter.Ge{"f": 3.14},
			wantErr: nil,
		},
		{
			name:    "",
			ex:      filter.Expression{},
			want:    nil,
			wantErr: nil,
		},
		{
			name:    "",
			ex:      nil,
			want:    nil,
			wantErr: nil,
		},
		{
			name:    "",
			ex:      filter.Expression{filter.Expression{filter.Expression{}}},
			want:    nil,
			wantErr: filter.ErrEmptyExpression,
		},
		{
			name:    "",
			ex:      filter.Expression{nil},
			want:    nil,
			wantErr: filter.ErrUnknownExpression,
		},
		{
			name:    "",
			ex:      filter.Expression{filter.Operation{"f", "GE", 3.14}, filter.Operation{"f", "LE", 3.14}},
			want:    nil,
			wantErr: filter.ErrUnsuitableOperation,
		},
		{
			name:    "",
			ex:      filter.Expression{filter.Operation{"f", "GE", 3.14}, nil},
			want:    nil,
			wantErr: filter.ErrIllegalExpression,
		},
		{
			name:    "",
			ex:      filter.Expression{filter.Expression{nil}},
			want:    nil,
			wantErr: filter.ErrIllegalExpression,
		},
		{
			name:    "",
			ex:      filter.Expression{filter.Expression{filter.Expression{filter.Expression{}}}, nil},
			want:    nil,
			wantErr: filter.ErrIllegalExpression,
		},
		{
			name:    "",
			ex:      filter.Expression{filter.Expression{filter.Expression{filter.Expression{}, filter.Expression{}}}, nil},
			want:    nil,
			wantErr: filter.ErrMalformedExpression,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ex.Filter()
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func ExampleMarshalJSON() {
	a := filter.And{filter.Or{filter.Ge{"f": 0}, filter.Eq{"b": false}}, filter.In{"t": {"1970-01-01T00:00:00Z"}}}

	fmt.Printf("%#v\n", a)

	b, err := filter.MarshalJSON(a)
	if err != nil {
		return
	}

	fmt.Printf("%s\n", b)

	var e filter.Expression
	err = json.Unmarshal(b, &e)
	if err != nil {
		return
	}

	fmt.Printf("%#v\n", e)

	f, err := e.Filter()
	if err != nil {
		return
	}

	fmt.Printf("%#v\n", f)

	var o struct {
		E []filter.Expression
	}
	err = json.Unmarshal([]byte(`{"e":[null,[]]}`), &o)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%#v\n", o)

	// Output:
	//
	// filter.And{filter.Or{filter.Ge{"f":0}, filter.Eq{"b":false}}, filter.In{"t":filter.Array{"1970-01-01T00:00:00Z"}}}
	// [[[[[["f","GE",0]]],[[["b","EQ",false]]]],[["t","IN",["1970-01-01T00:00:00Z"]]]]]
	// filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Expression{filter.Operation{"f", "GE", 0}}}, filter.Expression{filter.Expression{filter.Operation{"b", "EQ", false}}}}, filter.Expression{filter.Operation{"t", "IN", []interface {}{"1970-01-01T00:00:00Z"}}}}}
	// filter.And{filter.Or{filter.Ge{"f":0}, filter.Eq{"b":false}}, filter.In{"t":filter.Array{"1970-01-01T00:00:00Z"}}}
	// struct { E []filter.Expression }{E:[]filter.Expression{filter.Expression(nil), filter.Expression(nil)}}
}
