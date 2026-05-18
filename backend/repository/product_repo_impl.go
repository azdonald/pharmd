package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/azdonald/pharmd/backend/models"
)

type ProductRepoImpl struct {
	db *sql.DB
}

func NewProductRepositoryImpl(db *sql.DB) ProductRepository {
	return &ProductRepoImpl{db: db}
}

func (r *ProductRepoImpl) ListProducts(ctx context.Context, page, limit int, query, categoryID string) ([]models.Product, int, error) {
	orgID := ctx.Value("organisation_id").(string)
	offset := (page - 1) * limit

	where := "organisation_id = ? AND deleted_at IS NULL"
	args := []interface{}{orgID}

	if query != "" {
		where += " AND (name LIKE ? OR brand_name LIKE ? OR generic_name LIKE ? OR barcode LIKE ? OR ndc LIKE ?)"
		q := "%" + query + "%"
		args = append(args, q, q, q, q, q)
	}
	if categoryID != "" {
		where += " AND category_id = ?"
		args = append(args, categoryID)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM products WHERE " + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT id, organisation_id, name, description, category_id, classification,
		       brand_name, generic_name, manufacturer, barcode, ndc, unit_of_measure,
		       strength, form, reorder_level, is_active, created_at, updated_at
		 FROM products WHERE %s ORDER BY name ASC LIMIT ? OFFSET ?`, where),
		append(args, limit, offset)...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var catID sql.NullString
		if err := rows.Scan(&p.ID, &p.OrganisationID, &p.Name, &p.Description, &catID,
			&p.Classification, &p.BrandName, &p.GenericName, &p.Manufacturer,
			&p.Barcode, &p.NDC, &p.UnitOfMeasure, &p.Strength, &p.Form,
			&p.ReorderLevel, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		p.CategoryID = catID.String
		products = append(products, p)
	}
	return products, total, nil
}

func (r *ProductRepoImpl) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	orgID := ctx.Value("organisation_id").(string)
	p := &models.Product{}
	var catID sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, organisation_id, name, description, category_id, classification,
		        brand_name, generic_name, manufacturer, barcode, ndc, unit_of_measure,
		        strength, form, reorder_level, is_active, created_at, updated_at
		 FROM products WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&p.ID, &p.OrganisationID, &p.Name, &p.Description, &catID,
		&p.Classification, &p.BrandName, &p.GenericName, &p.Manufacturer,
		&p.Barcode, &p.NDC, &p.UnitOfMeasure, &p.Strength, &p.Form,
		&p.ReorderLevel, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	p.CategoryID = catID.String
	return p, nil
}

func (r *ProductRepoImpl) GetProductByBarcode(ctx context.Context, barcode string) (*models.Product, error) {
	orgID := ctx.Value("organisation_id").(string)
	p := &models.Product{}
	var catID sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, organisation_id, name, description, category_id, classification,
		        brand_name, generic_name, manufacturer, barcode, ndc, unit_of_measure,
		        strength, form, reorder_level, is_active, created_at, updated_at
		 FROM products WHERE barcode = ? AND organisation_id = ? AND deleted_at IS NULL
		 LIMIT 1`,
		barcode, orgID,
	).Scan(&p.ID, &p.OrganisationID, &p.Name, &p.Description, &catID,
		&p.Classification, &p.BrandName, &p.GenericName, &p.Manufacturer,
		&p.Barcode, &p.NDC, &p.UnitOfMeasure, &p.Strength, &p.Form,
		&p.ReorderLevel, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	p.CategoryID = catID.String
	return p, nil
}

func (r *ProductRepoImpl) CreateProduct(ctx context.Context, product models.Product) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO products (id, organisation_id, name, description, category_id, classification,
		                       brand_name, generic_name, manufacturer, barcode, ndc,
		                       unit_of_measure, strength, form, reorder_level,
		                       is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, NULLIF(?, ''), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		product.ID, product.OrganisationID, product.Name, product.Description,
		product.CategoryID, product.Classification, product.BrandName,
		product.GenericName, product.Manufacturer, product.Barcode, product.NDC,
		product.UnitOfMeasure, product.Strength, product.Form, product.ReorderLevel,
		product.IsActive, product.CreatedAt, product.UpdatedAt,
	)
	return err
}

