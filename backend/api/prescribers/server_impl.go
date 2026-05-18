package prescribers

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
	manager service.PrescriberServiceManager
}

func NewServer(manager service.PrescriberServiceManager) ServerInterface {
	return &ServerImpl{manager: manager}
}

func prescriberToResponse(p models.Prescriber, createdAt, updatedAt string) Prescriber {
	isActive := p.IsActive
	return Prescriber{
		Id:            &p.ID,
		Name:          &p.Name,
		LicenseNumber: &p.LicenseNumber,
		Phone:         &p.Phone,
		Email:         &p.Email,
		Specialty:     &p.Specialty,
		DeaNumber:     &p.DEANumber,
		NpiNumber:     &p.NPINumber,
		Address:       &p.Address,
		IsActive:      &isActive,
		CreatedAt:     &createdAt,
		UpdatedAt:     &updatedAt,
	}
}

func (s *ServerImpl) GetPrescribers(w http.ResponseWriter, r *http.Request, params GetPrescribersParams) {
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

	prescribers, total, err := s.manager.List(ctx, query, page, limit)
	if err != nil {
		http.Error(w, "Failed to list prescribers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]Prescriber, len(prescribers))
	for i, p := range prescribers {
		ca := p.CreatedAt.Format("2006-01-02T15:04:05Z")
		ua := p.UpdatedAt.Format("2006-01-02T15:04:05Z")
		resp[i] = prescriberToResponse(p, ca, ua)
	}

	utils.WriteResponse(ctx, w, PrescriberListResponse{
		Data:  &resp,
		Page:  &page,
		Limit: &limit,
		Total: &total,
	}, http.StatusOK)
}

func (s *ServerImpl) PostPrescribers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreatePrescriberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	p := models.Prescriber{
		Name:          req.Name,
		LicenseNumber: utils.PtrOrZero(req.LicenseNumber),
		Phone:         utils.PtrOrZero(req.Phone),
		Email:         utils.PtrOrZero(req.Email),
		Specialty:     utils.PtrOrZero(req.Specialty),
		DEANumber:     utils.PtrOrZero(req.DeaNumber),
		NPINumber:     utils.PtrOrZero(req.NpiNumber),
		Address:       utils.PtrOrZero(req.Address),
	}

	created, err := s.manager.Create(ctx, p)
	if err != nil {
		http.Error(w, "Failed to create prescriber: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ca := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := created.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, prescriberToResponse(*created, ca, ua), http.StatusCreated)
}

func (s *ServerImpl) GetPrescribersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	p, err := s.manager.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Prescriber not found", http.StatusNotFound)
		return
	}

	ca := p.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := p.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, prescriberToResponse(*p, ca, ua), http.StatusOK)
}

func (s *ServerImpl) PutPrescribersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req UpdatePrescriberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	p := models.Prescriber{
		Name:          utils.PtrOrZero(req.Name),
		LicenseNumber: utils.PtrOrZero(req.LicenseNumber),
		Phone:         utils.PtrOrZero(req.Phone),
		Email:         utils.PtrOrZero(req.Email),
		Specialty:     utils.PtrOrZero(req.Specialty),
		DEANumber:     utils.PtrOrZero(req.DeaNumber),
		NPINumber:     utils.PtrOrZero(req.NpiNumber),
		Address:       utils.PtrOrZero(req.Address),
	}
	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}

	updated, err := s.manager.Update(ctx, id, p)
	if err != nil {
		http.Error(w, "Failed to update prescriber: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ca := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, prescriberToResponse(*updated, ca, ua), http.StatusOK)
}

func (s *ServerImpl) DeletePrescribersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := s.manager.Delete(ctx, id); err != nil {
		http.Error(w, "Failed to delete prescriber: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (wrapper *ServerInterfaceWrapper) RegisterPrescribersRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermPrescribersRead)).Get("/prescribers", wrapper.GetPrescribers)
	r.With(middleware.RequirePermission(utils.PermPrescribersCreate)).Post("/prescribers", wrapper.PostPrescribers)
	r.With(middleware.RequirePermission(utils.PermPrescribersRead)).Get("/prescribers/{id}", wrapper.GetPrescribersId)
	r.With(middleware.RequirePermission(utils.PermPrescribersUpdate)).Put("/prescribers/{id}", wrapper.PutPrescribersId)
	r.With(middleware.RequirePermission(utils.PermPrescribersDelete)).Delete("/prescribers/{id}", wrapper.DeletePrescribersId)
	return r
}
