package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1769000001_create_products_table",
		Up: func() error {
			return olympian.Table("products").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("name")
				olympian.Text("description").Nullable()
				olympian.String("category_id").Nullable()
				olympian.String("classification").Nullable()
				olympian.String("brand_name").Nullable()
				olympian.String("generic_name").Nullable()
				olympian.String("manufacturer").Nullable()
				olympian.String("barcode").Nullable()
				olympian.String("ndc").Nullable()
				olympian.String("unit_of_measure").Nullable()
				olympian.String("strength").Nullable()
				olympian.String("form").Nullable()
				olympian.Integer("reorder_level").Default("10")
				olympian.Boolean("is_active").Default("true")
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("products").Drop()
		},
	})
}
