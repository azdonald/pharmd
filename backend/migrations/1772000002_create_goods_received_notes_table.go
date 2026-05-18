package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1772000002_create_goods_received_notes_table",
		Up: func() error {
			return olympian.Table("goods_received_notes").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("purchase_order_id")
				olympian.String("received_date")
				olympian.Text("notes").Nullable()
				olympian.String("created_by")
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("goods_received_notes").Drop()
		},
	})
}
