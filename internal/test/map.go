package test

import (
	"sync"

	"database/sql/driver"

	"github.com/pshvedko/dbx/filter"
)

type Arrays map[string]*Array

type Map struct {
	sync.Mutex
	Arrays
}

func (m *Map) Get(k string) *Array {
	m.Lock()
	defer m.Unlock()
	if m.Arrays == nil {
		m.Arrays = Arrays{}
	}
	a, ok := m.Arrays[k]
	if !ok {
		a = &Array{}
		m.Arrays[k] = a
	}
	return a
}

func (m *Map) Add(k string, v any) (string, any) {
	m.Get(k).Add(v)
	return k, v
}

func (m *Map) Range(f func(k string, v *Array)) {
	m.Lock()
	defer m.Unlock()
	for k, v := range m.Arrays {
		f(k, v)
	}
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
