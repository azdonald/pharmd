package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1769000000_create_product_categories_table",
		Up: func() error {
			return olympian.Table("product_categories").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("name")
				olympian.Text("description").Nullable()
				olympian.String("parent_id").Nullable()
				olympian.Boolean("is_active").Default("true")
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("product_categories").Drop()
		},
	})
}
