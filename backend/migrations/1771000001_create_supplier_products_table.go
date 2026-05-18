package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1771000001_create_supplier_products_table",
		Up: func() error {
			return olympian.Table("supplier_products").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("supplier_id")
				olympian.String("product_id")
				olympian.Float("unit_price").Default("0")
				olympian.Integer("min_order_qty").Default("0")
				olympian.Integer("lead_time_days").Default("0")
				olympian.Text("notes").Nullable()
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("supplier_products").Drop()
		},
	})
}
