package model

import (
	"context"
	"github.com/google/uuid"
	"github.com/pshvedko/db/request"
)

type User struct {
	ID uuid.UUID `json:"id"`
}

func GetUser(ctx context.Context) (*User, error) {
	r := request.New(ctx)
	var obj User
	err := r.Get(&obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}
