package model

import (
	"github.com/google/uuid"
	"time"
)

const UserTable = "users"

type User struct {
	ID        uuid.UUID  `db:"id,primary"`
	DomainID  uuid.UUID  `db:"domain_id"`
	Login     string     `db:"login"`
	Password  []byte     `db:"password"`
	Active    *bool      `db:"active,auto"`
	UserRoles []UserRole `db:"user_roles"`
	Created   time.Time  `db:"created"`
	Updated   time.Time  `db:"updated"`
	Deleted   *time.Time `db:"deleted"`
}

type UserPK0 = uuid.UUID
