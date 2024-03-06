package help

import (
	"context"
	"github.com/google/uuid"
	"github.com/pshvedko/dbx/filter"
	"log/slog"
	"testing"
	"time"
)

type ObjectList []Object

func (o ObjectList) Get() filter.Projector {
	return &Object{}
}

func (o *ObjectList) Put(j filter.Projector) {
	switch v := j.(type) {
	case *Object:
		*o = append(*o, *v)
	default:
		panic("invalid injection")
	}
}

type Object struct {
	ID      uint32     `json:"id"`
	Bool    *bool      `json:"o_bool,omitempty"`
	Float32 float32    `json:"o_float_32,omitempty"`
	Float64 *float64   `json:"o_float_64,omitempty"`
	Int     int        `json:"o_int,omitempty"`
	Int16   *int16     `json:"o_int_16,omitempty"`
	Null    any        `json:"o_null,omitempty"`
	String1 *string    `json:"o_string_1,omitempty"`
	String2 string     `json:"o_string_2,omitempty"`
	String3 string     `json:"o_string_3,omitempty"`
	Uint64  *uint64    `json:"o_uint_64,omitempty"`
	UUID1   uuid.UUID  `json:"o_uuid_1,omitempty"`
	UUID2   *uuid.UUID `json:"o_uuid_2,omitempty"`
	UUID3   *uuid.UUID `json:"o_uuid_3,omitempty"`
	UUID4   uuid.UUID  `json:"o_uuid_4,omitempty"`
	Time1   time.Time  `json:"o_time_1,omitempty"`
	Time2   *time.Time `json:"o_time_2,omitempty"`
	Time3   *time.Time `json:"o_time_3,omitempty"`
	Time4   time.Time  `json:"o_time_4,omitempty"`
}

func (o Object) Table() string {
	return "objects"
}

func (o Object) Names() []string {
	return []string{
		"id", "o_bool", "o_float_32", "o_float_64", "o_int", "o_int_16", "o_null", "o_string_1",
		"o_string_2", "o_string_3", "o_uint_64",
		"o_uuid_1", "o_uuid_2", "o_uuid_3", "o_uuid_4",
		"o_time_1", "o_time_2", "o_time_3", "o_time_4",
	}
}

func (o *Object) Values() []any {
	return []any{
		&o.ID, &o.Bool, &o.Float32, &o.Float64, &o.Int, &o.Int16, &o.Null, &o.String1,
		&String{x: &o.String2}, &String{x: &o.String3}, &o.Uint64,
		&o.UUID1, &o.UUID2, &o.UUID3, &o.UUID4,
		&o.Time1, &o.Time2, &o.Time3, &Time{x: &o.Time4},
	}
}

type String struct {
	x *string
}

func (s *String) Scan(v any) error {
	switch x := v.(type) {
	case nil:
	case string:
		*s.x = x
	}
	return nil
}

type Time struct {
	x *time.Time
}

func (t *Time) Scan(v any) error {
	switch x := v.(type) {
	case nil:
	case time.Time:
		*t.x = x
	}
	return nil
}

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

type logHandler testing.T

func (h *logHandler) Enabled(context.Context, slog.Level) bool { return true }

func (h *logHandler) Handle(_ context.Context, r slog.Record) error {
	h.Log(r.Message)
	r.Attrs(func(a slog.Attr) bool {
		h.Log(a)
		return true
	})
	return nil
}

func (h *logHandler) WithAttrs([]slog.Attr) slog.Handler { return h }

func (h *logHandler) WithGroup(string) slog.Handler { return h }

func LogHandler(t *testing.T) slog.Handler { return (*logHandler)(t) }
