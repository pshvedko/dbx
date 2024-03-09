package help

import (
	"context"
	"database/sql/driver"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/request"
)

type ObjectList = filter.Injectable[*Object]

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

func (o *Object) Copy() filter.Projector {
	if o == nil {
		return &Object{}
	}
	x := *o
	return &x
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

func (o Object) Value(i int) (any, bool, bool) {
	v := o.Get(i)
	switch i {
	case 0:
		return v, v == nil, true
	default:
		return v, v == nil, false
	}
}

func (o Object) Get(i int) any {
	switch i {
	case 0:
		return filter.Nil(o.ID)
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
		return o.String1
	case 8:
		return filter.Nil(o.String2)
	case 9:
		return filter.Nil(o.String3)
	case 10:
		return o.Uint64
	case 11:
		return o.UUID1
	case 12:
		return o.UUID2
	case 13:
		return o.UUID3
	case 14:
		return filter.Nil(o.UUID4)
	case 15:
		return o.Time1
	case 16:
		return o.Time2
	case 17:
		return o.Time3
	case 18:
		return filter.Nil(o.Time4)
	default:
		panic("illegal index")
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
	h.Log(r.Level, r.Message)
	r.Attrs(func(a slog.Attr) bool {
		h.Log(r.Level, a)
		return true
	})
	return nil
}

func (h *logHandler) WithAttrs([]slog.Attr) slog.Handler { return h }

func (h *logHandler) WithGroup(string) slog.Handler { return h }

func LogHandler(t *testing.T) slog.Handler { return (*logHandler)(t) }

type DB struct {
	*sqlx.DB
	*slog.Logger
}

func (DB) Connect(context.Context) (*sqlx.Conn, error) { return nil, nil }

func (DB) Option() []request.Option { return nil }

type Map struct {
	sync.Mutex
	M map[string]*Array
}

func (m *Map) Get(k string) *Array {
	m.Lock()
	defer m.Unlock()
	if m.M == nil {
		m.M = map[string]*Array{}
	}
	a, ok := m.M[k]
	if !ok {
		a = &Array{}
		m.M[k] = a
	}
	return a
}

func (m *Map) Add(k string, v any) (string, any) {
	m.Get(k).Add(v)
	return k, v
}

type Array struct {
	sync.Mutex
	filter.Array
}

func (a *Array) Add(v any) {
	a.Lock()
	defer a.Unlock()
	a.Array = append(a.Array, v)
}

func (a *Array) Value() (driver.Value, error) {
	a.Lock()
	defer a.Unlock()
	return a.Array.Value()
}
