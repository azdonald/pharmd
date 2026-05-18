package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/azdonald/pharmd/backend/models"
)

type LocationRepoImpl struct {
	db *sql.DB
}

func NewLocationRepositoryImpl(db *sql.DB) LocationRepository {
	return &LocationRepoImpl{db: db}
}

func (r *LocationRepoImpl) ListLocations(ctx context.Context, page, limit int) ([]models.Location, error) {
	orgID := ctx.Value("organisation_id").(string)
	offset := (page - 1) * limit

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, organisation_id, name, address, city, state, country, phone, email,
		        tax_rate, timezone, is_active, created_at, updated_at
		 FROM locations WHERE organisation_id = ? AND deleted_at IS NULL
		 ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var l models.Location
		if err := rows.Scan(&l.ID, &l.OrganisationID, &l.Name, &l.Address, &l.City, &l.State,
			&l.Country, &l.Phone, &l.Email, &l.TaxRate, &l.Timezone, &l.IsActive,
			&l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	return locations, nil
}

func (r *LocationRepoImpl) GetLocationByID(ctx context.Context, id string) (*models.Location, error) {
	orgID := ctx.Value("organisation_id").(string)
	l := &models.Location{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, organisation_id, name, address, city, state, country, phone, email,
		        tax_rate, timezone, is_active, created_at, updated_at
		 FROM locations WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&l.ID, &l.OrganisationID, &l.Name, &l.Address, &l.City, &l.State,
		&l.Country, &l.Phone, &l.Email, &l.TaxRate, &l.Timezone, &l.IsActive,
		&l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (r *LocationRepoImpl) CreateLocation(ctx context.Context, location models.Location) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO locations (id, organisation_id, name, address, city, state, country, phone, email,
		                        tax_rate, timezone, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		location.ID, location.OrganisationID, location.Name, location.Address, location.City,
		location.State, location.Country, location.Phone, location.Email,
		location.TaxRate, location.Timezone, location.IsActive,
		location.CreatedAt, location.UpdatedAt,
	)
	return err
}

func (r *LocationRepoImpl) UpdateLocation(ctx context.Context, id string, location models.Location) error {
	orgID := ctx.Value("organisation_id").(string)
	query := "UPDATE locations SET updated_at = ?"
	args := []interface{}{time.Now()}

	if location.Name != "" {
		query += ", name = ?"
		args = append(args, location.Name)
	}
	if location.Address != "" {
		query += ", address = ?"
		args = append(args, location.Address)
	}
	if location.City != "" {
		query += ", city = ?"
		args = append(args, location.City)
	}
	if location.State != "" {
		query += ", state = ?"
		args = append(args, location.State)
	}
	if location.Country != "" {
		query += ", country = ?"
		args = append(args, location.Country)
	}
	if location.Phone != "" {
		query += ", phone = ?"
		args = append(args, location.Phone)
	}
	if location.Email != "" {
		query += ", email = ?"
		args = append(args, location.Email)
	}
	if location.TaxRate != 0 {
		query += ", tax_rate = ?"
		args = append(args, location.TaxRate)
	}
	if location.Timezone != "" {
		query += ", timezone = ?"
		args = append(args, location.Timezone)
	}

	query += fmt.Sprintf(" WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL")
	args = append(args, id, orgID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *LocationRepoImpl) DeleteLocation(ctx context.Context, id string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		"UPDATE locations SET deleted_at = NOW() WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL",
		id, orgID,
	)
	return err
}
