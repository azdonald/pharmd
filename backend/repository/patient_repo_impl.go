package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/azdonald/pharmd/backend/models"
)

type PatientRepoImpl struct {
	db *sql.DB
}

func NewPatientRepositoryImpl(db *sql.DB) PatientRepository {
	return &PatientRepoImpl{db: db}
}

func (r *PatientRepoImpl) ListPatients(ctx context.Context, page, limit int, query string) ([]models.Patient, int, error) {
	orgID := ctx.Value("organisation_id").(string)
	offset := (page - 1) * limit

	where := "organisation_id = ? AND deleted_at IS NULL"
	args := []interface{}{orgID}

	if query != "" {
		where += " AND (first_name LIKE ? OR last_name LIKE ? OR phone LIKE ? OR email LIKE ?)"
		q := "%" + query + "%"
		args = append(args, q, q, q, q)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM patients WHERE " + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT id, organisation_id, first_name, last_name, date_of_birth, gender,
		       phone, email, address, city, state, country, blood_group, genotype, notes,
		       emergency_contact_name, emergency_contact_phone, is_active, created_at, updated_at
		 FROM patients WHERE %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, where),
		append(args, limit, offset)...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var patients []models.Patient
	for rows.Next() {
		var p models.Patient
		if err := rows.Scan(&p.ID, &p.OrganisationID, &p.FirstName, &p.LastName, &p.DateOfBirth,
			&p.Gender, &p.Phone, &p.Email, &p.Address, &p.City, &p.State, &p.Country,
			&p.BloodGroup, &p.Genotype, &p.Notes, &p.EmergencyContactName,
			&p.EmergencyContactPhone, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		patients = append(patients, p)
	}
	return patients, total, nil
}

func (r *PatientRepoImpl) GetPatientByID(ctx context.Context, id string) (*models.Patient, error) {
	orgID := ctx.Value("organisation_id").(string)
	p := &models.Patient{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, organisation_id, first_name, last_name, date_of_birth, gender,
		        phone, email, address, city, state, country, blood_group, genotype, notes,
		        emergency_contact_name, emergency_contact_phone, is_active, created_at, updated_at
		 FROM patients WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&p.ID, &p.OrganisationID, &p.FirstName, &p.LastName, &p.DateOfBirth,
		&p.Gender, &p.Phone, &p.Email, &p.Address, &p.City, &p.State, &p.Country,
		&p.BloodGroup, &p.Genotype, &p.Notes, &p.EmergencyContactName,
		&p.EmergencyContactPhone, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PatientRepoImpl) CreatePatient(ctx context.Context, patient models.Patient) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO patients (id, organisation_id, first_name, last_name, date_of_birth, gender,
		                       phone, email, address, city, state, country, blood_group, genotype,
		                       notes, emergency_contact_name, emergency_contact_phone,
		                       is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		patient.ID, patient.OrganisationID, patient.FirstName, patient.LastName,
		patient.DateOfBirth, patient.Gender, patient.Phone, patient.Email,
		patient.Address, patient.City, patient.State, patient.Country,
		patient.BloodGroup, patient.Genotype, patient.Notes,
		patient.EmergencyContactName, patient.EmergencyContactPhone,
		patient.IsActive, patient.CreatedAt, patient.UpdatedAt,
	)
	return err
}

func (r *PatientRepoImpl) UpdatePatient(ctx context.Context, id string, patient models.Patient) error {
	orgID := ctx.Value("organisation_id").(string)
	query := "UPDATE patients SET updated_at = ?"
	args := []interface{}{time.Now()}

	addField(&query, &args, "first_name", patient.FirstName)
	addField(&query, &args, "last_name", patient.LastName)
	addField(&query, &args, "date_of_birth", patient.DateOfBirth)
	addField(&query, &args, "gender", patient.Gender)
	addField(&query, &args, "phone", patient.Phone)
	addField(&query, &args, "email", patient.Email)
	addField(&query, &args, "address", patient.Address)
	addField(&query, &args, "city", patient.City)
	addField(&query, &args, "state", patient.State)
	addField(&query, &args, "country", patient.Country)
	addField(&query, &args, "blood_group", patient.BloodGroup)
	addField(&query, &args, "genotype", patient.Genotype)
	addField(&query, &args, "notes", patient.Notes)
	addField(&query, &args, "emergency_contact_name", patient.EmergencyContactName)
	addField(&query, &args, "emergency_contact_phone", patient.EmergencyContactPhone)
	if !patient.IsActive {
		query += ", is_active = ?"
		args = append(args, patient.IsActive)
	}

	query += fmt.Sprintf(" WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL")
	args = append(args, id, orgID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PatientRepoImpl) DeletePatient(ctx context.Context, id string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		"UPDATE patients SET deleted_at = NOW() WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL",
		id, orgID,
	)
	return err
}

func (r *PatientRepoImpl) ListPatientAllergies(ctx context.Context, patientID string) ([]models.PatientAllergy, error) {
	orgID := ctx.Value("organisation_id").(string)
	rows, err := r.db.QueryContext(ctx,
		`SELECT a.id, a.patient_id, a.allergy, a.severity, a.notes, a.created_at
		 FROM patient_allergies a
		 JOIN patients p ON p.id = a.patient_id
		 WHERE a.patient_id = ? AND p.organisation_id = ? AND p.deleted_at IS NULL
		 ORDER BY a.created_at DESC`,
		patientID, orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allergies []models.PatientAllergy
	for rows.Next() {
		var a models.PatientAllergy
		if err := rows.Scan(&a.ID, &a.PatientID, &a.Allergy, &a.Severity, &a.Notes, &a.CreatedAt); err != nil {
			return nil, err
		}
		allergies = append(allergies, a)
	}
	return allergies, nil
}

func (r *PatientRepoImpl) AddPatientAllergy(ctx context.Context, allergy models.PatientAllergy) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO patient_allergies (id, patient_id, allergy, severity, notes, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		allergy.ID, allergy.PatientID, allergy.Allergy, allergy.Severity, allergy.Notes,
		allergy.CreatedAt, allergy.CreatedAt,
	)
	return err
}

func (r *PatientRepoImpl) RemovePatientAllergy(ctx context.Context, patientID, allergyID string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		`DELETE a FROM patient_allergies a
		 JOIN patients p ON p.id = a.patient_id
		 WHERE a.id = ? AND a.patient_id = ? AND p.organisation_id = ?`,
		allergyID, patientID, orgID,
	)
	return err
}

