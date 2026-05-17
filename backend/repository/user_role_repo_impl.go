package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/oklog/ulid/v2"
)

type UserRoleRepoImpl struct {
	db *sql.DB
}

func NewUserRoleRepositoryImpl(db *sql.DB) UserRoleRepository {
	return &UserRoleRepoImpl{db: db}
}

func (r *UserRoleRepoImpl) GetUserPermissions(ctx context.Context, userID, orgID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DISTINCT p.slug
		 FROM permissions p
		 JOIN role_permissions rp ON rp.permission_id = p.id
		 JOIN user_roles ur ON ur.role_id = rp.role_id
		 WHERE ur.user_id = ? AND ur.organisation_id = ? AND ur.deleted_at IS NULL
		 AND rp.deleted_at IS NULL AND p.deleted_at IS NULL`,
		userID, orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	if permissions == nil {
		return []string{}, nil
	}

	return permissions, nil
}

func (r *UserRoleRepoImpl) AssignRoleToUser(ctx context.Context, userID, roleID, orgID string) error {
	id := ulid.Make().String()
	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_roles (id, user_id, role_id, organisation_id, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE deleted_at = NULL, updated_at = ?`,
		id, userID, roleID, orgID, now, now, now,
	)
	return err
}
