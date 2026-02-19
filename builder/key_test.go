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
	f := filter.And{
		filter.Eq{
			"int_8":  "\\x0f01",
			"bool_2": true},
		filter.In{
			"string_1": []any{"yellow", "green"},
			"int_16":   []any{1., 2., 3.},
			"float_32": []any{.0}},
		filter.Or{
			filter.Gt{
				"int_64": 0.},
			filter.Na{
				"string_2": `%ing`}},
	}

	err := f.To(&k, &j)
	if err != nil {
		return
	}

	t.Logf("%s", &k)
	t.Log(k.String())
}

func TestUnmarshalJSON(t *testing.T) {
	f := filter.And{
		filter.Eq{
			"int_8":  "\\x0f01",
			"bool_2": true},
		filter.In{
			"string_1": []any{"yellow", "green"},
			"int_16":   []any{1., 2., 3.},
			"float_32": []any{.0}},
		filter.Or{
			filter.Gt{
				"int_64": 0.},
			filter.Na{
				"string_2": `%ing`}},
	}

	j, err := filter.MarshalJSON(f)
	require.NoError(t, err)

	t.Logf("%s", j)

	var m []any
	err = json.Unmarshal(j, &m)
	require.NoError(t, err)

	t.Logf("%+v", m)

	x, err := filter.UnmarshalJSON(j)
	require.NoError(t, err)
	require.Equal(t, f, x)
}