func (r *ProductRepoImpl) UpdateProduct(ctx context.Context, id string, product models.Product) error {
	orgID := ctx.Value("organisation_id").(string)
	query := "UPDATE products SET updated_at = ?"
	args := []interface{}{time.Now()}

	if product.Name != "" {
		query += ", name = ?"
		args = append(args, product.Name)
	}
	if product.Description != "" {
		query += ", description = ?"
		args = append(args, product.Description)
	}
	if product.CategoryID != "" {
		query += ", category_id = ?"
		args = append(args, product.CategoryID)
	}
	if product.Classification != "" {
		query += ", classification = ?"
		args = append(args, product.Classification)
	}
	if product.BrandName != "" {
		query += ", brand_name = ?"
		args = append(args, product.BrandName)
	}
	if product.GenericName != "" {
		query += ", generic_name = ?"
		args = append(args, product.GenericName)
	}
	if product.Manufacturer != "" {
		query += ", manufacturer = ?"
		args = append(args, product.Manufacturer)
	}
	if product.Barcode != "" {
		query += ", barcode = ?"
		args = append(args, product.Barcode)
	}
	if product.NDC != "" {
		query += ", ndc = ?"
		args = append(args, product.NDC)
	}
	if product.UnitOfMeasure != "" {
		query += ", unit_of_measure = ?"
		args = append(args, product.UnitOfMeasure)
	}
	if product.Strength != "" {
		query += ", strength = ?"
		args = append(args, product.Strength)
	}
	if product.Form != "" {
		query += ", form = ?"
		args = append(args, product.Form)
	}
	if product.ReorderLevel != 0 {
		query += ", reorder_level = ?"
		args = append(args, product.ReorderLevel)
	}
	if !product.IsActive {
		query += ", is_active = ?"
		args = append(args, product.IsActive)
	}

	query += fmt.Sprintf(" WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL")
	args = append(args, id, orgID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *ProductRepoImpl) DeleteProduct(ctx context.Context, id string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		"UPDATE products SET deleted_at = NOW() WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL",
		id, orgID,
	)
	return err
}

func (r *ProductRepoImpl) ListSubstitutes(ctx context.Context, productID string) ([]models.GenericSubstitution, error) {
	orgID := ctx.Value("organisation_id").(string)
	rows, err := r.db.QueryContext(ctx,
		`SELECT s.id, s.product_id, s.substitute_product_id, p.name, s.notes, s.created_at
		 FROM generic_substitutions s
		 JOIN products p ON p.id = s.substitute_product_id
		 WHERE s.product_id = ? AND s.organisation_id = ? AND p.deleted_at IS NULL
		 ORDER BY p.name ASC`,
		productID, orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.GenericSubstitution
	for rows.Next() {
		var s models.GenericSubstitution
		var subName string
		if err := rows.Scan(&s.ID, &s.ProductID, &s.SubstituteProductID, &subName, &s.Notes, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}

func (r *ProductRepoImpl) AddSubstitute(ctx context.Context, sub models.GenericSubstitution) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO generic_substitutions (id, organisation_id, product_id, substitute_product_id, notes, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		sub.ID, sub.OrganisationID, sub.ProductID, sub.SubstituteProductID, sub.Notes, sub.CreatedAt, sub.CreatedAt,
	)
	return err
}

func (r *ProductRepoImpl) RemoveSubstitute(ctx context.Context, productID, substituteID string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM generic_substitutions
		 WHERE id = ? AND product_id = ? AND organisation_id = ?`,
		substituteID, productID, orgID,
	)
	return err
}

var _ ProductRepository = (*ProductRepoImpl)(nil)
