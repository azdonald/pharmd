package service

import (
	"context"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
	"github.com/google/uuid"
)

type PricingServiceManager interface {
	ListPrices(ctx context.Context, productID, locationID string, page, limit int) ([]models.ProductPrice, int, error)
	GetPriceByID(ctx context.Context, id string) (*models.ProductPrice, error)
	UpsertPrice(ctx context.Context, price models.ProductPrice) (*models.ProductPrice, error)
	DeletePrice(ctx context.Context, id string) error

	ListDiscountRules(ctx context.Context, page, limit int) ([]models.DiscountRule, int, error)
	GetDiscountRuleByID(ctx context.Context, id string) (*models.DiscountRule, error)
	CreateDiscountRule(ctx context.Context, rule models.DiscountRule) (*models.DiscountRule, error)
	UpdateDiscountRule(ctx context.Context, id string, rule models.DiscountRule) (*models.DiscountRule, error)
	DeleteDiscountRule(ctx context.Context, id string) error
}

type PricingService struct {
	pricingRepo repository.PricingRepository
}

func NewPricingService(pricingRepo repository.PricingRepository) PricingServiceManager {
	return &PricingService{pricingRepo: pricingRepo}
}

func (s *PricingService) ListPrices(ctx context.Context, productID, locationID string, page, limit int) ([]models.ProductPrice, int, error) {
	return s.pricingRepo.ListPrices(ctx, productID, locationID, page, limit)
}

func (s *PricingService) GetPriceByID(ctx context.Context, id string) (*models.ProductPrice, error) {
	return s.pricingRepo.GetPriceByID(ctx, id)
}

func (s *PricingService) UpsertPrice(ctx context.Context, price models.ProductPrice) (*models.ProductPrice, error) {
	now := time.Now()
	price.ID = uuid.New().String()
	price.CreatedAt = now
	price.UpdatedAt = now
	if err := s.pricingRepo.UpsertPrice(ctx, price); err != nil {
		return nil, err
	}
	return &price, nil
}

func (s *PricingService) DeletePrice(ctx context.Context, id string) error {
	return s.pricingRepo.DeletePrice(ctx, id)
}

func (s *PricingService) ListDiscountRules(ctx context.Context, page, limit int) ([]models.DiscountRule, int, error) {
	return s.pricingRepo.ListDiscountRules(ctx, page, limit)
}

func (s *PricingService) GetDiscountRuleByID(ctx context.Context, id string) (*models.DiscountRule, error) {
	return s.pricingRepo.GetDiscountRuleByID(ctx, id)
}

func (s *PricingService) CreateDiscountRule(ctx context.Context, rule models.DiscountRule) (*models.DiscountRule, error) {
	now := time.Now()
	rule.ID = uuid.New().String()
	rule.CreatedAt = now
	rule.UpdatedAt = now
	if err := s.pricingRepo.CreateDiscountRule(ctx, rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

func (s *PricingService) UpdateDiscountRule(ctx context.Context, id string, rule models.DiscountRule) (*models.DiscountRule, error) {
	rule.UpdatedAt = time.Now()
	if err := s.pricingRepo.UpdateDiscountRule(ctx, id, rule); err != nil {
		return nil, err
	}
	return s.pricingRepo.GetDiscountRuleByID(ctx, id)
}

func (s *PricingService) DeleteDiscountRule(ctx context.Context, id string) error {
	return s.pricingRepo.DeleteDiscountRule(ctx, id)
}
