package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1776000001_create_sale_items_table",
		Up: func() error {
			return olympian.Table("sale_items").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("sale_id")
				olympian.String("product_id")
				olympian.Integer("quantity")
				olympian.Float("unit_price")
				olympian.Float("discount").Default("0")
				olympian.Float("line_total")
			})
		},
		Down: func() error {
			return olympian.Table("sale_items").Drop()
		},
	})
}
