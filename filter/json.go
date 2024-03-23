package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Operation [3]any

type Operator interface {
	Filter
	Append(string, any)
}

func (o Operation) Filter() (Operator, error) {
	switch k := o[0].(type) {
	case string:
		switch v := o[1].(type) {
		case string:
			switch v {
			case "EQ":
				return Eq{k: o[2]}, nil
			case "NE":
				return Ne{k: o[2]}, nil
			case "GE":
				return Ge{k: o[2]}, nil
			case "GT":
				return Gt{k: o[2]}, nil
			case "LE":
				return Le{k: o[2]}, nil
			case "LT":
				return Lt{k: o[2]}, nil
			}
			return nil, fmt.Errorf("unknown operation")
		}
		return nil, fmt.Errorf("illegal operation")
	}
	return nil, fmt.Errorf("malformed operation")
}

type Filterer interface {
	Filter() (Operator, error)
}

type Expression []Filterer

var _ = Expression{
	Expression{
		Expression{
			Expression{
				Expression{Operation{"f", "GE", 0}},
			},
			Expression{
				Expression{Operation{"b", "EQ", false}},
			},
		},
		Expression{Operation{"f", "LE", 0}},
	},
}

func (e Expression) Filter() (Operator, error) {
	if len(e) == 0 {
		return nil, fmt.Errorf("empty expression")
	}
	switch x := e[0].(type) {
	case Expression:
		panic(2)
	case Operation:
		f, err := x.Filter()
		if err != nil {
			return nil, err
		}
		for _, v := range e[1:] {
			switch o := v.(type) {
			case Operation:
				if x[1] != o[1] {
					return nil, fmt.Errorf("illegal operation")
				}
				switch k := o[0].(type) {
				case string:
					f.Append(k, o[2])
				default:
					return nil, fmt.Errorf("malformed operation")
				}
			default:
				return nil, fmt.Errorf("illegal expression")
			}
		}
		return f, nil
	default:
		return nil, fmt.Errorf("unknown expression")
	}
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
