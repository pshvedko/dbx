package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
)

type Operation [3]any

func (o Operation) Filter() (Filter, error) {
	switch k := o[0].(type) {
	case string:
		switch x := o[1].(type) {
		case string:
			switch v := o[2].(type) {
			case []any:
				switch x {
				case "IN":
					return In{k: v}, nil
				case "NI":
					return Ni{k: v}, nil
				}
			default:
				switch x {
				case "EQ":
					return Eq{k: v}, nil
				case "NE":
					return Ne{k: v}, nil
				case "GE":
					return Ge{k: v}, nil
				case "GT":
					return Gt{k: v}, nil
				case "LE":
					return Le{k: v}, nil
				case "LT":
					return Lt{k: v}, nil
				case "AS":
					return As{k: v}, nil
				}
			}
			return nil, fmt.Errorf("unknown operation")
		}
		return nil, fmt.Errorf("illegal operation")
	}
	return nil, fmt.Errorf("malformedd operation")
}

func Append[M ~map[K]V, K comparable, V any](t M, f M) error {
	maps.Copy(t, f)
	return nil
}

func (o Operation) Join(t Filter) error {
	f, err := o.Filter()
	if err != nil {
		return err
	}

	switch a := t.(type) {
	case Eq:
		switch b := f.(type) {
		case Eq:
			return Append(a, b)
		}
	case Ne:
		switch b := f.(type) {
		case Ne:
			return Append(a, b)
		}
	case Ge:
		switch b := f.(type) {
		case Ge:
			return Append(a, b)
		}
	case Gt:
		switch b := f.(type) {
		case Gt:
			return Append(a, b)
		}
	case Le:
		switch b := f.(type) {
		case Le:
			return Append(a, b)
		}
	case Lt:
		switch b := f.(type) {
		case Lt:
			return Append(a, b)
		}
	case As:
		switch b := f.(type) {
		case As:
			return Append(a, b)
		}
	case In:
		switch b := f.(type) {
		case In:
			return Append(a, b)
		}
	case Ni:
		switch b := f.(type) {
		case Ni:
			return Append(a, b)
		}
	}
	return fmt.Errorf("inopportune operation")
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
		return nil, nil
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
				err = o.Join(f)
				if err != nil {
					return nil, err
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
	case As:
		return OperationJSON(x, "AS")
	case In:
		return OperationJSON(x, "IN")
	case Ni:
		return OperationJSON(x, "NI")
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

func OperationJSON[T any](x map[string]T, o string) ([]byte, error) {
	a := make([][]any, 0, len(x))
	for k, v := range x {
		a = append(a, []any{k, o, v})
	}
	return json.Marshal(a)
}
