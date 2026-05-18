package models

import "time"

type DispenseRecord struct {
	ID                 string    `json:"id"`
	OrganisationID     string    `json:"organisation_id"`
	PrescriptionID     string    `json:"prescription_id"`
	PrescriptionItemID string    `json:"prescription_item_id"`
	PatientID          string    `json:"patient_id"`
	PatientName        string    `json:"patient_name"`
	ProductID          string    `json:"product_id"`
	ProductName        string    `json:"product_name"`
	LocationID         string    `json:"location_id"`
	QuantityDispensed  int       `json:"quantity_dispensed"`
	QuantityPrescribed int       `json:"quantity_prescribed"`
	PharmacistID       string    `json:"pharmacist_id"`
	PharmacistName     string    `json:"pharmacist_name"`
	TechnicianID       string    `json:"technician_id"`
	Status             string    `json:"status"`
	Notes              string    `json:"notes"`
	WitnessName        string    `json:"witness_name"`
	IsControlled       bool      `json:"is_controlled"`
	DispensedAt        string    `json:"dispensed_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type DispenseLog struct {
	ID            string    `json:"id"`
	OrganisationID string   `json:"organisation_id"`
	DispenseID    string    `json:"dispense_id"`
	Action        string    `json:"action"`
	ActorID       string    `json:"actor_id"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}
