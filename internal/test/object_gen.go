package test

import (
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test/model"
)

type ObjectX model.Object

func (o ObjectX) Self() filter.Copier {
	return o
}

func (ObjectX) PK() filter.PK {
	return []string{"id"}
}

func (o ObjectX) Copy() filter.Projector {
	return &o
}

func (ObjectX) Table() string {
	return "objects"
}

func (ObjectX) Names() []string {
	return []string{
		"id",                         // 1
		"uuid_2", "uuid_3", "uuid_4", // 2 3 4
		"bool_1", "bool_2", "bool_3", "bool_4", // 5 6 7 8
		"float_32", "float_64", // 9 10
		"int_8", "int_16", "int_32", "int_64", // 11 12 13 14
		"string_1", "string_2", "string_3", "string_4", // 15 16 17 18
		"time_1", "time_2", "time_3", "time_4", // 19 20 21 22
	}
}

func (o *ObjectX) Values() []any {
	return []any{
		&o.ID,
		&o.UUID2, &o.UUID3, &o.UUID4,
		&o.Bool1, &o.Bool2, &o.Bool3, &o.Bool4,
		&o.Float32, &o.Float64,
		&o.Int8, &o.Int16, &o.Int32, &o.Int64,
		&o.String1, &o.String2, &o.String3, &o.String4,
		&o.Time1, &o.Time2, &o.Time3, &o.Time4,
	}
}

func (o ObjectX) Value(i int) (any, bool, bool) {
	v, ok := o.Get(i)
	switch i {
	case 0, 2, 10, 16, 20: // id uuid_3 int_8 string_3 time_3
		return v, ok, true
	default:
		return v, ok, false
	}
}

func (o ObjectX) Get(i int) (any, bool) {
	switch i {
	case 0:
		return filter.NilIfZero(o.ID)
	case 1:
		return o.UUID2, false
	case 2:
		return filter.NilIfZero(o.UUID3)
	case 3:
		return filter.NilIfZero(o.UUID4)
	case 4:
		return o.Bool1, false
	case 5:
		return o.Bool2, false
	case 6:
		return filter.NilIfZero(o.Bool3)
	case 7:
		return filter.NilIfZero(o.Bool4)
	case 8:
		return o.Float32, false
	case 9:
		return filter.NilIfZero(o.Float64)
	case 10:
		return filter.NilIfZero(o.Int8)
	case 11:
		return o.Int16, false
	case 12:
		return filter.NilIfZero(o.Int32)
	case 13:
		return filter.NilIfZero(o.Int64)
	case 14:
		return o.String1, false
	case 15:
		return o.String2, false
	case 16:
		return filter.NilIfZero(o.String3)
	case 17:
		return filter.NilIfZero(o.String4)
	case 18:
		return o.Time1, false
	case 19:
		return o.Time2, false
	case 20:
		return filter.NilIfZero(o.Time3)
	case 21:
		return filter.NilIfZero(o.Time4)
	default:
		panic(i)
	}
}
