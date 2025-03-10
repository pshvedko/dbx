package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID  `db:"id"`
	DomainID  uuid.UUID  `db:"domain_id"`
	Login     string     `db:"login"`
	Password  []byte     `db:"password"`
	Active    bool       `db:"active"`
	Created   time.Time  `db:"created"`
	Updated   time.Time  `db:"updated"`
	Deleted   *time.Time `db:"deleted"`
	UserRoles []UserRole `db:"user_roles"`
}
