package filter

import (
	"bytes"
	"encoding/json"
)

type Operation [3]any

type Expression []any

func (f *Expression) UnmarshalJSON(b []byte) error {
	switch {
	case len(b) > 2 && bytes.Equal(b[:3], []byte{'[', '[', '"'}):
		return UnmarshalJSON(f, b, []Operation{})
	default:
		return UnmarshalJSON(f, b, []Expression{})
	}
}

func UnmarshalJSON[T interface{ Expression | Operation }](f *Expression, b []byte, a []T) error {
	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}
	for _, x := range a {
		*f = append(*f, x)
	}
	return nil
}

func (f Expression) Filter() Filter {

	return nil
}

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
		// [[a,b]]
		a := make([]any, 0, len(x))
		for _, v := range x {
			a = append(a, v)
		}
		return json.Marshal([][]any{a})
	case Or:
		// [[a][b]]
		a := make([][]any, 0, len(x))
		for _, v := range x {
			a = append(a, []any{v})
		}
		return json.Marshal(a)
	default:
		panic(x)
	}
}

func OperationJSON(x map[string]any, o string) ([]byte, error) {
	a := make([][]any, 0, len(x))
	for k, v := range x {
		a = append(a, []any{k, o, v})
	}
	return json.Marshal(a)
}
