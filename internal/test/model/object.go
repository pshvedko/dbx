package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
)

const ObjectTable = "objects"

type Object struct {
	ID      uuid.UUID   `db:"id,primary"`      //0
	UUID2   uuid.UUID   `db:"uuid_2,key=2"`    //1
	UUID3   *uuid.UUID  `db:"uuid_3,auto"`     //2
	UUID4   *uuid.UUID  `db:"uuid_4,key=1"`    //3
	Bool1   bool        `db:"bool_1"`          //4
	Bool2   bool        `db:"bool_2"`          //5
	Bool3   *bool       `db:"bool_3"`          //6
	Bool4   *bool       `db:"bool_4"`          //7
	Float32 float32     `db:"float_32"`        //8
	Float64 *float64    `db:"float_64"`        //9
	Int8    pgtype.Bits `db:"int_8,zero,auto"` //10
	Int16   int16       `db:"int_16"`          //11
	Int32   *int32      `db:"int_32"`          //12
	Int64   *int64      `db:"int_64"`          //13
	String1 string      `db:"string_1"`        //14
	String2 string      `db:"string_2"`        //15
	String3 *string     `db:"string_3,auto"`   //16
	String4 *string     `db:"string_4"`        //17
	Time1   time.Time   `db:"time_1"`          //18
	Time2   time.Time   `db:"time_2"`          //19
	Time3   *time.Time  `db:"time_3,auto"`     //20
	Time4   *time.Time  `db:"time_4"`          //21
}
