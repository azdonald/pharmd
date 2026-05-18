package pricing

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
	pricingManager service.PricingServiceManager
}

func NewServer(pricingManager service.PricingServiceManager) ServerInterface {
	return &ServerImpl{pricingManager: pricingManager}
}

func priceToResponse(p models.ProductPrice, createdAt, updatedAt string) ProductPrice {
	isActive := p.IsActive
	return ProductPrice{
		Id:           &p.ID,
		ProductId:    &p.ProductID,
		ProductName:  &p.ProductName,
		LocationId:   &p.LocationID,
		LocationName: &p.LocationName,
		SellingPrice: utils.Ptr(float32(p.SellingPrice)),
		CostPrice:    utils.Ptr(float32(p.CostPrice)),
		MinPrice:     utils.Ptr(float32(p.MinPrice)),
		MaxDiscount:  utils.Ptr(float32(p.MaxDiscount)),
		IsActive:     &isActive,
		CreatedAt:    &createdAt,
		UpdatedAt:    &updatedAt,
	}
}

func ruleToResponse(r models.DiscountRule, createdAt, updatedAt string) DiscountRule {
	ruleType := DiscountRuleType(r.Type)
	appliesTo := DiscountRuleAppliesTo(r.AppliesTo)
	isActive := r.IsActive
	value := float32(r.Value)
	var minOrderValue, maxDiscountAmount float32
	if r.MinOrderValue != 0 {
		minOrderValue = float32(r.MinOrderValue)
	}
	if r.MaxDiscountAmount != 0 {
		maxDiscountAmount = float32(r.MaxDiscountAmount)
	}
	return DiscountRule{
		Id:                &r.ID,
		Name:              &r.Name,
		Type:              &ruleType,
		Value:             &value,
		MinOrderValue:     &minOrderValue,
		MaxDiscountAmount: &maxDiscountAmount,
		AppliesTo:         &appliesTo,
		AppliesToId:       &r.AppliesToID,
		IsActive:          &isActive,
		StartDate:         &r.StartDate,
		EndDate:           &r.EndDate,
		CreatedAt:         &createdAt,
		UpdatedAt:         &updatedAt,
	}
}

