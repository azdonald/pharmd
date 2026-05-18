package service

import (
	"context"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
)

type OrganisationService struct {
	authRepo repository.AuthRepository
}

func NewOrganisationService(authRepo repository.AuthRepository) OrganisationServiceManager {
	return &OrganisationService{authRepo: authRepo}
}

func (s *OrganisationService) GetOrganisationByID(ctx context.Context, id string) (*models.Organisation, error) {
	return s.authRepo.GetOrganisationByID(ctx, id)
}

func (s *OrganisationService) UpdateOrganisation(ctx context.Context, id string, org models.Organisation) error {
	// Stub for future implementation
	return nil
}
