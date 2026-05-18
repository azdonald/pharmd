package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1770000000_create_stock_batches_table",
		Up: func() error {
			return olympian.Table("stock_batches").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("product_id")
				olympian.String("location_id")
				olympian.String("batch_number").Nullable()
				olympian.Integer("quantity")
				olympian.Integer("remaining_qty")
				olympian.Float("unit_cost").Default("0")
				olympian.Float("selling_price").Default("0")
				olympian.String("manufacturing_date").Nullable()
				olympian.String("expiry_date").Nullable()
				olympian.String("received_date")
				olympian.Boolean("is_active").Default("true")
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("stock_batches").Drop()
		},
	})
}
