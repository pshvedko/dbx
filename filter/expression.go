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
	ErrMalformedOperation  = errors.New("malformed operation")
	ErrIllegalOperation    = errors.New("illegal operation")
	ErrUnknownOperation    = errors.New("unknown operation")
)

type Operation [3]any

func (o Operation) Filter() (Filter, error) {
	switch k := o[0].(type) {
	case nil:
		switch o[1].(type) {
		case nil:
			switch v := o[2].(type) {
			case bool:
				if v {
					return True{}, nil
				}
				return False{}, nil
			}
		}
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
			if len(a) == 1 {
				return a[0], nil
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
			if len(a) == 1 {
				return a[0], nil
			}
			return a, nil
		}
	case Operation:
		var a And
		for _, v := range e {
			switch o := v.(type) {
			case Operation:
				t, err := o.Filter()
				if err != nil {
					return nil, err
				}
				a = append(a, t)
			default:
				return nil, ErrIllegalExpression
			}
		}
		return Collapse(a), nil
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
		return OperationJSON(x)
	case Ne:
		return OperationJSON(x)
	case Ge:
		return OperationJSON(x)
	case Gt:
		return OperationJSON(x)
	case Le:
		return OperationJSON(x)
	case Lt:
		return OperationJSON(x)
	case As:
		return OperationJSON(x)
	case Na:
		return OperationJSON(x)
	case In:
		return OperationJSON(x)
	case Ni:
		return OperationJSON(x)
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
	case True:
		return json.Marshal([1][3]any{{nil, nil, true}})
	case False:
		return json.Marshal([1][3]any{{nil, nil, false}})
	default:
		panic(x)
	}
}

func OperationJSON[T any, M interface {
	~map[string]T
	Type() Type
}](x M) ([]byte, error) {
	a := make([][]any, 0, len(x))
	for k, v := range x {
		a = append(a, []any{k, x.Type().String(), v})
	}
	return json.Marshal(a)
}

func UnmarshalJSON(j []byte) (Filter, error) {
	var e Expression
	err := json.Unmarshal(j, &e)
	if err != nil {
		return nil, err
	}
	return e.Filter()
}

//type Instruction struct {
//	Op    OpCode
//	Field int
//	Value any // Уже типизированный (int, time.Time...)
//}
//
//func (obj *Object) Match(prog []Instruction) bool {
//	// 32 уровня вложенности хватит для любого безумного фильтра с Vue
//	var stack [32]bool
//	top := -1
//
//	for _, inst := range prog {
//		switch inst.Op {
//		case OpEQ, OpGE, OpGT, OpLE, OpLT, OpNE, OpIN, OpNI:
//			top++
//			// Кодогенератор: прямой вызов сравнения без рефлексии
//			stack[top] = obj.Compare(inst.Field, inst.Op, inst.Value)
//
//		case OpAND:
//			stack[top-1] = stack[top-1] && stack[top]
//			top--
//
//		case OpOR:
//			stack[top-1] = stack[top-1] || stack[top]
//			top--
//		}
//	}
//	return stack[0]
//}
//Поскольку при subscribe=true ты один раз парсишь матрешку, тебе выгодно в этот же момент превратить IN [1000 элементов] в map[T]struct{}
//
//func (obj *Object) Compare(field int, op OpCode, val any) bool {
//    switch field {
//    case 5: // Bool2
//        return obj.Bool2 == val.(bool)
//    case 11: // Int16
//        target := val.(int16)
//        switch op {
//        case OpGE: return obj.Int16 >= target
//        case OpEQ: return obj.Int16 == target
//        // ...
//        }
//    case 14: // String1 (для IN)
//        if op == OpIN {
//            // val уже причесан к map[string]struct{} в момент подписки
//            m := val.(map[string]struct{})
//            _, ok := m[obj.String1]
//            return ok
//        }
//        return obj.String1 == val.(string)
//    }
//    return false
//}
//
