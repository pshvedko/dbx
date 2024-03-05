package db

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func Open(name string) (*DB, error) {
	db, err := sqlx.Open("pgx", name)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}
