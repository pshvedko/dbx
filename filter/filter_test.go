package filter_test

import (
	"github.com/google/uuid"
	"testing"

	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test/model"
)

func TestJoin_Names(t *testing.T) {
	a := model.Object{ID: uuid.UUID{1}}
	b := model.Object{ID: uuid.UUID{2}}
	c := model.Object{ID: uuid.UUID{3}}

	t.Log(a.Len())
	t.Log(a.Name())
	t.Log(a.Names())
	t.Log(a.Places())

	j0 := filter.NewJoin(&a, &b, filter.Eq{"id": filter.Column("id")})

	t.Log(j0.Len())
	t.Log(j0.Name())
	t.Log(j0.Names())
	t.Log(j0.Places())

	j1 := filter.NewJoin(&a, filter.NewJoin(&b, &c, filter.Eq{"id": filter.Column("id")}), filter.Eq{"id": filter.Column("id")})

	t.Log(j1.Len())
	t.Log(j1.Name())
	t.Log(j1.Names())
	t.Log(j1.Places())

	j2 := filter.NewJoin(filter.NewJoin(&b, &c, filter.Eq{"id": filter.Column("id")}), &a, filter.Eq{"id": filter.Column("id")})

	t.Log(j2.Len())
	t.Log(j2.Name())
	t.Log(j2.Names())
	t.Log(j2.Places())
}
