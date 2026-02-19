package builder_test

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test/model"
	"github.com/pshvedko/dbx/request"
)

func ExampleNewPermanent() {
	o := model.Object{}
	q := filter.Eq{"string_1": "red", "time_4": nil}

	p, err := builder.NewPermanent(&o, q)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(p)
	fmt.Println(filter.IsEmpty(p))

	for _, f := range []filter.Filter{
		filter.And{filter.Eq{"id": uuid.UUID{}}, p},
		filter.And{filter.Eq{"id": uuid.UUID{}}, p, filter.In{"string_2": []any{"green", "yellow"}}},
		filter.And{filter.Eq{"id": uuid.UUID{}}, filter.Or{filter.In{"string_2": []any{"green", "yellow"}}, p}},
		filter.And{p, p},
	} {
		b := builder.NewBuilder(110)
		err = f.To(b, &o)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(b)

		for i, v := range b.Values() {
			fmt.Print("$", i+1, " = ", v, "\n")
		}
	}

	r, err := request.NewWithOption([]request.Option{request.WithField{"id"}})
	if err != nil {
		fmt.Println(err)
		return
	}

	_, s, a, _, err := r.Constructor().Select(&o, p)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(s)
	fmt.Println(a)

	// Output:
	//
	// ( "objects"."string_1" = $1 AND "objects"."time_4" IS NULL )
	// false
	// ( "objects"."id" = $1 AND ( "objects"."string_1" = $2 AND "objects"."time_4" IS NULL ) )
	// $1 = 00000000-0000-0000-0000-000000000000
	// $2 = red
	// ( "objects"."id" = $1 AND ( "objects"."string_1" = $2 AND "objects"."time_4" IS NULL ) AND "objects"."string_2" = ANY($3) )
	// $1 = 00000000-0000-0000-0000-000000000000
	// $2 = red
	// $3 = [green yellow]
	// ( "objects"."id" = $1 AND ( "objects"."string_2" = ANY($2) OR ( "objects"."string_1" = $3 AND "objects"."time_4" IS NULL ) ) )
	// $1 = 00000000-0000-0000-0000-000000000000
	// $2 = [green yellow]
	// $3 = red
	// ( ( "objects"."string_1" = $1 AND "objects"."time_4" IS NULL ) AND ( "objects"."string_1" = $2 AND "objects"."time_4" IS NULL ) )
	// $1 = red
	// $2 = red
	// SELECT "o"."id" FROM "objects" AS "o" WHERE ( "objects"."string_1" = $1 AND "objects"."time_4" IS NULL )
	// [red]
}
