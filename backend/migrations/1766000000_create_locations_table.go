package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1766000000_create_locations_table",
		Up: func() error {
			return olympian.Table("locations").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("name")
				olympian.Text("address").Nullable()
				olympian.String("city").Nullable()
				olympian.String("state").Nullable()
				olympian.String("country").Nullable()
				olympian.String("phone").Nullable()
				olympian.String("email").Nullable()
				olympian.Float("tax_rate").Default("0")
				olympian.String("timezone").Default("UTC")
				olympian.Boolean("is_active").Default("true")
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("locations").Drop()
		},
	})
}
