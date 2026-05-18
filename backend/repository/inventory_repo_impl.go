package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/azdonald/pharmd/backend/models"
)

type InventoryRepoImpl struct {
	db *sql.DB
}

func NewInventoryRepositoryImpl(db *sql.DB) InventoryRepository {
	return &InventoryRepoImpl{db: db}
}

func (r *InventoryRepoImpl) CreateBatch(ctx context.Context, batch models.StockBatch) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO stock_batches (id, organisation_id, product_id, location_id, batch_number,
		                            quantity, remaining_qty, unit_cost, selling_price,
		                            manufacturing_date, expiry_date, received_date,
		                            is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NULLIF(?, ''), NULLIF(?, ''), ?, ?, ?, ?)`,
		batch.ID, batch.OrganisationID, batch.ProductID, batch.LocationID, batch.BatchNumber,
		batch.Quantity, batch.RemainingQty, batch.UnitCost, batch.SellingPrice,
		batch.ManufacturingDate, batch.ExpiryDate, batch.ReceivedDate,
		batch.IsActive, batch.CreatedAt, batch.UpdatedAt,
	)
	return err
}

func (r *InventoryRepoImpl) GetBatchByID(ctx context.Context, id string) (*models.StockBatch, error) {
	orgID := ctx.Value("organisation_id").(string)
	b := &models.StockBatch{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, organisation_id, product_id, location_id, batch_number,
		        quantity, remaining_qty, unit_cost, selling_price,
		        manufacturing_date, expiry_date, received_date, is_active, created_at, updated_at
		 FROM stock_batches WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&b.ID, &b.OrganisationID, &b.ProductID, &b.LocationID, &b.BatchNumber,
		&b.Quantity, &b.RemainingQty, &b.UnitCost, &b.SellingPrice,
		&b.ManufacturingDate, &b.ExpiryDate, &b.ReceivedDate,
		&b.IsActive, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *InventoryRepoImpl) UpdateBatchQty(ctx context.Context, id string, remainingQty int) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		"UPDATE stock_batches SET remaining_qty = ?, updated_at = NOW() WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL",
		remainingQty, id, orgID,
	)
	return err
}

func (r *InventoryRepoImpl) ListStock(ctx context.Context, locationID string, page, limit int, query string) ([]InventoryBatchView, int, error) {
	orgID := ctx.Value("organisation_id").(string)
	offset := (page - 1) * limit

	where := "sb.organisation_id = ? AND sb.location_id = ? AND sb.deleted_at IS NULL"
	args := []interface{}{orgID, locationID}

	if query != "" {
		where += " AND (p.name LIKE ? OR p.brand_name LIKE ? OR sb.batch_number LIKE ?)"
		q := "%" + query + "%"
		args = append(args, q, q, q)
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM stock_batches sb JOIN products p ON p.id = sb.product_id WHERE %s", where)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT sb.id, sb.organisation_id, sb.product_id, sb.location_id, sb.batch_number,
		       sb.quantity, sb.remaining_qty, sb.unit_cost, sb.selling_price,
		       sb.manufacturing_date, sb.expiry_date, sb.received_date, sb.is_active,
		       sb.created_at, sb.updated_at,
		       p.name, p.brand_name, p.generic_name, p.classification, p.reorder_level
		 FROM stock_batches sb
		 JOIN products p ON p.id = sb.product_id
		 WHERE %s ORDER BY p.name ASC, sb.expiry_date ASC LIMIT ? OFFSET ?`, where),
		append(args, limit, offset)...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []InventoryBatchView
	for rows.Next() {
		var v InventoryBatchView
		if err := rows.Scan(&v.StockBatch.ID, &v.StockBatch.OrganisationID,
			&v.StockBatch.ProductID, &v.StockBatch.LocationID, &v.StockBatch.BatchNumber,
			&v.StockBatch.Quantity, &v.StockBatch.RemainingQty, &v.StockBatch.UnitCost,
			&v.StockBatch.SellingPrice, &v.StockBatch.ManufacturingDate,
			&v.StockBatch.ExpiryDate, &v.StockBatch.ReceivedDate, &v.StockBatch.IsActive,
			&v.StockBatch.CreatedAt, &v.StockBatch.UpdatedAt,
			&v.ProductName, &v.BrandName, &v.GenericName, &v.Classification, &v.ReorderLevel); err != nil {
			return nil, 0, err
		}
		results = append(results, v)
	}
	return results, total, nil
}

