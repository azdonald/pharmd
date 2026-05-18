package suppliers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/azdonald/pharmd/backend/middleware"
	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/service"
	"github.com/azdonald/pharmd/backend/utils"
	"github.com/go-chi/chi/v5"
)

type ServerImpl struct {
	supplierManager service.SupplierServiceManager
}

func NewServer(supplierManager service.SupplierServiceManager) ServerInterface {
	return &ServerImpl{supplierManager: supplierManager}
}

func supplierToResponse(s models.Supplier, createdAt, updatedAt string) Supplier {
	return Supplier{
		Id:            &s.ID,
		Name:          &s.Name,
		ContactPerson: &s.ContactPerson,
		Phone:         &s.Phone,
		Email:         &s.Email,
		Address:       &s.Address,
		City:          &s.City,
		State:         &s.State,
		Country:       &s.Country,
		PaymentTerms:  &s.PaymentTerms,
		Notes:         &s.Notes,
		IsActive:      &s.IsActive,
		CreatedAt:     &createdAt,
		UpdatedAt:     &updatedAt,
	}
}

func (s *ServerImpl) GetSuppliers(w http.ResponseWriter, r *http.Request, params GetSuppliersParams) {
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

	suppliers, total, err := s.supplierManager.ListSuppliers(ctx, page, limit, query)
	if err != nil {
		http.Error(w, "Failed to list suppliers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseSuppliers := make([]Supplier, len(suppliers))
	for i, sp := range suppliers {
		createdAt := sp.CreatedAt.Format("2006-01-02T15:04:05Z")
		updatedAt := sp.UpdatedAt.Format("2006-01-02T15:04:05Z")
		responseSuppliers[i] = supplierToResponse(sp, createdAt, updatedAt)
	}

	utils.WriteResponse(ctx, w, SupplierListResponse{
		Data:  &responseSuppliers,
		Page:  &page,
		Limit: &limit,
		Total: &total,
	}, http.StatusOK)
}

func (s *ServerImpl) PostSuppliers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	supplier := models.Supplier{
		Name:          req.Name,
		ContactPerson: utils.PtrOrZero(req.ContactPerson),
		Phone:         utils.PtrOrZero(req.Phone),
		Email:         utils.PtrOrZero(req.Email),
		Address:       utils.PtrOrZero(req.Address),
		City:          utils.PtrOrZero(req.City),
		State:         utils.PtrOrZero(req.State),
		Country:       utils.PtrOrZero(req.Country),
		PaymentTerms:  utils.PtrOrZero(req.PaymentTerms),
		Notes:         utils.PtrOrZero(req.Notes),
	}

	created, err := s.supplierManager.CreateSupplier(ctx, supplier)
	if err != nil {
		http.Error(w, "Failed to create supplier: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := created.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, supplierToResponse(*created, createdAt, updatedAt), http.StatusCreated)
}

func (s *ServerImpl) GetSuppliersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	sp, err := s.supplierManager.GetSupplierByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Supplier not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get supplier: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := sp.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := sp.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, supplierToResponse(*sp, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) PutSuppliersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req UpdateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	supplier := models.Supplier{
		Name:          utils.PtrOrZero(req.Name),
		ContactPerson: utils.PtrOrZero(req.ContactPerson),
		Phone:         utils.PtrOrZero(req.Phone),
		Email:         utils.PtrOrZero(req.Email),
		Address:       utils.PtrOrZero(req.Address),
		City:          utils.PtrOrZero(req.City),
		State:         utils.PtrOrZero(req.State),
		Country:       utils.PtrOrZero(req.Country),
		PaymentTerms:  utils.PtrOrZero(req.PaymentTerms),
		Notes:         utils.PtrOrZero(req.Notes),
	}
	if req.IsActive != nil {
		supplier.IsActive = *req.IsActive
	}

	updated, err := s.supplierManager.UpdateSupplier(ctx, id, supplier)
	if err != nil {
		http.Error(w, "Failed to update supplier: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, supplierToResponse(*updated, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) DeleteSuppliersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := s.supplierManager.DeleteSupplier(ctx, id); err != nil {
		http.Error(w, "Failed to delete supplier: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerImpl) GetSuppliersIdProducts(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	products, err := s.supplierManager.ListSupplierProducts(ctx, id)
	if err != nil {
		http.Error(w, "Failed to list products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]SupplierProduct, len(products))
	for i, p := range products {
		createdAt := p.CreatedAt.Format("2006-01-02T15:04:05Z")
		updatedAt := p.UpdatedAt.Format("2006-01-02T15:04:05Z")
		response[i] = SupplierProduct{
			Id:           &p.ID,
			SupplierId:   &p.SupplierID,
			ProductId:    &p.ProductID,
			UnitPrice:    utils.Ptr(float32(p.UnitPrice)),
			MinOrderQty:  &p.MinOrderQty,
			LeadTimeDays: &p.LeadTimeDays,
			Notes:        &p.Notes,
			CreatedAt:    &createdAt,
			UpdatedAt:    &updatedAt,
		}
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PutSuppliersIdProducts(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req PutSuppliersIdProductsJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	products := make([]models.SupplierProduct, len(req))
	for i, p := range req {
		products[i] = models.SupplierProduct{
			ProductID:    p.ProductId,
			UnitPrice:    float64(p.UnitPrice),
			MinOrderQty:  utils.PtrOrZeroInt(p.MinOrderQty),
			LeadTimeDays: utils.PtrOrZeroInt(p.LeadTimeDays),
			Notes:        utils.PtrOrZero(p.Notes),
		}
	}

	updated, err := s.supplierManager.SetSupplierProducts(ctx, id, products)
	if err != nil {
		http.Error(w, "Failed to set products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]SupplierProduct, len(updated))
	for i, p := range updated {
		createdAt := p.CreatedAt.Format("2006-01-02T15:04:05Z")
		updatedAt := p.UpdatedAt.Format("2006-01-02T15:04:05Z")
		response[i] = SupplierProduct{
			Id:           &p.ID,
			SupplierId:   &p.SupplierID,
			ProductId:    &p.ProductID,
			UnitPrice:    utils.Ptr(float32(p.UnitPrice)),
			MinOrderQty:  &p.MinOrderQty,
			LeadTimeDays: &p.LeadTimeDays,
			Notes:        &p.Notes,
			CreatedAt:    &createdAt,
			UpdatedAt:    &updatedAt,
		}
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (wrapper *ServerInterfaceWrapper) RegisterSuppliersRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermSuppliersRead)).Get("/suppliers", wrapper.GetSuppliers)
	r.With(middleware.RequirePermission(utils.PermSuppliersCreate)).Post("/suppliers", wrapper.PostSuppliers)
	r.With(middleware.RequirePermission(utils.PermSuppliersRead)).Get("/suppliers/{id}", wrapper.GetSuppliersId)
	r.With(middleware.RequirePermission(utils.PermSuppliersUpdate)).Put("/suppliers/{id}", wrapper.PutSuppliersId)
	r.With(middleware.RequirePermission(utils.PermSuppliersDelete)).Delete("/suppliers/{id}", wrapper.DeleteSuppliersId)
	r.With(middleware.RequirePermission(utils.PermSuppliersRead)).Get("/suppliers/{id}/products", wrapper.GetSuppliersIdProducts)
	r.With(middleware.RequirePermission(utils.PermSuppliersUpdate)).Put("/suppliers/{id}/products", wrapper.PutSuppliersIdProducts)
	return r
}
