package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/oklog/ulid/v2"
)

type RoleRepoImpl struct {
	db *sql.DB
}

func NewRoleRepositoryImpl(db *sql.DB) RoleRepository {
	return &RoleRepoImpl{db: db}
}

func (r *RoleRepoImpl) ListRoles(ctx context.Context, page, limit int) ([]models.Role, error) {
	orgID := ctx.Value("organisation_id").(string)
	offset := (page - 1) * limit

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, slug, description, organisation_id, is_system, created_at, updated_at
		 FROM roles
		 WHERE deleted_at IS NULL
		   AND (organisation_id = ?
		     OR (is_system = true AND organisation_id IS NULL
		         AND NOT EXISTS (SELECT 1 FROM roles sub WHERE sub.slug = roles.slug AND sub.organisation_id = ?)))
		 ORDER BY name ASC LIMIT ? OFFSET ?`,
		orgID, orgID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Slug, &role.Description,
			&role.OrganisationID, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepoImpl) GetRoleByID(ctx context.Context, id string) (*models.Role, error) {
	orgID := ctx.Value("organisation_id").(string)
	role := &models.Role{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, slug, description, organisation_id, is_system, created_at, updated_at
		 FROM roles WHERE id = ? AND (organisation_id = ? OR is_system = true) AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&role.ID, &role.Name, &role.Slug, &role.Description,
		&role.OrganisationID, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepoImpl) CreateRole(ctx context.Context, role models.Role) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO roles (id, name, slug, description, organisation_id, is_system, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		role.ID, role.Name, role.Slug, role.Description, role.OrganisationID,
		false, role.CreatedAt, role.UpdatedAt,
	)
	return err
}

func (r *RoleRepoImpl) UpdateRole(ctx context.Context, id string, role models.Role) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE roles SET name = ?, description = ?, updated_at = ? WHERE id = ? AND is_system = false AND deleted_at IS NULL`,
		role.Name, role.Description, time.Now(), id,
	)
	return err
}

func (r *RoleRepoImpl) DeleteRole(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE roles SET deleted_at = NOW() WHERE id = ? AND is_system = false AND deleted_at IS NULL`, id,
	)
	return err
}

func (r *RoleRepoImpl) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT permission_id FROM role_permissions WHERE role_id = ? AND deleted_at IS NULL`, roleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *RoleRepoImpl) CloneSystemRolesForOrg(ctx context.Context, orgID string) (map[string]string, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, slug, description FROM roles WHERE is_system = true AND organisation_id IS NULL AND deleted_at IS NULL`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type systemRole struct {
		ID, Name, Slug, Desc string
	}
	var systemRoles []systemRole
	for rows.Next() {
		var sr systemRole
		if err := rows.Scan(&sr.ID, &sr.Name, &sr.Slug, &sr.Desc); err != nil {
			return nil, err
		}
		systemRoles = append(systemRoles, sr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()
	mapping := make(map[string]string, len(systemRoles))
	for _, sr := range systemRoles {
		newID := ulid.Make().String()
		mapping[sr.ID] = newID
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO roles (id, name, slug, description, organisation_id, is_system, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, false, ?, ?)`,
			newID, sr.Name, sr.Slug, sr.Desc, orgID, now, now,
		); err != nil {
			return nil, fmt.Errorf("clone role %s: %w", sr.Name, err)
		}

		permRows, err := r.db.QueryContext(ctx,
			`SELECT permission_id FROM role_permissions WHERE role_id = ? AND deleted_at IS NULL`, sr.ID,
		)
		if err != nil {
			return nil, err
		}
		var permIDs []string
		for permRows.Next() {
			var pid string
			if err := permRows.Scan(&pid); err != nil {
				permRows.Close()
				return nil, err
			}
			permIDs = append(permIDs, pid)
		}
		permRows.Close()

		for _, pid := range permIDs {
			rpID := ulid.Make().String()
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO role_permissions (id, role_id, permission_id, created_at, updated_at)
				 VALUES (?, ?, ?, ?, ?)`,
				rpID, newID, pid, now, now,
			); err != nil {
				return nil, fmt.Errorf("clone permission %s for role %s: %w", pid, sr.Name, err)
			}
		}
	}

	return mapping, tx.Commit()
}

func (r *RoleRepoImpl) SetRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		`UPDATE role_permissions SET deleted_at = NOW() WHERE role_id = ? AND deleted_at IS NULL`, roleID,
	); err != nil {
		return err
	}

	now := time.Now()
	for _, permID := range permissionIDs {
		id := ulid.Make().String()
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO role_permissions (id, role_id, permission_id, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?)`,
			id, roleID, permID, now, now,
		); err != nil {
			return fmt.Errorf("failed to assign permission %s: %w", permID, err)
		}
	}

	return tx.Commit()
}
