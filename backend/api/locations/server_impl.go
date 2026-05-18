package locations

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
	locationManager service.LocationServiceManager
}

func NewServer(locationManager service.LocationServiceManager) ServerInterface {
	return &ServerImpl{locationManager: locationManager}
}

func floatToPtr(f float64) *float32 {
	v := float32(f)
	return &v
}

func ptrToFloat(p *float32) float64 {
	if p != nil {
		return float64(*p)
	}
	return 0
}

func (s *ServerImpl) GetLocations(w http.ResponseWriter, r *http.Request, params GetLocationsParams) {
	ctx := r.Context()

	page := 1
	limit := 20
	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	locations, err := s.locationManager.ListLocations(ctx, page, limit)
	if err != nil {
		http.Error(w, "Failed to list locations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseLocations := make([]Location, len(locations))
	for i, l := range locations {
		createdAt := l.CreatedAt.Format("2006-01-02T15:04:05Z")
		updatedAt := l.UpdatedAt.Format("2006-01-02T15:04:05Z")
		responseLocations[i] = Location{
			Id:        &l.ID,
			Name:      &l.Name,
			Address:   &l.Address,
			City:      &l.City,
			State:     &l.State,
			Country:   &l.Country,
			Phone:     &l.Phone,
			Email:     &l.Email,
			TaxRate:   floatToPtr(l.TaxRate),
			Timezone:  &l.Timezone,
			IsActive:  &l.IsActive,
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		}
	}

	total := len(responseLocations)
	response := LocationListResponse{
		Data:  &responseLocations,
		Page:  &page,
		Limit: &limit,
		Total: &total,
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PostLocations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	location := models.Location{
		Name:     req.Name,
		Address:  utils.PtrOrZero(req.Address),
		City:     utils.PtrOrZero(req.City),
		State:    utils.PtrOrZero(req.State),
		Country:  utils.PtrOrZero(req.Country),
		Phone:    utils.PtrOrZero(req.Phone),
		Email:    utils.PtrOrZero(req.Email),
		TaxRate:  ptrToFloat(req.TaxRate),
		Timezone: utils.PtrOrZero(req.Timezone),
	}

	created, err := s.locationManager.CreateLocation(ctx, location)
	if err != nil {
		http.Error(w, "Failed to create location: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := created.UpdatedAt.Format("2006-01-02T15:04:05Z")

	response := Location{
		Id:        &created.ID,
		Name:      &created.Name,
		Address:   &created.Address,
		City:      &created.City,
		State:     &created.State,
		Country:   &created.Country,
		Phone:     &created.Phone,
		Email:     &created.Email,
		TaxRate:   floatToPtr(created.TaxRate),
		Timezone:  &created.Timezone,
		IsActive:  &created.IsActive,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	utils.WriteResponse(ctx, w, response, http.StatusCreated)
}

func (s *ServerImpl) GetLocationsId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	l, err := s.locationManager.GetLocationByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Location not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get location: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := l.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := l.UpdatedAt.Format("2006-01-02T15:04:05Z")

	response := Location{
		Id:        &l.ID,
		Name:      &l.Name,
		Address:   &l.Address,
		City:      &l.City,
		State:     &l.State,
		Country:   &l.Country,
		Phone:     &l.Phone,
		Email:     &l.Email,
		TaxRate:   floatToPtr(l.TaxRate),
		Timezone:  &l.Timezone,
		IsActive:  &l.IsActive,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PutLocationsId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req UpdateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	location := models.Location{
		Name:     utils.PtrOrZero(req.Name),
		Address:  utils.PtrOrZero(req.Address),
		City:     utils.PtrOrZero(req.City),
		State:    utils.PtrOrZero(req.State),
		Country:  utils.PtrOrZero(req.Country),
		Phone:    utils.PtrOrZero(req.Phone),
		Email:    utils.PtrOrZero(req.Email),
		TaxRate:  ptrToFloat(req.TaxRate),
		Timezone: utils.PtrOrZero(req.Timezone),
	}

	updated, err := s.locationManager.UpdateLocation(ctx, id, location)
	if err != nil {
		http.Error(w, "Failed to update location: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")

	response := Location{
		Id:        &updated.ID,
		Name:      &updated.Name,
		Address:   &updated.Address,
		City:      &updated.City,
		State:     &updated.State,
		Country:   &updated.Country,
		Phone:     &updated.Phone,
		Email:     &updated.Email,
		TaxRate:   floatToPtr(updated.TaxRate),
		Timezone:  &updated.Timezone,
		IsActive:  &updated.IsActive,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) DeleteLocationsId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := s.locationManager.DeleteLocation(ctx, id); err != nil {
		http.Error(w, "Failed to delete location: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (wrapper *ServerInterfaceWrapper) RegisterLocationsRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermLocationsRead)).Get("/locations", wrapper.GetLocations)
	r.With(middleware.RequirePermission(utils.PermLocationsCreate)).Post("/locations", wrapper.PostLocations)
	r.With(middleware.RequirePermission(utils.PermLocationsRead)).Get("/locations/{id}", wrapper.GetLocationsId)
	r.With(middleware.RequirePermission(utils.PermLocationsUpdate)).Put("/locations/{id}", wrapper.PutLocationsId)
	r.With(middleware.RequirePermission(utils.PermLocationsDelete)).Delete("/locations/{id}", wrapper.DeleteLocationsId)
	return r
}
