package filter

import (
	"encoding/json"
)

func MarshalJSON(f Filter) ([]byte, error) {
	switch x := f.(type) {
	case Eq:
		return OperationJSON(x, "EQ")
	case Ne:
		return OperationJSON(x, "NE")
	case Ge:
		return OperationJSON(x, "GE")
	case Gt:
		return OperationJSON(x, "GT")
	case Le:
		return OperationJSON(x, "LE")
	case Lt:
		return OperationJSON(x, "LT")
	case And:
		// a&&b [[a,b]]
		var a []any
		for _, v := range x {
			a = append(a, v)
		}
		return json.Marshal([][]any{a})
	case Or:
		// a||b [[a][b]]
		var a [][]any
		for _, v := range x {
			a = append(a, []any{v})
		}
		return json.Marshal(a)
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
