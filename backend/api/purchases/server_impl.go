package purchases

import (
	"encoding/json"
	"net/http"

	"github.com/azdonald/pharmd/backend/middleware"
	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/service"
	"github.com/azdonald/pharmd/backend/utils"
	"github.com/go-chi/chi/v5"
)

type ServerImpl struct {
	purchaseManager service.PurchaseOrderServiceManager
}

func NewServer(purchaseManager service.PurchaseOrderServiceManager) ServerInterface {
	return &ServerImpl{purchaseManager: purchaseManager}
}

func purchaseToResponse(po models.PurchaseOrder, createdAt, updatedAt string) PurchaseOrder {
	resp := PurchaseOrder{
		Id:           &po.ID,
		PoNumber:     &po.PONumber,
		SupplierId:   &po.SupplierID,
		LocationId:   &po.LocationID,
		Status:       &po.Status,
		OrderDate:    utils.Ptr(po.OrderDate.Format("2006-01-02T15:04:05Z")),
		ExpectedDate: &po.ExpectedDate,
		Notes:        &po.Notes,
		Subtotal:     utils.Ptr(float32(po.Subtotal)),
		TaxTotal:     utils.Ptr(float32(po.TaxTotal)),
		GrandTotal:   utils.Ptr(float32(po.GrandTotal)),
		CreatedBy:    &po.CreatedBy,
		ApprovedBy:   &po.ApprovedBy,
		CreatedAt:    &createdAt,
		UpdatedAt:    &updatedAt,
	}
	if po.ApprovedAt != nil {
		aa := po.ApprovedAt.Format("2006-01-02T15:04:05Z")
		resp.ApprovedAt = &aa
	}
	return resp
}

func itemToResponse(item models.PurchaseOrderItem) PurchaseOrderItem {
	productName := ""
	ls := float32(item.LineTotal)
	qo := item.QuantityOrdered
	qr := item.QuantityReceived
	uc := float32(item.UnitCost)
	return PurchaseOrderItem{
		Id:               &item.ID,
		ProductId:        &item.ProductID,
		ProductName:      &productName,
		QuantityOrdered:  &qo,
		QuantityReceived: &qr,
		UnitCost:         &uc,
		LineTotal:        &ls,
	}
}

func (s *ServerImpl) GetPurchases(w http.ResponseWriter, r *http.Request, params GetPurchasesParams) {
	ctx := r.Context()

	page := 1
	limit := 20
	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}
	status := ""
	if params.Status != nil {
		status = *params.Status
	}

	orders, total, err := s.purchaseManager.ListPurchaseOrders(ctx, page, limit, status)
	if err != nil {
		http.Error(w, "Failed to list purchase orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseOrders := make([]PurchaseOrder, len(orders))
	for i, po := range orders {
		createdAt := po.CreatedAt.Format("2006-01-02T15:04:05Z")
		updatedAt := po.UpdatedAt.Format("2006-01-02T15:04:05Z")
		responseOrders[i] = purchaseToResponse(po, createdAt, updatedAt)
	}

	utils.WriteResponse(ctx, w, PurchaseOrderListResponse{
		Data:  &responseOrders,
		Page:  &page,
		Limit: &limit,
		Total: &total,
	}, http.StatusOK)
}

func (s *ServerImpl) PostPurchases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreatePurchaseOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	po := models.PurchaseOrder{
		SupplierID:   req.SupplierId,
		LocationID:   req.LocationId,
		ExpectedDate: utils.PtrOrZero(req.ExpectedDate),
		Notes:        utils.PtrOrZero(req.Notes),
	}

	items := make([]models.PurchaseOrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = models.PurchaseOrderItem{
			ProductID:       item.ProductId,
			QuantityOrdered: item.QuantityOrdered,
			UnitCost:        float64(item.UnitCost),
		}
	}

	created, err := s.purchaseManager.CreatePurchaseOrder(ctx, po, items)
	if err != nil {
		http.Error(w, "Failed to create purchase order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := created.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, purchaseToResponse(*created, createdAt, updatedAt), http.StatusCreated)
}

func (s *ServerImpl) GetPurchasesId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	po, err := s.purchaseManager.GetPurchaseOrderByID(ctx, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "Purchase order not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get purchase order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := po.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := po.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, purchaseToResponse(*po, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) PutPurchasesIdApprove(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	po, err := s.purchaseManager.ApprovePurchaseOrder(ctx, id)
	if err != nil {
		http.Error(w, "Failed to approve purchase order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := po.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := po.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, purchaseToResponse(*po, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) PostPurchasesIdReceive(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req ReceiveGoodsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	items := make([]models.PurchaseOrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = models.PurchaseOrderItem{
			ID:               item.ItemId,
			QuantityReceived: item.QuantityReceived,
		}
	}

	po, err := s.purchaseManager.ReceiveGoods(ctx, id, items, utils.PtrOrZero(req.Notes))
	if err != nil {
		http.Error(w, "Failed to receive goods: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := po.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := po.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, purchaseToResponse(*po, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) PutPurchasesIdReject(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	po, err := s.purchaseManager.RejectPurchaseOrder(ctx, id)
	if err != nil {
		http.Error(w, "Failed to reject purchase order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := po.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := po.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, purchaseToResponse(*po, createdAt, updatedAt), http.StatusOK)
}

func (wrapper *ServerInterfaceWrapper) RegisterPurchasesRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermPurchasesRead)).Get("/purchases", wrapper.GetPurchases)
	r.With(middleware.RequirePermission(utils.PermPurchasesCreate)).Post("/purchases", wrapper.PostPurchases)
	r.With(middleware.RequirePermission(utils.PermPurchasesRead)).Get("/purchases/{id}", wrapper.GetPurchasesId)
	r.With(middleware.RequirePermission(utils.PermPurchasesApprove)).Put("/purchases/{id}/approve", wrapper.PutPurchasesIdApprove)
	r.With(middleware.RequirePermission(utils.PermPurchasesReceive)).Post("/purchases/{id}/receive", wrapper.PostPurchasesIdReceive)
	r.With(middleware.RequirePermission(utils.PermPurchasesApprove)).Put("/purchases/{id}/reject", wrapper.PutPurchasesIdReject)
	return r
}
