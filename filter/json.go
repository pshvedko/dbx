package filter

import (
	"encoding/json"
)

type ProtoFilter []any

func (f *ProtoFilter) UnmarshalJSON(b []byte) error {
	var a []any
	err := json.Unmarshal(b, &a)
	switch e := err.(type) {
	case nil:
		*f = append(*f, a...)
	case *json.UnmarshalTypeError:
		if e.Type != nil {
			return &json.UnmarshalTypeError{}
		}
		var v [3]any
		err = json.Unmarshal(b, &v)
		switch e := err.(type) {
		case nil:
			*f = append(*f, v)
		case *json.UnmarshalTypeError:
			return &json.UnsupportedTypeError{Type: e.Type}
		default:
			return err
		}
	default:
		return err
	}
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
