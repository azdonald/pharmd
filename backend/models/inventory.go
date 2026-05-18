package models

import "time"

type StockBatch struct {
	ID               string     `json:"id"`
	OrganisationID   string     `json:"organisation_id"`
	ProductID        string     `json:"product_id"`
	LocationID       string     `json:"location_id"`
	BatchNumber      string     `json:"batch_number"`
	Quantity         int        `json:"quantity"`
	RemainingQty     int        `json:"remaining_qty"`
	UnitCost         float64    `json:"unit_cost"`
	SellingPrice     float64    `json:"selling_price"`
	ManufacturingDate string    `json:"manufacturing_date"`
	ExpiryDate       string     `json:"expiry_date"`
	ReceivedDate     time.Time  `json:"received_date"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type StockMovement struct {
	ID             string    `json:"id"`
	OrganisationID string    `json:"organisation_id"`
	LocationID     string    `json:"location_id"`
	ProductID      string    `json:"product_id"`
	BatchID        string    `json:"batch_id"`
	MovementType   string    `json:"movement_type"`
	Quantity       int       `json:"quantity"`
	ReferenceType  string    `json:"reference_type"`
	ReferenceID    string    `json:"reference_id"`
	Notes          string    `json:"notes"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

type StockCountItem struct {
	ProductID  string `json:"product_id"`
	LocationID string `json:"location_id"`
	BatchID    string `json:"batch_id"`
	CountedQty int    `json:"counted_qty"`
	Notes      string `json:"notes"`
}
