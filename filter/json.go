package filter

import (
	"encoding/json"
)

func MarshalJSON(f Filter) ([]byte, error) {
	switch x := f.(type) {
	case Eq:
		return OperationJSON(x, "==")
	case Ne:
		return OperationJSON(x, "<>")
	default:
		panic(x)
	}
}

func OperationJSON(m map[string]any, o string) ([]byte, error) {
	var a [][]any
	for k, v := range m {
		a = append(a, []any{k, o, v})
	}
	return json.Marshal(a)
}
