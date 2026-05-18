package service

import (
	"context"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
	"github.com/google/uuid"
)

type ProductCategoryService struct {
	categoryRepo repository.ProductCategoryRepository
}

func NewProductCategoryService(categoryRepo repository.ProductCategoryRepository) ProductCategoryServiceManager {
	return &ProductCategoryService{categoryRepo: categoryRepo}
}

func (s *ProductCategoryService) ListCategories(ctx context.Context) ([]models.ProductCategory, error) {
	return s.categoryRepo.ListCategories(ctx)
}

func (s *ProductCategoryService) GetCategoryByID(ctx context.Context, id string) (*models.ProductCategory, error) {
	return s.categoryRepo.GetCategoryByID(ctx, id)
}

func (s *ProductCategoryService) CreateCategory(ctx context.Context, category models.ProductCategory) (*models.ProductCategory, error) {
	orgID := ctx.Value("organisation_id").(string)
	now := time.Now()

	category.ID = uuid.New().String()
	category.OrganisationID = orgID
	category.IsActive = true
	category.CreatedAt = now
	category.UpdatedAt = now

	if err := s.categoryRepo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *ProductCategoryService) UpdateCategory(ctx context.Context, id string, category models.ProductCategory) (*models.ProductCategory, error) {
	existing, err := s.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if category.Name != "" {
		existing.Name = category.Name
	}
	if category.Description != "" {
		existing.Description = category.Description
	}
	if category.ParentID != "" {
		existing.ParentID = category.ParentID
	}

	if err := s.categoryRepo.UpdateCategory(ctx, id, *existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *ProductCategoryService) DeleteCategory(ctx context.Context, id string) error {
	return s.categoryRepo.DeleteCategory(ctx, id)
}

var _ ProductCategoryServiceManager = (*ProductCategoryService)(nil)
