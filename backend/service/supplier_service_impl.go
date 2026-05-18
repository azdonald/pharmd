package service

import (
	"context"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
	"github.com/google/uuid"
)

type SupplierService struct {
	supplierRepo repository.SupplierRepository
}

func NewSupplierService(supplierRepo repository.SupplierRepository) SupplierServiceManager {
	return &SupplierService{supplierRepo: supplierRepo}
}

func (s *SupplierService) ListSuppliers(ctx context.Context, page, limit int, query string) ([]models.Supplier, int, error) {
	return s.supplierRepo.ListSuppliers(ctx, page, limit, query)
}

func (s *SupplierService) GetSupplierByID(ctx context.Context, id string) (*models.Supplier, error) {
	return s.supplierRepo.GetSupplierByID(ctx, id)
}

func (s *SupplierService) CreateSupplier(ctx context.Context, supplier models.Supplier) (*models.Supplier, error) {
	orgID := ctx.Value("organisation_id").(string)
	now := time.Now()

	supplier.ID = uuid.New().String()
	supplier.OrganisationID = orgID
	supplier.IsActive = true
	supplier.CreatedAt = now
	supplier.UpdatedAt = now

	if err := s.supplierRepo.CreateSupplier(ctx, supplier); err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (s *SupplierService) UpdateSupplier(ctx context.Context, id string, supplier models.Supplier) (*models.Supplier, error) {
	existing, err := s.supplierRepo.GetSupplierByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if supplier.Name != "" {
		existing.Name = supplier.Name
	}
	if supplier.ContactPerson != "" {
		existing.ContactPerson = supplier.ContactPerson
	}
	if supplier.Phone != "" {
		existing.Phone = supplier.Phone
	}
	if supplier.Email != "" {
		existing.Email = supplier.Email
	}
	if supplier.Address != "" {
		existing.Address = supplier.Address
	}
	if supplier.City != "" {
		existing.City = supplier.City
	}
	if supplier.State != "" {
		existing.State = supplier.State
	}
	if supplier.Country != "" {
		existing.Country = supplier.Country
	}
	if supplier.PaymentTerms != "" {
		existing.PaymentTerms = supplier.PaymentTerms
	}
	if supplier.Notes != "" {
		existing.Notes = supplier.Notes
	}

	if err := s.supplierRepo.UpdateSupplier(ctx, id, *existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *SupplierService) DeleteSupplier(ctx context.Context, id string) error {
	return s.supplierRepo.DeleteSupplier(ctx, id)
}

func (s *SupplierService) ListSupplierProducts(ctx context.Context, supplierID string) ([]models.SupplierProduct, error) {
	return s.supplierRepo.ListSupplierProducts(ctx, supplierID)
}

func (s *SupplierService) SetSupplierProducts(ctx context.Context, supplierID string, products []models.SupplierProduct) ([]models.SupplierProduct, error) {
	orgID := ctx.Value("organisation_id").(string)
	now := time.Now()

	for i := range products {
		products[i].ID = uuid.New().String()
		products[i].OrganisationID = orgID
		products[i].SupplierID = supplierID
		products[i].CreatedAt = now
		products[i].UpdatedAt = now
	}

	if err := s.supplierRepo.SetSupplierProducts(ctx, supplierID, products); err != nil {
		return nil, err
	}
	return products, nil
}

var _ SupplierServiceManager = (*SupplierService)(nil)
