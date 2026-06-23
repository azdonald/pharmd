package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func deductFEFO(ctx context.Context, tx *sql.Tx, orgID, locationID, productID string, quantity int, movementType, referenceType, referenceID, userID string) error {
	rows, err := tx.QueryContext(ctx,
		`SELECT id, remaining_qty FROM stock_batches
		 WHERE organisation_id = ? AND location_id = ? AND product_id = ?
		   AND remaining_qty > 0 AND deleted_at IS NULL
		 ORDER BY expiry_date ASC
		 FOR UPDATE`,
		orgID, locationID, productID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	toDeduct := quantity
	for rows.Next() {
		var batchID string
		var remainingQty int
		if err := rows.Scan(&batchID, &remainingQty); err != nil {
			return err
		}

		deductQty := toDeduct
		if deductQty > remainingQty {
			deductQty = remainingQty
		}

		newQty := remainingQty - deductQty
		if _, err := tx.ExecContext(ctx,
			"UPDATE stock_batches SET remaining_qty = ?, updated_at = NOW() WHERE id = ? AND organisation_id = ?",
			newQty, batchID, orgID,
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx,
			`INSERT INTO stock_movements (id, organisation_id, location_id, product_id, batch_id,
			                               movement_type, quantity, reference_type, reference_id,
			                               notes, created_by, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, '', ?, NOW(), NOW())`,
			uuid.New().String(), orgID, locationID, productID, batchID,
			movementType, deductQty, referenceType, referenceID, userID,
		); err != nil {
			return err
		}

		toDeduct -= deductQty
		if toDeduct == 0 {
			break
		}
	}

	if toDeduct > 0 {
		return fmt.Errorf("insufficient stock for product %s: short %d units", productID, toDeduct)
	}
	return nil
}

func restoreMovements(ctx context.Context, tx *sql.Tx, orgID string, movementType, referenceType, referenceID, userID string) error {
	rows, err := tx.QueryContext(ctx,
		`SELECT location_id, product_id, batch_id, quantity FROM stock_movements
		 WHERE organisation_id = ? AND reference_type = ? AND reference_id = ?
		   AND movement_type = ? AND batch_id IS NOT NULL AND batch_id != ''
		 FOR UPDATE`,
		orgID, referenceType, referenceID, movementType,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	reversalType := movementType + "_reversal"
	for rows.Next() {
		var locationID, productID, batchID string
		var qty int
		if err := rows.Scan(&locationID, &productID, &batchID, &qty); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx,
			"UPDATE stock_batches SET remaining_qty = remaining_qty + ?, updated_at = NOW() WHERE id = ? AND organisation_id = ?",
			qty, batchID, orgID,
		); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx,
			`INSERT INTO stock_movements (id, organisation_id, location_id, product_id, batch_id,
			                               movement_type, quantity, reference_type, reference_id,
			                               notes, created_by, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, '', ?, NOW(), NOW())`,
			uuid.New().String(), orgID, locationID, productID, batchID,
			reversalType, qty, referenceType, referenceID, userID,
		); err != nil {
			return err
		}
	}
	return nil
}
