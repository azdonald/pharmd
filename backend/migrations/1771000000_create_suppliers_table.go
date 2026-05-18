package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1771000000_create_suppliers_table",
		Up: func() error {
			return olympian.Table("suppliers").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("name")
				olympian.String("contact_person").Nullable()
				olympian.String("phone").Nullable()
				olympian.String("email").Nullable()
				olympian.Text("address").Nullable()
				olympian.String("city").Nullable()
				olympian.String("state").Nullable()
				olympian.String("country").Nullable()
				olympian.String("payment_terms").Nullable()
				olympian.Text("notes").Nullable()
				olympian.Boolean("is_active").Default("true")
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("suppliers").Drop()
		},
	})
}
