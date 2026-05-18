package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/azdonald/pharmd/backend/models"
)

type PrescriberRepository interface {
	List(ctx context.Context, query string, page, limit int) ([]models.Prescriber, int, error)
	GetByID(ctx context.Context, id string) (*models.Prescriber, error)
	Create(ctx context.Context, p models.Prescriber) error
	Update(ctx context.Context, id string, p models.Prescriber) error
	Delete(ctx context.Context, id string) error
}

type PrescriberRepoImpl struct {
	db *sql.DB
}

func NewPrescriberRepositoryImpl(db *sql.DB) PrescriberRepository {
	return &PrescriberRepoImpl{db: db}
}

func (r *PrescriberRepoImpl) List(ctx context.Context, query string, page, limit int) ([]models.Prescriber, int, error) {
	orgID := ctx.Value("organisation_id").(string)
	offset := (page - 1) * limit

	where := "organisation_id = ? AND deleted_at IS NULL"
	args := []interface{}{orgID}
	if query != "" {
		where += " AND (name LIKE ? OR license_number LIKE ? OR specialty LIKE ?)"
		q := "%" + query + "%"
		args = append(args, q, q, q)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM prescribers WHERE %s", where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf("SELECT id, organisation_id, name, COALESCE(license_number,''), COALESCE(phone,''), COALESCE(email,''), COALESCE(specialty,''), COALESCE(dea_number,''), COALESCE(npi_number,''), COALESCE(address,''), is_active, created_at, updated_at FROM prescribers WHERE %s ORDER BY name ASC LIMIT ? OFFSET ?", where),
		append(args, limit, offset)...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var prescribers []models.Prescriber
	for rows.Next() {
		var p models.Prescriber
		if err := rows.Scan(&p.ID, &p.OrganisationID, &p.Name, &p.LicenseNumber, &p.Phone, &p.Email, &p.Specialty, &p.DEANumber, &p.NPINumber, &p.Address, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		prescribers = append(prescribers, p)
	}
	return prescribers, total, nil
}

func (r *PrescriberRepoImpl) GetByID(ctx context.Context, id string) (*models.Prescriber, error) {
	orgID := ctx.Value("organisation_id").(string)
	p := &models.Prescriber{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, organisation_id, name, COALESCE(license_number,''), COALESCE(phone,''), COALESCE(email,''), COALESCE(specialty,''), COALESCE(dea_number,''), COALESCE(npi_number,''), COALESCE(address,''), is_active, created_at, updated_at
		 FROM prescribers WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&p.ID, &p.OrganisationID, &p.Name, &p.LicenseNumber, &p.Phone, &p.Email, &p.Specialty, &p.DEANumber, &p.NPINumber, &p.Address, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PrescriberRepoImpl) Create(ctx context.Context, p models.Prescriber) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO prescribers (id, organisation_id, name, license_number, phone, email, specialty, dea_number, npi_number, address, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, NULLIF(?,''), NULLIF(?,''), NULLIF(?,''), NULLIF(?,''), NULLIF(?,''), NULLIF(?,''), NULLIF(?,''), ?, ?, ?)`,
		p.ID, orgID, p.Name, p.LicenseNumber, p.Phone, p.Email, p.Specialty, p.DEANumber, p.NPINumber, p.Address, p.IsActive, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *PrescriberRepoImpl) Update(ctx context.Context, id string, p models.Prescriber) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		`UPDATE prescribers SET name=?, license_number=NULLIF(?,''), phone=NULLIF(?,''), email=NULLIF(?,''), specialty=NULLIF(?,''), dea_number=NULLIF(?,''), npi_number=NULLIF(?,''), address=NULLIF(?,''), is_active=?, updated_at=? WHERE id=? AND organisation_id=? AND deleted_at IS NULL`,
		p.Name, p.LicenseNumber, p.Phone, p.Email, p.Specialty, p.DEANumber, p.NPINumber, p.Address, p.IsActive, p.UpdatedAt, id, orgID,
	)
	return err
}

func (r *PrescriberRepoImpl) Delete(ctx context.Context, id string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx, "UPDATE prescribers SET deleted_at = NOW() WHERE id = ? AND organisation_id = ?", id, orgID)
	return err
}

var _ PrescriberRepository = (*PrescriberRepoImpl)(nil)
