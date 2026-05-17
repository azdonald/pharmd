package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1767000003_create_user_roles_table",
		Up: func() error {
			return olympian.Table("user_roles").Create(func() {
				olympian.String("id").Primary()
				olympian.String("user_id")
				olympian.String("role_id")
				olympian.String("organisation_id")
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("user_roles").Drop()
		},
	})
}
