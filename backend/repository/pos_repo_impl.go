package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/azdonald/pharmd/backend/models"
)

type POSRepository interface {
	CreateSale(ctx context.Context, sale models.Sale, items []models.SaleItem) error
	GetSaleByID(ctx context.Context, id string) (*models.Sale, []models.SaleItem, []models.Payment, error)
	ListSales(ctx context.Context, status string, page, limit int) ([]models.Sale, int, error)
	UpdateSale(ctx context.Context, id, status, voidedBy, notes string) error
	RecordPayments(ctx context.Context, saleID string, payments []models.Payment) error
	UpdateSalePaid(ctx context.Context, saleID string, paidAmount, changeAmount float64) error
	CompleteSale(ctx context.Context, saleID string, sale models.Sale, items []models.SaleItem, payments []models.Payment, userID string, totalPaid, changeAmount float64) error
	RestoreSaleInventory(ctx context.Context, saleID, status, voidedBy string) error
	GetDailySummary(ctx context.Context, locationID, date string) (*models.DailySummary, []struct{ Method string; Total float64 }, error)
	UpsertDailySummary(ctx context.Context, summary models.DailySummary) error
	CloseDay(ctx context.Context, id, closedBy, notes string) error
}

type POSRepoImpl struct {
	db *sql.DB
}

func NewPOSRepositoryImpl(db *sql.DB) POSRepository {
	return &POSRepoImpl{db: db}
}

