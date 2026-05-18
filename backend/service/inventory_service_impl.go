package service

import (
	"context"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
	"github.com/google/uuid"
)

type InventoryService struct {
	inventoryRepo repository.InventoryRepository
}

func NewInventoryService(inventoryRepo repository.InventoryRepository) InventoryServiceManager {
	return &InventoryService{inventoryRepo: inventoryRepo}
}

func (s *InventoryService) CreateBatch(ctx context.Context, batch models.StockBatch) (*models.StockBatch, error) {
	orgID := ctx.Value("organisation_id").(string)
	userID := ctx.Value("user_id").(string)
	now := time.Now()

	batch.ID = uuid.New().String()
	batch.OrganisationID = orgID
	batch.RemainingQty = batch.Quantity
	batch.IsActive = true
	batch.ReceivedDate = now
	batch.CreatedAt = now
	batch.UpdatedAt = now

	if err := s.inventoryRepo.CreateBatch(ctx, batch); err != nil {
		return nil, err
	}

	// record movement
	movement := models.StockMovement{
		OrganisationID: orgID,
		LocationID:     batch.LocationID,
		ProductID:      batch.ProductID,
		BatchID:        batch.ID,
		MovementType:   "receipt",
		Quantity:       batch.Quantity,
		ReferenceType:  "batch",
		ReferenceID:    batch.ID,
		CreatedBy:      userID,
	}
	movement.ID = uuid.New().String()
	movement.CreatedAt = now
	if err := s.inventoryRepo.CreateMovement(ctx, movement); err != nil {
		return nil, err
	}

	return &batch, nil
}

func (s *InventoryService) ListStock(ctx context.Context, locationID string, page, limit int, query string) ([]repository.InventoryBatchView, int, error) {
	return s.inventoryRepo.ListStock(ctx, locationID, page, limit, query)
}

func (s *InventoryService) CreateAdjustment(ctx context.Context, movement models.StockMovement) (*models.StockMovement, error) {
	orgID := ctx.Value("organisation_id").(string)
	userID := ctx.Value("user_id").(string)
	now := time.Now()

	movement.ID = uuid.New().String()
	movement.OrganisationID = orgID
	movement.CreatedBy = userID
	movement.CreatedAt = now

	// validate and deduct from batch
	if movement.BatchID != "" {
		batch, err := s.inventoryRepo.GetBatchByID(ctx, movement.BatchID)
		if err != nil {
			return nil, err
		}
		newQty := batch.RemainingQty - movement.Quantity
		if newQty < 0 {
			newQty = 0
		}
		if err := s.inventoryRepo.UpdateBatchQty(ctx, movement.BatchID, newQty); err != nil {
			return nil, err
		}
	}

	if err := s.inventoryRepo.CreateMovement(ctx, movement); err != nil {
		return nil, err
	}

	return &movement, nil
}

func (s *InventoryService) ListAlerts(ctx context.Context, locationID string) ([]repository.InventoryAlertView, error) {
	return s.inventoryRepo.ListAlerts(ctx, locationID)
}

func (s *InventoryService) ListExpiring(ctx context.Context, locationID string, days int) ([]repository.InventoryExpiringView, error) {
	return s.inventoryRepo.ListExpiring(ctx, locationID, days)
}

func (s *InventoryService) StockCount(ctx context.Context, items []models.StockCountItem) (int, error) {
	userID := ctx.Value("user_id").(string)
	return s.inventoryRepo.StockCount(ctx, items, userID)
}

var _ InventoryServiceManager = (*InventoryService)(nil)
