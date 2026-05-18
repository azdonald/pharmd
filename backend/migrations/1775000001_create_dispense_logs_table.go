package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1775000001_create_dispense_logs_table",
		Up: func() error {
			return olympian.Table("dispense_logs").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("dispense_id")
				olympian.String("action")
				olympian.String("actor_id")
				olympian.Text("notes").Nullable()
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("dispense_logs").Drop()
		},
	})
}
