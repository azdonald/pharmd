package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1773000001_create_discount_rules_table",
		Up: func() error {
			return olympian.Table("discount_rules").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("name")
				olympian.String("type")
				olympian.Float("value")
				olympian.Float("min_order_value").Default("0")
				olympian.Float("max_discount_amount").Default("0")
				olympian.String("applies_to").Default("all")
				olympian.String("applies_to_id").Nullable()
				olympian.Boolean("is_active").Default("true")
				olympian.String("start_date").Nullable()
				olympian.String("end_date").Nullable()
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("discount_rules").Drop()
		},
	})
}
