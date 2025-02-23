package help

import (
	"time"

	"github.com/google/uuid"

	"github.com/pshvedko/dbx/filter"
)

//type Profile struct {
//}
//
//func (p Profile) PK() filter.PK {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (p Profile) Names() []string {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (p *Profile) Values() []any {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (p Profile) Value(i int) (any, bool, bool) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (p Profile) Get(i int) any {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (p Profile) Copy() filter.Projector {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (p Profile) Table() string {
//	//TODO implement me
//	panic("implement me")
//}
//
//var _ filter.Projector = &Profile{}

type ObjectList = filter.Injectable[*Object]

type Object struct {
	ID      uint32     `json:"id"`                   // 0 pk auto
	Bool    *bool      `json:"o_bool,omitempty"`     // 1
	Float32 float32    `json:"o_float_32,omitempty"` // 2
	Float64 *float64   `json:"o_float_64,omitempty"` // 3
	Int     int        `json:"o_int,omitempty"`      // 4
	Int16   *int16     `json:"o_int_16,omitempty"`   // 5
	Null    any        `json:"o_null,omitempty"`     // 6
	String1 *string    `json:"o_string_1,omitempty"` // 7 auto
	String2 string     `json:"o_string_2,omitempty"` // 8 null
	String3 string     `json:"o_string_3,omitempty"` // 9 null
	Uint64  *uint64    `json:"o_uint_64,omitempty"`  // 0
	UUID1   uuid.UUID  `json:"o_uuid_1,omitempty"`   // 1
	UUID2   *uuid.UUID `json:"o_uuid_2,omitempty"`   // 2
	UUID3   *uuid.UUID `json:"o_uuid_3,omitempty"`   // 3
	UUID4   uuid.UUID  `json:"o_uuid_4,omitempty"`   // 4 auto
	Time0   time.Time  `json:"o_time_0,omitempty"`   // 5 null
	Time1   time.Time  `json:"o_time_1,omitempty"`   // 6
	Time2   *time.Time `json:"o_time_2,omitempty"`   // 7
	Time3   *time.Time `json:"o_time_3,omitempty"`   // 8
	Time4   time.Time  `json:"o_time_4,omitempty"`   // 9 null
	//Parent  *Object    `json:"o_parent,omitempty"`   // 0 join,on=o_uuid_2
	//Profile Profile    `json:"o_profile,omitempty"`  // 0 join,on=o_uuid_1
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
	return "objects"
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
