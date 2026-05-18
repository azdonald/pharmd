package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/azdonald/pharmd/backend/models"
)

type SupplierRepoImpl struct {
	db *sql.DB
}

func NewSupplierRepositoryImpl(db *sql.DB) SupplierRepository {
	return &SupplierRepoImpl{db: db}
}

func (r *SupplierRepoImpl) ListSuppliers(ctx context.Context, page, limit int, query string) ([]models.Supplier, int, error) {
	orgID := ctx.Value("organisation_id").(string)
	offset := (page - 1) * limit

	where := "organisation_id = ? AND deleted_at IS NULL"
	args := []interface{}{orgID}

	if query != "" {
		where += " AND (name LIKE ? OR phone LIKE ? OR email LIKE ?)"
		q := "%" + query + "%"
		args = append(args, q, q, q)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM suppliers WHERE " + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT id, organisation_id, name, contact_person, phone, email,
		       address, city, state, country, payment_terms, notes, is_active,
		       created_at, updated_at
		 FROM suppliers WHERE %s ORDER BY name ASC LIMIT ? OFFSET ?`, where),
		append(args, limit, offset)...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var suppliers []models.Supplier
	for rows.Next() {
		var s models.Supplier
		if err := rows.Scan(&s.ID, &s.OrganisationID, &s.Name, &s.ContactPerson,
			&s.Phone, &s.Email, &s.Address, &s.City, &s.State, &s.Country,
			&s.PaymentTerms, &s.Notes, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		suppliers = append(suppliers, s)
	}
	return suppliers, total, nil
}

func (r *SupplierRepoImpl) GetSupplierByID(ctx context.Context, id string) (*models.Supplier, error) {
	orgID := ctx.Value("organisation_id").(string)
	s := &models.Supplier{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, organisation_id, name, contact_person, phone, email,
		        address, city, state, country, payment_terms, notes, is_active,
		        created_at, updated_at
		 FROM suppliers WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&s.ID, &s.OrganisationID, &s.Name, &s.ContactPerson,
		&s.Phone, &s.Email, &s.Address, &s.City, &s.State, &s.Country,
		&s.PaymentTerms, &s.Notes, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SupplierRepoImpl) CreateSupplier(ctx context.Context, supplier models.Supplier) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO suppliers (id, organisation_id, name, contact_person, phone, email,
		                        address, city, state, country, payment_terms, notes,
		                        is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		supplier.ID, supplier.OrganisationID, supplier.Name, supplier.ContactPerson,
		supplier.Phone, supplier.Email, supplier.Address, supplier.City,
		supplier.State, supplier.Country, supplier.PaymentTerms, supplier.Notes,
		supplier.IsActive, supplier.CreatedAt, supplier.UpdatedAt,
	)
	return err
}

func (r *SupplierRepoImpl) UpdateSupplier(ctx context.Context, id string, supplier models.Supplier) error {
	orgID := ctx.Value("organisation_id").(string)
	query := "UPDATE suppliers SET updated_at = ?"
	args := []interface{}{time.Now()}

	if supplier.Name != "" {
		query += ", name = ?"
		args = append(args, supplier.Name)
	}
	if supplier.ContactPerson != "" {
		query += ", contact_person = ?"
		args = append(args, supplier.ContactPerson)
	}
	if supplier.Phone != "" {
		query += ", phone = ?"
		args = append(args, supplier.Phone)
	}
	if supplier.Email != "" {
		query += ", email = ?"
		args = append(args, supplier.Email)
	}
	if supplier.Address != "" {
		query += ", address = ?"
		args = append(args, supplier.Address)
	}
	if supplier.City != "" {
		query += ", city = ?"
		args = append(args, supplier.City)
	}
	if supplier.State != "" {
		query += ", state = ?"
		args = append(args, supplier.State)
	}
	if supplier.Country != "" {
		query += ", country = ?"
		args = append(args, supplier.Country)
	}
	if supplier.PaymentTerms != "" {
		query += ", payment_terms = ?"
		args = append(args, supplier.PaymentTerms)
	}
	if supplier.Notes != "" {
		query += ", notes = ?"
		args = append(args, supplier.Notes)
	}
	if !supplier.IsActive {
		query += ", is_active = ?"
		args = append(args, supplier.IsActive)
	}

	query += fmt.Sprintf(" WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL")
	args = append(args, id, orgID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *SupplierRepoImpl) DeleteSupplier(ctx context.Context, id string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		"UPDATE suppliers SET deleted_at = NOW() WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL",
		id, orgID,
	)
	return err
}

func (r *SupplierRepoImpl) ListSupplierProducts(ctx context.Context, supplierID string) ([]models.SupplierProduct, error) {
	orgID := ctx.Value("organisation_id").(string)
	rows, err := r.db.QueryContext(ctx,
		`SELECT sp.id, sp.organisation_id, sp.supplier_id, sp.product_id, p.name,
		        sp.unit_price, sp.min_order_qty, sp.lead_time_days, sp.notes,
		        sp.created_at, sp.updated_at
		 FROM supplier_products sp
		 JOIN products p ON p.id = sp.product_id
		 WHERE sp.supplier_id = ? AND sp.organisation_id = ?
		 ORDER BY p.name ASC`,
		supplierID, orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.SupplierProduct
	for rows.Next() {
		var sp models.SupplierProduct
		var productName string
		if err := rows.Scan(&sp.ID, &sp.OrganisationID, &sp.SupplierID, &sp.ProductID,
			&productName, &sp.UnitPrice, &sp.MinOrderQty, &sp.LeadTimeDays,
			&sp.Notes, &sp.CreatedAt, &sp.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, sp)
	}
	return products, nil
}

func (r *SupplierRepoImpl) SetSupplierProducts(ctx context.Context, supplierID string, products []models.SupplierProduct) error {
	orgID := ctx.Value("organisation_id").(string)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		"DELETE FROM supplier_products WHERE supplier_id = ? AND organisation_id = ?",
		supplierID, orgID,
	); err != nil {
		return err
	}

	for _, p := range products {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO supplier_products (id, organisation_id, supplier_id, product_id,
			                                unit_price, min_order_qty, lead_time_days, notes,
			                                created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			p.ID, p.OrganisationID, p.SupplierID, p.ProductID,
			p.UnitPrice, p.MinOrderQty, p.LeadTimeDays, p.Notes,
			p.CreatedAt, p.UpdatedAt,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

var _ SupplierRepository = (*SupplierRepoImpl)(nil)
