package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/azdonald/pharmd/backend/models"
)

type UserRepoImpl struct {
	db *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) UserRepository {
	return &UserRepoImpl{db: db}
}

func (r *UserRepoImpl) ListUsers(ctx context.Context, page, limit int) ([]models.User, error) {
	offset := (page - 1) * limit
	orgID := ctx.Value("organisation_id").(string)

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, first_name, last_name, email, organisation_id, location_id, is_active, created_at, updated_at
		 FROM users WHERE organisation_id = ? AND deleted_at IS NULL
		 ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var locID sql.NullString
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email,
			&u.OrganisationID, &locID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		u.LocationID = locID.String
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepoImpl) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	orgID := ctx.Value("organisation_id").(string)
	user := &models.User{}
	var locID sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, first_name, last_name, email, password, organisation_id, location_id, is_active, created_at, updated_at
		 FROM users WHERE id = ? AND organisation_id = ? AND deleted_at IS NULL`, id, orgID,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password,
		&user.OrganisationID, &locID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	user.LocationID = locID.String
	return user, nil
}

func (r *UserRepoImpl) CreateUser(ctx context.Context, user models.User) error {
	orgID := ctx.Value("organisation_id").(string)
	user.OrganisationID = orgID

	var locID interface{}
	if user.LocationID != "" {
		locID = user.LocationID
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, first_name, last_name, email, password, organisation_id, location_id, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.FirstName, user.LastName, user.Email, user.Password,
		user.OrganisationID, locID, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepoImpl) UpdateUser(ctx context.Context, id string, user models.User) error {
	orgID := ctx.Value("organisation_id").(string)
	query := "UPDATE users SET updated_at = NOW()"
	args := []interface{}{}

	if user.FirstName != "" {
		query += ", first_name = ?"
		args = append(args, user.FirstName)
	}
	if user.LastName != "" {
		query += ", last_name = ?"
		args = append(args, user.LastName)
	}
	if user.LocationID != "" {
		query += ", location_id = ?"
		args = append(args, user.LocationID)
	}
	query += fmt.Sprintf(" WHERE id = ? AND organisation_id = ?")
	args = append(args, id, orgID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *UserRepoImpl) DeleteUser(ctx context.Context, id string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET deleted_at = NOW() WHERE id = ? AND organisation_id = ?", id, orgID,
	)
	return err
}
