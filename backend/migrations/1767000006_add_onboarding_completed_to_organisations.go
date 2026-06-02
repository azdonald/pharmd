package migrations

import (
	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1767000006_add_onboarding_completed_to_organisations",
		Up: func() error {
			return olympian.Raw("ALTER TABLE organisations ADD COLUMN onboarding_completed BOOLEAN NOT NULL DEFAULT FALSE")
		},
		Down: func() error {
			return olympian.Raw("ALTER TABLE organisations DROP COLUMN onboarding_completed")
		},
	})
}
