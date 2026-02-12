package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
)

const ObjectTable = "objects"

type Object struct {
	ID      uuid.UUID   `db:"id,primary"`      // 1
	UUID2   uuid.UUID   `db:"uuid_2,key=2"`    // 2
	UUID3   *uuid.UUID  `db:"uuid_3,auto"`     // 3
	UUID4   *uuid.UUID  `db:"uuid_4,key=1"`    // 4
	Bool1   bool        `db:"bool_1"`          // 5
	Bool2   bool        `db:"bool_2"`          // 6
	Bool3   *bool       `db:"bool_3"`          // 7
	Bool4   *bool       `db:"bool_4"`          // 8
	Float32 float32     `db:"float_32"`        // 9
	Float64 *float64    `db:"float_64"`        // 10
	Int8    pgtype.Bits `db:"int_8,null,auto"` // 11
	Int16   int16       `db:"int_16"`          // 12
	Int32   *int32      `db:"int_32"`          // 13
	Int64   *int64      `db:"int_64"`          // 14
	String1 string      `db:"string_1"`        // 15
	String2 string      `db:"string_2"`        // 16
	String3 *string     `db:"string_3,auto"`   // 17
	String4 *string     `db:"string_4"`        // 18
	Time1   time.Time   `db:"time_1"`          // 19 created
	Time2   time.Time   `db:"time_2"`          // 20 updated
	Time3   *time.Time  `db:"time_3,auto"`     // 21
	Time4   *time.Time  `db:"time_4"`          // 22 deleted
}
