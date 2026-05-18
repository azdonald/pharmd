package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/azdonald/pharmd/backend/models"
)

type PrescriptionRepository interface {
	List(ctx context.Context, status, patientID string, page, limit int) ([]struct {
		models.Prescription
		PatientName    string
		PrescriberName string
	}, int, error)
	GetByID(ctx context.Context, id string) (*models.Prescription, []models.PrescriptionItem, error)
	Create(ctx context.Context, rx models.Prescription, items []models.PrescriptionItem) error
	Update(ctx context.Context, id string, rx models.Prescription) error
	Delete(ctx context.Context, id string) error
	RecordRefill(ctx context.Context, refill models.Refill) error
	IncrementRefillUsed(ctx context.Context, itemID string) error
}

type PrescriptionRepoImpl struct {
	db *sql.DB
}

func NewPrescriptionRepositoryImpl(db *sql.DB) PrescriptionRepository {
	return &PrescriptionRepoImpl{db: db}
}

func (r *PrescriptionRepoImpl) List(ctx context.Context, status, patientID string, page, limit int) ([]struct {
	models.Prescription
	PatientName    string
	PrescriberName string
}, int, error) {
	orgID := ctx.Value("organisation_id").(string)
	offset := (page - 1) * limit

	where := "rx.organisation_id = ? AND rx.deleted_at IS NULL"
	args := []interface{}{orgID}
	if status != "" {
		where += " AND rx.status = ?"
		args = append(args, status)
	}
	if patientID != "" {
		where += " AND rx.patient_id = ?"
		args = append(args, patientID)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM prescriptions rx WHERE %s", where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT rx.id, rx.organisation_id, rx.patient_id, p.first_name, rx.prescriber_id, pr.name,
		       rx.location_id, rx.status, COALESCE(rx.diagnosis,''), COALESCE(rx.notes,''),
		       COALESCE(rx.issued_date,''), COALESCE(rx.expiry_date,''),
		       rx.created_by, rx.created_at, rx.updated_at
		 FROM prescriptions rx
		 JOIN patients p ON p.id = rx.patient_id
		 JOIN prescribers pr ON pr.id = rx.prescriber_id
		 WHERE %s ORDER BY rx.created_at DESC LIMIT ? OFFSET ?`, where),
		append(args, limit, offset)...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []struct {
		models.Prescription
		PatientName    string
		PrescriberName string
	}
	for rows.Next() {
		var r struct {
			models.Prescription
			PatientName    string
			PrescriberName string
		}
		if err := rows.Scan(&r.ID, &r.OrganisationID, &r.PatientID, &r.PatientName, &r.PrescriberID, &r.PrescriberName,
			&r.LocationID, &r.Status, &r.Diagnosis, &r.Notes, &r.IssuedDate, &r.ExpiryDate,
			&r.CreatedBy, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, r)
	}

	// convert to return type
	out := make([]struct {
		models.Prescription
		PatientName    string
		PrescriberName string
	}, len(result))
	for i, r := range result {
		out[i] = struct {
			models.Prescription
			PatientName    string
			PrescriberName string
		}{r.Prescription, r.PatientName, r.PrescriberName}
	}
	return out, total, nil
}

func (r *PrescriptionRepoImpl) GetByID(ctx context.Context, id string) (*models.Prescription, []models.PrescriptionItem, error) {
	orgID := ctx.Value("organisation_id").(string)
	rx := &models.Prescription{}
	items := []models.PrescriptionItem{}

	err := r.db.QueryRowContext(ctx,
		`SELECT rx.id, rx.organisation_id, rx.patient_id, COALESCE(p.first_name,''), rx.prescriber_id, COALESCE(pr.name,''),
		        rx.location_id, rx.status, COALESCE(rx.diagnosis,''), COALESCE(rx.notes,''),
		        COALESCE(rx.issued_date,''), COALESCE(rx.expiry_date,''),
		        rx.created_by, rx.created_at, rx.updated_at
		 FROM prescriptions rx
		 JOIN patients p ON p.id = rx.patient_id
		 JOIN prescribers pr ON pr.id = rx.prescriber_id
		 WHERE rx.id = ? AND rx.organisation_id = ? AND rx.deleted_at IS NULL`,
		id, orgID,
	).Scan(&rx.ID, &rx.OrganisationID, &rx.PatientID, &rx.PatientName, &rx.PrescriberID, &rx.PrescriberName,
		&rx.LocationID, &rx.Status, &rx.Diagnosis, &rx.Notes, &rx.IssuedDate, &rx.ExpiryDate,
		&rx.CreatedBy, &rx.CreatedAt, &rx.UpdatedAt)
	if err != nil {
		return nil, nil, err
	}

	itemRows, err := r.db.QueryContext(ctx,
		`SELECT i.id, i.organisation_id, i.prescription_id, i.product_id, COALESCE(pd.name,''),
		        i.dosage, i.frequency, COALESCE(i.duration,''),
		        i.quantity, i.refills_authorized, i.refills_used, COALESCE(i.notes,'')
		 FROM prescription_items i
		 JOIN products pd ON pd.id = i.product_id
		 WHERE i.prescription_id = ? AND i.organisation_id = ?`,
		id, orgID,
	)
	if err != nil {
		return nil, nil, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item models.PrescriptionItem
		if err := itemRows.Scan(&item.ID, &item.OrganisationID, &item.PrescriptionID, &item.ProductID, &item.ProductName,
			&item.Dosage, &item.Frequency, &item.Duration,
			&item.Quantity, &item.RefillsAuthorized, &item.RefillsUsed, &item.Notes); err != nil {
			return nil, nil, err
		}
		items = append(items, item)
	}

	return rx, items, nil
}

func (r *PrescriptionRepoImpl) Create(ctx context.Context, rx models.Prescription, items []models.PrescriptionItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	orgID := ctx.Value("organisation_id").(string)
	_, err = tx.ExecContext(ctx,
		`INSERT INTO prescriptions (id, organisation_id, patient_id, prescriber_id, location_id, status, diagnosis, notes, issued_date, expiry_date, created_by, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, NULLIF(?,''), NULLIF(?,''), NULLIF(?,''), NULLIF(?,''), ?, ?, ?)`,
		rx.ID, orgID, rx.PatientID, rx.PrescriberID, rx.LocationID, rx.Status,
		rx.Diagnosis, rx.Notes, rx.IssuedDate, rx.ExpiryDate, rx.CreatedBy, rx.CreatedAt, rx.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, item := range items {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO prescription_items (id, organisation_id, prescription_id, product_id, dosage, frequency, duration, quantity, refills_authorized, refills_used, notes)
			 VALUES (?, ?, ?, ?, ?, ?, NULLIF(?,''), ?, ?, 0, NULLIF(?,''))`,
			item.ID, orgID, item.PrescriptionID, item.ProductID, item.Dosage, item.Frequency,
			item.Duration, item.Quantity, item.RefillsAuthorized, item.Notes,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PrescriptionRepoImpl) Update(ctx context.Context, id string, rx models.Prescription) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		`UPDATE prescriptions SET status=NULLIF(?,''), diagnosis=NULLIF(?,''), notes=NULLIF(?,''), expiry_date=NULLIF(?,''), updated_at=? WHERE id=? AND organisation_id=? AND deleted_at IS NULL`,
		rx.Status, rx.Diagnosis, rx.Notes, rx.ExpiryDate, rx.UpdatedAt, id, orgID,
	)
	return err
}

func (r *PrescriptionRepoImpl) Delete(ctx context.Context, id string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx, "UPDATE prescriptions SET deleted_at = NOW() WHERE id = ? AND organisation_id = ?", id, orgID)
	return err
}

func (r *PrescriptionRepoImpl) RecordRefill(ctx context.Context, refill models.Refill) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO refills (id, organisation_id, prescription_id, item_id, refilled_by, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		refill.ID, orgID, refill.PrescriptionID, refill.ItemID, refill.RefilledBy, refill.CreatedAt, refill.CreatedAt,
	)
	return err
}

func (r *PrescriptionRepoImpl) IncrementRefillUsed(ctx context.Context, itemID string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE prescription_items SET refills_used = refills_used + 1 WHERE id = ?", itemID,
	)
	return err
}

var _ PrescriptionRepository = (*PrescriptionRepoImpl)(nil)
