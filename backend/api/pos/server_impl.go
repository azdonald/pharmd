package pos

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
	manager service.POSServiceManager
}

func NewServer(manager service.POSServiceManager) ServerInterface {
	return &ServerImpl{manager: manager}
}

func saleToResponse(s models.Sale, createdAt, updatedAt string) Sale {
	st := SaleSaleType(s.SaleType)
	status := SaleStatus(s.Status)
	sub := float32(s.Subtotal)
	tax := float32(s.TaxTotal)
	disc := float32(s.DiscountTotal)
	grand := float32(s.GrandTotal)
	paid := float32(s.PaidAmount)
	chg := float32(s.ChangeAmount)
	return Sale{
		Id:            &s.ID,
		LocationId:    &s.LocationID,
		PatientId:     &s.PatientID,
		PatientName:   &s.PatientName,
		PrescriptionId: &s.PrescriptionID,
		SaleType:      &st,
		Status:        &status,
		Subtotal:      &sub,
		TaxTotal:      &tax,
		DiscountTotal: &disc,
		GrandTotal:    &grand,
		PaidAmount:    &paid,
		ChangeAmount:  &chg,
		Notes:         &s.Notes,
		CreatedBy:     &s.CreatedBy,
		VoidedBy:      &s.VoidedBy,
		VoidedAt:      &s.VoidedAt,
		CreatedAt:     &createdAt,
		UpdatedAt:     &updatedAt,
	}
}

func saleItemToResponse(item models.SaleItem) SaleItem {
	qty := item.Quantity
	up := float32(item.UnitPrice)
	disc := float32(item.Discount)
	lt := float32(item.LineTotal)
	return SaleItem{
		Id:          &item.ID,
		ProductId:   &item.ProductID,
		ProductName: &item.ProductName,
		Quantity:    &qty,
		UnitPrice:   &up,
		Discount:    &disc,
		LineTotal:   &lt,
	}
}

func paymentToResponse(p models.Payment) Payment {
	m := PaymentMethod(p.Method)
	amt := float32(p.Amount)
	return Payment{
		Id:        &p.ID,
		SaleId:    &p.SaleID,
		Method:    &m,
		Amount:    &amt,
		Reference: &p.Reference,
	}
}

