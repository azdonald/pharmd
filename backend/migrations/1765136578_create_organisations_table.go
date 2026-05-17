package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1765136578_create_organisations_table",
		Up: func() error {
			return olympian.Table("organisations").Create(func() {
				olympian.String("id").Primary()
				olympian.String("name").Unique()
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("organisations").Drop()
		},
	})
}
