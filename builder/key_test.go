package builder_test

import (
	"github.com/pshvedko/dbx/request"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test/model"
)

func TestHash_Places(t *testing.T) {
	var o model.Object
	type args struct {
		j filter.Projector
		L string
	}
	tests := []struct {
		name    string
		args    args
		want    []any
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				j: &o,
				L: "0-21",
			},
			want: o.Places(),
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "",
			},
			want: []any{},
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "0",
			},
			want: []any{&o.ID},
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "1",
			},
			want: []any{&o.UUID2},
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "0,1",
			},
			want: []any{&o.ID, &o.UUID2},
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "0-1",
			},
			want: []any{&o.ID, &o.UUID2},
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "0,1,3",
			},
			want: []any{&o.ID, &o.UUID2, &o.UUID4},
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "0,1,3-5,7",
			},
			want: []any{&o.ID, &o.UUID2, &o.UUID4, &o.Bool1, &o.Bool2, &o.Bool4},
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "0,1,3-5-7",
			},
			wantErr: builder.ErrInvalidRange,
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "0-",
			},
			wantErr: builder.ErrInvalidRange,
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "1-1",
			},
			wantErr: builder.ErrInvalidRange,
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "1-",
			},
			wantErr: builder.ErrInvalidRange,
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "1-2,",
			},
			wantErr: builder.ErrInvalidRange,
		},
		{
			name: "",
			args: args{
				j: &o,
				L: ",",
			},
			wantErr: builder.ErrInvalidRange,
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "-",
			},
			wantErr: builder.ErrInvalidRange,
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "2-1",
			},
			wantErr: builder.ErrInvalidRange,
		},
		{
			name: "",
			args: args{
				j: &o,
				L: "1,3,2",
			},
			wantErr: builder.ErrInvalidRange,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := builder.Hash{Static: builder.Static{List: tt.args.L}}
			got, err := h.Places(tt.args.j)
			require.ErrorIs(t, tt.wantErr, err)
			require.Equal(t, tt.want, got)
			for i, v := range got {
				require.Same(t, tt.want[i], v)
			}
		})
	}
}

func TestKey_WriteIndices(t *testing.T) {
	type args struct {
		j filter.Projector
		f builder.Fielder
	}
	o := model.Object{}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				j: &o,
				f: &builder.ExcludedColumn{},
			},
			want: "0-21",
		},
		{
			name: "",
			args: args{
				j: &o,
				f: &builder.IncludedColumn{},
			},
			want: "",
		},
		{
			name: "",
			args: args{
				j: &o,
				f: &builder.IncludedColumn{{"id": 0}},
			},
			want: "0",
		},
		{
			name: "",
			args: args{
				j: &o,
				f: &builder.IncludedColumn{{"id": 0, "uuid_2": 1}},
			},
			want: "0,1",
		},
		{
			name: "",
			args: args{
				j: &o,
				f: &builder.IncludedColumn{{"id": 0, "uuid_2": 1}, {}},
			},
			want: "",
		},
		{
			name: "",
			args: args{
				j: &o,
				f: &builder.IncludedColumn{{"id": 0}, {"uuid_2": 1}},
			},
			want: "1",
		},
		{
			name: "",
			args: args{
				j: &o,
				f: &builder.IncludedColumn{{"id": 0, "uuid_2": 1}, o.Columns()},
			},
			want: "0-21",
		},
		{
			name: "",
			args: args{
				j: &o,
				f: IndexedColumn(t, o, 1, 1, 1),
			},
			want: "0-2",
		},
		{
			name: "",
			args: args{
				j: &o,
				f: IndexedColumn(t, o, 1, 0, 1),
			},
			want: "0,2",
		},
		{
			name: "",
			args: args{
				j: &o,
				f: IndexedColumn(t, o, 1, 1, 0, 1, 1, 1, 1, 0, 1),
			},
			want: "0,1,3-6,8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := builder.Key{Constructor: &builder.Constructor{Fielder: tt.args.f}}
			err := k.WriteIndices(tt.args.j)
			require.ErrorIs(t, err, tt.wantErr)
			key := k.String()
			require.Equal(t, tt.want, key)
		})
	}
}

type TestedColumn []string

func (t TestedColumn) Used(n string) bool { return slices.Contains(t, n) }

