package model

import (
	"github.com/google/uuid"
	"time"
)

type Permission struct {
	ID       uuid.UUID  `db:"id"`
	DomainID uuid.UUID  `db:"domain_id"`
	RoleID   uuid.UUID  `db:"role_id"`
	UserID   uuid.UUID  `db:"user_id"`
	Entity   string     `db:"entity"`
	Allow    byte       `db:"allow"`
	Created  time.Time  `db:"created"`
	Updated  time.Time  `db:"updated"`
	Deleted  *time.Time `db:"deleted"`
}
