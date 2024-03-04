package db

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
	"github.com/pshvedko/db/filter"
)

type DB struct {
	*sqlx.DB
}

var database DB

func Open(name string) error {
	db, err := sqlx.Open("pgx", name)
	if err == nil {
		database.DB = db
	}
	return err
}

func Close() {
	if database.DB != nil {
		_ = database.DB.Close()
	}
}

type Object interface {
	Fields() []any
}

type Performer interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Connection struct {
	Performer
	error
}

//func (db *DB) GetConnection(ctx context.Context) (*Connection, error) {
//
//}

func Get(ctx context.Context, o Object, f filter.Filter) error {
	return database.Get(ctx, o, f)
}

func (db *DB) Get(ctx context.Context, o Object, f filter.Filter) error {
	//c, err := db.GetConnection(ctx)
	//if err != nil {
	//	return err
	//}
	//defer c.Close()
	////q, a, err := makeSelectQuery(o, f)
	////
	////fmt.Println(q)
	////fmt.Println(a)
	////
	////tx, err := db.BeginTxx(ctx, nil)
	////tx.GetContext()
	////tx.Close()
	//if err != nil {
	//	return err
	//}
	//return c.GetContext(ctx, o, q, a...)
	return nil
}

//func makeSelectQuery(o Object, f filter.Filter) (string, []any, error) {
//	w, a, err := buildFilter(f)
//	if err != nil {
//		return "", nil, err
//	}
//	return "SELECT * FROM files" + w, a, err
//}

//func buildFilter(f filter.Filter, a ...any) (string, []any, error) {
//	b := Builder{v: a}
//	err := f.To(&b)
//	if err != nil || b.b.Len() == 0 {
//		return "", nil, err
//	}
//	return " WHERE " + b.b.String(), b.v, nil
//}
