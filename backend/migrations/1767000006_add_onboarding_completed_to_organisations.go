package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1767000006_add_onboarding_completed_to_organisations",
		Up: func() error {
			return olympian.Table("organisations").Modify(func() {
				olympian.Boolean("onboarding_completed").Default("false")
			})
		},
		Down: func() error {
			return olympian.Table("organisations").DropColumn("onboarding_completed")
		},
	})
}
