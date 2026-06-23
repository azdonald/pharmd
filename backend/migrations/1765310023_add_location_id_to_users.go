package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1765310023_add_location_id_to_users",
		Up: func() error {
			return olympian.Table("users").Modify(func() {
				olympian.String("location_id").Nullable()
			})
		},
		Down: func() error {
			return olympian.Table("users").DropColumn("location_id")
		},
	})
}
