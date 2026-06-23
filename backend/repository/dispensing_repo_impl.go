package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/azdonald/pharmd/backend/models"
)

type DispensingRepository interface {
	List(ctx context.Context, status, prescriptionID string, page, limit int) ([]models.DispenseRecord, int, error)
	GetByID(ctx context.Context, id string) (*models.DispenseRecord, error)
	Create(ctx context.Context, dr models.DispenseRecord) error
	Update(ctx context.Context, id string, dr models.DispenseRecord) error
	UpdateStatus(ctx context.Context, id, status string) error
}

type DispensingRepoImpl struct {
	db *sql.DB
}

func NewDispensingRepositoryImpl(db *sql.DB) DispensingRepository {
	return &DispensingRepoImpl{db: db}
}

func (r *DispensingRepoImpl) List(ctx context.Context, status, prescriptionID string, page, limit int) ([]models.DispenseRecord, int, error) {
	orgID := ctx.Value("organisation_id").(string)
	offset := (page - 1) * limit

	where := "d.organisation_id = ?"
	args := []interface{}{orgID}
	if status != "" {
		where += " AND d.status = ?"
		args = append(args, status)
	}
	if prescriptionID != "" {
		where += " AND d.prescription_id = ?"
		args = append(args, prescriptionID)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM dispensing_records d WHERE %s", where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT d.id, d.organisation_id, d.prescription_id, d.prescription_item_id,
		       d.patient_id, COALESCE(pa.first_name,''), d.product_id, COALESCE(p.name,''),
		       d.location_id, d.quantity_dispensed, d.quantity_prescribed,
		       d.pharmacist_id, COALESCE(u.first_name,''), COALESCE(d.technician_id,''),
		       d.status, COALESCE(d.notes,''), COALESCE(d.witness_name,''), d.is_controlled,
		       COALESCE(d.dispensed_at,''), d.created_at, d.updated_at
		 FROM dispensing_records d
		 JOIN patients pa ON pa.id = d.patient_id
		 JOIN products p ON p.id = d.product_id
		 JOIN users u ON u.id = d.pharmacist_id
		 WHERE %s ORDER BY d.created_at DESC LIMIT ? OFFSET ?`, where),
		append(args, limit, offset)...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []models.DispenseRecord
	for rows.Next() {
		var dr models.DispenseRecord
		if err := rows.Scan(&dr.ID, &dr.OrganisationID, &dr.PrescriptionID, &dr.PrescriptionItemID,
			&dr.PatientID, &dr.PatientName, &dr.ProductID, &dr.ProductName,
			&dr.LocationID, &dr.QuantityDispensed, &dr.QuantityPrescribed,
			&dr.PharmacistID, &dr.PharmacistName, &dr.TechnicianID,
			&dr.Status, &dr.Notes, &dr.WitnessName, &dr.IsControlled,
			&dr.DispensedAt, &dr.CreatedAt, &dr.UpdatedAt); err != nil {
			return nil, 0, err
		}
		records = append(records, dr)
	}
	return records, total, nil
}

func (r *DispensingRepoImpl) GetByID(ctx context.Context, id string) (*models.DispenseRecord, error) {
	orgID := ctx.Value("organisation_id").(string)
	dr := &models.DispenseRecord{}
	err := r.db.QueryRowContext(ctx,
		`SELECT d.id, d.organisation_id, d.prescription_id, d.prescription_item_id,
		        d.patient_id, COALESCE(pa.first_name,''), d.product_id, COALESCE(p.name,''),
		        d.location_id, d.quantity_dispensed, d.quantity_prescribed,
		        d.pharmacist_id, COALESCE(u.first_name,''), COALESCE(d.technician_id,''),
		        d.status, COALESCE(d.notes,''), COALESCE(d.witness_name,''), d.is_controlled,
		        COALESCE(d.dispensed_at,''), d.created_at, d.updated_at
		 FROM dispensing_records d
		 JOIN patients pa ON pa.id = d.patient_id
		 JOIN products p ON p.id = d.product_id
		 JOIN users u ON u.id = d.pharmacist_id
		 WHERE d.id = ? AND d.organisation_id = ?`,
		id, orgID,
	).Scan(&dr.ID, &dr.OrganisationID, &dr.PrescriptionID, &dr.PrescriptionItemID,
		&dr.PatientID, &dr.PatientName, &dr.ProductID, &dr.ProductName,
		&dr.LocationID, &dr.QuantityDispensed, &dr.QuantityPrescribed,
		&dr.PharmacistID, &dr.PharmacistName, &dr.TechnicianID,
		&dr.Status, &dr.Notes, &dr.WitnessName, &dr.IsControlled,
		&dr.DispensedAt, &dr.CreatedAt, &dr.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return dr, nil
}

func (r *DispensingRepoImpl) Create(ctx context.Context, dr models.DispenseRecord) error {
	orgID := ctx.Value("organisation_id").(string)
	userID := ctx.Value("user_id").(string)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO dispensing_records (id, organisation_id, prescription_id, prescription_item_id,
		                                  patient_id, product_id, location_id,
		                                  quantity_dispensed, quantity_prescribed,
		                                  pharmacist_id, technician_id, status, notes,
		                                  witness_name, is_controlled, dispensed_at,
		                                  created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULLIF(?,''), ?, NULLIF(?,''), NULLIF(?,''), ?, ?, ?, ?)`,
		dr.ID, orgID, dr.PrescriptionID, dr.PrescriptionItemID,
		dr.PatientID, dr.ProductID, dr.LocationID,
		dr.QuantityDispensed, dr.QuantityPrescribed,
		dr.PharmacistID, dr.TechnicianID, dr.Status, dr.Notes,
		dr.WitnessName, dr.IsControlled, dr.DispensedAt,
		dr.CreatedAt, dr.UpdatedAt,
	); err != nil {
		return err
	}

	if err := deductFEFO(ctx, tx, orgID, dr.LocationID, dr.ProductID, dr.QuantityDispensed, "dispense", "dispensing", dr.ID, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *DispensingRepoImpl) Update(ctx context.Context, id string, dr models.DispenseRecord) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		`UPDATE dispensing_records SET technician_id=NULLIF(?,''), witness_name=NULLIF(?,''),
		        notes=NULLIF(?,''), updated_at=? WHERE id=? AND organisation_id=?`,
		dr.TechnicianID, dr.WitnessName, dr.Notes, dr.UpdatedAt, id, orgID,
	)
	return err
}

func (r *DispensingRepoImpl) UpdateStatus(ctx context.Context, id, status string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		"UPDATE dispensing_records SET status=?, updated_at=NOW() WHERE id=? AND organisation_id=?", status, id, orgID,
	)
	return err
}

var _ DispensingRepository = (*DispensingRepoImpl)(nil)
