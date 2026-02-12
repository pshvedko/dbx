package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
)

const ObjectTable = "objects"

type Object struct {
	ID      uuid.UUID   `db:"id,primary"`
	UUID2   uuid.UUID   `db:"uuid_2,key=2"`
	UUID3   *uuid.UUID  `db:"uuid_3,auto"`
	UUID4   *uuid.UUID  `db:"uuid_4,key=1"`
	Bool1   bool        `db:"bool_1"`
	Bool2   bool        `db:"bool_2"`
	Bool3   *bool       `db:"bool_3"`
	Bool4   *bool       `db:"bool_4"`
	Float32 float32     `db:"float_32"`
	Float64 *float64    `db:"float_64"`
	Int8    pgtype.Bits `db:"int_8,zero,auto"`
	Int16   int16       `db:"int_16"`
	Int32   *int32      `db:"int_32"`
	Int64   *int64      `db:"int_64"`
	String1 string      `db:"string_1"`
	String2 string      `db:"string_2"`
	String3 *string     `db:"string_3,auto"`
	String4 *string     `db:"string_4"`
	Time1   time.Time   `db:"time_1"`
	Time2   time.Time   `db:"time_2"`
	Time3   *time.Time  `db:"time_3,auto"`
	Time4   *time.Time  `db:"time_4"`
}
