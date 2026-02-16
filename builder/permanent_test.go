package builder_test

import (
	"fmt"
	"github.com/pshvedko/dbx/internal/test/model"

	"github.com/google/uuid"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
)

func ExampleNewPermanent() {
	o := model.Object{}
	q := filter.Eq{"string_1": "red", "time_4": nil}

	p, err := builder.NewPermanent(q, &o)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(p)

	for _, f := range []filter.Filter{
		filter.And{filter.Eq{"id": uuid.UUID{}}, p},
		filter.And{filter.Eq{"id": uuid.UUID{}}, p, filter.In{"string_2": []any{"green", "yellow"}}},
		filter.And{filter.Eq{"id": uuid.UUID{}}, filter.Or{filter.In{"string_2": []any{"green", "yellow"}}, p}},
		filter.And{p, p},
	} {
		b := builder.NewBuilder()
		err = f.To(b, &o)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(b)

		for i, v := range b.Places() {
			fmt.Print("$", i+1, " = ", v, "\n")
		}
	}

	// Output:
	//
	// ( "objects"."string_1" = $1 AND "objects"."time_4" IS NULL )
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
}
