package service

import (
	"context"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
	"github.com/google/uuid"
)

type PrescriberServiceManager interface {
	List(ctx context.Context, query string, page, limit int) ([]models.Prescriber, int, error)
	GetByID(ctx context.Context, id string) (*models.Prescriber, error)
	Create(ctx context.Context, p models.Prescriber) (*models.Prescriber, error)
	Update(ctx context.Context, id string, p models.Prescriber) (*models.Prescriber, error)
	Delete(ctx context.Context, id string) error
}

type PrescriberService struct {
	repo repository.PrescriberRepository
}

func NewPrescriberService(repo repository.PrescriberRepository) PrescriberServiceManager {
	return &PrescriberService{repo: repo}
}

func (s *PrescriberService) List(ctx context.Context, query string, page, limit int) ([]models.Prescriber, int, error) {
	return s.repo.List(ctx, query, page, limit)
}

func (s *PrescriberService) GetByID(ctx context.Context, id string) (*models.Prescriber, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PrescriberService) Create(ctx context.Context, p models.Prescriber) (*models.Prescriber, error) {
	now := time.Now()
	p.ID = uuid.New().String()
	p.IsActive = true
	p.CreatedAt = now
	p.UpdatedAt = now
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *PrescriberService) Update(ctx context.Context, id string, p models.Prescriber) (*models.Prescriber, error) {
	p.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, id, p); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *PrescriberService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

type PrescriptionServiceManager interface {
	List(ctx context.Context, status, patientID string, page, limit int) ([]models.Prescription, int, error)
	GetByID(ctx context.Context, id string) (*models.Prescription, []models.PrescriptionItem, error)
	Create(ctx context.Context, rx models.Prescription, items []models.PrescriptionItem) (*models.Prescription, error)
	Update(ctx context.Context, id string, rx models.Prescription) (*models.Prescription, error)
	Delete(ctx context.Context, id string) error
	RecordRefill(ctx context.Context, rxID, itemID, userID string) (*models.Prescription, error)
}

type PrescriptionService struct {
	repo repository.PrescriptionRepository
}

func NewPrescriptionService(repo repository.PrescriptionRepository) PrescriptionServiceManager {
	return &PrescriptionService{repo: repo}
}

func (s *PrescriptionService) List(ctx context.Context, status, patientID string, page, limit int) ([]models.Prescription, int, error) {
	rows, total, err := s.repo.List(ctx, status, patientID, page, limit)
	if err != nil {
		return nil, 0, err
	}
	out := make([]models.Prescription, len(rows))
	for i, r := range rows {
		r.Prescription.PatientName = r.PatientName
		r.Prescription.PrescriberName = r.PrescriberName
		out[i] = r.Prescription
	}
	return out, total, nil
}

func (s *PrescriptionService) GetByID(ctx context.Context, id string) (*models.Prescription, []models.PrescriptionItem, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PrescriptionService) Create(ctx context.Context, rx models.Prescription, items []models.PrescriptionItem) (*models.Prescription, error) {
	orgID := ctx.Value("organisation_id").(string)
	userID := ctx.Value("user_id").(string)
	now := time.Now()

	rx.ID = uuid.New().String()
	rx.OrganisationID = orgID
	rx.Status = "active"
	rx.CreatedBy = userID
	rx.CreatedAt = now
	rx.UpdatedAt = now

	for i := range items {
		items[i].ID = uuid.New().String()
		items[i].OrganisationID = orgID
		items[i].PrescriptionID = rx.ID
	}

	if err := s.repo.Create(ctx, rx, items); err != nil {
		return nil, err
	}
	return &rx, nil
}

func (s *PrescriptionService) Update(ctx context.Context, id string, rx models.Prescription) (*models.Prescription, error) {
	rx.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, id, rx); err != nil {
		return nil, err
	}
	rx2, _, err := s.repo.GetByID(ctx, id)
	return rx2, err
}

func (s *PrescriptionService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *PrescriptionService) RecordRefill(ctx context.Context, rxID, itemID, userID string) (*models.Prescription, error) {
	refill := models.Refill{
		ID:            uuid.New().String(),
		PrescriptionID: rxID,
		ItemID:        itemID,
		RefilledBy:    userID,
		CreatedAt:     time.Now(),
	}
	if err := s.repo.RecordRefill(ctx, refill); err != nil {
		return nil, err
	}
	if err := s.repo.IncrementRefillUsed(ctx, itemID); err != nil {
		return nil, err
	}
	rx, _, err := s.repo.GetByID(ctx, rxID)
	return rx, err
}
