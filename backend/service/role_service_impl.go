package service

import (
	"context"
	"strings"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
	"github.com/oklog/ulid/v2"
)

type RoleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) RoleServiceManager {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) ListRoles(ctx context.Context, page, limit int) ([]models.Role, error) {
	return s.roleRepo.ListRoles(ctx, page, limit)
}

func (s *RoleService) GetRoleByID(ctx context.Context, id string) (*models.Role, error) {
	return s.roleRepo.GetRoleByID(ctx, id)
}

func (s *RoleService) CreateRole(ctx context.Context, role models.Role) (*models.Role, error) {
	orgID := ctx.Value("organisation_id").(string)
	role.ID = ulid.Make().String()
	role.Slug = strings.ToLower(strings.ReplaceAll(role.Name, " ", "_"))
	role.OrganisationID = orgID
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	if err := s.roleRepo.CreateRole(ctx, role); err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, id string, role models.Role) (*models.Role, error) {
	if err := s.roleRepo.UpdateRole(ctx, id, role); err != nil {
		return nil, err
	}
	return s.roleRepo.GetRoleByID(ctx, id)
}

func (s *RoleService) DeleteRole(ctx context.Context, id string) error {
	return s.roleRepo.DeleteRole(ctx, id)
}

func (s *RoleService) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	return s.roleRepo.GetRolePermissions(ctx, roleID)
}

func (s *RoleService) SetRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	_, err := s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	return s.roleRepo.SetRolePermissions(ctx, roleID, permissionIDs)
}

type PermissionService struct {
	permRepo repository.PermissionRepository
}

func NewPermissionService(permRepo repository.PermissionRepository) PermissionServiceManager {
	return &PermissionService{permRepo: permRepo}
}

func (s *PermissionService) ListPermissions(ctx context.Context) ([]models.Permission, error) {
	return s.permRepo.ListPermissions(ctx)
}

type LocationService struct {
	locationRepo repository.LocationRepository
}

func NewLocationService(locationRepo repository.LocationRepository) LocationServiceManager {
	return &LocationService{locationRepo: locationRepo}
}

func (s *LocationService) ListLocations(ctx context.Context, page, limit int) ([]models.Location, error) {
	return s.locationRepo.ListLocations(ctx, page, limit)
}

func (s *LocationService) GetLocationByID(ctx context.Context, id string) (*models.Location, error) {
	return s.locationRepo.GetLocationByID(ctx, id)
}

func (s *LocationService) CreateLocation(ctx context.Context, location models.Location) (*models.Location, error) {
	orgID := ctx.Value("organisation_id").(string)
	location.ID = ulid.Make().String()
	location.OrganisationID = orgID
	location.IsActive = true
	location.CreatedAt = time.Now()
	location.UpdatedAt = time.Now()

	if err := s.locationRepo.CreateLocation(ctx, location); err != nil {
		return nil, err
	}
	return &location, nil
}

func (s *LocationService) UpdateLocation(ctx context.Context, id string, location models.Location) (*models.Location, error) {
	existing, err := s.locationRepo.GetLocationByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if location.Name != "" {
		existing.Name = location.Name
	}
	if location.Address != "" {
		existing.Address = location.Address
	}
	if location.City != "" {
		existing.City = location.City
	}
	if location.State != "" {
		existing.State = location.State
	}
	if location.Country != "" {
		existing.Country = location.Country
	}
	if location.Phone != "" {
		existing.Phone = location.Phone
	}
	if location.Email != "" {
		existing.Email = location.Email
	}
	if location.TaxRate != 0 {
		existing.TaxRate = location.TaxRate
	}
	if location.Timezone != "" {
		existing.Timezone = location.Timezone
	}

	if err := s.locationRepo.UpdateLocation(ctx, id, *existing); err != nil {
		return nil, err
	}
	return s.locationRepo.GetLocationByID(ctx, id)
}

func (s *LocationService) DeleteLocation(ctx context.Context, id string) error {
	return s.locationRepo.DeleteLocation(ctx, id)
}

// Ensure locationRepo satisfies the interface
var _ LocationServiceManager = (*LocationService)(nil)
