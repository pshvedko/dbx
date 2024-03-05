package request

import (
	"reflect"
	"testing"

	"github.com/pshvedko/db/filter"
)

type Object struct {
	Bool    bool    `json:"bool,omitempty"`
	Int     int     `json:"int,omitempty"`
	Int16   int16   `json:"int_16,omitempty"`
	Float32 float32 `json:"float_32,omitempty"`
	Float64 float64 `json:"float_64,omitempty"`
	String  string  `json:"string,omitempty"`
}

func (o Object) Table() string {
	return "objects"
}

func (o Object) Names() []string {
	return []string{"bool", "int", "int_16", "float_32", "float_64", "string"}
}

func (o *Object) Values() []any {
	return []any{&o.Bool, &o.Int, &o.Int16, &o.Float32, &o.Float64, &o.String}
}

func TestRequest_makeSelect(t *testing.T) {
	type args struct {
		j filter.Projector
		f filter.Filter
		o []Option
	}
	object := Object{}
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
				j: &object,
				f: filter.Eq{"int": 1, "bool": true},
				o: nil,
			},
			want:    `SELECT "bool", "int", "int_16", "float_32", "float_64", "string" FROM "objects" WHERE ( bool IS TRUE AND int = $1 )`,
			want1:   []any{1},
			want2:   object.Values(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := New(nil, nil, tt.args.o...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, got1, got2, err := r.makeSelect(tt.args.j, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeSelect() error = %v, wantErr %v", err, tt.wantErr)
				return
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
