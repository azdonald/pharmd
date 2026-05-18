package service

import (
	"context"
	"time"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
	"github.com/google/uuid"
)

type PatientService struct {
	patientRepo repository.PatientRepository
}

func NewPatientService(patientRepo repository.PatientRepository) PatientServiceManager {
	return &PatientService{patientRepo: patientRepo}
}

func (s *PatientService) ListPatients(ctx context.Context, page, limit int, query string) ([]models.Patient, int, error) {
	return s.patientRepo.ListPatients(ctx, page, limit, query)
}

func (s *PatientService) GetPatientByID(ctx context.Context, id string) (*models.Patient, error) {
	return s.patientRepo.GetPatientByID(ctx, id)
}

func (s *PatientService) CreatePatient(ctx context.Context, patient models.Patient) (*models.Patient, error) {
	orgID := ctx.Value("organisation_id").(string)
	now := time.Now()

	patient.ID = uuid.New().String()
	patient.OrganisationID = orgID
	patient.IsActive = true
	patient.CreatedAt = now
	patient.UpdatedAt = now

	if err := s.patientRepo.CreatePatient(ctx, patient); err != nil {
		return nil, err
	}
	return &patient, nil
}

func (s *PatientService) UpdatePatient(ctx context.Context, id string, patient models.Patient) (*models.Patient, error) {
	existing, err := s.patientRepo.GetPatientByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if patient.FirstName != "" {
		existing.FirstName = patient.FirstName
	}
	if patient.LastName != "" {
		existing.LastName = patient.LastName
	}
	if patient.DateOfBirth != "" {
		existing.DateOfBirth = patient.DateOfBirth
	}
	if patient.Gender != "" {
		existing.Gender = patient.Gender
	}
	if patient.Phone != "" {
		existing.Phone = patient.Phone
	}
	if patient.Email != "" {
		existing.Email = patient.Email
	}
	if patient.Address != "" {
		existing.Address = patient.Address
	}
	if patient.City != "" {
		existing.City = patient.City
	}
	if patient.State != "" {
		existing.State = patient.State
	}
	if patient.Country != "" {
		existing.Country = patient.Country
	}
	if patient.BloodGroup != "" {
		existing.BloodGroup = patient.BloodGroup
	}
	if patient.Genotype != "" {
		existing.Genotype = patient.Genotype
	}
	if patient.Notes != "" {
		existing.Notes = patient.Notes
	}
	if patient.EmergencyContactName != "" {
		existing.EmergencyContactName = patient.EmergencyContactName
	}
	if patient.EmergencyContactPhone != "" {
		existing.EmergencyContactPhone = patient.EmergencyContactPhone
	}

	if err := s.patientRepo.UpdatePatient(ctx, id, *existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *PatientService) DeletePatient(ctx context.Context, id string) error {
	return s.patientRepo.DeletePatient(ctx, id)
}

func (s *PatientService) ListPatientAllergies(ctx context.Context, patientID string) ([]models.PatientAllergy, error) {
	return s.patientRepo.ListPatientAllergies(ctx, patientID)
}

func (s *PatientService) AddPatientAllergy(ctx context.Context, patientID string, allergy models.PatientAllergy) (*models.PatientAllergy, error) {
	allergy.ID = uuid.New().String()
	allergy.PatientID = patientID
	allergy.CreatedAt = time.Now()

	if err := s.patientRepo.AddPatientAllergy(ctx, allergy); err != nil {
		return nil, err
	}
	return &allergy, nil
}

func (s *PatientService) RemovePatientAllergy(ctx context.Context, patientID, allergyID string) error {
	return s.patientRepo.RemovePatientAllergy(ctx, patientID, allergyID)
}

func (s *PatientService) ListPatientConditions(ctx context.Context, patientID string) ([]models.PatientCondition, error) {
	return s.patientRepo.ListPatientConditions(ctx, patientID)
}

func (s *PatientService) AddPatientCondition(ctx context.Context, patientID string, condition models.PatientCondition) (*models.PatientCondition, error) {
	condition.ID = uuid.New().String()
	condition.PatientID = patientID
	condition.CreatedAt = time.Now()

	if err := s.patientRepo.AddPatientCondition(ctx, condition); err != nil {
		return nil, err
	}
	return &condition, nil
}

var _ PatientServiceManager = (*PatientService)(nil)
