package models

import "time"

type Supplier struct {
	ID            string     `json:"id"`
	OrganisationID string    `json:"organisation_id"`
	Name          string     `json:"name"`
	ContactPerson string     `json:"contact_person"`
	Phone         string     `json:"phone"`
	Email         string     `json:"email"`
	Address       string     `json:"address"`
	City          string     `json:"city"`
	State         string     `json:"state"`
	Country       string     `json:"country"`
	PaymentTerms  string     `json:"payment_terms"`
	Notes         string     `json:"notes"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

type SupplierProduct struct {
	ID           string    `json:"id"`
	OrganisationID string   `json:"organisation_id"`
	SupplierID   string    `json:"supplier_id"`
	ProductID    string    `json:"product_id"`
	UnitPrice    float64   `json:"unit_price"`
	MinOrderQty  int       `json:"min_order_qty"`
	LeadTimeDays int       `json:"lead_time_days"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
