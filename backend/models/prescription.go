package models

import "time"

type Prescriber struct {
	ID            string    `json:"id"`
	OrganisationID string   `json:"organisation_id"`
	Name          string    `json:"name"`
	LicenseNumber string    `json:"license_number"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	Specialty     string    `json:"specialty"`
	DEANumber     string    `json:"dea_number"`
	NPINumber     string    `json:"npi_number"`
	Address       string    `json:"address"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

type Prescription struct {
	ID              string    `json:"id"`
	OrganisationID  string    `json:"organisation_id"`
	PatientID       string    `json:"patient_id"`
	PatientName     string    `json:"patient_name"`
	PrescriberID    string    `json:"prescriber_id"`
	PrescriberName  string    `json:"prescriber_name"`
	LocationID      string    `json:"location_id"`
	Status          string    `json:"status"`
	Diagnosis       string    `json:"diagnosis"`
	Notes           string    `json:"notes"`
	IssuedDate      string    `json:"issued_date"`
	ExpiryDate      string    `json:"expiry_date"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type PrescriptionItem struct {
	ID               string `json:"id"`
	OrganisationID   string `json:"organisation_id"`
	PrescriptionID   string `json:"prescription_id"`
	ProductID        string `json:"product_id"`
	ProductName      string `json:"product_name"`
	Dosage           string `json:"dosage"`
	Frequency        string `json:"frequency"`
	Duration         string `json:"duration"`
	Quantity         int    `json:"quantity"`
	RefillsAuthorized int   `json:"refills_authorized"`
	RefillsUsed      int    `json:"refills_used"`
	Notes            string `json:"notes"`
}

type Refill struct {
	ID             string    `json:"id"`
	OrganisationID string    `json:"organisation_id"`
	PrescriptionID string    `json:"prescription_id"`
	ItemID         string    `json:"item_id"`
	RefilledBy     string    `json:"refilled_by"`
	CreatedAt      time.Time `json:"created_at"`
}
