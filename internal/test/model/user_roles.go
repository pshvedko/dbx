package model

import (
	"github.com/google/uuid"
)

type UserRole struct {
	DomainID    uuid.UUID    `db:"domain_id"`
	UserID      uuid.UUID    `db:"user_id"`
	RoleID      uuid.UUID    `db:"role_id"`
	Permissions []Permission `db:"permissions"`
}
