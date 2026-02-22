package builder_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test/model"
	"github.com/pshvedko/dbx/request"
)

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
	require.Equal(t, "0-7,10-13,18-21[[5EQT&10EQ]&[8IN&11IN&14IN]&[T|13GT|15NA]]+21", q)
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
