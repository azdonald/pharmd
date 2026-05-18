package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/azdonald/pharmd/backend/models"
)

type ProductCategoryRepoImpl struct {
	db *sql.DB
}

func NewProductCategoryRepositoryImpl(db *sql.DB) ProductCategoryRepository {
	return &ProductCategoryRepoImpl{db: db}
}

func (r *ProductCategoryRepoImpl) ListCategories(ctx context.Context) ([]models.ProductCategory, error) {
	orgID := ctx.Value("organisation_id").(string)
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, organisation_id, name, description, parent_id, is_active, created_at, updated_at
		 FROM product_categories WHERE organisation_id = ? AND deleted_at IS NULL
		 ORDER BY name ASC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.ProductCategory
	for rows.Next() {
		var c models.ProductCategory
		var parentID sql.NullString
		if err := rows.Scan(&c.ID, &c.OrganisationID, &c.Name, &c.Description, &parentID,
			&c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.ParentID = parentID.String
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *ProductCategoryRepoImpl) GetCategoryByID(ctx context.Context, id string) (*models.ProductCategory, error) {
	orgID := ctx.Value("organisation_id").(string)
	c := &models.ProductCategory{}
	var parentID sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, organisation_id, name, description, parent_id, is_active, created_at, updated_at
		 FROM product_categories WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&c.ID, &c.OrganisationID, &c.Name, &c.Description, &parentID, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	c.ParentID = parentID.String
	return c, nil
}

func (r *ProductCategoryRepoImpl) CreateCategory(ctx context.Context, category models.ProductCategory) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO product_categories (id, organisation_id, name, description, parent_id, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, NULLIF(?, ''), ?, ?, ?)`,
		category.ID, category.OrganisationID, category.Name, category.Description,
		category.ParentID, category.IsActive, category.CreatedAt, category.UpdatedAt,
	)
	return err
}

func (r *ProductCategoryRepoImpl) UpdateCategory(ctx context.Context, id string, category models.ProductCategory) error {
	orgID := ctx.Value("organisation_id").(string)
	query := "UPDATE product_categories SET updated_at = ?"
	args := []interface{}{time.Now()}

	if category.Name != "" {
		query += ", name = ?"
		args = append(args, category.Name)
	}
	if category.Description != "" {
		query += ", description = ?"
		args = append(args, category.Description)
	}
	if category.ParentID != "" {
		query += ", parent_id = ?"
		args = append(args, category.ParentID)
	}
	if !category.IsActive {
		query += ", is_active = ?"
		args = append(args, category.IsActive)
	}

	query += fmt.Sprintf(" WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL")
	args = append(args, id, orgID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *ProductCategoryRepoImpl) DeleteCategory(ctx context.Context, id string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		"UPDATE product_categories SET deleted_at = NOW() WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL",
		id, orgID,
	)
	return err
}

var _ ProductCategoryRepository = (*ProductCategoryRepoImpl)(nil)
