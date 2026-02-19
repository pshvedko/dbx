package builder_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test/model"
)

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
			//			b.Fatal(q)
		}
	}
	//b.Logf(q)
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
