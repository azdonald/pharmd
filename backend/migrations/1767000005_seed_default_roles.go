package migrations

import (
	"fmt"

	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1767000005_seed_default_roles",
		Up: func() error {
			db, _ := olympian.GetDB()

			roles := []struct {
				ID, Name, Slug, Desc string
				IsSystem             bool
				Permissions          []string
			}{
				{
					ID: "R001", Name: "Admin", Slug: "admin",
					Desc: "Full access to all features", IsSystem: true,
					Permissions: []string{
						"P001", "P002", "P003", "P004", "P005", "P006", "P007", "P008",
						"P009", "P010", "P011", "P012", "P013", "P014", "P015", "P016",
						"P017", "P018", "P019", "P020", "P021", "P022", "P023", "P024",
						"P025", "P026", "P027", "P028", "P029", "P030", "P031", "P032",
						"P033", "P034", "P035", "P036", "P037", "P038", "P039", "P040",
					},
				},
				{
					ID: "R002", Name: "Pharmacist-in-Charge", Slug: "pharmacist_in_charge",
					Desc: "Full pharmacy operations", IsSystem: true,
					Permissions: []string{
						"P006", "P008", "P010", "P012", "P014",
						"P017", "P018", "P019", "P021", "P022", "P023",
						"P025", "P026", "P027", "P028", "P029",
						"P031", "P033", "P034", "P035", "P036",
						"P037", "P039", "P040",
					},
				},
				{
					ID: "R003", Name: "Pharmacist", Slug: "pharmacist",
					Desc: "Dispensing and patient care", IsSystem: true,
					Permissions: []string{
						"P010", "P014",
						"P017", "P018", "P019",
						"P022", "P023", "P026",
						"P028",
						"P033", "P034", "P035",
						"P037", "P039",
					},
				},
				{
					ID: "R004", Name: "Technician", Slug: "technician",
					Desc: "Inventory and stock support", IsSystem: true,
					Permissions: []string{
						"P010", "P018", "P022", "P025", "P026",
						"P028", "P029", "P031", "P033",
					},
				},
				{
					ID: "R005", Name: "Cashier", Slug: "cashier",
					Desc: "Front-desk sales only", IsSystem: true,
					Permissions: []string{
						"P010", "P018", "P022", "P026", "P037",
					},
				},
			}

			for _, role := range roles {
				isSystem := 0
				if role.IsSystem {
					isSystem = 1
				}
				_, err := db.Exec(
					`INSERT INTO roles (id, name, slug, description, is_system, created_at, updated_at)
					 VALUES (?, ?, ?, ?, ?, NOW(), NOW())
					 ON DUPLICATE KEY UPDATE name=name`,
					role.ID, role.Name, role.Slug, role.Desc, isSystem,
				)
				if err != nil {
					return fmt.Errorf("failed to insert role %s: %w", role.ID, err)
				}

				for _, permID := range role.Permissions {
					rpID := role.ID + "_" + permID
					_, err := db.Exec(
						`INSERT INTO role_permissions (id, role_id, permission_id, created_at, updated_at)
						 VALUES (?, ?, ?, NOW(), NOW())
						 ON DUPLICATE KEY UPDATE id=id`,
						rpID, role.ID, permID,
					)
					if err != nil {
						return fmt.Errorf("failed to assign permission %s to role %s: %w", permID, role.ID, err)
					}
				}
			}
			return nil
		},
		Down: func() error {
			db, _ := olympian.GetDB()
			if _, err := db.Exec("DELETE FROM role_permissions WHERE role_id LIKE 'R%'"); err != nil {
				return err
			}
			_, err := db.Exec("DELETE FROM roles WHERE id LIKE 'R%'")
			return err
		},
	})
}
