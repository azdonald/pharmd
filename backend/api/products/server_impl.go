package products

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
	productManager service.ProductServiceManager
}

func NewServer(productManager service.ProductServiceManager) ServerInterface {
	return &ServerImpl{productManager: productManager}
}

func (s *ServerImpl) GetProducts(w http.ResponseWriter, r *http.Request, params GetProductsParams) {
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
	categoryID := ""
	if params.CategoryId != nil {
		categoryID = *params.CategoryId
	}

	products, total, err := s.productManager.ListProducts(ctx, page, limit, query, categoryID)
	if err != nil {
		http.Error(w, "Failed to list products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseProducts := make([]Product, len(products))
	for i, p := range products {
		createdAt := p.CreatedAt.Format("2006-01-02T15:04:05Z")
		updatedAt := p.UpdatedAt.Format("2006-01-02T15:04:05Z")
		responseProducts[i] = productToResponse(p, createdAt, updatedAt)
	}

	response := ProductListResponse{
		Data:  &responseProducts,
		Page:  &page,
		Limit: &limit,
		Total: &total,
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PostProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product := models.Product{
		Name:           req.Name,
		Description:    utils.PtrOrZero(req.Description),
		CategoryID:     utils.PtrOrZero(req.CategoryId),
		Classification: utils.PtrOrZero(req.Classification),
		BrandName:      utils.PtrOrZero(req.BrandName),
		GenericName:    utils.PtrOrZero(req.GenericName),
		Manufacturer:   utils.PtrOrZero(req.Manufacturer),
		Barcode:        utils.PtrOrZero(req.Barcode),
		NDC:            utils.PtrOrZero(req.Ndc),
		UnitOfMeasure:  utils.PtrOrZero(req.UnitOfMeasure),
		Strength:       utils.PtrOrZero(req.Strength),
		Form:           utils.PtrOrZero(req.Form),
	}
	if req.ReorderLevel != nil {
		product.ReorderLevel = *req.ReorderLevel
	}

	created, err := s.productManager.CreateProduct(ctx, product)
	if err != nil {
		http.Error(w, "Failed to create product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := created.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, productToResponse(*created, createdAt, updatedAt), http.StatusCreated)
}

func (s *ServerImpl) GetProductsId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	p, err := s.productManager.GetProductByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := p.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := p.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, productToResponse(*p, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) PutProductsId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product := models.Product{
		Name:           utils.PtrOrZero(req.Name),
		Description:    utils.PtrOrZero(req.Description),
		CategoryID:     utils.PtrOrZero(req.CategoryId),
		Classification: utils.PtrOrZero(req.Classification),
		BrandName:      utils.PtrOrZero(req.BrandName),
		GenericName:    utils.PtrOrZero(req.GenericName),
		Manufacturer:   utils.PtrOrZero(req.Manufacturer),
		Barcode:        utils.PtrOrZero(req.Barcode),
		NDC:            utils.PtrOrZero(req.Ndc),
		UnitOfMeasure:  utils.PtrOrZero(req.UnitOfMeasure),
		Strength:       utils.PtrOrZero(req.Strength),
		Form:           utils.PtrOrZero(req.Form),
	}
	if req.ReorderLevel != nil {
		product.ReorderLevel = *req.ReorderLevel
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	updated, err := s.productManager.UpdateProduct(ctx, id, product)
	if err != nil {
		http.Error(w, "Failed to update product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, productToResponse(*updated, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) DeleteProductsId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := s.productManager.DeleteProduct(ctx, id); err != nil {
		http.Error(w, "Failed to delete product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerImpl) PostProductsBarcodeLookup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req BarcodeLookupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	p, err := s.productManager.GetProductByBarcode(ctx, req.Barcode)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to lookup product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := p.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := p.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, productToResponse(*p, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) GetProductsIdSubstitutes(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	subs, err := s.productManager.ListSubstitutes(ctx, id)
	if err != nil {
		http.Error(w, "Failed to list substitutes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]GenericSubstitution, len(subs))
	for i, sub := range subs {
		createdAt := sub.CreatedAt.Format("2006-01-02T15:04:05Z")
		response[i] = GenericSubstitution{
			CreatedAt:           &createdAt,
			Id:                  &sub.ID,
			Notes:               &sub.Notes,
			ProductId:           &sub.ProductID,
			SubstituteProductId: &sub.SubstituteProductID,
		}
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PostProductsIdSubstitutes(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req CreateGenericSubstitutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sub := models.GenericSubstitution{
		SubstituteProductID: req.SubstituteProductId,
		Notes:               utils.PtrOrZero(req.Notes),
	}

	created, err := s.productManager.AddSubstitute(ctx, id, sub)
	if err != nil {
		http.Error(w, "Failed to add substitute: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, GenericSubstitution{
		CreatedAt:           &createdAt,
		Id:                  &created.ID,
		Notes:               &created.Notes,
		ProductId:           &created.ProductID,
		SubstituteProductId: &created.SubstituteProductID,
	}, http.StatusCreated)
}

func (s *ServerImpl) DeleteProductsIdSubstitutes(w http.ResponseWriter, r *http.Request, id string, params DeleteProductsIdSubstitutesParams) {
	ctx := r.Context()

	if err := s.productManager.RemoveSubstitute(ctx, id, params.SubstituteId); err != nil {
		http.Error(w, "Failed to remove substitute: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func productToResponse(p models.Product, createdAt, updatedAt string) Product {
	return Product{
		Id:             &p.ID,
		Name:           &p.Name,
		Description:    &p.Description,
		CategoryId:     &p.CategoryID,
		Classification: &p.Classification,
		BrandName:      &p.BrandName,
		GenericName:    &p.GenericName,
		Manufacturer:   &p.Manufacturer,
		Barcode:        &p.Barcode,
		Ndc:            &p.NDC,
		UnitOfMeasure:  &p.UnitOfMeasure,
		Strength:       &p.Strength,
		Form:           &p.Form,
		ReorderLevel:   &p.ReorderLevel,
		IsActive:       &p.IsActive,
		CreatedAt:      &createdAt,
		UpdatedAt:      &updatedAt,
	}
}

func (wrapper *ServerInterfaceWrapper) RegisterProductsRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermProductsRead)).Get("/products", wrapper.GetProducts)
	r.With(middleware.RequirePermission(utils.PermProductsCreate)).Post("/products", wrapper.PostProducts)
	r.With(middleware.RequirePermission(utils.PermProductsRead)).Get("/products/{id}", wrapper.GetProductsId)
	r.With(middleware.RequirePermission(utils.PermProductsUpdate)).Put("/products/{id}", wrapper.PutProductsId)
	r.With(middleware.RequirePermission(utils.PermProductsDelete)).Delete("/products/{id}", wrapper.DeleteProductsId)
	r.With(middleware.RequirePermission(utils.PermProductsRead)).Get("/products/{id}/substitutes", wrapper.GetProductsIdSubstitutes)
	r.With(middleware.RequirePermission(utils.PermProductsCreate)).Post("/products/{id}/substitutes", wrapper.PostProductsIdSubstitutes)
	r.With(middleware.RequirePermission(utils.PermProductsDelete)).Delete("/products/{id}/substitutes", wrapper.DeleteProductsIdSubstitutes)
	r.With(middleware.RequirePermission(utils.PermProductsRead)).Post("/products/barcode-lookup", wrapper.PostProductsBarcodeLookup)
	return r
}