func (r *PatientRepoImpl) ListPatientConditions(ctx context.Context, patientID string) ([]models.PatientCondition, error) {
	orgID := ctx.Value("organisation_id").(string)
	rows, err := r.db.QueryContext(ctx,
		`SELECT c.id, c.patient_id, c.condition, c.notes, c.created_at
		 FROM patient_conditions c
		 JOIN patients p ON p.id = c.patient_id
		 WHERE c.patient_id = ? AND p.organisation_id = ? AND p.deleted_at IS NULL
		 ORDER BY c.created_at DESC`,
		patientID, orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conditions []models.PatientCondition
	for rows.Next() {
		var c models.PatientCondition
		if err := rows.Scan(&c.ID, &c.PatientID, &c.Condition, &c.Notes, &c.CreatedAt); err != nil {
			return nil, err
		}
		conditions = append(conditions, c)
	}
	return conditions, nil
}

func (r *PatientRepoImpl) AddPatientCondition(ctx context.Context, condition models.PatientCondition) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO patient_conditions (id, patient_id, condition, notes, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		condition.ID, condition.PatientID, condition.Condition, condition.Notes,
		condition.CreatedAt, condition.CreatedAt,
	)
	return err
}

func addField(query *string, args *[]interface{}, field, value string) {
	if value != "" {
		*query += ", " + field + " = ?"
		*args = append(*args, value)
	}
}

var _ PatientRepository = (*PatientRepoImpl)(nil)
