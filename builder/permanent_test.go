package builder_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/t"
)

func ExampleNewPermanent() {
	o := help.Object{}
	q := filter.Eq{"o_string_1": "red", "o_time_4": nil}

	p, err := builder.NewPermanent(q, &o)
	if err != nil {
		return
	}

	fmt.Println(p)

	for _, f := range []filter.Filter{
		filter.And{filter.Eq{"id": uuid.UUID{}}, p},
		filter.And{filter.Eq{"id": uuid.UUID{}}, p, filter.In{"o_string_2": []any{"green", "yellow"}}},
		filter.And{filter.Eq{"id": uuid.UUID{}}, filter.Or{filter.In{"o_string_2": []any{"green", "yellow"}}, p}},
		filter.And{p, p},
	} {
		b := builder.NewBuilder()
		err = f.To(b, &o)
		if err != nil {
			return
		}

		fmt.Println(b)

		for i, v := range b.Values() {
			fmt.Print("$", i+1, " = ", v, "\n")
		}
	}

	// Output:
	//
	// ( "objects"."o_string_1" = $1 AND "objects"."o_time_4" IS NULL )
	// ( "objects"."id" = $1 AND ( "objects"."o_string_1" = $2 AND "objects"."o_time_4" IS NULL ) )
	// $1 = 00000000-0000-0000-0000-000000000000
	// $2 = red
	// ( "objects"."id" = $1 AND ( "objects"."o_string_1" = $2 AND "objects"."o_time_4" IS NULL ) AND "objects"."o_string_2" = ANY($3) )
	// $1 = 00000000-0000-0000-0000-000000000000
	// $2 = red
	// $3 = [green yellow]
	// ( "objects"."id" = $1 AND ( "objects"."o_string_2" = ANY($2) OR ( "objects"."o_string_1" = $3 AND "objects"."o_time_4" IS NULL ) ) )
	// $1 = 00000000-0000-0000-0000-000000000000
	// $2 = [green yellow]
	// $3 = red
	// ( ( "objects"."o_string_1" = $1 AND "objects"."o_time_4" IS NULL ) AND ( "objects"."o_string_1" = $2 AND "objects"."o_time_4" IS NULL ) )
	// $1 = red
	// $2 = red
}
