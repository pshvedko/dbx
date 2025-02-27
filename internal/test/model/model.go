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

type SubjectList = filter.Injectable[*Subject]

type ObjectList = filter.Injectable[*Object]

type Object struct {
	ID      uuid.UUID  `db:"id,primary"`         //
	UUID2   uuid.UUID  `db:"o_uuid_2,key=2"`     //
	UUID3   *uuid.UUID `db:"o_uuid_3"`           //
	UUID4   *uuid.UUID `db:"o_uuid_4,key=1"`     //
	Bool1   bool       `db:"o_bool_1"`           //
	Bool2   bool       `db:"o_bool_2"`           //
	Bool3   *bool      `db:"o_bool_3"`           //
	Bool4   *bool      `db:"o_bool_4"`           //
	Float32 float32    `db:"o_float_32"`         //
	Float64 *float64   `db:"o_float_64"`         //
	Int8    int8       `db:"o_int_8"`            //
	Int16   int16      `db:"o_int_16"`           //
	Int32   *int32     `db:"o_int"`              //
	Int64   *int64     `db:"o_int_64"`           //
	String1 string     `db:"o_string_1"`         //
	String2 string     `db:"o_string_2"`         //
	String3 *string    `db:"o_string_3"`         //
	String4 *string    `db:"o_string_4"`         //
	Time1   time.Time  `db:"o_time_1"`           //
	Time2   time.Time  `db:"o_time_2"`           //
	Time3   *time.Time `db:"o_time_3"`           //
	Time4   *time.Time `db:"o_time_4"`           //
	Parent  *Object    `db:"parent,join,key=1"`  // left join parent on parent.id = object.o_uuid_4
	Child   []Object   `db:"child,join,key=1"`   // where child.o_uuid_4 = $object.id
	Subject Subject    `db:"subject,join,key=2"` // join subject on subject.id = object.o_uuid_2
}

func (Object) PK() filter.PK {
	return []string{"id"}
}

func (o *Object) Copy() filter.Projector {
	if o == nil {
		return &Object{}
	}
	x := *o
	return &x
}

func (Object) Table() string {
	return "object"
}

func (Object) Names() []string {
	return []string{
		"id", "o_bool", "o_float_32", "o_float_64", "o_int", "o_int_16", "o_null", "o_string_1",
		"o_string_2", "o_string_3", "o_uint_64",
		"o_uuid_1", "o_uuid_2", "o_uuid_3", "o_uuid_4",
		"o_time_0", "o_time_1", "o_time_2", "o_time_3", "o_time_4",
		//		"o_parent",
	}
}

func (o *Object) Values() []any {
	return []any{
		&o.ID, &o.Bool, &o.Float32, &o.Float64, &o.Int, &o.Int16, &o.Null, &o.String1,
		filter.Nil(&o.String2), filter.Nil(&o.String3), &o.Uint64,
		&o.UUID1, &o.UUID2, &o.UUID3, &o.UUID4,
		&o.Time0, &o.Time1, &o.Time2, &o.Time3, filter.Nil(&o.Time4),
		//filter.JoinAny(&o.Parent, "o_uuid_2"),
		//filter.JoinTo(&o.Profile, "o_uuid_1"),
	}
}

func (o Object) Value(i int) (any, bool, bool) {
	v, ok := o.Get(i)
	switch i {
	case 0, 7, 14:
		return v, ok, true
	default:
		return v, ok, false
	}
}

func (o Object) Get(i int) (any, bool) {
	switch i {
	case 0:
		return filter.NilIfZero(o.ID)
	case 1:
		return o.Bool, o.Bool == nil
	case 2:
		return o.Float32, false
	case 3:
		return o.Float64, o.Float64 == nil
	case 4:
		return o.Int, false
	case 5:
		return o.Int16, o.Int16 == nil
	case 6:
		return o.Null, o.Null == nil
	case 7:
		return filter.NilIfZero(o.String1)
	case 8:
		return filter.NilIfZero(o.String2)
	case 9:
		return filter.NilIfZero(o.String3)
	case 10:
		return o.Uint64, o.Uint64 == nil
	case 11:
		return o.UUID1, false
	case 12:
		return o.UUID2, o.UUID2 == nil
	case 13:
		return o.UUID3, o.UUID3 == nil
	case 14:
		return filter.NilIfZero(o.UUID4)
	case 15:
		return filter.NilIfZero(o.Time0)
	case 16:
		return o.Time1, true
	case 17:
		return o.Time2, o.Time2 == nil
	case 18:
		return o.Time3, o.Time3 == nil
	case 19:
		return filter.NilIfZero(o.Time4)
	default:
		panic("illegal index")
	}
}
