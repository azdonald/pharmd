package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1769000002_create_generic_substitutions_table",
		Up: func() error {
			return olympian.Table("generic_substitutions").Create(func() {
				olympian.String("id").Primary()
				olympian.String("organisation_id")
				olympian.String("product_id")
				olympian.String("substitute_product_id")
				olympian.Text("notes").Nullable()
				olympian.Timestamps()
			})
		},
		Down: func() error {
			return olympian.Table("generic_substitutions").Drop()
		},
	})
}
