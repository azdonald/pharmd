package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1768000000_create_patients_table",
		Up: func() error {
			return olympian.Table("patients").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("first_name")
				olympian.String("last_name")
				olympian.String("date_of_birth").Nullable()
				olympian.String("gender").Nullable()
				olympian.String("phone").Nullable()
				olympian.String("email").Nullable()
				olympian.Text("address").Nullable()
				olympian.String("city").Nullable()
				olympian.String("state").Nullable()
				olympian.String("country").Nullable()
				olympian.String("blood_group").Nullable()
				olympian.String("genotype").Nullable()
				olympian.Text("notes").Nullable()
				olympian.String("emergency_contact_name").Nullable()
				olympian.String("emergency_contact_phone").Nullable()
				olympian.Boolean("is_active").Default("true")
				olympian.Timestamps()
				olympian.SoftDeletes()
			})
		},
		Down: func() error {
			return olympian.Table("patients").Drop()
		},
	})
}
