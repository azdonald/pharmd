package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1772000000_create_purchase_orders_table",
		Up: func() error {
			return olympian.Table("purchase_orders").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("po_number")
				olympian.String("supplier_id")
				olympian.String("location_id")
				olympian.String("status").Default("draft")
				olympian.String("order_date")
				olympian.String("expected_date").Nullable()
				olympian.Text("notes").Nullable()
				olympian.Float("subtotal").Default("0")
				olympian.Float("tax_total").Default("0")
				olympian.Float("grand_total").Default("0")
				olympian.String("created_by")
				olympian.String("approved_by").Nullable()
				olympian.String("approved_at").Nullable()
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("purchase_orders").Drop()
		},
	})
}
