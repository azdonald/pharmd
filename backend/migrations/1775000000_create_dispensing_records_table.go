package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1775000000_create_dispensing_records_table",
		Up: func() error {
			return olympian.Table("dispensing_records").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("prescription_id")
				olympian.String("prescription_item_id")
				olympian.String("patient_id")
				olympian.String("product_id")
				olympian.String("location_id")
				olympian.Integer("quantity_dispensed")
				olympian.Integer("quantity_prescribed")
				olympian.String("pharmacist_id")
				olympian.String("technician_id").Nullable()
				olympian.String("status").Default("completed")
				olympian.Text("notes").Nullable()
				olympian.String("witness_name").Nullable()
				olympian.Boolean("is_controlled").Default("false")
				olympian.String("dispensed_at")
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("dispensing_records").Drop()
		},
	})
}
