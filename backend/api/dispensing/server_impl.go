package dispensing

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
	manager service.DispensingServiceManager
}

func NewServer(manager service.DispensingServiceManager) ServerInterface {
	return &ServerImpl{manager: manager}
}

func drToResponse(dr models.DispenseRecord, createdAt, updatedAt string) DispenseRecord {
	qtyDisp := dr.QuantityDispensed
	qtyPres := dr.QuantityPrescribed
	isCtrl := dr.IsControlled
	return DispenseRecord{
		Id:                 &dr.ID,
		PrescriptionId:     &dr.PrescriptionID,
		PrescriptionItemId: &dr.PrescriptionItemID,
		PatientId:          &dr.PatientID,
		PatientName:        &dr.PatientName,
		ProductId:          &dr.ProductID,
		ProductName:        &dr.ProductName,
		LocationId:         &dr.LocationID,
		QuantityDispensed:  &qtyDisp,
		QuantityPrescribed: &qtyPres,
		PharmacistId:       &dr.PharmacistID,
		PharmacistName:     &dr.PharmacistName,
		TechnicianId:       &dr.TechnicianID,
		Status:             &dr.Status,
		Notes:              &dr.Notes,
		WitnessName:        &dr.WitnessName,
		IsControlled:       &isCtrl,
		DispensedAt:        &dr.DispensedAt,
		CreatedAt:          &createdAt,
		UpdatedAt:          &updatedAt,
	}
}

func (s *ServerImpl) GetDispensing(w http.ResponseWriter, r *http.Request, params GetDispensingParams) {
	ctx := r.Context()
	page := 1
	limit := 20
	if params.Page != nil { page = *params.Page }
	if params.Limit != nil { limit = *params.Limit }
	status := ""
	if params.Status != nil { status = *params.Status }
	prescriptionID := ""
	if params.PrescriptionId != nil { prescriptionID = *params.PrescriptionId }

	records, total, err := s.manager.List(ctx, status, prescriptionID, page, limit)
	if err != nil {
		http.Error(w, "Failed to list dispensing records: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]DispenseRecord, len(records))
	for i, dr := range records {
		ca := dr.CreatedAt.Format("2006-01-02T15:04:05Z")
		ua := dr.UpdatedAt.Format("2006-01-02T15:04:05Z")
		resp[i] = drToResponse(dr, ca, ua)
	}

	utils.WriteResponse(ctx, w, DispenseListResponse{
		Data:  &resp,
		Page:  &page,
		Limit: &limit,
		Total: &total,
	}, http.StatusOK)
}

func (s *ServerImpl) PostDispensing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateDispenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dr := models.DispenseRecord{
		PrescriptionItemID: req.PrescriptionItemId,
		QuantityDispensed:  req.QuantityDispensed,
		PharmacistID:       req.PharmacistId,
		TechnicianID:       utils.PtrOrZero(req.TechnicianId),
		Notes:              utils.PtrOrZero(req.Notes),
		WitnessName:        utils.PtrOrZero(req.WitnessName),
	}
	if req.IsControlled != nil {
		dr.IsControlled = *req.IsControlled
	}

	created, err := s.manager.Create(ctx, dr)
	if err != nil {
		http.Error(w, "Failed to create dispense record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ca := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := created.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, drToResponse(*created, ca, ua), http.StatusCreated)
}

func (s *ServerImpl) GetDispensingId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	dr, err := s.manager.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Dispense record not found", http.StatusNotFound)
		return
	}

	ca := dr.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := dr.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, drToResponse(*dr, ca, ua), http.StatusOK)
}

func (s *ServerImpl) PutDispensingId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req UpdateDispenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dr := models.DispenseRecord{
		TechnicianID: utils.PtrOrZero(req.TechnicianId),
		WitnessName:  utils.PtrOrZero(req.WitnessName),
		Notes:        utils.PtrOrZero(req.Notes),
	}

	updated, err := s.manager.Update(ctx, id, dr)
	if err != nil {
		http.Error(w, "Failed to update dispense record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ca := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, drToResponse(*updated, ca, ua), http.StatusOK)
}

func (s *ServerImpl) PutDispensingIdStatus(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req StatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := s.manager.UpdateStatus(ctx, id, string(req.Status))
	if err != nil {
		http.Error(w, "Failed to update status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ca := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, drToResponse(*updated, ca, ua), http.StatusOK)
}

func (s *ServerImpl) GetDispensingIdLabel(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	data, err := s.manager.GetLabelData(ctx, id)
	if err != nil {
		http.Error(w, "Failed to get label data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(ctx, w, data, http.StatusOK)
}

func (s *ServerImpl) GetDrugsProductIdInteractions(w http.ResponseWriter, r *http.Request, productId string, params GetDrugsProductIdInteractionsParams) {
	ctx := r.Context()

	warnings, err := s.manager.CheckInteractions(ctx, productId, params.PatientId)
	if err != nil {
		http.Error(w, "Failed to check interactions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]InteractionWarning, len(warnings))
	for i, w := range warnings {
		msg := w
		sev := Moderate
		resp[i] = InteractionWarning{
			Message:  &msg,
			Severity: &sev,
		}
	}

	utils.WriteResponse(ctx, w, resp, http.StatusOK)
}

func (s *ServerImpl) GetDrugsProductIdAllergyCheck(w http.ResponseWriter, r *http.Request, productId string, params GetDrugsProductIdAllergyCheckParams) {
	ctx := r.Context()

	warnings, err := s.manager.CheckAllergies(ctx, productId, params.PatientId)
	if err != nil {
		http.Error(w, "Failed to check allergies: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]AllergyWarning, len(warnings))
	for i, w := range warnings {
		allergen := w
		sev := "moderate"
		resp[i] = AllergyWarning{
			Allergen: &allergen,
			Severity: &sev,
		}
	}

	utils.WriteResponse(ctx, w, resp, http.StatusOK)
}

func (wrapper *ServerInterfaceWrapper) RegisterDispensingRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermDispensingRead)).Get("/dispensing", wrapper.GetDispensing)
	r.With(middleware.RequirePermission(utils.PermDispensingManage)).Post("/dispensing", wrapper.PostDispensing)
	r.With(middleware.RequirePermission(utils.PermDispensingRead)).Get("/dispensing/{id}", wrapper.GetDispensingId)
	r.With(middleware.RequirePermission(utils.PermDispensingManage)).Put("/dispensing/{id}", wrapper.PutDispensingId)
	r.With(middleware.RequirePermission(utils.PermDispensingManage)).Put("/dispensing/{id}/status", wrapper.PutDispensingIdStatus)
	r.With(middleware.RequirePermission(utils.PermDispensingRead)).Get("/dispensing/{id}/label", wrapper.GetDispensingIdLabel)
	r.With(middleware.RequirePermission(utils.PermDispensingRead)).Get("/drugs/{productId}/interactions", wrapper.GetDrugsProductIdInteractions)
	r.With(middleware.RequirePermission(utils.PermDispensingRead)).Get("/drugs/{productId}/allergy-check", wrapper.GetDrugsProductIdAllergyCheck)
	return r
}
