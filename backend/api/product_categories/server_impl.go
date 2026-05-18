package product_categories

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
	categoryManager service.ProductCategoryServiceManager
}

func NewServer(categoryManager service.ProductCategoryServiceManager) ServerInterface {
	return &ServerImpl{categoryManager: categoryManager}
}

func (s *ServerImpl) GetProductCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := s.categoryManager.ListCategories(ctx)
	if err != nil {
		http.Error(w, "Failed to list categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]ProductCategory, len(categories))
	for i, c := range categories {
		createdAt := c.CreatedAt.Format("2006-01-02T15:04:05Z")
		updatedAt := c.UpdatedAt.Format("2006-01-02T15:04:05Z")
		response[i] = ProductCategory{
			Id:          &c.ID,
			Name:        &c.Name,
			Description: &c.Description,
			ParentId:    &c.ParentID,
			IsActive:    &c.IsActive,
			CreatedAt:   &createdAt,
			UpdatedAt:   &updatedAt,
		}
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PostProductCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateProductCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	category := models.ProductCategory{
		Name:        req.Name,
		Description: utils.PtrOrZero(req.Description),
		ParentID:    utils.PtrOrZero(req.ParentId),
	}

	created, err := s.categoryManager.CreateCategory(ctx, category)
	if err != nil {
		http.Error(w, "Failed to create category: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := created.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, ProductCategory{
		Id:          &created.ID,
		Name:        &created.Name,
		Description: &created.Description,
		ParentId:    &created.ParentID,
		IsActive:    &created.IsActive,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}, http.StatusCreated)
}

func (s *ServerImpl) GetProductCategoriesId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	c, err := s.categoryManager.GetCategoryByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get category: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := c.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := c.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, ProductCategory{
		Id:          &c.ID,
		Name:        &c.Name,
		Description: &c.Description,
		ParentId:    &c.ParentID,
		IsActive:    &c.IsActive,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}, http.StatusOK)
}

func (s *ServerImpl) PutProductCategoriesId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req UpdateProductCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	category := models.ProductCategory{
		Name:        utils.PtrOrZero(req.Name),
		Description: utils.PtrOrZero(req.Description),
		ParentID:    utils.PtrOrZero(req.ParentId),
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	updated, err := s.categoryManager.UpdateCategory(ctx, id, category)
	if err != nil {
		http.Error(w, "Failed to update category: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, ProductCategory{
		Id:          &updated.ID,
		Name:        &updated.Name,
		Description: &updated.Description,
		ParentId:    &updated.ParentID,
		IsActive:    &updated.IsActive,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}, http.StatusOK)
}

func (s *ServerImpl) DeleteProductCategoriesId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := s.categoryManager.DeleteCategory(ctx, id); err != nil {
		http.Error(w, "Failed to delete category: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (wrapper *ServerInterfaceWrapper) RegisterProductCategoriesRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermProductsRead)).Get("/product-categories", wrapper.GetProductCategories)
	r.With(middleware.RequirePermission(utils.PermProductsCreate)).Post("/product-categories", wrapper.PostProductCategories)
	r.With(middleware.RequirePermission(utils.PermProductsRead)).Get("/product-categories/{id}", wrapper.GetProductCategoriesId)
	r.With(middleware.RequirePermission(utils.PermProductsUpdate)).Put("/product-categories/{id}", wrapper.PutProductCategoriesId)
	r.With(middleware.RequirePermission(utils.PermProductsDelete)).Delete("/product-categories/{id}", wrapper.DeleteProductCategoriesId)
	return r
}
