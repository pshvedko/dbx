package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
)

type Permission struct {
	ID       uuid.UUID   `db:"id,primary"`
	DomainID uuid.UUID   `db:"domain_id"`
	RoleID   uuid.UUID   `db:"role_id"`
	UserID   uuid.UUID   `db:"user_id"`
	Entity   string      `db:"entity"`
	Allow    pgtype.Bits `db:"allow"`
	Created  time.Time   `db:"created"`
	Updated  time.Time   `db:"updated"`
	Deleted  *time.Time  `db:"deleted"`
}
