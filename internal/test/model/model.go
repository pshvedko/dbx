package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/pshvedko/dbx/filter"
)

//go:generate go run ../../../cmd/dbx -x Object User

func UUIDNilIfZero[T uuid.UUID | *uuid.UUID](v T) (any, bool) {
	return filter.NilIfZero(v)
}

func BoolNilIfZero[T bool | *bool](v T) (any, bool) {
	return filter.NilIfZero(v)
}

func StringNilIfZero[T string | *string](v T) (any, bool) {
	return filter.NilIfZero(v)
}

func Int64NilIfZero[T int64 | *int64](v T) (any, bool) {
	return filter.NilIfZero(v)
}

func Int32NilIfZero[T int32 | *int32](v T) (any, bool) {
	return filter.NilIfZero(v)
}

func Float64NilIfZero[T float64 | *float64](v T) (any, bool) {
	return filter.NilIfZero(v)
}

func TimeNilIfZero[T time.Time | *time.Time](v T) (any, bool) {
	return filter.NilIfZero(v)
}

func BitsNilIfZero(v pgtype.Bits) (any, bool) {
	if !v.Valid && v.Len == 0 && v.Bytes == nil {
		return nil, true
	}
	return nil, true
}

func UserRoleNilIfZero(v []UserRole) (any, bool) {
	if v == nil {
		return nil, true
	}
	return v, false
}
