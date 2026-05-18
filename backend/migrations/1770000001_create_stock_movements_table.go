package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1770000001_create_stock_movements_table",
		Up: func() error {
			return olympian.Table("stock_movements").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("location_id")
				olympian.String("product_id")
				olympian.String("batch_id").Nullable()
				olympian.String("movement_type")
				olympian.Integer("quantity")
				olympian.String("reference_type").Nullable()
				olympian.String("reference_id").Nullable()
				olympian.Text("notes").Nullable()
				olympian.String("created_by")
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("stock_movements").Drop()
		},
	})
}
