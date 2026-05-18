package migrations

import (
	"fmt"

	"github.com/ichtrojan/olympian"
)

func init() {
	olympian.RegisterMigration(olympian.Migration{
		Name: "1767000004_seed_permissions",
		Up: func() error {
			db, _ := olympian.GetDB()

			perms := []struct {
				ID, Name, Slug, Desc string
			}{
				{"P001", "Create Organisation", "organisation:create", "Create a new organisation"},
				{"P002", "Read Organisation", "organisation:read", "View organisation details"},
				{"P003", "Update Organisation", "organisation:update", "Update organisation settings"},
				{"P004", "Delete Organisation", "organisation:delete", "Delete organisation"},
				{"P005", "Create Users", "users:create", "Create new users"},
				{"P006", "Read Users", "users:read", "View user list and details"},
				{"P007", "Update Users", "users:update", "Update user profiles"},
				{"P008", "Delete Users", "users:delete", "Deactivate/delete users"},
				{"P009", "Create Locations", "locations:create", "Create new locations"},
				{"P010", "Read Locations", "locations:read", "View location list and details"},
				{"P011", "Update Locations", "locations:update", "Update location settings"},
				{"P012", "Delete Locations", "locations:delete", "Delete locations"},
				{"P013", "Create Roles", "roles:create", "Create new roles"},
				{"P014", "Read Roles", "roles:read", "View role list and details"},
				{"P015", "Update Roles", "roles:update", "Update role permissions"},
				{"P016", "Delete Roles", "roles:delete", "Delete roles"},
				{"P017", "Create Patients", "patients:create", "Register new patients"},
				{"P018", "Read Patients", "patients:read", "View patient records"},
				{"P019", "Update Patients", "patients:update", "Update patient information"},
				{"P020", "Delete Patients", "patients:delete", "Delete patient records"},
				{"P021", "Create Products", "products:create", "Add new products"},
				{"P022", "Read Products", "products:read", "View product catalog"},
				{"P023", "Update Products", "products:update", "Update product details"},
				{"P024", "Delete Products", "products:delete", "Remove products"},
				{"P025", "Manage Inventory", "inventory:manage", "Receive, adjust, transfer stock"},
				{"P026", "Read Inventory", "inventory:read", "View stock levels and movements"},
				{"P027", "Create Suppliers", "suppliers:create", "Add new suppliers"},
				{"P028", "Read Suppliers", "suppliers:read", "View supplier list"},
				{"P029", "Update Suppliers", "suppliers:update", "Update supplier details"},
				{"P030", "Delete Suppliers", "suppliers:delete", "Remove suppliers"},
				{"P031", "Create Purchase Orders", "purchases:create", "Create purchase orders"},
				{"P032", "Approve Purchase Orders", "purchases:approve", "Approve purchase orders"},
				{"P033", "Receive Goods", "purchases:receive", "Record goods received"},
				{"P034", "Manage Prescriptions", "prescriptions:manage", "Create and manage prescriptions"},
				{"P035", "Dispense Medications", "dispensing:dispense", "Dispense medications"},
				{"P036", "Override Pricing", "pricing:override", "Override sale prices"},
				{"P037", "Process Sales", "pos:sale", "Process POS transactions"},
				{"P038", "Void Sales", "pos:void", "Void or refund transactions"},
				{"P039", "View Reports", "reports:read", "Access reports and analytics"},
				{"P040", "Manage Settings", "settings:manage", "Update system settings"},
			}

			for _, p := range perms {
				query := `INSERT INTO permissions (id, name, slug, description, created_at, updated_at)
						  VALUES (?, ?, ?, ?, NOW(), NOW())
						  ON DUPLICATE KEY UPDATE name=name`
				_, err := db.Exec(query, p.ID, p.Name, p.Slug, p.Desc)
				if err != nil {
					return fmt.Errorf("failed to insert permission %s: %w", p.ID, err)
				}
			}
			return nil
		},
		Down: func() error {
			db, _ := olympian.GetDB()
			if _, err := db.Exec("DELETE FROM permissions WHERE id LIKE 'P%'"); err != nil {
				return err
			}
			return nil
		},
	})
}
