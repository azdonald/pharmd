package models

import "time"

type Patient struct {
	ID                   string     `json:"id"`
	OrganisationID       string     `json:"organisation_id"`
	FirstName            string     `json:"first_name"`
	LastName             string     `json:"last_name"`
	DateOfBirth          string     `json:"date_of_birth"`
	Gender               string     `json:"gender"`
	Phone                string     `json:"phone"`
	Email                string     `json:"email"`
	Address              string     `json:"address"`
	City                 string     `json:"city"`
	State                string     `json:"state"`
	Country              string     `json:"country"`
	BloodGroup           string     `json:"blood_group"`
	Genotype             string     `json:"genotype"`
	Notes                string     `json:"notes"`
	EmergencyContactName string     `json:"emergency_contact_name"`
	EmergencyContactPhone string    `json:"emergency_contact_phone"`
	IsActive             bool       `json:"is_active"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty"`
}

type PatientAllergy struct {
	ID        string    `json:"id"`
	PatientID string    `json:"patient_id"`
	Allergy   string    `json:"allergy"`
	Severity  string    `json:"severity"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

type PatientCondition struct {
	ID        string    `json:"id"`
	PatientID string    `json:"patient_id"`
	Condition string    `json:"condition"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}
