package model

import (
	"github.com/google/uuid"
	"github.com/pshvedko/dbx/filter"
)

type UserRole struct {
	DomainID    uuid.UUID    `db:"domain_id"`
	UserID      uuid.UUID    `db:"user_id"`
	RoleID      uuid.UUID    `db:"role_id"`
	Permissions []Permission `db:"permissions"`
}

func (u UserRole) Self() filter.Copier {
	//TODO implement me
	panic("implement me")
}

func (u UserRole) PK() filter.PK {
	//TODO implement me
	panic("implement me")
}

func (u UserRole) Names() []string {
	//TODO implement me
	panic("implement me")
}

func (u *UserRole) Values() []any {
	//TODO implement me
	panic("implement me")
}

func (u UserRole) Value(i int) (any, bool, bool) {
	//TODO implement me
	panic("implement me")
}

func (u UserRole) Get(i int) (any, bool) {
	//TODO implement me
	panic("implement me")
}

func (u UserRole) Copy() filter.Projector {
	//TODO implement me
	panic("implement me")
}

func (u UserRole) Table() string {
	//TODO implement me
	panic("implement me")
}
