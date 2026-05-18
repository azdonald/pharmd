package service

import (
	"context"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
	"github.com/google/uuid"
)

type DispensingServiceManager interface {
	List(ctx context.Context, status, prescriptionID string, page, limit int) ([]models.DispenseRecord, int, error)
	GetByID(ctx context.Context, id string) (*models.DispenseRecord, error)
	Create(ctx context.Context, dr models.DispenseRecord) (*models.DispenseRecord, error)
	Update(ctx context.Context, id string, dr models.DispenseRecord) (*models.DispenseRecord, error)
	UpdateStatus(ctx context.Context, id, status string) (*models.DispenseRecord, error)
	CheckInteractions(ctx context.Context, productID, patientID string) ([]string, error)
	CheckAllergies(ctx context.Context, productID, patientID string) ([]string, error)
	GetLabelData(ctx context.Context, id string) (*map[string]interface{}, error)
}

type DispensingService struct {
	repo repository.DispensingRepository
}

func NewDispensingService(repo repository.DispensingRepository) DispensingServiceManager {
	return &DispensingService{repo: repo}
}

func (s *DispensingService) List(ctx context.Context, status, prescriptionID string, page, limit int) ([]models.DispenseRecord, int, error) {
	return s.repo.List(ctx, status, prescriptionID, page, limit)
}

func (s *DispensingService) GetByID(ctx context.Context, id string) (*models.DispenseRecord, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DispensingService) Create(ctx context.Context, dr models.DispenseRecord) (*models.DispenseRecord, error) {
	orgID := ctx.Value("organisation_id").(string)
	now := time.Now()

	dr.ID = uuid.New().String()
	dr.OrganisationID = orgID
	dr.Status = "completed"
	dr.DispensedAt = now.Format("2006-01-02T15:04:05Z")
	dr.CreatedAt = now
	dr.UpdatedAt = now

	if err := s.repo.Create(ctx, dr); err != nil {
		return nil, err
	}
	return &dr, nil
}

func (s *DispensingService) Update(ctx context.Context, id string, dr models.DispenseRecord) (*models.DispenseRecord, error) {
	dr.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, id, dr); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *DispensingService) UpdateStatus(ctx context.Context, id, status string) (*models.DispenseRecord, error) {
	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *DispensingService) CheckInteractions(ctx context.Context, productID, patientID string) ([]string, error) {
	return []string{}, nil
}

func (s *DispensingService) CheckAllergies(ctx context.Context, productID, patientID string) ([]string, error) {
	return []string{}, nil
}

func (s *DispensingService) GetLabelData(ctx context.Context, id string) (*map[string]interface{}, error) {
	dr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{
		"patient_name": dr.PatientName,
		"drug_name":    dr.ProductName,
		"dosage":       "",
		"frequency":    "",
		"quantity":     dr.QuantityDispensed,
		"pharmacist":   dr.PharmacistName,
		"date":         dr.DispensedAt,
		"instructions": dr.Notes,
		"warning_labels": []string{},
	}
	return &data, nil
}
