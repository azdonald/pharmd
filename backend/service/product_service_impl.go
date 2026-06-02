package service

import (
	"context"
	"strconv"
	"strings"
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

func (s *ProductService) ImportProductsCSV(ctx context.Context, orgID string, records [][]string) (int, int, []string) {
	var products []models.Product
	var errors []string
	now := time.Now()

	for i, row := range records {
		if i == 0 && strings.EqualFold(row[0], "name") {
			continue
		}
		if len(row) < 1 || strings.TrimSpace(row[0]) == "" {
			continue
		}
		p := models.Product{
			ID:             uuid.New().String(),
			OrganisationID: orgID,
			Name:           strings.TrimSpace(row[0]),
			IsActive:       true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if len(row) > 1 {
			p.Description = strings.TrimSpace(row[1])
		}
		if len(row) > 2 {
			p.CategoryID = strings.TrimSpace(row[2])
		}
		if len(row) > 3 {
			p.Classification = strings.TrimSpace(row[3])
		}
		if len(row) > 4 {
			p.BrandName = strings.TrimSpace(row[4])
		}
		if len(row) > 5 {
			p.GenericName = strings.TrimSpace(row[5])
		}
		if len(row) > 6 {
			p.Manufacturer = strings.TrimSpace(row[6])
		}
		if len(row) > 7 {
			p.Barcode = strings.TrimSpace(row[7])
		}
		if len(row) > 8 {
			p.NDC = strings.TrimSpace(row[8])
		}
		if len(row) > 9 {
			p.UnitOfMeasure = strings.TrimSpace(row[9])
		}
		if len(row) > 10 {
			p.Strength = strings.TrimSpace(row[10])
		}
		if len(row) > 11 {
			p.Form = strings.TrimSpace(row[11])
		}
		if len(row) > 12 {
			if v, err := strconv.Atoi(strings.TrimSpace(row[12])); err == nil {
				p.ReorderLevel = v
			}
		}
		products = append(products, p)
	}

	if len(products) == 0 {
		return 0, 0, []string{"no valid products found"}
	}

	if err := s.productRepo.BulkCreateProducts(ctx, products); err != nil {
		return 0, len(products), []string{err.Error()}
	}

	return len(products), 0, errors
}

var _ ProductServiceManager = (*ProductService)(nil)
