package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1776000000_create_sales_table",
		Up: func() error {
			return olympian.Table("sales").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("location_id")
				olympian.String("patient_id").Nullable()
				olympian.String("prescription_id").Nullable()
				olympian.String("sale_number")
				olympian.String("sale_type").Default("otc")
				olympian.String("status").Default("completed")
				olympian.Float("subtotal").Default("0")
				olympian.Float("tax_total").Default("0")
				olympian.Float("discount_total").Default("0")
				olympian.Float("grand_total").Default("0")
				olympian.Float("paid_amount").Default("0")
				olympian.Float("change_amount").Default("0")
				olympian.Text("notes").Nullable()
				olympian.String("created_by")
				olympian.String("voided_by").Nullable()
				olympian.String("voided_at").Nullable()
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("sales").Drop()
		},
	})
}
