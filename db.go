package db

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pshvedko/db/filter"
)

var db *sqlx.DB

func Open(name string) (err error) {
	db, err = sqlx.Open("pgx", name)
	return
}

func Close() error {
	return db.Close()
}

type Object interface {
	Fields() []any
}

func Get(ctx context.Context, o Object, f filter.Filter) error {
	c, err := db.Connx(ctx)
	if err != nil {
		return err
	}
	defer c.Close()
	q, a, err := makeSelectQuery(o, f)

	fmt.Println(q)
	fmt.Println(a)

	if err != nil {
		return err
	}
	return c.GetContext(ctx, o, q, a...)
}

func makeSelectQuery(o Object, f filter.Filter) (string, []any, error) {
	w, a, err := buildFilter(f)
	if err != nil {
		return "", nil, err
	}
	return "SELECT * FROM files" + w, a, err
}

func buildFilter(f filter.Filter, a ...any) (string, []any, error) {
	b := Builder{v: a}
	err := f.To(&b)
	if err != nil || b.b.Len() == 0 {
		return "", nil, err
	}
	return " WHERE " + b.b.String(), b.v, nil
}
