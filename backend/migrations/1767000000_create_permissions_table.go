package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1767000000_create_permissions_table",
		Up: func() error {
			return olympian.Table("permissions").Create(func() {
				olympian.String("id").Primary()
				olympian.String("name")
				olympian.String("slug").Unique()
				olympian.String("description").Nullable()
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("permissions").Drop()
		},
	})
}
