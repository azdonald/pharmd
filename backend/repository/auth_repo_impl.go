package repository

import (
	"context"
	"database/sql"

	"github.com/azdonald/pharmd/backend/models"
)

type AuthRepoImpl struct {
	db *sql.DB
}

func NewAuthRepositoryImpl(db *sql.DB) AuthRepository {
	return &AuthRepoImpl{db: db}
}

func (r *AuthRepoImpl) CreateOrganisation(ctx context.Context, org models.Organisation) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO organisations (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)",
		org.ID, org.Name, org.CreatedAt, org.UpdatedAt,
	)
	return err
}

func (r *AuthRepoImpl) CreateUser(ctx context.Context, user models.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, first_name, last_name, email, password, organisation_id, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.FirstName, user.LastName, user.Email, user.Password,
		user.OrganisationID, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *AuthRepoImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, first_name, last_name, email, password, organisation_id, is_active, created_at, updated_at
		 FROM users WHERE email = ? AND deleted_at IS NULL`, email,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password,
		&user.OrganisationID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *AuthRepoImpl) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, first_name, last_name, email, password, organisation_id, is_active, created_at, updated_at
		 FROM users WHERE id = ? AND deleted_at IS NULL`, id,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password,
		&user.OrganisationID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *AuthRepoImpl) GetOrganisationByID(ctx context.Context, id string) (*models.Organisation, error) {
	org := &models.Organisation{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, created_at, updated_at FROM organisations WHERE id = ? AND deleted_at IS NULL", id,
	).Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (r *AuthRepoImpl) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET password = ?, updated_at = NOW() WHERE id = ?",
		hashedPassword, userID,
	)
	return err
}
