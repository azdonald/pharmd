package models

import "time"

type ProductPrice struct {
	ID             string    `json:"id"`
	OrganisationID string    `json:"organisation_id"`
	ProductID      string    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	LocationID     string    `json:"location_id"`
	LocationName   string    `json:"location_name"`
	SellingPrice   float64   `json:"selling_price"`
	CostPrice      float64   `json:"cost_price"`
	MinPrice       float64   `json:"min_price"`
	MaxDiscount    float64   `json:"max_discount"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DiscountRule struct {
	ID                string    `json:"id"`
	OrganisationID    string    `json:"organisation_id"`
	Name              string    `json:"name"`
	Type              string    `json:"type"`
	Value             float64   `json:"value"`
	MinOrderValue     float64   `json:"min_order_value"`
	MaxDiscountAmount float64   `json:"max_discount_amount"`
	AppliesTo         string    `json:"applies_to"`
	AppliesToID       string    `json:"applies_to_id"`
	IsActive          bool      `json:"is_active"`
	StartDate         string    `json:"start_date"`
	EndDate           string    `json:"end_date"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
