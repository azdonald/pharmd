package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1767000002_create_role_permissions_table",
		Up: func() error {
			return olympian.Table("role_permissions").Create(func() {
				olympian.String("id").Primary()
				olympian.String("role_id")
				olympian.String("permission_id")
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("role_permissions").Drop()
		},
	})
}
