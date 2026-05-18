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
	ListRoles(ctx context.Context, page, limit int) ([]models.Role, error)
	GetRoleByID(ctx context.Context, id string) (*models.Role, error)
	CreateRole(ctx context.Context, role models.Role) error
	UpdateRole(ctx context.Context, id string, role models.Role) error
	DeleteRole(ctx context.Context, id string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]string, error)
	SetRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error
}

type PermissionRepository interface {
	ListPermissions(ctx context.Context) ([]models.Permission, error)
}

type LocationRepository interface {
	ListLocations(ctx context.Context, page, limit int) ([]models.Location, error)
	GetLocationByID(ctx context.Context, id string) (*models.Location, error)
	CreateLocation(ctx context.Context, location models.Location) error
	UpdateLocation(ctx context.Context, id string, location models.Location) error
	DeleteLocation(ctx context.Context, id string) error
}

type PatientRepository interface {
	ListPatients(ctx context.Context, page, limit int, query string) ([]models.Patient, int, error)
	GetPatientByID(ctx context.Context, id string) (*models.Patient, error)
	CreatePatient(ctx context.Context, patient models.Patient) error
	UpdatePatient(ctx context.Context, id string, patient models.Patient) error
	DeletePatient(ctx context.Context, id string) error
	ListPatientAllergies(ctx context.Context, patientID string) ([]models.PatientAllergy, error)
	AddPatientAllergy(ctx context.Context, allergy models.PatientAllergy) error
	RemovePatientAllergy(ctx context.Context, patientID, allergyID string) error
	ListPatientConditions(ctx context.Context, patientID string) ([]models.PatientCondition, error)
	AddPatientCondition(ctx context.Context, condition models.PatientCondition) error
}
