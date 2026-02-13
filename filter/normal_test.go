package filter

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestExpand(t *testing.T) {
	type args struct {
		f Filter
	}
	tests := []struct {
		name string
		args args
		want Filter
	}{
		// TODO: Add test cases.
		{
			args: args{
				f: Eq{"a": 1},
			},
			want: And{Eq{"a": 1}},
		},
		{
			args: args{
				f: Eq{"a": 1, "b": 2},
			},
			want: And{Eq{"a": 1}, Eq{"b": 2}},
		},
		{
			args: args{
				f: Ne{"a": 1, "b": 2},
			},
			want: And{Ne{"a": 1}, Ne{"b": 2}},
		},
		{
			args: args{
				f: In{"a": {1, 2, 3}},
			},
			want: Or{Eq{"a": 1}, Eq{"a": 2}, Eq{"a": 3}},
		},
		{
			args: args{
				f: Ni{"a": {time.Duration(1), time.Duration(2), time.Duration(3)}},
			},
			want: And{Ne{"a": time.Duration(3)}, Ne{"a": time.Duration(2)}, Ne{"a": time.Duration(1)}},
		},
		{
			args: args{
				f: And{Eq{"a": 1}, Eq{"b": 2}},
			},
			want: And{Eq{"a": 1}, Eq{"b": 2}},
		},
		{
			args: args{
				f: And{Eq{"a": 1}, And{Eq{"b": 2}}},
			},
			want: And{Eq{"a": 1}, Eq{"b": 2}},
		},
		{
			args: args{
				f: And{Eq{"a": 1}, Eq{"b": 2, "c": 3}},
			},
			want: And{Eq{"a": 1}, Eq{"b": 2}, Eq{"c": 3}},
		},
		{
			args: args{
				f: And{Eq{"a": 1}, Or{Eq{"b": 2}}},
			},
			want: Or{And{Eq{"a": 1}, Eq{"b": 2}}},
		},
		{
			args: args{
				f: And{Eq{"a": 1}, Or{Eq{"b": 2}, Eq{"c": 3}}},
			},
			want: Or{And{Eq{"a": 1}, Eq{"b": 2}}, And{Eq{"a": 1}, Eq{"c": 3}}},
		},
		{
			args: args{
				f: And{Or{Eq{"a": 1}, Eq{"b": 2}}, Or{Eq{"c": 3}, Eq{"d": 4}}},
			},
			want: Or{
				And{Eq{"a": 1}, Eq{"c": 3}},
				And{Eq{"a": 1}, Eq{"d": 4}},
				And{Eq{"b": 2}, Eq{"c": 3}},
				And{Eq{"b": 2}, Eq{"d": 4}}},
		},
		{
			args: args{
				f: And{Eq{"e": 0}, Or{Eq{"a": 1}, Eq{"b": 2}}, Or{Eq{"c": 3}, Eq{"d": 4}}},
			},
			want: Or{
				And{Eq{"e": 0}, Eq{"a": 1}, Eq{"c": 3}},
				And{Eq{"e": 0}, Eq{"a": 1}, Eq{"d": 4}},
				And{Eq{"e": 0}, Eq{"b": 2}, Eq{"c": 3}},
				And{Eq{"e": 0}, Eq{"b": 2}, Eq{"d": 4}}},
		},
		{
			args: args{
				f: And{Or{Eq{"e": 0}, Eq{"f": 0}}, Or{Eq{"a": 1}, Eq{"b": 2}}, Or{Eq{"c": 3}, Eq{"d": 4}}},
			},
			want: Or{
				And{Eq{"e": 0}, Eq{"a": 1}, Eq{"c": 3}},
				And{Eq{"e": 0}, Eq{"a": 1}, Eq{"d": 4}},
				And{Eq{"e": 0}, Eq{"b": 2}, Eq{"c": 3}},
				And{Eq{"e": 0}, Eq{"b": 2}, Eq{"d": 4}},
				And{Eq{"f": 0}, Eq{"a": 1}, Eq{"c": 3}},
				And{Eq{"f": 0}, Eq{"a": 1}, Eq{"d": 4}},
				And{Eq{"f": 0}, Eq{"b": 2}, Eq{"c": 3}},
				And{Eq{"f": 0}, Eq{"b": 2}, Eq{"d": 4}}},
		},
		{
			args: args{
				f: Or{Eq{"a": 1}, And{Eq{"b": 2}, Eq{"c": 3}}},
			},
			want: Or{And{Eq{"a": 1}}, And{Eq{"b": 2}, Eq{"c": 3}}},
		},
		{
			args: args{
				f: Or{Eq{"a": 1}, Or{Eq{"b": 2}, Eq{"c": 3}}},
			},
			want: Or{And{Eq{"a": 1}}, And{Eq{"b": 2}}, And{Eq{"c": 3}}},
		},
		{
			args: args{
				f: Or{In{"a": {1, 2, 3}}, In{"b": {1, 2, 3}}},
			},
			want: Or{Eq{"a": 1}, Eq{"a": 2}, Eq{"a": 3}, Eq{"b": 1}, Eq{"b": 2}, Eq{"b": 3}},
		},
		{
			args: args{
				f: In{"a": {}},
			},
			want: False{},
		},
		{
			args: args{
				f: And{In{"a": {}}},
			},
			want: And{False{}},
		},
		{
			args: args{
				f: And{In{"a": {}}, In{"b": {}}},
			},
			want: And{False{}, False{}},
		},
		{
			args: args{
				f: Or{In{"a": {}}, In{"b": {}}},
			},
			want: Or{False{}, False{}},
		},
		{
			args: args{
				f: Or{},
			},
			want: Or{},
		},
		{
			args: args{
				f: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Expand(tt.args.f); !assert.IsType(t, tt.want, got) || !assert.ElementsMatch(t, got, tt.want) {
				t.Errorf("Expand() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func BenchmarkExpand(b *testing.B) {
	a := And{Or{Eq{"e": 0}, Eq{"f": 0}}, Or{Eq{"a": 1}, Eq{"b": 2}}, Or{Eq{"c": 3}, Eq{"d": 4}}}
	o := Or{And{Eq{"e": 0}, Eq{"a": 1}, Eq{"c": 3}}, And{Eq{"e": 0}, Eq{"a": 1}, Eq{"d": 4}}, And{Eq{"e": 0}, Eq{"b": 2}, Eq{"c": 3}}, And{Eq{"e": 0}, Eq{"b": 2}, Eq{"d": 4}}, And{Eq{"f": 0}, Eq{"a": 1}, Eq{"c": 3}}, And{Eq{"f": 0}, Eq{"a": 1}, Eq{"d": 4}}, And{Eq{"f": 0}, Eq{"b": 2}, Eq{"c": 3}}, And{Eq{"f": 0}, Eq{"b": 2}, Eq{"d": 4}}}
	for i := 0; i < b.N; i++ {
		require.ElementsMatch(b, o, Expand(a))
	}
}

func TestCollapse(t *testing.T) {
	tests := []struct {
		name string
		args And
		want Filter
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: And{Eq{"f": 3.14}, Eq{"b": true}, Eq{"n": nil}},
			want: Eq{"f": 3.14, "b": true, "n": nil},
		},
		{
			name: "",
			args: And{Eq{"f": 3.14, "x": 1}, Eq{"x": 2, "b": true}, Eq{"n": nil}},
			want: And{Eq{"x": 2}, Eq{"b": true, "f": 3.14, "n": nil, "x": 1}},
		},
		{
			name: "",
			args: And{Eq{"f": 3.14, "x": 1}, Eq{"x": 2, "b": true}, Eq{"x": 3, "n": nil}},
			want: And{Eq{"x": 1}, Eq{"x": 2}, Eq{"b": true, "f": 3.14, "n": nil, "x": 3}},
		},
		{
			name: "",
			args: And{True{}, True{}, False{}},
			want: And{True{}, False{}},
		},
		{
			name: "",
			args: And{True{}, True{}},
			want: True{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Collapse(tt.args), "Collapse(%v)", tt.args)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	type args struct {
		f Filter
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{},
			want: true,
		},
		{
			name: "",
			args: args{f: And{Ni{}}},
			want: true,
		},
		{
			name: "",
			args: args{f: And{Or{Eq{}, Gt{}}, In{"a": {}}}},
			want: true,
		},
		{
			name: "",
			args: args{f: True{}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsEmpty(tt.args.f), "IsEmpty(%v)", tt.args.f)
		})
	}
}