func (s *ServerImpl) GetPosSales(w http.ResponseWriter, r *http.Request, params GetPosSalesParams) {
	ctx := r.Context()
	page := 1
	limit := 20
	if params.Page != nil { page = *params.Page }
	if params.Limit != nil { limit = *params.Limit }
	status := ""
	if params.Status != nil { status = *params.Status }

	sales, total, err := s.manager.ListSales(ctx, status, page, limit)
	if err != nil {
		http.Error(w, "Failed to list sales: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]Sale, len(sales))
	for i, sale := range sales {
		ca := sale.CreatedAt.Format("2006-01-02T15:04:05Z")
		ua := sale.UpdatedAt.Format("2006-01-02T15:04:05Z")
		resp[i] = saleToResponse(sale, ca, ua)
	}

	utils.WriteResponse(ctx, w, SaleListResponse{Data: &resp, Page: &page, Limit: &limit, Total: &total}, http.StatusOK)
}

func (s *ServerImpl) PostPosSales(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sale := models.Sale{
		LocationID:     req.LocationId,
		PatientID:      utils.PtrOrZero(req.PatientId),
		PrescriptionID: utils.PtrOrZero(req.PrescriptionId),
		SaleType:       utils.PtrOrZero(req.SaleType),
		Notes:          utils.PtrOrZero(req.Notes),
	}

	items := make([]models.SaleItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = models.SaleItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
			UnitPrice: float64(item.UnitPrice),
			Discount:  float64(utils.PtrOrZero32(item.Discount)),
		}
	}

	created, err := s.manager.CreateSale(ctx, sale, items)
	if err != nil {
		http.Error(w, "Failed to create sale: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ca := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := created.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, saleToResponse(*created, ca, ua), http.StatusCreated)
}

func (s *ServerImpl) GetPosSalesId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	sale, items, payments, err := s.manager.GetSaleByID(ctx, id)
	if err != nil {
		http.Error(w, "Sale not found", http.StatusNotFound)
		return
	}

	ca := sale.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := sale.UpdatedAt.Format("2006-01-02T15:04:05Z")
	resp := saleToResponse(*sale, ca, ua)

	respItems := make([]SaleItem, len(items))
	for i, item := range items {
		respItems[i] = saleItemToResponse(item)
	}
	resp.Items = &respItems

	respPayments := make([]Payment, len(payments))
	for i, p := range payments {
		respPayments[i] = paymentToResponse(p)
	}
	resp.Payments = &respPayments

	utils.WriteResponse(ctx, w, resp, http.StatusOK)
}

func (s *ServerImpl) PutPosSalesId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(string)

	var req UpdateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var updated *models.Sale
	var err error
	switch *req.Status {
	case UpdateSaleRequestStatusVoided:
		updated, err = s.manager.VoidSale(ctx, id, userID)
	case UpdateSaleRequestStatusRefunded:
		updated, err = s.manager.RefundSale(ctx, id, userID)
	case UpdateSaleRequestStatusHeld:
		updated, err = s.manager.HoldSale(ctx, id)
	default:
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Failed to update sale: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ca := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, saleToResponse(*updated, ca, ua), http.StatusOK)
}

func (s *ServerImpl) GetPosSalesIdReceipt(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	receipt, err := s.manager.GetReceipt(ctx, id)
	if err != nil {
		http.Error(w, "Failed to get receipt: "+err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteResponse(ctx, w, receipt, http.StatusOK)
}

func (s *ServerImpl) PostPosPayments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req RecordPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	payments := make([]models.Payment, len(req.Payments))
	for i, p := range req.Payments {
		payments[i] = models.Payment{
			Method:    string(p.Method),
			Amount:    float64(p.Amount),
			Reference: utils.PtrOrZero(p.Reference),
		}
	}

	updated, err := s.manager.RecordPayments(ctx, req.SaleId, payments)
	if err != nil {
		http.Error(w, "Failed to record payments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ca := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, saleToResponse(*updated, ca, ua), http.StatusOK)
}

func (s *ServerImpl) GetPosSummary(w http.ResponseWriter, r *http.Request, params GetPosSummaryParams) {
	ctx := r.Context()
	locationID := ""
	if params.LocationId != nil { locationID = *params.LocationId }
	date := ""
	if params.Date != nil { date = *params.Date }

	summary, err := s.manager.GetDailySummary(ctx, locationID, date)
	if err != nil {
		http.Error(w, "Failed to get summary: "+err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteResponse(ctx, w, summary, http.StatusOK)
}

func (s *ServerImpl) PostPosCloseDay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(string)

	var req CloseDayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := s.manager.CloseDay(ctx, req.LocationId, req.Date, userID, utils.PtrOrZero(req.Notes))
	if err != nil {
		http.Error(w, "Failed to close day: "+err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteResponse(ctx, w, result, http.StatusOK)
}

func (wrapper *ServerInterfaceWrapper) RegisterPOSRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermPosRead)).Get("/pos/sales", wrapper.GetPosSales)
	r.With(middleware.RequirePermission(utils.PermPosCreate)).Post("/pos/sales", wrapper.PostPosSales)
	r.With(middleware.RequirePermission(utils.PermPosRead)).Get("/pos/sales/{id}", wrapper.GetPosSalesId)
	r.With(middleware.RequirePermission(utils.PermPosUpdate)).Put("/pos/sales/{id}", wrapper.PutPosSalesId)
	r.With(middleware.RequirePermission(utils.PermPosRead)).Get("/pos/sales/{id}/receipt", wrapper.GetPosSalesIdReceipt)
	r.With(middleware.RequirePermission(utils.PermPosCreate)).Post("/pos/payments", wrapper.PostPosPayments)
	r.With(middleware.RequirePermission(utils.PermPosReports)).Get("/pos/summary", wrapper.GetPosSummary)
	r.With(middleware.RequirePermission(utils.PermPosCloseout)).Post("/pos/close-day", wrapper.PostPosCloseDay)
	return r
}