func (r *POSRepoImpl) CreateSale(ctx context.Context, sale models.Sale, items []models.SaleItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	orgID := ctx.Value("organisation_id").(string)
	_, err = tx.ExecContext(ctx,
		`INSERT INTO sales (id, organisation_id, location_id, patient_id, prescription_id, sale_number,
		                     sale_type, status, subtotal, tax_total, discount_total, grand_total,
		                     paid_amount, change_amount, notes, created_by, created_at, updated_at)
		 VALUES (?, ?, ?, NULLIF(?,''), NULLIF(?,''), ?, ?, ?, ?, ?, ?, ?, 0, 0, NULLIF(?,''), ?, ?, ?)`,
		sale.ID, orgID, sale.LocationID, sale.PatientID, sale.PrescriptionID, sale.SaleNumber,
		sale.SaleType, sale.Status, sale.Subtotal, sale.TaxTotal, sale.DiscountTotal, sale.GrandTotal,
		sale.Notes, sale.CreatedBy, sale.CreatedAt, sale.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, item := range items {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO sale_items (id, organisation_id, sale_id, product_id, quantity, unit_price, discount, line_total)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			item.ID, orgID, item.SaleID, item.ProductID, item.Quantity, item.UnitPrice, item.Discount, item.LineTotal,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *POSRepoImpl) GetSaleByID(ctx context.Context, id string) (*models.Sale, []models.SaleItem, []models.Payment, error) {
	orgID := ctx.Value("organisation_id").(string)
	s := &models.Sale{}

	err := r.db.QueryRowContext(ctx,
		`SELECT s.id, s.organisation_id, s.location_id, COALESCE(s.patient_id,''), COALESCE(pa.first_name,''),
		        COALESCE(s.prescription_id,''), s.sale_number, s.sale_type, s.status,
		        s.subtotal, s.tax_total, s.discount_total, s.grand_total,
		        s.paid_amount, s.change_amount, COALESCE(s.notes,''), s.created_by,
		        COALESCE(s.voided_by,''), COALESCE(s.voided_at,''),
		        s.created_at, s.updated_at
		 FROM sales s
		 LEFT JOIN patients pa ON pa.id = s.patient_id
		 WHERE s.id = ? AND s.organisation_id = ?`,
		id, orgID,
	).Scan(&s.ID, &s.OrganisationID, &s.LocationID, &s.PatientID, &s.PatientName,
		&s.PrescriptionID, &s.SaleNumber, &s.SaleType, &s.Status,
		&s.Subtotal, &s.TaxTotal, &s.DiscountTotal, &s.GrandTotal,
		&s.PaidAmount, &s.ChangeAmount, &s.Notes, &s.CreatedBy,
		&s.VoidedBy, &s.VoidedAt, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, nil, nil, err
	}

	itemRows, err := r.db.QueryContext(ctx,
		`SELECT i.id, i.organisation_id, i.sale_id, i.product_id, COALESCE(p.name,''),
		        i.quantity, i.unit_price, i.discount, i.line_total
		 FROM sale_items i
		 JOIN products p ON p.id = i.product_id
		 WHERE i.sale_id = ? AND i.organisation_id = ?`,
		id, orgID,
	)
	if err != nil {
		return nil, nil, nil, err
	}
	defer itemRows.Close()

	var items []models.SaleItem
	for itemRows.Next() {
		var item models.SaleItem
		if err := itemRows.Scan(&item.ID, &item.OrganisationID, &item.SaleID, &item.ProductID, &item.ProductName,
			&item.Quantity, &item.UnitPrice, &item.Discount, &item.LineTotal); err != nil {
			return nil, nil, nil, err
		}
		items = append(items, item)
	}

	payRows, err := r.db.QueryContext(ctx,
		`SELECT id, organisation_id, sale_id, method, amount, COALESCE(reference,'')
		 FROM payments WHERE sale_id = ? AND organisation_id = ?`,
		id, orgID,
	)
	if err != nil {
		return nil, nil, nil, err
	}
	defer payRows.Close()

	var payments []models.Payment
	for payRows.Next() {
		var p models.Payment
		if err := payRows.Scan(&p.ID, &p.OrganisationID, &p.SaleID, &p.Method, &p.Amount, &p.Reference); err != nil {
			return nil, nil, nil, err
		}
		payments = append(payments, p)
	}

	return s, items, payments, nil
}

func (r *POSRepoImpl) ListSales(ctx context.Context, status string, page, limit int) ([]models.Sale, int, error) {
	orgID := ctx.Value("organisation_id").(string)
	offset := (page - 1) * limit

	where := "s.organisation_id = ?"
	args := []interface{}{orgID}
	if status != "" {
		where += " AND s.status = ?"
		args = append(args, status)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM sales s WHERE %s", where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT s.id, s.organisation_id, s.location_id, COALESCE(s.patient_id,''), COALESCE(pa.first_name,''),
		        COALESCE(s.prescription_id,''), s.sale_number, s.sale_type, s.status,
		        s.subtotal, s.tax_total, s.discount_total, s.grand_total,
		        s.paid_amount, s.change_amount, COALESCE(s.notes,''), s.created_by,
		        COALESCE(s.voided_by,''), COALESCE(s.voided_at,''),
		        s.created_at, s.updated_at
		 FROM sales s
		 LEFT JOIN patients pa ON pa.id = s.patient_id
		 WHERE %s ORDER BY s.created_at DESC LIMIT ? OFFSET ?`, where),
		append(args, limit, offset)...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sales []models.Sale
	for rows.Next() {
		var s models.Sale
		if err := rows.Scan(&s.ID, &s.OrganisationID, &s.LocationID, &s.PatientID, &s.PatientName,
			&s.PrescriptionID, &s.SaleNumber, &s.SaleType, &s.Status,
			&s.Subtotal, &s.TaxTotal, &s.DiscountTotal, &s.GrandTotal,
			&s.PaidAmount, &s.ChangeAmount, &s.Notes, &s.CreatedBy,
			&s.VoidedBy, &s.VoidedAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		sales = append(sales, s)
	}
	return sales, total, nil
}

func (r *POSRepoImpl) UpdateSale(ctx context.Context, id, status, voidedBy, notes string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		`UPDATE sales SET status=?, voided_by=NULLIF(?,''), voided_at=CASE WHEN ? IN ('voided','refunded') THEN NOW() ELSE NULL END,
		        notes=COALESCE(NULLIF(?,''), notes), updated_at=NOW()
		 WHERE id=? AND organisation_id=?`,
		status, voidedBy, status, notes, id, orgID,
	)
	return err
}

func (r *POSRepoImpl) RecordPayments(ctx context.Context, saleID string, payments []models.Payment) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	orgID := ctx.Value("organisation_id").(string)
	for _, p := range payments {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO payments (id, organisation_id, sale_id, method, amount, reference)
			 VALUES (?, ?, ?, ?, ?, NULLIF(?,''))`,
			p.ID, orgID, saleID, p.Method, p.Amount, p.Reference,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *POSRepoImpl) UpdateSalePaid(ctx context.Context, saleID string, paidAmount, changeAmount float64) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		"UPDATE sales SET paid_amount=?, change_amount=?, status='completed', updated_at=NOW() WHERE id=? AND organisation_id=?",
		paidAmount, changeAmount, saleID, orgID,
	)
	return err
}

func (r *POSRepoImpl) CompleteSale(ctx context.Context, saleID string, sale models.Sale, items []models.SaleItem, payments []models.Payment, userID string, totalPaid, changeAmount float64) error {
	orgID := ctx.Value("organisation_id").(string)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, p := range payments {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO payments (id, organisation_id, sale_id, method, amount, reference)
			 VALUES (?, ?, ?, ?, ?, NULLIF(?,''))`,
			p.ID, orgID, saleID, p.Method, p.Amount, p.Reference,
		); err != nil {
			return err
		}
	}

	for _, item := range items {
		if err := deductFEFO(ctx, tx, orgID, sale.LocationID, item.ProductID, item.Quantity, "sale", "sale", saleID, userID); err != nil {
			return err
		}
	}

	if _, err := tx.ExecContext(ctx,
		"UPDATE sales SET paid_amount=?, change_amount=?, status='completed', updated_at=NOW() WHERE id=? AND organisation_id=?",
		totalPaid, changeAmount, saleID, orgID,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *POSRepoImpl) RestoreSaleInventory(ctx context.Context, saleID, status, voidedBy string) error {
	orgID := ctx.Value("organisation_id").(string)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		`UPDATE sales SET status=?, voided_by=NULLIF(?,''), voided_at=NOW(), updated_at=NOW()
		 WHERE id=? AND organisation_id=? AND status='completed'`,
		status, voidedBy, saleID, orgID,
	); err != nil {
		return err
	}

	if err := restoreMovements(ctx, tx, orgID, "sale", "sale", saleID, voidedBy); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *POSRepoImpl) GetDailySummary(ctx context.Context, locationID, date string) (*models.DailySummary, []struct{ Method string; Total float64 }, error) {
	orgID := ctx.Value("organisation_id").(string)
	ds := &models.DailySummary{}

	err := r.db.QueryRowContext(ctx,
		`SELECT id, organisation_id, location_id, date, total_sales, total_revenue,
		        total_tax, total_discounts, is_closed, COALESCE(closed_by,''),
		        COALESCE(closed_at,''), COALESCE(notes,''), created_at, updated_at
		 FROM daily_summaries
		 WHERE organisation_id=? AND location_id=? AND date=?
		 ORDER BY created_at DESC LIMIT 1`,
		orgID, locationID, date,
	).Scan(&ds.ID, &ds.OrganisationID, &ds.LocationID, &ds.Date,
		&ds.TotalSales, &ds.TotalRevenue, &ds.TotalTax, &ds.TotalDiscounts,
		&ds.IsClosed, &ds.ClosedBy, &ds.ClosedAt, &ds.Notes,
		&ds.CreatedAt, &ds.UpdatedAt)
	if err == sql.ErrNoRows {
		// build summary from sales
		row := r.db.QueryRowContext(ctx,
			`SELECT COUNT(*), COALESCE(SUM(grand_total),0), COALESCE(SUM(tax_total),0), COALESCE(SUM(discount_total),0)
			 FROM sales WHERE organisation_id=? AND location_id=? AND DATE(created_at)=? AND status='completed'`,
			orgID, locationID, date,
		)
		if err := row.Scan(&ds.TotalSales, &ds.TotalRevenue, &ds.TotalTax, &ds.TotalDiscounts); err != nil {
			return nil, nil, err
		}
		ds.OrganisationID = orgID
		ds.LocationID = locationID
		ds.Date = date
		ds.IsClosed = false
	} else if err != nil {
		return nil, nil, err
	}

	// get totals by payment method
	rows, err := r.db.QueryContext(ctx,
		`SELECT p.method, COALESCE(SUM(p.amount),0)
		 FROM payments p
		 JOIN sales s ON s.id = p.sale_id
		 WHERE s.organisation_id=? AND s.location_id=? AND DATE(s.created_at)=? AND s.status='completed'
		 GROUP BY p.method`,
		orgID, locationID, date,
	)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var methods []struct{ Method string; Total float64 }
	for rows.Next() {
		var m struct{ Method string; Total float64 }
		if err := rows.Scan(&m.Method, &m.Total); err != nil {
			return nil, nil, err
		}
		methods = append(methods, m)
	}

	return ds, methods, nil
}

func (r *POSRepoImpl) UpsertDailySummary(ctx context.Context, ds models.DailySummary) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO daily_summaries (id, organisation_id, location_id, date, total_sales, total_revenue, total_tax, total_discounts, is_closed, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)
		 ON DUPLICATE KEY UPDATE total_sales=VALUES(total_sales), total_revenue=VALUES(total_revenue),
		         total_tax=VALUES(total_tax), total_discounts=VALUES(total_discounts), updated_at=VALUES(updated_at)`,
		ds.ID, orgID, ds.LocationID, ds.Date, ds.TotalSales, ds.TotalRevenue, ds.TotalTax, ds.TotalDiscounts,
		ds.CreatedAt, ds.UpdatedAt,
	)
	return err
}

func (r *POSRepoImpl) CloseDay(ctx context.Context, id, closedBy, notes string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE daily_summaries SET is_closed=1, closed_by=?, closed_at=NOW(), notes=COALESCE(NULLIF(?,''), notes), updated_at=NOW() WHERE id=?`,
		closedBy, notes, id,
	)
	return err
}

var _ POSRepository = (*POSRepoImpl)(nil)
