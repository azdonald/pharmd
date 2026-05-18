package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1774000000_create_prescribers_table",
		Up: func() error {
			return olympian.Table("prescribers").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("name")
				olympian.String("license_number").Nullable()
				olympian.String("phone").Nullable()
				olympian.String("email").Nullable()
				olympian.String("specialty").Nullable()
				olympian.String("dea_number").Nullable()
				olympian.String("npi_number").Nullable()
				olympian.Text("address").Nullable()
				olympian.Boolean("is_active").Default("true")
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("prescribers").Drop()
		},
	})
}
