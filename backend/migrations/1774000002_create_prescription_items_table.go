package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1774000002_create_prescription_items_table",
		Up: func() error {
			return olympian.Table("prescription_items").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("prescription_id")
				olympian.String("product_id")
				olympian.String("dosage")
				olympian.String("frequency")
				olympian.String("duration").Nullable()
				olympian.Integer("quantity")
				olympian.Integer("refills_authorized").Default("0")
				olympian.Integer("refills_used").Default("0")
				olympian.Text("notes").Nullable()
			})
		},
		Down: func() error {
			return olympian.Table("prescription_items").Drop()
		},
	})
}
