package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1768000002_create_patient_conditions_table",
		Up: func() error {
			return olympian.Table("patient_conditions").Create(func() {
				olympian.String("id").Primary()
				olympian.String("patient_id")
				olympian.String("condition_name")
				olympian.Text("notes").Nullable()
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("patient_conditions").Drop()
		},
	})
}
