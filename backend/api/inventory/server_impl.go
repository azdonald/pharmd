package inventory

import (
	"encoding/json"
	"net/http"

	"github.com/azdonald/pharmd/backend/middleware"
	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/service"
	"github.com/azdonald/pharmd/backend/utils"
	"github.com/go-chi/chi/v5"
)

type stockDataItem struct {
	Batches        *[]BatchSummary `json:"batches,omitempty"`
	BrandName      *string         `json:"brand_name,omitempty"`
	Classification *string         `json:"classification,omitempty"`
	GenericName    *string         `json:"generic_name,omitempty"`
	ProductId      *string         `json:"product_id,omitempty"`
	ProductName    *string         `json:"product_name,omitempty"`
	ReorderLevel   *int            `json:"reorder_level,omitempty"`
	TotalQuantity  *int            `json:"total_quantity,omitempty"`
}

type ServerImpl struct {
	inventoryManager service.InventoryServiceManager
}

func NewServer(inventoryManager service.InventoryServiceManager) ServerInterface {
	return &ServerImpl{inventoryManager: inventoryManager}
}

func floatTo32(f float64) *float32 {
	v := float32(f)
	return &v
}

func (s *ServerImpl) GetInventory(w http.ResponseWriter, r *http.Request, params GetInventoryParams) {
	ctx := r.Context()

	page := 1
	limit := 20
	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}
	query := ""
	if params.Query != nil {
		query = *params.Query
	}

	batches, total, err := s.inventoryManager.ListStock(ctx, params.LocationId, page, limit, query)
	if err != nil {
		http.Error(w, "Failed to list stock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	grouped := make(map[string]*stockDataItem)
	var order []string
	for _, b := range batches {
		pid := b.StockBatch.ProductID
		if _, ok := grouped[pid]; !ok {
			totalQ := 0
			order = append(order, pid)
			grouped[pid] = &stockDataItem{
				ProductId:      &pid,
				ProductName:    &b.ProductName,
				BrandName:      &b.BrandName,
				GenericName:    &b.GenericName,
				Classification: &b.Classification,
				ReorderLevel:   &b.ReorderLevel,
				TotalQuantity:  &totalQ,
				Batches:        &[]BatchSummary{},
			}
		}
		item := grouped[pid]
		*item.TotalQuantity += b.StockBatch.RemainingQty
		bs := *item.Batches
		expiry := b.StockBatch.ExpiryDate
		bs = append(bs, BatchSummary{
			Id:           &b.StockBatch.ID,
			BatchNumber:  &b.StockBatch.BatchNumber,
			Quantity:     &b.StockBatch.RemainingQty,
			UnitCost:     floatTo32(b.StockBatch.UnitCost),
			SellingPrice: floatTo32(b.StockBatch.SellingPrice),
			ExpiryDate:   &expiry,
			IsActive:     &b.StockBatch.IsActive,
		})
		item.Batches = &bs
	}

	data := make([]stockDataItem, len(order))
	for i, pid := range order {
		data[i] = *grouped[pid]
	}

	utils.WriteResponse(ctx, w, map[string]interface{}{
		"data":  data,
		"page":  page,
		"limit": limit,
		"total": total,
	}, http.StatusOK)
}

func (s *ServerImpl) PostInventoryBatches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	batch := models.StockBatch{
		ProductID:        req.ProductId,
		LocationID:       req.LocationId,
		BatchNumber:      utils.PtrOrZero(req.BatchNumber),
		Quantity:         req.Quantity,
		UnitCost:         float64(utils.PtrOrZero32(req.UnitCost)),
		SellingPrice:     float64(utils.PtrOrZero32(req.SellingPrice)),
		ManufacturingDate: utils.PtrOrZero(req.ManufacturingDate),
		ExpiryDate:       utils.PtrOrZero(req.ExpiryDate),
	}

	created, err := s.inventoryManager.CreateBatch(ctx, batch)
	if err != nil {
		http.Error(w, "Failed to create batch: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := created.UpdatedAt.Format("2006-01-02T15:04:05Z")
	receivedDate := created.ReceivedDate.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, batchToResponse(*created, receivedDate, createdAt, updatedAt), http.StatusCreated)
}

func (s *ServerImpl) PostInventoryAdjustments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateAdjustmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	movement := models.StockMovement{
		LocationID:   req.LocationId,
		ProductID:    req.ProductId,
		BatchID:      utils.PtrOrZero(req.BatchId),
		MovementType: req.MovementType,
		Quantity:     req.Quantity,
		Notes:        utils.PtrOrZero(req.Notes),
	}

	created, err := s.inventoryManager.CreateAdjustment(ctx, movement)
	if err != nil {
		http.Error(w, "Failed to record adjustment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, StockMovement{
		Id:           &created.ID,
		BatchId:      &created.BatchID,
		CreatedAt:    &createdAt,
		CreatedBy:    &created.CreatedBy,
		MovementType: &created.MovementType,
		Notes:        &created.Notes,
		ProductId:    &created.ProductID,
		Quantity:     &created.Quantity,
	}, http.StatusOK)
}

