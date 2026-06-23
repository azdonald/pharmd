package service

import (
	"context"
	"fmt"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
	"github.com/google/uuid"
)

type POSServiceManager interface {
	ListSales(ctx context.Context, status string, page, limit int) ([]models.Sale, int, error)
	GetSaleByID(ctx context.Context, id string) (*models.Sale, []models.SaleItem, []models.Payment, error)
	CreateSale(ctx context.Context, sale models.Sale, items []models.SaleItem) (*models.Sale, error)
	VoidSale(ctx context.Context, id, userID string) (*models.Sale, error)
	RefundSale(ctx context.Context, id, userID string) (*models.Sale, error)
	HoldSale(ctx context.Context, id string) (*models.Sale, error)
	RecordPayments(ctx context.Context, saleID string, payments []models.Payment) (*models.Sale, error)
	GetDailySummary(ctx context.Context, locationID, date string) (*map[string]interface{}, error)
	CloseDay(ctx context.Context, locationID, date, userID, notes string) (*map[string]interface{}, error)
	GetReceipt(ctx context.Context, id string) (*map[string]interface{}, error)
}

type POSService struct {
	repo repository.POSRepository
}

func NewPOSService(repo repository.POSRepository) POSServiceManager {
	return &POSService{repo: repo}
}

func (s *POSService) ListSales(ctx context.Context, status string, page, limit int) ([]models.Sale, int, error) {
	return s.repo.ListSales(ctx, status, page, limit)
}

func (s *POSService) GetSaleByID(ctx context.Context, id string) (*models.Sale, []models.SaleItem, []models.Payment, error) {
	return s.repo.GetSaleByID(ctx, id)
}

func (s *POSService) CreateSale(ctx context.Context, sale models.Sale, items []models.SaleItem) (*models.Sale, error) {
	orgID := ctx.Value("organisation_id").(string)
	userID := ctx.Value("user_id").(string)
	now := time.Now()

	sale.ID = uuid.New().String()
	sale.OrganisationID = orgID
	sale.SaleNumber = fmt.Sprintf("SALE-%s-%d", orgID[:8], now.Unix())
	if sale.SaleType == "" {
		sale.SaleType = "otc"
	}
	sale.Status = "active"
	sale.CreatedBy = userID
	sale.CreatedAt = now
	sale.UpdatedAt = now

	var subtotal float64
	for i := range items {
		items[i].ID = uuid.New().String()
		items[i].OrganisationID = orgID
		items[i].SaleID = sale.ID
		items[i].LineTotal = float64(items[i].Quantity)*items[i].UnitPrice - items[i].Discount
		subtotal += items[i].LineTotal
	}
	sale.Subtotal = subtotal
	sale.GrandTotal = subtotal + sale.TaxTotal - sale.DiscountTotal

	if err := s.repo.CreateSale(ctx, sale, items); err != nil {
		return nil, err
	}
	return &sale, nil
}

func (s *POSService) VoidSale(ctx context.Context, id, userID string) (*models.Sale, error) {
	existing, _, _, err := s.repo.GetSaleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing.Status == "completed" {
		if err := s.repo.RestoreSaleInventory(ctx, id, "voided", userID); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.UpdateSale(ctx, id, "voided", userID, ""); err != nil {
			return nil, err
		}
	}

	rx, _, _, err := s.repo.GetSaleByID(ctx, id)
	return rx, err
}

func (s *POSService) RefundSale(ctx context.Context, id, userID string) (*models.Sale, error) {
	existing, _, _, err := s.repo.GetSaleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing.Status == "completed" {
		if err := s.repo.RestoreSaleInventory(ctx, id, "refunded", userID); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.UpdateSale(ctx, id, "refunded", userID, ""); err != nil {
			return nil, err
		}
	}

	rx, _, _, err := s.repo.GetSaleByID(ctx, id)
	return rx, err
}

func (s *POSService) HoldSale(ctx context.Context, id string) (*models.Sale, error) {
	if err := s.repo.UpdateSale(ctx, id, "held", "", ""); err != nil {
		return nil, err
	}
	rx, _, _, err := s.repo.GetSaleByID(ctx, id)
	return rx, err
}

