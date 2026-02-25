package builder_test

import (
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pshvedko/dbx/builder"
	"github.com/pshvedko/dbx/filter"
	"github.com/pshvedko/dbx/internal/test/model"
	"github.com/pshvedko/dbx/request"
)

func TestConstructor_Select_Join(t *testing.T) {
	a := model.Object{ID: uuid.UUID{1}}
	b := model.Object{ID: uuid.UUID{2}}
	c := model.Object{ID: uuid.UUID{3}}
	type args struct {
		j filter.Projector
		f filter.Filter
		o []request.Option
		b *uint
		l *uint
		y builder.Order
	}
	tests := []struct {
		name     string
		args     args
		want     string
		want1    []any
		want2    []any
		want3    string
		want4    int
		wantErr  error
		wantErr1 error
		wantErr2 error
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				j: filter.NewJoin(&a, &b, filter.Eq{"id": filter.Column("id")}),
				f: filter.Eq{},
				o: []request.Option{
					request.WithField{"id"},
					request.WithDeleted("time_4"),
				},
			},
			want:  `SELECT "o"."id", "o1"."id" FROM "objects" AS "o" JOIN "objects" AS "o1" ON ( "o1"."id" = "o"."id" AND "o1"."time_4" IS NULL ) WHERE ( "o"."time_4" IS NULL )`,
			want1: []any{},
			want2: []any{&a.ID, &b.ID},
		},
		{
			name: "",
			args: args{
				j: filter.NewJoin(&a, filter.NewJoin(&b, &c, filter.Eq{"id": filter.Column("id")}), filter.Eq{"id": filter.Column("id")}),
				f: filter.Eq{},
				o: []request.Option{
					request.WithField{"id"},
					request.WithDeleted("time_4"),
				},
			},
			want:  `SELECT "o"."id", "o1"."id", "o2"."id" FROM "objects" AS "o" JOIN "objects" AS "o1" ON ( "o1"."id" = "o"."id" AND "o1"."time_4" IS NULL ) JOIN "objects" AS "o2" ON ( "o2"."id" = "o1"."id" AND "o2"."time_4" IS NULL ) WHERE ( "o"."time_4" IS NULL )`,
			want1: []any{},
			want2: []any{&a.ID, &b.ID, &c.ID},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := request.NewWithOption(tt.args.o)
			require.ErrorIs(t, err, tt.wantErr)
			z, got, got1, got2, err := r.Constructor().Range(tt.args.b, tt.args.l).Sort(tt.args.y).Select(tt.args.j, tt.args.f)
			t.Log(got)
			require.Equal(t, err, tt.wantErr1)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want1, got1)
			require.Equal(t, tt.want2, got2)
			for i, v := range got2 {
				require.Same(t, tt.want2[i], v, i)
			}
			require.Equal(t, z == nil, tt.want3 == "")
			if z == nil {
				return
			}
			got3, got4, err := z.Count()
			t.Log(got3)
			require.Equal(t, err, tt.wantErr2)
			require.Equal(t, tt.want3, got3)
			require.Equal(t, tt.want4, got4)
		})
	}
}
