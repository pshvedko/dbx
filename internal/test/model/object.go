package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/pshvedko/dbx/filter"
)

type Subject struct {
	ID uint32 `db:"id,pk"`
}

func (p Subject) PK() filter.PK {
	//TODO implement me
	panic("implement me")
}

func (p Subject) Names() []string {
	//TODO implement me
	panic("implement me")
}

func (p *Subject) Values() []any {
	//TODO implement me
	panic("implement me")
}

func (p Subject) Value(i int) (any, bool, bool) {
	//TODO implement me
	panic("implement me")
}

func (p Subject) Get(i int) (any, bool) {
	//TODO implement me
	panic("implement me")
}

func (p Subject) Copy() filter.Projector {
	//TODO implement me
	panic("implement me")
}

func (p Subject) Table() string {
	//TODO implement me
	panic("implement me")
}

type ObjectList = filter.Injectable[Object]

type Object struct {
	ID      uuid.UUID  `db:"id,primary"`         //
	UUID2   uuid.UUID  `db:"uuid_2,key=2"`       //
	UUID3   *uuid.UUID `db:"uuid_3"`             //
	UUID4   *uuid.UUID `db:"uuid_4,key=1"`       //
	Bool1   bool       `db:"bool_1"`             //
	Bool2   bool       `db:"bool_2"`             //
	Bool3   *bool      `db:"bool_3"`             //
	Bool4   *bool      `db:"bool_4"`             //
	Float32 float32    `db:"float_32"`           //
	Float64 *float64   `db:"float_64"`           //
	Int8    int8       `db:"int_8"`              //
	Int16   int16      `db:"int_16"`             //
	Int32   *int32     `db:"int_32"`             //
	Int64   *int64     `db:"int_64"`             //
	String1 string     `db:"string_1"`           //
	String2 string     `db:"string_2"`           //
	String3 *string    `db:"string_3"`           //
	String4 *string    `db:"string_4"`           //
	Time1   time.Time  `db:"time_1"`             //
	Time2   time.Time  `db:"time_2"`             //
	Time3   *time.Time `db:"time_3"`             //
	Time4   *time.Time `db:"time_4"`             //
	Parent  *Object    `db:"parent,join,key=1"`  // left join parent on parent.id = object.uuid_4
	Child   []Object   `db:"child,join,key=1"`   // where child.uuid_4 = $object.id
	Subject Subject    `db:"subject,join,key=2"` // join subject on subject.id = object.uuid_2
}

func (o Object) Self() filter.Copier {
	return o
}

func (Object) PK() filter.PK {
	return []string{"id"}
}

func (o Object) Copy() filter.Projector {
	return &o
}

func (Object) Table() string {
	return "objects"
}

func (Object) Names() []string {
	return []string{
		"id",
		"uuid_2", "uuid_3", "uuid_4",
		"bool_1", "bool_2", "bool_3", "bool_4",
		"float_32", "float_64",
		"int_8", "int_16", "int_32", "int_64",
		"string_1", "string_2", "string_3", "string_4",
		"time_1", "time_2", "time_3", "time_4",
	}
}

func (o *Object) Values() []any {
	return []any{
		&o.ID,
		&o.UUID2, &o.UUID3, &o.UUID4,
		&o.Bool1, &o.Bool2, &o.Bool3, &o.Bool4,
		&o.Float32, &o.Float64,
		&o.Int8, &o.Int16, &o.Int32, &o.Int64,
		&o.String1, filter.Nil(&o.String2), &o.String3, &o.String4,
		&o.Time1, &o.Time2, &o.Time3, &o.Time4,
	}
}

func (o Object) Value(i int) (any, bool, bool) {
	v, ok := o.Get(i)
	switch i {
	case 0, 7, 14:
		fallthrough
	default:
		return v, ok, false
	}
}

func (o Object) Get(i int) (any, bool) {
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
	case 7:
		return filter.NilIfZero(o.Bool3)
	case 8:
		return filter.NilIfZero(o.Bool4)
	case 9:
		return o.Float32, false
	case 10:
		return filter.NilIfZero(o.Float64)
	case 11:
		return o.Int8, false
	case 12:
		return o.Int16, false
	case 13:
		return filter.NilIfZero(o.Int32)
	case 14:
		return filter.NilIfZero(o.Int64)
	case 15:
		return o.String1, false
	case 16:
		return o.String2, false
	case 17:
		return filter.NilIfZero(o.String3)
	case 18:
		return filter.NilIfZero(o.String4)
	case 19:
		return o.Time1, false
	case 20:
		return o.Time2, false
	case 21:
		return filter.NilIfZero(o.Time3)
	case 22:
		return filter.NilIfZero(o.Time4)
	default:
		panic("illegal index")
	}
}
