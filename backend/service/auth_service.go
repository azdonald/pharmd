package service

import (
	"context"

	"github.com/azdonald/pharmd/backend/models"
)

type AuthServiceManager interface {
	Register(ctx context.Context, org models.Organisation, user models.User) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetOrganisationByID(ctx context.Context, id string) (*models.Organisation, error)
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
}

type UserServiceManager interface {
	ListUsers(ctx context.Context, page, limit int) ([]models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	UpdateUser(ctx context.Context, id string, user models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserRoleServiceManager interface {
	GetUserPermissions(ctx context.Context, userID, orgID string) ([]string, error)
	AssignRoleToUser(ctx context.Context, userID, roleID, orgID string) error
}