func (t TestedColumn) Names() [2]map[string]int { return [2]map[string]int{} }

func (t TestedColumn) Returned(n string) bool { return t.Used(n) }

func IndexedColumn(t *testing.T, j filter.Fielder, idx ...byte) (c TestedColumn) {
	t.Helper()
	n := j.Names()
	for i, b := range idx {
		if b != 0 {
			c = append(c, n[i])
		}
	}
	return
}

func TestConstructor_Select_WithCache(t *testing.T) {
	j := model.Object{}

	r, err := request.NewWithOption([]request.Option{
		request.WithCreated("time_1"),
		request.WithUpdated("time_2"),
		request.WithDeleted("time_4"),
		request.WithoutField{"uuid_3", "bool_1", "float_32", "float_64", "string_1", "string_2", "string_3", "string_4"},
		request.WithCache(&__c),
	})
	require.NoError(t, err)

	f := __f

	z0, q0, a0, p0, err := r.Constructor().Range(&__o, &__l).Sort(__y).Select(&j, f)
	require.NoError(t, err)
	t.Log(q0)
	t.Log(a0)
	t.Log(p0)
	t.Log(z0)

	z1, q1, a1, p1, err := r.Constructor().Range(&__o, &__l).Sort(__y).Select(&j, f)
	require.NoError(t, err)
	t.Log(q1)
	t.Log(a1)
	t.Log(p1)
	t.Log(z1)
	require.Equal(t, q0, q1)
	require.Equal(t, a0, a1)
	require.Equal(t, p0, p1)
	require.Equal(t, z0, z1)
}

func TestKey_WriteHash(t *testing.T) {
	j := model.Object{}

	r, err := request.NewWithOption([]request.Option{
		request.WithCreated("time_1"),
		request.WithUpdated("time_2"),
		request.WithDeleted("time_4"),
		request.WithoutCount(),
		request.WithoutField{"float_32", "float_64", "string_1", "string_2", "string_3", "string_4"},
		request.WithCache(&__c),
	})
	require.NoError(t, err)

	k := builder.Key{Constructor: r.Constructor().Range(&__o, &__l).Sort(__y)}
	n, err := k.WriteHash(&j, __f)
	require.NoError(t, err)
	q := k.String()
	t.Log(q)
	require.Equal(t, 15, n)
	require.Equal(t, "0-7,10-13,18-21[[5EQT&10EQ]&[8IN&11IN&14IN]&[T|13GT|15NA]]+21-1+18OL", q)
}

func TestKey_To(t *testing.T) {
	j := model.Object{}
	k := builder.Key{Constructor: &builder.Constructor{}}
	f := __f

	err := f.To(&k, &j)
	require.NoError(t, err)
	t.Logf("%s", &k)
	require.Equal(t, `[[5EQT&10EQ]&[8IN&11IN&14IN]&[T|13GT|15NA]]`, k.String())
}

const (
	KEY  = `[[5EQT&10EQ]&[8IN&11IN&14IN]&[T|13GT|15NA]]`
	HASH = `0-21` + KEY + `+21`
)

func BenchmarkKey(b *testing.B) {
	f := __f
	j := model.Object{}
	c := builder.Constructor{T: true}
	var q string
	for i := 0; i < b.N; i++ {
		k := builder.Key{Constructor: &c}
		k.Alloc(16)
		k.Grow(128)
		err := f.To(&k, &j)
		if err != nil {
			b.Fatal(err)
		}
		q = k.String()
		switch q {
		case KEY:
		default:
			b.Fatal(q)
		}
		k.Reset()
	}
}

func BenchmarkKey_WriteHash(b *testing.B) {
	f := __f
	j := model.Object{}
	c := builder.Constructor{
		Fielder: builder.ExcludedColumn{},
		Modify: builder.Modify{
			Created: "time_1",
			Updated: "time_2",
			Deleted: builder.DeletedNone("time_4"),
		},
		T: true,
	}
	var q string
	for i := 0; i < b.N; i++ {
		k := builder.Key{Constructor: &c}
		k.Alloc(16)
		k.Grow(128)
		_, err := k.WriteHash(&j, f)
		if err != nil {
			b.Fatal(err)
		}
		q = k.String()
		switch q {
		case HASH:
		default:
			b.Fatal(q)
		}
		k.Reset()
	}
}
