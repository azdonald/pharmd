package patients

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
	patientManager service.PatientServiceManager
}

func NewServer(patientManager service.PatientServiceManager) ServerInterface {
	return &ServerImpl{patientManager: patientManager}
}

func (s *ServerImpl) GetPatients(w http.ResponseWriter, r *http.Request, params GetPatientsParams) {
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

	patients, total, err := s.patientManager.ListPatients(ctx, page, limit, query)
	if err != nil {
		http.Error(w, "Failed to list patients: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responsePatients := make([]Patient, len(patients))
	for i, p := range patients {
		createdAt := p.CreatedAt.Format("2006-01-02T15:04:05Z")
		updatedAt := p.UpdatedAt.Format("2006-01-02T15:04:05Z")
		responsePatients[i] = patientToResponse(p, createdAt, updatedAt)
	}

	response := PatientListResponse{
		Data:  &responsePatients,
		Page:  &page,
		Limit: &limit,
		Total: &total,
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PostPatients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreatePatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	patient := models.Patient{
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		DateOfBirth:          utils.PtrOrZero(req.DateOfBirth),
		Gender:               utils.PtrOrZero(req.Gender),
		Phone:                utils.PtrOrZero(req.Phone),
		Email:                utils.PtrOrZero(req.Email),
		Address:              utils.PtrOrZero(req.Address),
		City:                 utils.PtrOrZero(req.City),
		State:                utils.PtrOrZero(req.State),
		Country:              utils.PtrOrZero(req.Country),
		BloodGroup:           utils.PtrOrZero(req.BloodGroup),
		Genotype:             utils.PtrOrZero(req.Genotype),
		Notes:                utils.PtrOrZero(req.Notes),
		EmergencyContactName: utils.PtrOrZero(req.EmergencyContactName),
		EmergencyContactPhone: utils.PtrOrZero(req.EmergencyContactPhone),
	}

	created, err := s.patientManager.CreatePatient(ctx, patient)
	if err != nil {
		http.Error(w, "Failed to create patient: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := created.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, patientToResponse(*created, createdAt, updatedAt), http.StatusCreated)
}

func (s *ServerImpl) GetPatientsId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	p, err := s.patientManager.GetPatientByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Patient not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get patient: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := p.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := p.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, patientToResponse(*p, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) PutPatientsId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req UpdatePatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	patient := models.Patient{
		FirstName:            utils.PtrOrZero(req.FirstName),
		LastName:             utils.PtrOrZero(req.LastName),
		DateOfBirth:          utils.PtrOrZero(req.DateOfBirth),
		Gender:               utils.PtrOrZero(req.Gender),
		Phone:                utils.PtrOrZero(req.Phone),
		Email:                utils.PtrOrZero(req.Email),
		Address:              utils.PtrOrZero(req.Address),
		City:                 utils.PtrOrZero(req.City),
		State:                utils.PtrOrZero(req.State),
		Country:              utils.PtrOrZero(req.Country),
		BloodGroup:           utils.PtrOrZero(req.BloodGroup),
		Genotype:             utils.PtrOrZero(req.Genotype),
		Notes:                utils.PtrOrZero(req.Notes),
		EmergencyContactName: utils.PtrOrZero(req.EmergencyContactName),
		EmergencyContactPhone: utils.PtrOrZero(req.EmergencyContactPhone),
	}
	if req.IsActive != nil {
		patient.IsActive = *req.IsActive
	}

	updated, err := s.patientManager.UpdatePatient(ctx, id, patient)
	if err != nil {
		http.Error(w, "Failed to update patient: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, patientToResponse(*updated, createdAt, updatedAt), http.StatusOK)
}

func (s *ServerImpl) DeletePatientsId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := s.patientManager.DeletePatient(ctx, id); err != nil {
		http.Error(w, "Failed to delete patient: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerImpl) GetPatientsIdAllergies(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	allergies, err := s.patientManager.ListPatientAllergies(ctx, id)
	if err != nil {
		http.Error(w, "Failed to list allergies: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]PatientAllergy, len(allergies))
	for i, a := range allergies {
		createdAt := a.CreatedAt.Format("2006-01-02T15:04:05Z")
		response[i] = PatientAllergy{
			Allergy:   &a.Allergy,
			CreatedAt: &createdAt,
			Id:        &a.ID,
			Notes:     &a.Notes,
			PatientId: &a.PatientID,
			Severity:  &a.Severity,
		}
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PostPatientsIdAllergies(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req CreatePatientAllergyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	allergy := models.PatientAllergy{
		Allergy:  req.Allergy,
		Severity: utils.PtrOrZero(req.Severity),
		Notes:    utils.PtrOrZero(req.Notes),
	}

	created, err := s.patientManager.AddPatientAllergy(ctx, id, allergy)
	if err != nil {
		http.Error(w, "Failed to add allergy: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	response := PatientAllergy{
		Allergy:   &created.Allergy,
		CreatedAt: &createdAt,
		Id:        &created.ID,
		Notes:     &created.Notes,
		PatientId: &created.PatientID,
		Severity:  &created.Severity,
	}

	utils.WriteResponse(ctx, w, response, http.StatusCreated)
}

func (s *ServerImpl) DeletePatientsIdAllergies(w http.ResponseWriter, r *http.Request, id string, params DeletePatientsIdAllergiesParams) {
	ctx := r.Context()

	if err := s.patientManager.RemovePatientAllergy(ctx, id, params.AllergyId); err != nil {
		http.Error(w, "Failed to remove allergy: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerImpl) GetPatientsIdConditions(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	conditions, err := s.patientManager.ListPatientConditions(ctx, id)
	if err != nil {
		http.Error(w, "Failed to list conditions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]PatientCondition, len(conditions))
	for i, c := range conditions {
		createdAt := c.CreatedAt.Format("2006-01-02T15:04:05Z")
		response[i] = PatientCondition{
			Condition: &c.Condition,
			CreatedAt: &createdAt,
			Id:        &c.ID,
			Notes:     &c.Notes,
			PatientId: &c.PatientID,
		}
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PostPatientsIdConditions(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req CreatePatientConditionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	condition := models.PatientCondition{
		Condition: req.Condition,
		Notes:     utils.PtrOrZero(req.Notes),
	}

	created, err := s.patientManager.AddPatientCondition(ctx, id, condition)
	if err != nil {
		http.Error(w, "Failed to add condition: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	response := PatientCondition{
		Condition: &created.Condition,
		CreatedAt: &createdAt,
		Id:        &created.ID,
		Notes:     &created.Notes,
		PatientId: &created.PatientID,
	}

	utils.WriteResponse(ctx, w, response, http.StatusCreated)
}

func patientToResponse(p models.Patient, createdAt, updatedAt string) Patient {
	return Patient{
		Id:                    &p.ID,
		FirstName:             &p.FirstName,
		LastName:              &p.LastName,
		DateOfBirth:           &p.DateOfBirth,
		Gender:                &p.Gender,
		Phone:                 &p.Phone,
		Email:                 &p.Email,
		Address:               &p.Address,
		City:                  &p.City,
		State:                 &p.State,
		Country:               &p.Country,
		BloodGroup:            &p.BloodGroup,
		Genotype:              &p.Genotype,
		Notes:                 &p.Notes,
		EmergencyContactName:  &p.EmergencyContactName,
		EmergencyContactPhone: &p.EmergencyContactPhone,
		IsActive:              &p.IsActive,
		CreatedAt:             &createdAt,
		UpdatedAt:             &updatedAt,
	}
}

func (wrapper *ServerInterfaceWrapper) RegisterPatientsRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermPatientsRead)).Get("/patients", wrapper.GetPatients)
	r.With(middleware.RequirePermission(utils.PermPatientsCreate)).Post("/patients", wrapper.PostPatients)
	r.With(middleware.RequirePermission(utils.PermPatientsRead)).Get("/patients/{id}", wrapper.GetPatientsId)
	r.With(middleware.RequirePermission(utils.PermPatientsUpdate)).Put("/patients/{id}", wrapper.PutPatientsId)
	r.With(middleware.RequirePermission(utils.PermPatientsDelete)).Delete("/patients/{id}", wrapper.DeletePatientsId)
	r.With(middleware.RequirePermission(utils.PermPatientsRead)).Get("/patients/{id}/allergies", wrapper.GetPatientsIdAllergies)
	r.With(middleware.RequirePermission(utils.PermPatientsCreate)).Post("/patients/{id}/allergies", wrapper.PostPatientsIdAllergies)
	r.With(middleware.RequirePermission(utils.PermPatientsDelete)).Delete("/patients/{id}/allergies", wrapper.DeletePatientsIdAllergies)
	r.With(middleware.RequirePermission(utils.PermPatientsRead)).Get("/patients/{id}/conditions", wrapper.GetPatientsIdConditions)
	r.With(middleware.RequirePermission(utils.PermPatientsCreate)).Post("/patients/{id}/conditions", wrapper.PostPatientsIdConditions)
	return r
}
