package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/azdonald/pharmd/backend/models"
)

type PurchaseOrderRepoImpl struct {
	db *sql.DB
}

func NewPurchaseOrderRepositoryImpl(db *sql.DB) PurchaseOrderRepository {
	return &PurchaseOrderRepoImpl{db: db}
}

func (r *PurchaseOrderRepoImpl) ListPurchaseOrders(ctx context.Context, page, limit int, status string) ([]models.PurchaseOrder, int, error) {
	orgID := ctx.Value("organisation_id").(string)
	offset := (page - 1) * limit

	where := "po.organisation_id = ? AND po.deleted_at IS NULL"
	args := []interface{}{orgID}

	if status != "" {
		where += " AND po.status = ?"
		args = append(args, status)
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM purchase_orders po WHERE %s", where)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT po.id, po.organisation_id, po.po_number, po.supplier_id, s.name,
		       po.location_id, po.status, po.order_date, po.expected_date, po.notes,
		       po.subtotal, po.tax_total, po.grand_total, po.created_by,
		       po.approved_by, po.approved_at, po.created_at, po.updated_at
		 FROM purchase_orders po
		 JOIN suppliers s ON s.id = po.supplier_id
		 WHERE %s ORDER BY po.created_at DESC LIMIT ? OFFSET ?`, where),
		append(args, limit, offset)...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []models.PurchaseOrder
	for rows.Next() {
		var o models.PurchaseOrder
		var approvedBy, approvedAt sql.NullString
		if err := rows.Scan(&o.ID, &o.OrganisationID, &o.PONumber, &o.SupplierID,
			&o.Notes, &o.LocationID, &o.Status, &o.OrderDate, &o.ExpectedDate,
			&o.Notes, &o.Subtotal, &o.TaxTotal, &o.GrandTotal, &o.CreatedBy,
			&approvedBy, &approvedAt, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		o.ApprovedBy = approvedBy.String
		orders = append(orders, o)
	}
	return orders, total, nil
}

func (r *PurchaseOrderRepoImpl) GetPurchaseOrderByID(ctx context.Context, id string) (*models.PurchaseOrder, error) {
	orgID := ctx.Value("organisation_id").(string)
	o := &models.PurchaseOrder{}
	var approvedBy, approvedAt sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT po.id, po.organisation_id, po.po_number, po.supplier_id, s.name,
		        po.location_id, po.status, po.order_date, po.expected_date, po.notes,
		        po.subtotal, po.tax_total, po.grand_total, po.created_by,
		        po.approved_by, po.approved_at, po.created_at, po.updated_at
		 FROM purchase_orders po
		 JOIN suppliers s ON s.id = po.supplier_id
		 WHERE po.id = ? AND po.organisation_id = ? AND po.deleted_at IS NULL`,
		id, orgID,
	).Scan(&o.ID, &o.OrganisationID, &o.PONumber, &o.SupplierID,
		&o.Notes, &o.LocationID, &o.Status, &o.OrderDate, &o.ExpectedDate,
		&o.Notes, &o.Subtotal, &o.TaxTotal, &o.GrandTotal, &o.CreatedBy,
		&approvedBy, &approvedAt, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	o.ApprovedBy = approvedBy.String
	return o, nil
}

func (r *PurchaseOrderRepoImpl) CreatePurchaseOrder(ctx context.Context, po models.PurchaseOrder, items []models.PurchaseOrderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO purchase_orders (id, organisation_id, po_number, supplier_id, location_id,
		                               status, order_date, expected_date, notes,
		                               subtotal, tax_total, grand_total, created_by,
		                               created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULLIF(?, ''), ?, ?, ?, ?, ?, ?)`,
		po.ID, po.OrganisationID, po.PONumber, po.SupplierID, po.LocationID,
		po.Status, po.OrderDate, po.ExpectedDate, po.Notes,
		po.Subtotal, po.TaxTotal, po.GrandTotal, po.CreatedBy,
		po.CreatedAt, po.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, item := range items {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO purchase_order_items (id, organisation_id, purchase_order_id, product_id,
			                                    quantity_ordered, quantity_received, unit_cost, line_total,
			                                    created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?, ?)`,
			item.ID, item.OrganisationID, item.PurchaseOrderID, item.ProductID,
			item.QuantityOrdered, item.UnitCost, item.LineTotal,
			po.CreatedAt, po.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PurchaseOrderRepoImpl) UpdatePOStatus(ctx context.Context, id, status, approvedBy string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		"UPDATE purchase_orders SET status = ?, approved_by = NULLIF(?, ''), approved_at = CASE WHEN ? = 'approved' THEN NOW() ELSE NULL END, updated_at = NOW() WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL",
		status, approvedBy, status, id, orgID,
	)
	return err
}

func (r *PurchaseOrderRepoImpl) ReceiveGoods(ctx context.Context, poID string, items []models.PurchaseOrderItem, notes, userID string) error {
	orgID := ctx.Value("organisation_id").(string)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range items {
		_, err = tx.ExecContext(ctx,
			`UPDATE purchase_order_items
			 SET quantity_received = quantity_received + ?, updated_at = NOW()
			 WHERE id = ? AND purchase_order_id = ? AND organisation_id = ?`,
			item.QuantityReceived, item.ID, poID, orgID,
		)
		if err != nil {
			return err
		}
	}

	// record GRN
	_, err = tx.ExecContext(ctx,
		`INSERT INTO goods_received_notes (id, organisation_id, purchase_order_id, received_date, notes, created_by, created_at, updated_at)
		 VALUES (?, ?, ?, NOW(), ?, ?, NOW(), NOW())`,
		fmt.Sprintf("GRN-%s-%d", poID[:8], time.Now().UnixNano()), orgID, poID, notes, userID,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE purchase_orders SET status = 'received', updated_at = NOW()
		 WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL`,
		poID, orgID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

var _ PurchaseOrderRepository = (*PurchaseOrderRepoImpl)(nil)
