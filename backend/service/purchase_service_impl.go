package service

import (
	"context"
	"fmt"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
	"github.com/google/uuid"
)

type PurchaseOrderService struct {
	purchaseRepo repository.PurchaseOrderRepository
}

func NewPurchaseOrderService(purchaseRepo repository.PurchaseOrderRepository) PurchaseOrderServiceManager {
	return &PurchaseOrderService{purchaseRepo: purchaseRepo}
}

func (s *PurchaseOrderService) ListPurchaseOrders(ctx context.Context, page, limit int, status string) ([]models.PurchaseOrder, int, error) {
	return s.purchaseRepo.ListPurchaseOrders(ctx, page, limit, status)
}

func (s *PurchaseOrderService) GetPurchaseOrderByID(ctx context.Context, id string) (*models.PurchaseOrder, error) {
	return s.purchaseRepo.GetPurchaseOrderByID(ctx, id)
}

func (s *PurchaseOrderService) CreatePurchaseOrder(ctx context.Context, po models.PurchaseOrder, items []models.PurchaseOrderItem) (*models.PurchaseOrder, error) {
	orgID := ctx.Value("organisation_id").(string)
	userID := ctx.Value("user_id").(string)
	now := time.Now()

	po.ID = uuid.New().String()
	po.OrganisationID = orgID
	po.PONumber = fmt.Sprintf("PO-%s-%d", orgID[:8], now.Unix())
	po.Status = "draft"
	po.OrderDate = now
	po.CreatedBy = userID
	po.CreatedAt = now
	po.UpdatedAt = now

	var subtotal float64
	for i := range items {
		items[i].ID = uuid.New().String()
		items[i].OrganisationID = orgID
		items[i].PurchaseOrderID = po.ID
		items[i].LineTotal = float64(items[i].QuantityOrdered) * items[i].UnitCost
		subtotal += items[i].LineTotal
	}
	po.Subtotal = subtotal
	po.GrandTotal = subtotal + po.TaxTotal

	if err := s.purchaseRepo.CreatePurchaseOrder(ctx, po, items); err != nil {
		return nil, err
	}
	return &po, nil
}

func (s *PurchaseOrderService) ApprovePurchaseOrder(ctx context.Context, id string) (*models.PurchaseOrder, error) {
	userID := ctx.Value("user_id").(string)
	if err := s.purchaseRepo.UpdatePOStatus(ctx, id, "approved", userID); err != nil {
		return nil, err
	}
	return s.purchaseRepo.GetPurchaseOrderByID(ctx, id)
}

func (s *PurchaseOrderService) RejectPurchaseOrder(ctx context.Context, id string) (*models.PurchaseOrder, error) {
	if err := s.purchaseRepo.UpdatePOStatus(ctx, id, "rejected", ""); err != nil {
		return nil, err
	}
	return s.purchaseRepo.GetPurchaseOrderByID(ctx, id)
}

func (s *PurchaseOrderService) ReceiveGoods(ctx context.Context, id string, items []models.PurchaseOrderItem, notes string) (*models.PurchaseOrder, error) {
	userID := ctx.Value("user_id").(string)
	if err := s.purchaseRepo.ReceiveGoods(ctx, id, items, notes, userID); err != nil {
		return nil, err
	}
	return s.purchaseRepo.GetPurchaseOrderByID(ctx, id)
}

var _ PurchaseOrderServiceManager = (*PurchaseOrderService)(nil)
