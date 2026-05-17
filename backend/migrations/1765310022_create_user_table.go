package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1765310022_create_user_table",
		Up: func() error {
			return olympian.Table("users").Create(func() {
				olympian.String("id").Primary()
				olympian.String("first_name")
				olympian.String("last_name")
				olympian.String("email").Unique()
				olympian.String("password")
				olympian.String("organisation_id")
				olympian.Boolean("is_active").Default("true")
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("users").Drop()
		},
	})
}
