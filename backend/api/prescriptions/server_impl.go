package prescriptions

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
	rxManager   service.PrescriptionServiceManager
	prManager   service.PrescriberServiceManager
}

func NewServer(rxManager service.PrescriptionServiceManager, prManager service.PrescriberServiceManager) ServerInterface {
	return &ServerImpl{rxManager: rxManager, prManager: prManager}
}

func rxToResponse(rx models.Prescription, createdAt, updatedAt string) Prescription {
	status := rx.Status
	return Prescription{
		Id:             &rx.ID,
		PatientId:      &rx.PatientID,
		PatientName:    &rx.PatientName,
		PrescriberId:   &rx.PrescriberID,
		PrescriberName: &rx.PrescriberName,
		LocationId:     &rx.LocationID,
		Status:         &status,
		Diagnosis:      &rx.Diagnosis,
		Notes:          &rx.Notes,
		IssuedDate:     &rx.IssuedDate,
		ExpiryDate:     &rx.ExpiryDate,
		CreatedBy:      &rx.CreatedBy,
		CreatedAt:      &createdAt,
		UpdatedAt:      &updatedAt,
	}
}

func itemToResponse(item models.PrescriptionItem) PrescriptionItem {
	qty := item.Quantity
	refAuth := item.RefillsAuthorized
	refUsed := item.RefillsUsed
	return PrescriptionItem{
		Id:               &item.ID,
		ProductId:        &item.ProductID,
		ProductName:      &item.ProductName,
		Dosage:           &item.Dosage,
		Frequency:        &item.Frequency,
		Duration:         &item.Duration,
		Quantity:         &qty,
		RefillsAuthorized: &refAuth,
		RefillsUsed:      &refUsed,
		Notes:            &item.Notes,
	}
}

func (s *ServerImpl) GetPrescriptions(w http.ResponseWriter, r *http.Request, params GetPrescriptionsParams) {
	ctx := r.Context()
	page := 1
	limit := 20
	if params.Page != nil { page = *params.Page }
	if params.Limit != nil { limit = *params.Limit }
	status := ""
	if params.Status != nil { status = *params.Status }
	patientID := ""
	if params.PatientId != nil { patientID = *params.PatientId }

	rxs, total, err := s.rxManager.List(ctx, status, patientID, page, limit)
	if err != nil {
		http.Error(w, "Failed to list prescriptions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]Prescription, len(rxs))
	for i, rx := range rxs {
		ca := rx.CreatedAt.Format("2006-01-02T15:04:05Z")
		ua := rx.UpdatedAt.Format("2006-01-02T15:04:05Z")
		resp[i] = rxToResponse(rx, ca, ua)
	}

	utils.WriteResponse(ctx, w, PrescriptionListResponse{
		Data:  &resp,
		Page:  &page,
		Limit: &limit,
		Total: &total,
	}, http.StatusOK)
}

func (s *ServerImpl) PostPrescriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreatePrescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rx := models.Prescription{
		PatientID:    req.PatientId,
		PrescriberID: req.PrescriberId,
		LocationID:   req.LocationId,
		Diagnosis:    utils.PtrOrZero(req.Diagnosis),
		Notes:        utils.PtrOrZero(req.Notes),
		IssuedDate:   utils.PtrOrZero(req.IssuedDate),
		ExpiryDate:   utils.PtrOrZero(req.ExpiryDate),
	}

	items := make([]models.PrescriptionItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = models.PrescriptionItem{
			ProductID:        item.ProductId,
			Dosage:           item.Dosage,
			Frequency:        item.Frequency,
			Duration:         utils.PtrOrZero(item.Duration),
			Quantity:         item.Quantity,
			RefillsAuthorized: utils.PtrOrZeroInt(item.RefillsAuthorized),
			Notes:            utils.PtrOrZero(item.Notes),
		}
	}

	created, err := s.rxManager.Create(ctx, rx, items)
	if err != nil {
		http.Error(w, "Failed to create prescription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ca := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := created.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, rxToResponse(*created, ca, ua), http.StatusCreated)
}

func (s *ServerImpl) GetPrescriptionsId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	rx, items, err := s.rxManager.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Prescription not found", http.StatusNotFound)
		return
	}

	ca := rx.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := rx.UpdatedAt.Format("2006-01-02T15:04:05Z")
	resp := rxToResponse(*rx, ca, ua)
	responseItems := make([]PrescriptionItem, len(items))
	for i, item := range items {
		responseItems[i] = itemToResponse(item)
	}
	resp.Items = &responseItems

	utils.WriteResponse(ctx, w, resp, http.StatusOK)
}

func (s *ServerImpl) PutPrescriptionsId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req UpdatePrescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rx := models.Prescription{
		Diagnosis:  utils.PtrOrZero(req.Diagnosis),
		Notes:      utils.PtrOrZero(req.Notes),
		Status:     utils.PtrOrZero(req.Status),
		ExpiryDate: utils.PtrOrZero(req.ExpiryDate),
	}

	updated, err := s.rxManager.Update(ctx, id, rx)
	if err != nil {
		http.Error(w, "Failed to update prescription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ca := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, rxToResponse(*updated, ca, ua), http.StatusOK)
}

func (s *ServerImpl) DeletePrescriptionsId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	if err := s.rxManager.Delete(ctx, id); err != nil {
		http.Error(w, "Failed to delete prescription: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerImpl) PostPrescriptionsIdRefill(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(string)

	var req RefillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := s.rxManager.RecordRefill(ctx, id, req.ItemId, userID)
	if err != nil {
		http.Error(w, "Failed to record refill: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ca := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	ua := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")
	utils.WriteResponse(ctx, w, rxToResponse(*updated, ca, ua), http.StatusOK)
}

func (wrapper *ServerInterfaceWrapper) RegisterPrescriptionsRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermPrescriptionsRead)).Get("/prescriptions", wrapper.GetPrescriptions)
	r.With(middleware.RequirePermission(utils.PermPrescriptionsCreate)).Post("/prescriptions", wrapper.PostPrescriptions)
	r.With(middleware.RequirePermission(utils.PermPrescriptionsRead)).Get("/prescriptions/{id}", wrapper.GetPrescriptionsId)
	r.With(middleware.RequirePermission(utils.PermPrescriptionsUpdate)).Put("/prescriptions/{id}", wrapper.PutPrescriptionsId)
	r.With(middleware.RequirePermission(utils.PermPrescriptionsDelete)).Delete("/prescriptions/{id}", wrapper.DeletePrescriptionsId)
	r.With(middleware.RequirePermission(utils.PermPrescriptionsUpdate)).Post("/prescriptions/{id}/refill", wrapper.PostPrescriptionsIdRefill)
	return r
}
