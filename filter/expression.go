package filter

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrMalformedExpression = errors.New("malformed expression")
	ErrIllegalExpression   = errors.New("illegal expression")
	ErrUnknownExpression   = errors.New("unknown expression")
	ErrEmptyExpression     = errors.New("empty expression")
	ErrUnsuitableOperation = errors.New("unsuitable operation")
	ErrMalformedOperation  = errors.New("malformed operation")
	ErrIllegalOperation    = errors.New("illegal operation")
	ErrUnknownOperation    = errors.New("unknown operation")
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
				default:
					switch s := v.(type) {
					case string:
						switch x {
						case "AS":
							return As{k: s}, nil
						case "NA":
							return Na{k: s}, nil
						}
					}
				}
			}
			return nil, ErrUnknownOperation
		}
		return nil, ErrIllegalOperation
	}
	return nil, ErrMalformedOperation
}

func Append[M interface {
	~map[K]V
	Filter
}, K string, V any](t M, f M) M {
	for k, v := range f {
		t[k] = v
	}
	return t
}

func Join(t Filter, f ...Filter) (Filter, error) {
	if len(f) == 0 {
		return t, nil
	}
	switch a := t.(type) {
	case Eq:
		switch b := f[0].(type) {
		case Eq:
			return Join(Append(a, b), f[1:]...)
		}
	case Ne:
		switch b := f[0].(type) {
		case Ne:
			return Join(Append(a, b), f[1:]...)
		}
	case Ge:
		switch b := f[0].(type) {
		case Ge:
			return Join(Append(a, b), f[1:]...)
		}
	case Gt:
		switch b := f[0].(type) {
		case Gt:
			return Join(Append(a, b), f[1:]...)
		}
	case Le:
		switch b := f[0].(type) {
		case Le:
			return Join(Append(a, b), f[1:]...)
		}
	case Lt:
		switch b := f[0].(type) {
		case Lt:
			return Join(Append(a, b), f[1:]...)
		}
	case As:
		switch b := f[0].(type) {
		case As:
			return Join(Append(a, b), f[1:]...)
		}
	case In:
		switch b := f[0].(type) {
		case In:
			return Join(Append(a, b), f[1:]...)
		}
	case Ni:
		switch b := f[0].(type) {
		case Ni:
			return Join(Append(a, b), f[1:]...)
		}
	}
	return nil, ErrUnsuitableOperation
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
			for _, o := range x {
				switch o.(type) {
				case Expression:
					f, err := o.Filter()
					if err != nil {
						return nil, err
					}
					if f == nil {
						return nil, ErrEmptyExpression
					}
					a = append(a, f)
				default:
					return nil, ErrIllegalExpression
				}
			}
			return a, nil
		default:
			var a Or
			for _, v := range e {
				switch o := v.(type) {
				case Expression:
					if len(o) != 1 {
						return nil, ErrMalformedExpression
					}
					f, err := o[0].Filter()
					if err != nil {
						return nil, err
					}
					if f == nil {
						return nil, ErrEmptyExpression
					}
					a = append(a, f)
				default:
					return nil, ErrIllegalExpression
				}
			}
			return a, nil
		}
	case Operation:
		t, err := x.Filter()
		if err != nil {
			return nil, err
		}
		for _, v := range e[1:] {
			switch o := v.(type) {
			case Operation:
				var f Filter
				f, err = o.Filter()
				if err != nil {
					return nil, err
				}
				t, err = Join(t, f)
				if err != nil {
					return nil, err
				}
			default:
				return nil, ErrIllegalExpression
			}
		}
		return t, nil
	default:
		return nil, ErrUnknownExpression
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
