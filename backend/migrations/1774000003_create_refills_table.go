package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1774000003_create_refills_table",
		Up: func() error {
			return olympian.Table("refills").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("prescription_id")
				olympian.String("item_id")
				olympian.String("refilled_by")
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("refills").Drop()
		},
	})
}
