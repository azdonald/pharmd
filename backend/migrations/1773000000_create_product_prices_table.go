package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1773000000_create_product_prices_table",
		Up: func() error {
			return olympian.Table("product_prices").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("product_id")
				olympian.String("location_id")
				olympian.Float("selling_price").Default("0")
				olympian.Float("cost_price").Default("0")
				olympian.Float("min_price").Default("0")
				olympian.Float("max_discount").Default("0")
				olympian.Boolean("is_active").Default("true")
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("product_prices").Drop()
		},
	})
}
