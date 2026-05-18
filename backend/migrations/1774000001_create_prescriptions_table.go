package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1774000001_create_prescriptions_table",
		Up: func() error {
			return olympian.Table("prescriptions").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("patient_id")
				olympian.String("prescriber_id")
				olympian.String("location_id")
				olympian.String("status").Default("active")
				olympian.Text("diagnosis").Nullable()
				olympian.Text("notes").Nullable()
				olympian.String("issued_date").Nullable()
				olympian.String("expiry_date").Nullable()
				olympian.String("created_by")
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("prescriptions").Drop()
		},
	})
}
