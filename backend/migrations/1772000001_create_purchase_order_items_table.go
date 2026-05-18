package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1772000001_create_purchase_order_items_table",
		Up: func() error {
			return olympian.Table("purchase_order_items").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("purchase_order_id")
				olympian.String("product_id")
				olympian.Integer("quantity_ordered")
				olympian.Integer("quantity_received").Default("0")
				olympian.Float("unit_cost").Default("0")
				olympian.Float("line_total").Default("0")
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("purchase_order_items").Drop()
		},
	})
}
