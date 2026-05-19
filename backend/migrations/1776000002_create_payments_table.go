package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1776000002_create_payments_table",
		Up: func() error {
			return olympian.Table("payments").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("sale_id")
				olympian.String("method")
				olympian.Float("amount")
				olympian.String("reference").Nullable()
			})
		},
		Down: func() error {
			return olympian.Table("payments").Drop()
		},
	})
}
