package service

import (
	"context"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
	"github.com/google/uuid"
)

type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductServiceManager {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) ListProducts(ctx context.Context, page, limit int, query, categoryID string) ([]models.Product, int, error) {
	return s.productRepo.ListProducts(ctx, page, limit, query, categoryID)
}

func (s *ProductService) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	return s.productRepo.GetProductByID(ctx, id)
}

func (s *ProductService) GetProductByBarcode(ctx context.Context, barcode string) (*models.Product, error) {
	return s.productRepo.GetProductByBarcode(ctx, barcode)
}

func (s *ProductService) CreateProduct(ctx context.Context, product models.Product) (*models.Product, error) {
	orgID := ctx.Value("organisation_id").(string)
	now := time.Now()

	product.ID = uuid.New().String()
	product.OrganisationID = orgID
	product.IsActive = true
	product.CreatedAt = now
	product.UpdatedAt = now

	if err := s.productRepo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, product models.Product) (*models.Product, error) {
	existing, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if product.Name != "" {
		existing.Name = product.Name
	}
	if product.Description != "" {
		existing.Description = product.Description
	}
	if product.CategoryID != "" {
		existing.CategoryID = product.CategoryID
	}
	if product.Classification != "" {
		existing.Classification = product.Classification
	}
	if product.BrandName != "" {
		existing.BrandName = product.BrandName
	}
	if product.GenericName != "" {
		existing.GenericName = product.GenericName
	}
	if product.Manufacturer != "" {
		existing.Manufacturer = product.Manufacturer
	}
	if product.Barcode != "" {
		existing.Barcode = product.Barcode
	}
	if product.NDC != "" {
		existing.NDC = product.NDC
	}
	if product.UnitOfMeasure != "" {
		existing.UnitOfMeasure = product.UnitOfMeasure
	}
	if product.Strength != "" {
		existing.Strength = product.Strength
	}
	if product.Form != "" {
		existing.Form = product.Form
	}
	if product.ReorderLevel != 0 {
		existing.ReorderLevel = product.ReorderLevel
	}

	if err := s.productRepo.UpdateProduct(ctx, id, *existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	return s.productRepo.DeleteProduct(ctx, id)
}

func (s *ProductService) ListSubstitutes(ctx context.Context, productID string) ([]models.GenericSubstitution, error) {
	return s.productRepo.ListSubstitutes(ctx, productID)
}

func (s *ProductService) AddSubstitute(ctx context.Context, productID string, sub models.GenericSubstitution) (*models.GenericSubstitution, error) {
	orgID := ctx.Value("organisation_id").(string)
	sub.ID = uuid.New().String()
	sub.OrganisationID = orgID
	sub.ProductID = productID
	sub.CreatedAt = time.Now()

	if err := s.productRepo.AddSubstitute(ctx, sub); err != nil {
		return nil, err
	}
	return &sub, nil
}

func (s *ProductService) RemoveSubstitute(ctx context.Context, productID, substituteID string) error {
	return s.productRepo.RemoveSubstitute(ctx, productID, substituteID)
}

var _ ProductServiceManager = (*ProductService)(nil)