func (s *POSService) RecordPayments(ctx context.Context, saleID string, payments []models.Payment) (*models.Sale, error) {
	now := time.Now()
	userID := ctx.Value("user_id").(string)
	for i := range payments {
		payments[i].ID = uuid.New().String()
	}

	sale, items, _, err := s.repo.GetSaleByID(ctx, saleID)
	if err != nil {
		return nil, err
	}

	var totalPaid float64
	for _, p := range payments {
		totalPaid += p.Amount
	}

	change := totalPaid - sale.GrandTotal
	if change < 0 {
		change = 0
	}

	if err := s.repo.CompleteSale(ctx, saleID, *sale, items, payments, userID, totalPaid, change); err != nil {
		return nil, err
	}

	// update daily summary
	_ = s.updateDailySummary(ctx, sale.LocationID, now.Format("2006-01-02"))

	sale, _, _, err = s.repo.GetSaleByID(ctx, saleID)
	return sale, err
}

func (s *POSService) updateDailySummary(ctx context.Context, locationID, date string) error {
	ds, _, err := s.repo.GetDailySummary(ctx, locationID, date)
	if err != nil {
		return err
	}
	ds.ID = uuid.New().String()
	ds.CreatedAt = time.Now()
	ds.UpdatedAt = time.Now()
	return s.repo.UpsertDailySummary(ctx, *ds)
}

func (s *POSService) GetDailySummary(ctx context.Context, locationID, date string) (*map[string]interface{}, error) {
	ds, methods, err := s.repo.GetDailySummary(ctx, locationID, date)
	if err != nil {
		return nil, err
	}

	byMethod := map[string]float64{}
	for _, m := range methods {
		byMethod[m.Method] = m.Total
	}

	result := map[string]interface{}{
		"date":            ds.Date,
		"total_sales":     ds.TotalSales,
		"total_revenue":   ds.TotalRevenue,
		"total_tax":       ds.TotalTax,
		"total_discounts": ds.TotalDiscounts,
		"by_method":       byMethod,
		"is_closed":       ds.IsClosed,
		"closed_by":       ds.ClosedBy,
		"closed_at":       ds.ClosedAt,
	}
	return &result, nil
}

func (s *POSService) CloseDay(ctx context.Context, locationID, date, userID, notes string) (*map[string]interface{}, error) {
	ds, _, err := s.repo.GetDailySummary(ctx, locationID, date)
	if err != nil {
		return nil, err
	}
	if err := s.repo.CloseDay(ctx, ds.ID, userID, notes); err != nil {
		return nil, err
	}
	ds.IsClosed = true
	ds.ClosedBy = userID
	ds.Notes = notes

	byMethod := map[string]float64{}
	result := map[string]interface{}{
		"date":            ds.Date,
		"total_sales":     ds.TotalSales,
		"total_revenue":   ds.TotalRevenue,
		"total_tax":       ds.TotalTax,
		"total_discounts": ds.TotalDiscounts,
		"by_method":       byMethod,
		"is_closed":       true,
		"closed_by":       userID,
		"closed_at":       time.Now().Format("2006-01-02T15:04:05Z"),
	}
	return &result, nil
}

func (s *POSService) GetReceipt(ctx context.Context, id string) (*map[string]interface{}, error) {
	sale, items, payments, err := s.repo.GetSaleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	receiptItems := make([]map[string]interface{}, len(items))
	for i, item := range items {
		receiptItems[i] = map[string]interface{}{
			"name":  item.ProductName,
			"qty":   item.Quantity,
			"price": item.UnitPrice,
			"total": item.LineTotal,
		}
	}

	paymentStrs := make([]string, len(payments))
	for i, p := range payments {
		paymentStrs[i] = fmt.Sprintf("%s: $%.2f", p.Method, p.Amount)
	}

	result := map[string]interface{}{
		"sale_number": sale.SaleNumber,
		"date":        sale.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"items":       receiptItems,
		"subtotal":    sale.Subtotal,
		"tax":         sale.TaxTotal,
		"discount":    sale.DiscountTotal,
		"grand_total": sale.GrandTotal,
		"paid":        sale.PaidAmount,
		"change":      sale.ChangeAmount,
		"payments":    paymentStrs,
	}
	return &result, nil
}
