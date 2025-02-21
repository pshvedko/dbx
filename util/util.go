package util

import (
	"time"

	"github.com/google/uuid"
)

func PtrBool(v bool) *bool {
	return &v
}

func PtrFloat64(v float64) *float64 {
	return &v
}

func PtrInt16(v int16) *int16 {
	return &v
}

func PtrUint(v uint) *uint {
	return &v
}

func PtrString(v string) *string {
	return &v
}

func PtrUUID(v uuid.UUID) *uuid.UUID {
	return &v
}

func PtrTime(v time.Time) *time.Time {
	return &v
}