func (s *ServerImpl) GetInventoryAlerts(w http.ResponseWriter, r *http.Request, params GetInventoryAlertsParams) {
	ctx := r.Context()

	locationID := ""
	if params.LocationId != nil {
		locationID = *params.LocationId
	}

	alerts, err := s.inventoryManager.ListAlerts(ctx, locationID)
	if err != nil {
		http.Error(w, "Failed to get alerts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]StockAlert, len(alerts))
	for i, a := range alerts {
		response[i] = StockAlert{
			ProductId:     &a.ProductID,
			ProductName:   &a.ProductName,
			BrandName:     &a.BrandName,
			TotalQuantity: &a.TotalQuantity,
			ReorderLevel:  &a.ReorderLevel,
			LocationId:    &a.LocationID,
		}
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) GetInventoryExpiring(w http.ResponseWriter, r *http.Request, params GetInventoryExpiringParams) {
	ctx := r.Context()

	days := 30
	if params.Days != nil {
		days = *params.Days
	}
	locationID := ""
	if params.LocationId != nil {
		locationID = *params.LocationId
	}

	expiring, err := s.inventoryManager.ListExpiring(ctx, locationID, days)
	if err != nil {
		http.Error(w, "Failed to get expiring stock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]ExpiringBatch, len(expiring))
	for i, e := range expiring {
		response[i] = ExpiringBatch{
			Id:              &e.StockBatch.ID,
			BatchNumber:     &e.StockBatch.BatchNumber,
			ExpiryDate:      &e.StockBatch.ExpiryDate,
			ProductId:       &e.StockBatch.ProductID,
			ProductName:     &e.ProductName,
			RemainingQty:    &e.StockBatch.RemainingQty,
			DaysUntilExpiry: &e.DaysUntilExpiry,
		}
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PostInventoryCounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req PostInventoryCountsJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	items := make([]models.StockCountItem, len(req))
	for i, item := range req {
		items[i] = models.StockCountItem{
			ProductID:  item.ProductId,
			LocationID: item.LocationId,
			BatchID:    utils.PtrOrZero(item.BatchId),
			CountedQty: item.CountedQty,
			Notes:      utils.PtrOrZero(item.Notes),
		}
	}

	adjustments, err := s.inventoryManager.StockCount(ctx, items)
	if err != nil {
		http.Error(w, "Failed to record stock count: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(ctx, w, map[string]int{"adjustments": adjustments}, http.StatusOK)
}

func batchToResponse(b models.StockBatch, receivedDate, createdAt, updatedAt string) StockBatch {
	return StockBatch{
		Id:                &b.ID,
		ProductId:         &b.ProductID,
		LocationId:        &b.LocationID,
		BatchNumber:       &b.BatchNumber,
		Quantity:          &b.Quantity,
		RemainingQty:      &b.RemainingQty,
		UnitCost:          floatTo32(b.UnitCost),
		SellingPrice:      floatTo32(b.SellingPrice),
		ManufacturingDate: &b.ManufacturingDate,
		ExpiryDate:        &b.ExpiryDate,
		ReceivedDate:      &receivedDate,
		IsActive:          &b.IsActive,
		CreatedAt:         &createdAt,
		UpdatedAt:         &updatedAt,
	}
}

func (wrapper *ServerInterfaceWrapper) RegisterInventoryRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermInventoryRead)).Get("/inventory", wrapper.GetInventory)
	r.With(middleware.RequirePermission(utils.PermInventoryManage)).Post("/inventory/batches", wrapper.PostInventoryBatches)
	r.With(middleware.RequirePermission(utils.PermInventoryManage)).Post("/inventory/adjustments", wrapper.PostInventoryAdjustments)
	r.With(middleware.RequirePermission(utils.PermInventoryRead)).Get("/inventory/alerts", wrapper.GetInventoryAlerts)
	r.With(middleware.RequirePermission(utils.PermInventoryRead)).Get("/inventory/expiring", wrapper.GetInventoryExpiring)
	r.With(middleware.RequirePermission(utils.PermInventoryManage)).Post("/inventory/counts", wrapper.PostInventoryCounts)
	return r
}
