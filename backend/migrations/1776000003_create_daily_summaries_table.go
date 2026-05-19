package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1776000003_create_daily_summaries_table",
		Up: func() error {
			return olympian.Table("daily_summaries").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("location_id")
				olympian.String("date")
				olympian.Integer("total_sales").Default("0")
				olympian.Float("total_revenue").Default("0")
				olympian.Float("total_tax").Default("0")
				olympian.Float("total_discounts").Default("0")
				olympian.Boolean("is_closed").Default("false")
				olympian.String("closed_by").Nullable()
				olympian.String("closed_at").Nullable()
				olympian.Text("notes").Nullable()
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("daily_summaries").Drop()
		},
	})
}