func (r *InventoryRepoImpl) ListAlerts(ctx context.Context, locationID string) ([]InventoryAlertView, error) {
	orgID := ctx.Value("organisation_id").(string)

	where := "sb.organisation_id = ? AND sb.deleted_at IS NULL"
	args := []interface{}{orgID}
	if locationID != "" {
		where += " AND sb.location_id = ?"
		args = append(args, locationID)
	}

	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT sb.product_id, p.name, p.brand_name,
		       SUM(sb.remaining_qty) as total_qty, p.reorder_level, sb.location_id
		 FROM stock_batches sb
		 JOIN products p ON p.id = sb.product_id
		 WHERE %s
		 GROUP BY sb.product_id, sb.location_id
		 HAVING total_qty < p.reorder_level
		 ORDER BY total_qty ASC`, where),
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []InventoryAlertView
	for rows.Next() {
		var a InventoryAlertView
		if err := rows.Scan(&a.ProductID, &a.ProductName, &a.BrandName,
			&a.TotalQuantity, &a.ReorderLevel, &a.LocationID); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

func (r *InventoryRepoImpl) ListExpiring(ctx context.Context, locationID string, days int) ([]InventoryExpiringView, error) {
	orgID := ctx.Value("organisation_id").(string)

	where := "sb.organisation_id = ? AND sb.deleted_at IS NULL AND sb.remaining_qty > 0"
	args := []interface{}{orgID}
	if locationID != "" {
		where += " AND sb.location_id = ?"
		args = append(args, locationID)
	}

	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT sb.id, sb.organisation_id, sb.product_id, sb.location_id, sb.batch_number,
		       sb.quantity, sb.remaining_qty, sb.unit_cost, sb.selling_price,
		       sb.manufacturing_date, sb.expiry_date, sb.received_date, sb.is_active,
		       sb.created_at, sb.updated_at,
		       p.name,
		       DATEDIFF(STR_TO_DATE(sb.expiry_date, '%%Y-%%m-%%d'), CURDATE()) as days_until_expiry
		 FROM stock_batches sb
		 JOIN products p ON p.id = sb.product_id
		 WHERE %s AND sb.expiry_date IS NOT NULL AND sb.expiry_date != ''
		   AND STR_TO_DATE(sb.expiry_date, '%%Y-%%m-%%d') <= DATE_ADD(CURDATE(), INTERVAL ? DAY)
		   AND STR_TO_DATE(sb.expiry_date, '%%Y-%%m-%%d') >= CURDATE()
		 ORDER BY sb.expiry_date ASC`, where),
		append(args, days)...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []InventoryExpiringView
	for rows.Next() {
		var v InventoryExpiringView
		if err := rows.Scan(&v.StockBatch.ID, &v.StockBatch.OrganisationID,
			&v.StockBatch.ProductID, &v.StockBatch.LocationID, &v.StockBatch.BatchNumber,
			&v.StockBatch.Quantity, &v.StockBatch.RemainingQty, &v.StockBatch.UnitCost,
			&v.StockBatch.SellingPrice, &v.StockBatch.ManufacturingDate,
			&v.StockBatch.ExpiryDate, &v.StockBatch.ReceivedDate, &v.StockBatch.IsActive,
			&v.StockBatch.CreatedAt, &v.StockBatch.UpdatedAt,
			&v.ProductName, &v.DaysUntilExpiry); err != nil {
			return nil, err
		}
		results = append(results, v)
	}
	return results, nil
}

func (r *InventoryRepoImpl) CreateMovement(ctx context.Context, movement models.StockMovement) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO stock_movements (id, organisation_id, location_id, product_id, batch_id,
		                               movement_type, quantity, reference_type, reference_id,
		                               notes, created_by, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		movement.ID, movement.OrganisationID, movement.LocationID, movement.ProductID,
		movement.BatchID, movement.MovementType, movement.Quantity,
		movement.ReferenceType, movement.ReferenceID, movement.Notes,
		movement.CreatedBy, movement.CreatedAt, movement.CreatedAt,
	)
	return err
}

func (r *InventoryRepoImpl) StockCount(ctx context.Context, items []models.StockCountItem, userID string) (int, error) {
	orgID := ctx.Value("organisation_id").(string)
	adjustments := 0

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	for _, item := range items {
		// get current remaining qty
		var currentQty int
		var locationID string
		if item.BatchID != "" {
			err = tx.QueryRowContext(ctx,
				`SELECT remaining_qty, location_id FROM stock_batches
				 WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL FOR UPDATE`,
				item.BatchID, orgID,
			).Scan(&currentQty, &locationID)
			if err != nil {
				return 0, err
			}
		} else {
			err = tx.QueryRowContext(ctx,
				`SELECT COALESCE(SUM(remaining_qty), 0), ? FROM stock_batches
				 WHERE product_id = ? AND location_id = ? AND organisation_id = ? AND deleted_at IS NULL`,
				item.LocationID, item.ProductID, item.LocationID, orgID,
			).Scan(&currentQty)
			if err != nil {
				return 0, err
			}
			locationID = item.LocationID
		}

		if currentQty != item.CountedQty {
			diff := item.CountedQty - currentQty
			adjustments++

			movementID := fmt.Sprintf("MOV-%s-%d", orgID[:8], time.Now().UnixNano())
			notes := item.Notes
			if notes == "" {
				notes = fmt.Sprintf("Stock count adjustment (%d → %d)", currentQty, item.CountedQty)
			}

			_, err = tx.ExecContext(ctx,
				`INSERT INTO stock_movements (id, organisation_id, location_id, product_id, batch_id,
				                               movement_type, quantity, reference_type, reference_id,
				                               notes, created_by, created_at, updated_at)
				 VALUES (?, ?, ?, ?, ?, 'count_correction', ?, '', '', ?, ?, NOW(), NOW())`,
				movementID, orgID, locationID, item.ProductID, item.BatchID,
				diff, notes, userID,
			)
			if err != nil {
				return 0, err
			}

			// update batch qty (for batch-specific counts)
			if item.BatchID != "" {
				_, err = tx.ExecContext(ctx,
					`UPDATE stock_batches SET remaining_qty = ?, updated_at = NOW()
					 WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL`,
					item.CountedQty, item.BatchID, orgID,
				)
				if err != nil {
					return 0, err
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return adjustments, nil
}

var _ InventoryRepository = (*InventoryRepoImpl)(nil)
