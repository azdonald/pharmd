package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1768000001_create_patient_allergies_table",
		Up: func() error {
			return olympian.Table("patient_allergies").Create(func() {
				olympian.String("id").Primary()
				olympian.String("patient_id")
				olympian.String("allergy")
				olympian.String("severity").Nullable()
				olympian.Text("notes").Nullable()
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("patient_allergies").Drop()
		},
	})
}
