package help

import (
	"time"

	"github.com/google/uuid"

	"github.com/pshvedko/dbx/filter"
)

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
	}
}

func (o *Object) Values() []any {
	return []any{
		&o.ID, &o.Bool, &o.Float32, &o.Float64, &o.Int, &o.Int16, &o.Null, &o.String1,
		filter.Nil(&o.String2), filter.Nil(&o.String3), &o.Uint64,
		&o.UUID1, &o.UUID2, &o.UUID3, &o.UUID4,
		&o.Time0, &o.Time1, &o.Time2, &o.Time3, filter.Nil(&o.Time4),
	}
}

func (o Object) Value(i int) (any, bool, bool) {
	v := o.Get(i)
	switch i {
	case 0, 7, 14:
		return v, v == nil, true
	default:
		return v, v == nil, false
	}
}

func (o Object) Get(i int) any {
	switch i {
	case 0:
		return filter.NilIfZero(o.ID)
	case 1:
		return o.Bool
	case 2:
		return o.Float32
	case 3:
		return o.Float64
	case 4:
		return o.Int
	case 5:
		return o.Int16
	case 6:
		return o.Null
	case 7:
		return filter.NilIfZero(o.String1)
	case 8:
		return filter.NilIfZero(o.String2)
	case 9:
		return filter.NilIfZero(o.String3)
	case 10:
		return o.Uint64
	case 11:
		return o.UUID1
	case 12:
		return o.UUID2
	case 13:
		return o.UUID3
	case 14:
		return filter.NilIfZero(o.UUID4)
	case 15:
		return filter.NilIfZero(o.Time0)
	case 16:
		return o.Time1
	case 17:
		return o.Time2
	case 18:
		return o.Time3
	case 19:
		return filter.NilIfZero(o.Time4)
	default:
		panic("illegal index")
	}
}
