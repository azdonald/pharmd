package models

import "time"

type PurchaseOrder struct {
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	PONumber       string     `json:"po_number"`
	SupplierID     string     `json:"supplier_id"`
	LocationID     string     `json:"location_id"`
	Status         string     `json:"status"`
	OrderDate      time.Time  `json:"order_date"`
	ExpectedDate   string     `json:"expected_date"`
	Notes          string     `json:"notes"`
	Subtotal       float64    `json:"subtotal"`
	TaxTotal       float64    `json:"tax_total"`
	GrandTotal     float64    `json:"grand_total"`
	CreatedBy      string     `json:"created_by"`
	ApprovedBy     string     `json:"approved_by"`
	ApprovedAt     *time.Time `json:"approved_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type PurchaseOrderItem struct {
	ID               string `json:"id"`
	OrganisationID   string `json:"organisation_id"`
	PurchaseOrderID  string `json:"purchase_order_id"`
	ProductID        string `json:"product_id"`
	QuantityOrdered  int    `json:"quantity_ordered"`
	QuantityReceived int    `json:"quantity_received"`
	UnitCost         float64 `json:"unit_cost"`
	LineTotal        float64 `json:"line_total"`
}

type GoodsReceivedNote struct {
	ID              string    `json:"id"`
	OrganisationID  string    `json:"organisation_id"`
	PurchaseOrderID string    `json:"purchase_order_id"`
	ReceivedDate    time.Time `json:"received_date"`
	Notes           string    `json:"notes"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
}
