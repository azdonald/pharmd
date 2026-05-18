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

type RoleServiceManager interface {
	ListRoles(ctx context.Context, page, limit int) ([]models.Role, error)
	GetRoleByID(ctx context.Context, id string) (*models.Role, error)
	CreateRole(ctx context.Context, role models.Role) (*models.Role, error)
	UpdateRole(ctx context.Context, id string, role models.Role) (*models.Role, error)
	DeleteRole(ctx context.Context, id string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]string, error)
	SetRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error
}

type PermissionServiceManager interface {
	ListPermissions(ctx context.Context) ([]models.Permission, error)
}

type LocationServiceManager interface {
	ListLocations(ctx context.Context, page, limit int) ([]models.Location, error)
	GetLocationByID(ctx context.Context, id string) (*models.Location, error)
	CreateLocation(ctx context.Context, location models.Location) (*models.Location, error)
	UpdateLocation(ctx context.Context, id string, location models.Location) (*models.Location, error)
	DeleteLocation(ctx context.Context, id string) error
}

type OrganisationServiceManager interface {
	GetOrganisationByID(ctx context.Context, id string) (*models.Organisation, error)
	UpdateOrganisation(ctx context.Context, id string, org models.Organisation) error
}
