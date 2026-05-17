package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1767000001_create_roles_table",
		Up: func() error {
			return olympian.Table("roles").Create(func() {
				olympian.String("id").Primary()
				olympian.String("name")
				olympian.String("slug").Unique()
				olympian.String("description").Nullable()
				olympian.String("organisation_id")
				olympian.Boolean("is_system").Default("false")
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("roles").Drop()
		},
	})
}
