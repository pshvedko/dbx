package builder_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test/model"
	"github.com/pshvedko/dbx/request"
)

func TestConstructor_TryCache(t *testing.T) {
	var j model.Object
	r, err := request.NewWithOption([]request.Option{
		request.WithCreated("time_1"),
		request.WithUpdated("time_2"),
		request.WithDeleted("time_4"),
		request.WithoutCount(),
		request.WithCache(&__c),
	})
	require.NoError(t, err)

	k, err := r.Constructor().Range(&__o, &__l).Sort(__y).TryCache(&j, __f, 16, '$')
	require.NoError(t, err)
	require.NotZero(t, k)

	t.Logf(k.String())
}

func TestKey_To(t *testing.T) {
	j := model.Object{}
	k := builder.Key{}
	f := __f

	k.Alloc(16)
	k.Grow(128)

	err := f.To(&k, &j)
	require.NoError(t, err)

	t.Logf("%s", &k)
}

const KEY = `[[5EQT&10EQ]&[8IN&11IN&14IN]&[T|13GT|15NA]]`

func BenchmarkKey(b *testing.B) {
	j := model.Object{}
	f := __f
	var q string
	for i := 0; i < b.N; i++ {
		k := builder.Key{}
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
			//b.Fatal(q)
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	f := __f

	j, err := filter.MarshalJSON(f)
	require.NoError(t, err)

	t.Logf("%s", j)

	var m []any
	err = json.Unmarshal(j, &m)
	require.NoError(t, err)

	t.Logf("%v", m)

	x, err := filter.UnmarshalJSON(j)
	require.NoError(t, err)
	require.Equal(t, f, x)
}

func TestX(t *testing.T) {

	builder.WritePresentRanges([]int{1, 1, 2, 3, 4, 5, 0, 0, 0, 0, 10, 11, 12})
	//                               0  1  2  3  4  5  6  7  8  9  10  11  12
	builder.WritePresentRanges([]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 00, 00, 00})
	//                                           4     6  7  8  9
	builder.WritePresentRanges([]int{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 00, 00, 00})
	builder.WritePresentRanges([]int{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 00, 00, 00})
	builder.WritePresentRanges([]int{1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 00, 00, 00})
	builder.WritePresentRanges([]int{1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 00, 00, 00})
	builder.WritePresentRanges([]int{1, 1, 1, 0, 1, 5, 6, 7, 8, 0, 00, 00, 12})
	builder.WritePresentRanges([]int{1, 1, 1, 0, 1, 5, 6, 7, 8, 0, 00, 11, 00})
	builder.WritePresentRanges([]int{1, 1, 1, 1, 1, 5, 6, 7, 8, 1, 11, 11, 11})

}
