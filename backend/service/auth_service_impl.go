package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo     repository.AuthRepository
	userRepo     repository.UserRepository
	locationRepo repository.LocationRepository
	userRoleRepo repository.UserRoleRepository
}

func NewAuthService(
	authRepo repository.AuthRepository,
	userRepo repository.UserRepository,
	locationRepo repository.LocationRepository,
	userRoleRepo repository.UserRoleRepository,
) AuthServiceManager {
	return &AuthService{
		authRepo:     authRepo,
		userRepo:     userRepo,
		locationRepo: locationRepo,
		userRoleRepo: userRoleRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, org models.Organisation, user models.User) (*models.User, error) {
	org.ID = ulid.Make().String()
	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()

	if err := s.authRepo.CreateOrganisation(ctx, org); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.ID = ulid.Make().String()
	user.Password = hashedPassword
	user.OrganisationID = org.ID
	user.IsActive = true
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	defaultLocation := models.Location{
		ID:             ulid.Make().String(),
		OrganisationID: org.ID,
		Name:           org.Name + " - Main",
		Address:        "",
		City:           "",
		State:          "",
		Country:        "",
		Phone:          "",
		Email:          "",
		TaxRate:        0,
		Timezone:       "UTC",
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.locationRepo.CreateLocation(ctx, defaultLocation); err != nil {
		return nil, err
	}

	if err := s.userRoleRepo.AssignRoleToUser(ctx, user.ID, "R001", org.ID); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.authRepo.GetUserByID(ctx, id)
}

func (s *AuthService) GetOrganisationByID(ctx context.Context, id string) (*models.Organisation, error) {
	return s.authRepo.GetOrganisationByID(ctx, id)
}

func (s *AuthService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(oldPassword)); err != nil {
		return errors.New("incorrect old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.authRepo.UpdatePassword(ctx, userID, string(hashedPassword))
}

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserServiceManager {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) ListUsers(ctx context.Context, page, limit int) ([]models.User, error) {
	return s.userRepo.ListUsers(ctx, page, limit)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	user.ID = ulid.Make().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	hashedPassword, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user models.User) (*models.User, error) {
	user.UpdatedAt = time.Now()
	if err := s.userRepo.UpdateUser(ctx, id, user); err != nil {
		return nil, err
	}
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.DeleteUser(ctx, id)
}

type UserRoleService struct {
	userRoleRepo repository.UserRoleRepository
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
}

func NewUserRoleService(userRoleRepo repository.UserRoleRepository, userRepo repository.UserRepository, roleRepo repository.RoleRepository) UserRoleServiceManager {
	return &UserRoleService{
		userRoleRepo: userRoleRepo,
		userRepo:     userRepo,
		roleRepo:     roleRepo,
	}
}

func (s *UserRoleService) GetUserPermissions(ctx context.Context, userID, orgID string) ([]string, error) {
	permissions, err := s.userRoleRepo.GetUserPermissions(ctx, userID, orgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []string{}, nil
		}
		return nil, err
	}
	return permissions, nil
}

func (s *UserRoleService) AssignRoleToUser(ctx context.Context, userID, roleID, orgID string) error {
	return s.userRoleRepo.AssignRoleToUser(ctx, userID, roleID, orgID)
}
