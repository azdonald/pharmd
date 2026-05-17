package models

import "time"

type Location struct {
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Name           string     `json:"name"`
	Address        string     `json:"address"`
	City           string     `json:"city"`
	State          string     `json:"state"`
	Country        string     `json:"country"`
	Phone          string     `json:"phone"`
	Email          string     `json:"email"`
	TaxRate        float64    `json:"tax_rate"`
	Timezone       string     `json:"timezone"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