func (s *ServerImpl) GetPricing(w http.ResponseWriter, r *http.Request, params GetPricingParams) {
	ctx := r.Context()

	page := 1
	limit := 20
	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}
	productID := ""
	if params.ProductId != nil {
		productID = *params.ProductId
	}
	locationID := ""
	if params.LocationId != nil {
		locationID = *params.LocationId
	}

	prices, total, err := s.pricingManager.ListPrices(ctx, productID, locationID, page, limit)
	if err != nil {
		http.Error(w, "Failed to list prices: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responsePrices := make([]ProductPrice, len(prices))
	for i, p := range prices {
		createdAt := p.CreatedAt.Format("2006-01-02T15:04:05Z")
		updatedAt := p.UpdatedAt.Format("2006-01-02T15:04:05Z")
		responsePrices[i] = priceToResponse(p, createdAt, updatedAt)
	}

	utils.WriteResponse(ctx, w, PriceListResponse{
		Data:  &responsePrices,
		Page:  &page,
		Limit: &limit,
		Total: &total,
	}, http.StatusOK)
}

func (s *ServerImpl) PostPricing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req UpsertPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	price := models.ProductPrice{
		ProductID:    req.ProductId,
		LocationID:   req.LocationId,
		SellingPrice: float64(req.SellingPrice),
		CostPrice:    float64(utils.PtrOrZero32(req.CostPrice)),
		MinPrice:     float64(utils.PtrOrZero32(req.MinPrice)),
		MaxDiscount:  float64(utils.PtrOrZero32(req.MaxDiscount)),
		IsActive:     true,
	}

	created, err := s.pricingManager.UpsertPrice(ctx, price)
	if err != nil {
		http.Error(w, "Failed to save price: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := created.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, priceToResponse(*created, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) GetPricingId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	p, err := s.pricingManager.GetPriceByID(ctx, id)
	if err != nil {
		http.Error(w, "Price not found", http.StatusNotFound)
		return
	}

	createdAt := p.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := p.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, priceToResponse(*p, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) DeletePricingId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := s.pricingManager.DeletePrice(ctx, id); err != nil {
		http.Error(w, "Failed to delete price: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerImpl) GetPricingRules(w http.ResponseWriter, r *http.Request, params GetPricingRulesParams) {
	ctx := r.Context()

	page := 1
	limit := 20
	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	rules, total, err := s.pricingManager.ListDiscountRules(ctx, page, limit)
	if err != nil {
		http.Error(w, "Failed to list discount rules: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseRules := make([]DiscountRule, len(rules))
	for i, r := range rules {
		createdAt := r.CreatedAt.Format("2006-01-02T15:04:05Z")
		updatedAt := r.UpdatedAt.Format("2006-01-02T15:04:05Z")
		responseRules[i] = ruleToResponse(r, createdAt, updatedAt)
	}

	utils.WriteResponse(ctx, w, DiscountRuleListResponse{
		Data:  &responseRules,
		Page:  &page,
		Limit: &limit,
		Total: &total,
	}, http.StatusOK)
}

func (s *ServerImpl) PostPricingRules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateDiscountRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	appliesTo := "all"
	if req.AppliesTo != nil {
		appliesTo = string(*req.AppliesTo)
	}

	rule := models.DiscountRule{
		Name:              req.Name,
		Type:              string(req.Type),
		Value:             float64(req.Value),
		MinOrderValue:     float64(utils.PtrOrZero32(req.MinOrderValue)),
		MaxDiscountAmount: float64(utils.PtrOrZero32(req.MaxDiscountAmount)),
		AppliesTo:         appliesTo,
		AppliesToID:       utils.PtrOrZero(req.AppliesToId),
		IsActive:          true,
		StartDate:         utils.PtrOrZero(req.StartDate),
		EndDate:           utils.PtrOrZero(req.EndDate),
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}

	created, err := s.pricingManager.CreateDiscountRule(ctx, rule)
	if err != nil {
		http.Error(w, "Failed to create discount rule: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := created.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, ruleToResponse(*created, createdAt, updatedAt), http.StatusCreated)
}

func (s *ServerImpl) GetPricingRulesId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	rule, err := s.pricingManager.GetDiscountRuleByID(ctx, id)
	if err != nil {
		http.Error(w, "Discount rule not found", http.StatusNotFound)
		return
	}

	createdAt := rule.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := rule.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, ruleToResponse(*rule, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) PutPricingRulesId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req UpdateDiscountRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	appliesTo := "all"
	if req.AppliesTo != nil {
		appliesTo = string(*req.AppliesTo)
	}
	ruleType := "percentage"
	if req.Type != nil {
		ruleType = string(*req.Type)
	}

	rule := models.DiscountRule{
		Name:              utils.PtrOrZero(req.Name),
		Type:              ruleType,
		Value:             float64(utils.PtrOrZero32(req.Value)),
		MinOrderValue:     float64(utils.PtrOrZero32(req.MinOrderValue)),
		MaxDiscountAmount: float64(utils.PtrOrZero32(req.MaxDiscountAmount)),
		AppliesTo:         appliesTo,
		AppliesToID:       utils.PtrOrZero(req.AppliesToId),
		StartDate:         utils.PtrOrZero(req.StartDate),
		EndDate:           utils.PtrOrZero(req.EndDate),
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}

	updated, err := s.pricingManager.UpdateDiscountRule(ctx, id, rule)
	if err != nil {
		http.Error(w, "Failed to update discount rule: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, ruleToResponse(*updated, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) DeletePricingRulesId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := s.pricingManager.DeleteDiscountRule(ctx, id); err != nil {
		http.Error(w, "Failed to delete discount rule: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (wrapper *ServerInterfaceWrapper) RegisterPricingRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermPricingRead)).Get("/pricing", wrapper.GetPricing)
	r.With(middleware.RequirePermission(utils.PermPricingCreate)).Post("/pricing", wrapper.PostPricing)
	r.With(middleware.RequirePermission(utils.PermDiscountsRead)).Get("/pricing/rules", wrapper.GetPricingRules)
	r.With(middleware.RequirePermission(utils.PermDiscountsCreate)).Post("/pricing/rules", wrapper.PostPricingRules)
	r.With(middleware.RequirePermission(utils.PermDiscountsRead)).Get("/pricing/rules/{id}", wrapper.GetPricingRulesId)
	r.With(middleware.RequirePermission(utils.PermDiscountsUpdate)).Put("/pricing/rules/{id}", wrapper.PutPricingRulesId)
	r.With(middleware.RequirePermission(utils.PermDiscountsDelete)).Delete("/pricing/rules/{id}", wrapper.DeletePricingRulesId)
	r.With(middleware.RequirePermission(utils.PermPricingRead)).Get("/pricing/{id}", wrapper.GetPricingId)
	r.With(middleware.RequirePermission(utils.PermPricingDelete)).Delete("/pricing/{id}", wrapper.DeletePricingId)
	return r
}
