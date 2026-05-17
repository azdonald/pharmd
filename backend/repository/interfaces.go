package repository

import (
	"context"

	"github.com/azdonald/pharmd/backend/models"
)

type AuthRepository interface {
	CreateOrganisation(ctx context.Context, org models.Organisation) error
	CreateUser(ctx context.Context, user models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetOrganisationByID(ctx context.Context, id string) (*models.Organisation, error)
	UpdatePassword(ctx context.Context, userID, hashedPassword string) error
}

type UserRepository interface {
	ListUsers(ctx context.Context, page, limit int) ([]models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, user models.User) error
	UpdateUser(ctx context.Context, id string, user models.User) error
	DeleteUser(ctx context.Context, id string) error
}

type UserRoleRepository interface {
	GetUserPermissions(ctx context.Context, userID, orgID string) ([]string, error)
	AssignRoleToUser(ctx context.Context, userID, roleID, orgID string) error
}

type RoleRepository interface {
	// Placeholder for future role operations
}
