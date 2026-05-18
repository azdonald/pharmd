package repository

import (
	"context"
	"database/sql"

	"github.com/azdonald/pharmd/backend/models"
)

type PermissionRepoImpl struct {
	db *sql.DB
}

func NewPermissionRepositoryImpl(db *sql.DB) PermissionRepository {
	return &PermissionRepoImpl{db: db}
}

func (r *PermissionRepoImpl) ListPermissions(ctx context.Context) ([]models.Permission, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, slug, description, created_at, updated_at
		 FROM permissions WHERE deleted_at IS NULL ORDER BY id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}
