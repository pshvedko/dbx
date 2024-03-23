package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Operation [3]any

func (f Operation) Filter() (Filter, error) {
	switch k := f[0].(type) {
	case string:
		switch o := f[1].(type) {
		case string:
			switch o {
			case "EQ":
				return Eq{k: f[2]}, nil
			case "NE":
				return Ne{k: f[2]}, nil
			case "GE":
				return Ge{k: f[2]}, nil
			case "GT":
				return Gt{k: f[2]}, nil
			case "LE":
				return Le{k: f[2]}, nil
			case "LT":
				return Lt{k: f[2]}, nil
			}
			return nil, fmt.Errorf("unknown operation")
		}
		return nil, fmt.Errorf("illegal operation")
	}
	return nil, fmt.Errorf("illegal field")
}

type Filterer interface {
	Filter() (Filter, error)
}

type Expression []Filterer

func (e Expression) Filter() (Filter, error) {
	for _, v := range e {
		switch x := v.(type) {
		case Operation:
		case Expression:
			// FIXME
			return x.Filter()
		}
	}
	return nil, nil
}

func (e *Expression) UnmarshalJSON(b []byte) error {
	j := json.NewDecoder(bytes.NewReader(b))
	n := 0
	for n < 3 {
		t, err := j.Token()
		switch err {
		case nil:
			switch t {
			case json.Delim('['):
				n++
				continue
			}
		case io.EOF:
		default:
			return err
		}
		break
	}
	switch n {
	case 2:
		return ExpressionJSON(b, e, []Operation{})
	default:
		return ExpressionJSON(b, e, []Expression{})
	}
}

func ExpressionJSON[T Filterer](b []byte, e *Expression, a []T) error {
	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}
	for i := range a {
		*e = append(*e, a[i])
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
