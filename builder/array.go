package builder

import (
	"database/sql"
	"errors"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test/model"
)

type Reader struct {
	s string
	a []any
	e int
}

func NewReader(s string) *Reader {
	return &Reader{
		s: s,
		a: nil,
	}
}

func (r *Reader) CompareAndReadByte(b byte) bool {
	if len(r.s) == 0 || r.s[0] != b {
		return false
	}
	r.s = r.s[1:]
	return true
}

func (r *Reader) Array() []any {
	return r.a
}

// ReadArray ...
// {"(2eb40782-cf1d-4777-a5ac-fad81ea4f5a2,cd484cf9-0702-451d-8770-f70cc0861e56,9801a142-7d5b-4d27-9049-bc223cfe3a9f)","(2eb40782-cf1d-4777-a5ac-fad81ea4f5a2,cd484cf9-0702-451d-8770-f70cc0861e56,f83693ca-7449-43dc-a061-3e9c86a70638)"}
func (r *Reader) ReadArray() error {
	if !r.CompareAndReadByte('{') {
		return errors.New("not open bracket")
	}
	for {
		err := r.ReadRow()
		if err != nil {
			return err
		}
		if !r.CompareAndReadByte(',') {
			break
		}
	}
	if !r.CompareAndReadByte('}') {
		return errors.New("not close bracket")
	}
	return nil
}

// ReadRow ...
// "(2eb40782-cf1d-4777-a5ac-fad81ea4f5a2,cd484cf9-0702-451d-8770-f70cc0861e56,9801a142-7d5b-4d27-9049-bc223cfe3a9f)","(2eb40782-cf1d-4777-a5ac-fad81ea4f5a2,cd484cf9-0702-451d-8770-f70cc0861e56,f83693ca-7449-43dc-a061-3e9c86a70638)"
func (r *Reader) ReadRow() error {
	if !r.ReadQuote(1) {
		return errors.New("not open quote")
	}
	if !r.CompareAndReadByte('(') {
		return errors.New("not open tuple")
	}
	for {
		err := r.ReadColumn()
		if err != nil {
			return err
		}
		if !r.CompareAndReadByte(',') {
			break
		}
	}
	if !r.CompareAndReadByte(')') {
		return errors.New("not close tuple")
	}
	if !r.ReadQuote(-1) {
		return errors.New("not close quote")
	}
	return nil
}

// ReadColumn ...
// 2eb40782-cf1d-4777-a5ac-fad81ea4f5a2,cd484cf9-0702-451d-8770-f70cc0861e56,9801a142-7d5b-4d27-9049-bc223cfe3a9f,\"2025-03-09 15:01:22.010824+03\"
func (r *Reader) ReadColumn() error {
	q := r.ReadQuote(1)

	if q && !r.ReadQuote(-1) {
		return errors.New("not close quote")
	}
	return nil
}

func (r *Reader) ReadQuote(e int) bool {
	for i := 0; i < r.e; i++ {
		if !r.CompareAndReadByte('\\') {
			return false
		}
	}
	if !r.CompareAndReadByte('"') {
		return false
	}

	r.e += e
	return true
}

type ROW []string

type ARRAY []ROW

// SplitRowArray ...
//
//	+{
//	 +"
//	  +(
//	   +2eb40782-cf1d-4777-a5ac-fad81ea4f5a2,
//	   +cd484cf9-0702-451d-8770-f70cc0861e56,
//	   +,
//	   +9801a142-7d5b-4d27-9049-bc223cfe3a9f,
//	   +\"
//	     +2025-03-09 15:01:22.010824+03
//	   -\",
//	   +\"
//	     +{
//	      +\"
//	      -\"
//	      +(
//	       +8f23f2fa-41b2-4bda-91e9-7bb6f2474f02,
//	       +2eb40782-cf1d-4777-a5ac-fad81ea4f5a2,
//	       +9801a142-7d5b-4d27-9049-bc223cfe3a9f,
//	       +,
//	       +objects,
//	       +00011,
//	       +\\\\\"\"
//	               +2025-03-08 16:01:09.516008+03
//	       -\\\\\"\",
//	       +\\\\\"\"
//	               +2025-03-08 16:01:09.516008+03
//	       -\\\\\"\",
//	       +
//	      -)
//	      +\"
//	      -\",
//	      +\"
//	      -\"
//	      +(
//	       +9ebd7789-f738-4c89-b9d2-f9bac86de3dd,
//	       +2eb40782-cf1d-4777-a5ac-fad81ea4f5a2,9801a142-7d5b-4d27-9049-bc223cfe3a9f,,users,00001,\\\\\"\"2025-03-08 16:12:24.782876+03\\\\\"\",\\\\\"\"2025-03-08 16:12:24.782876+03\\\\\"\",)\"\",\"\"(4906e175-663b-4c13-82bd-296ee6c7d334,2eb40782-cf1d-4777-a5ac-fad81ea4f5a2,9801a142-7d5b-4d27-9049-bc223cfe3a9f,,domains,00001,\\\\\"\"2025-03-08 16:12:43.2469+03\\\\\"\",\\\\\"\"2025-03-08 16:12:43.2469+03\\\\\"\",)\"\",\"\"(1f6e5334-47c1-489d-a31c-296e6e0d5520,2eb40782-cf1d-4777-a5ac-fad81ea4f5a2,9801a142-7d5b-4d27-9049-bc223cfe3a9f,,roles,00001,\\\\\"\"2025-03-08 16:13:06.833611+03\\\\\"\",\\\\\"\"2025-03-08 16:13:06.833611+03\\\\\"\",)\"\",\"\"(dd450a0c-69e8-447c-a80f-5f8841f5f69e,2eb40782-cf1d-4777-a5ac-fad81ea4f5a2,9801a142-7d5b-4d27-9049-bc223cfe3a9f,,permissions,00001,\\\\\"\"2025-03-08 16:13:16.321388+03\\\\\"\",\\\\\"\"2025-03-08 16:13:16.321388+03\\\\\"\",)\"\"}\")","(2eb40782-cf1d-4777-a5ac-fad81ea4f5a2,cd484cf9-0702-451d-8770-f70cc0861e56,f83693ca-7449-43dc-a061-3e9c86a70638,{})"}
func SplitRowArray(array string) ([]any, error) {
	r := NewReader(array)
	err := r.ReadArray()
	if err != nil {
		return nil, err
	}
	return r.Array(), nil
}

type Tokenizer struct {
	s string
	q []string
}

func (t *Tokenizer) Token() (string, error) {
	return "", nil
}

func NewTokenizer(s string) *Tokenizer {
	return &Tokenizer{s: s}
}

func Decode(s string) ([]any, error) {

	return nil, nil
}

type RowArray[T filter.Copier] struct {
	rows *[]T
}

func (rr RowArray[T]) Scan(src any) error {
	switch s := src.(type) {
	case string:
		aa, err := Decode(s)
		if err != nil {
			return err
		}
		return rr.ScanArray(aa)
	default:
		return errors.New("illegal type")
	}
}

func (rr RowArray[T]) ScanArray(aa []any, dst ...any) error {

	return nil
}

func Rows[T filter.Copier](rows *[]T) sql.Scanner {
	return &RowArray[T]{
		rows: rows,
	}
}

var _ = Rows(&[]model.UserRole{})
