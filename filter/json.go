package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Operation [3]any

func (o Operation) Filter() (Filter, error) {
	switch k := o[0].(type) {
	case string:
		switch x := o[1].(type) {
		case string:
			switch x {
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
			case "AS":
				return As{k: o[2]}, nil
			}
			return nil, fmt.Errorf("unknown operation")
		}
		return nil, fmt.Errorf("illegal operation")
	}
	return nil, fmt.Errorf("malformed operation")
}

type Filterer interface {
	Filter() (Filter, error)
}

type Expression []Filterer

// Filter ...
//
//	( ( A || B ) && C )
//
//	     M=1
//	    +---+             2
//	    | A |         +---+---+            M
//	 OR +---+ 2       | A | B | N=1       +- NxM
//	    | B |         +---+---+         N |
//	    +---+            AND
//
//	Nx1 - OR:  A || B -> [[ A ] , [ B ]] = D
//
//	1xM - AND: D && C -> [[ D , C ]]
//
//	[[ [[ A ] , [ B ]] , C ]]
//	 \  \____OR_____/      /
//	  \______AND__________/
func (e Expression) Filter() (Filter, error) {
	if len(e) == 0 {
		return nil, fmt.Errorf("empty expression")
	}
	switch x := e[0].(type) {
	case Expression:
		switch len(e) {
		case 1:
			var a And
			switch v := e[0].(type) {
			case Expression:
				for _, o := range v {
					switch o.(type) {
					case Expression:
						f, err := o.Filter()
						if err != nil {
							return nil, err
						}
						a = append(a, f)
					default:
						return nil, fmt.Errorf("malformed expression")
					}
				}
			default:
				return nil, fmt.Errorf("illegal expression")
			}
			return a, nil
		default:
			var a Or
			for _, v := range e {
				switch o := v.(type) {
				case Expression:
					if len(o) != 1 {
						return nil, fmt.Errorf("malformed expression")
					}
					f, err := o[0].Filter()
					if err != nil {
						return nil, err
					}
					a = append(a, f)
				default:
					return nil, fmt.Errorf("illegal expression")
				}
			}
			return a, nil
		}
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
					f = Append(f, k, o[2])
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

func Append(f Filter, k string, v any) Filter {
	switch m := f.(type) {
	case Eq:
		m[k] = v
	case Ne:
		m[k] = v
	case Ge:
		m[k] = v
	case Gt:
		m[k] = v
	case Le:
		m[k] = v
	case Lt:
		m[k] = v
	case As:
		m[k] = v
	default:
		panic(m)
	}
	return f
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

//func normalize(v Filterer) Filterer {
//	switch o := v.(type) {
//	case Operation:
//		switch x := o[2].(type) {
//		case float64:
//			i, f := math.Modf(x)
//			if f == 0 {
//				o[2] = int(i)
//			}
//			return o
//		case string:
//			// 1970-01-01T00:00:00+00:00
//			if len(x) >= 20 && x[4] == '-' && x[7] == '-' && x[10] == 'T' && x[13] == ':' && x[16] == ':' {
//				if x[19] == 'Z' || x[19] == '+' {
//					t, err := time.Parse(time.RFC3339, x)
//					if err == nil {
//						o[2] = t
//						return o
//					}
//				}
//			}
//		}
//	}
//	return v
//}

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
