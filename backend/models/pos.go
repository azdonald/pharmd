package models

import "time"

type Sale struct {
	ID              string    `json:"id"`
	OrganisationID  string    `json:"organisation_id"`
	LocationID      string    `json:"location_id"`
	PatientID       string    `json:"patient_id"`
	PatientName     string    `json:"patient_name"`
	PrescriptionID  string    `json:"prescription_id"`
	SaleNumber      string    `json:"sale_number"`
	SaleType        string    `json:"sale_type"`
	Status          string    `json:"status"`
	Subtotal        float64   `json:"subtotal"`
	TaxTotal        float64   `json:"tax_total"`
	DiscountTotal   float64   `json:"discount_total"`
	GrandTotal      float64   `json:"grand_total"`
	PaidAmount      float64   `json:"paid_amount"`
	ChangeAmount    float64   `json:"change_amount"`
	Notes           string    `json:"notes"`
	CreatedBy       string    `json:"created_by"`
	VoidedBy        string    `json:"voided_by"`
	VoidedAt        string    `json:"voided_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type SaleItem struct {
	ID             string  `json:"id"`
	OrganisationID string  `json:"organisation_id"`
	SaleID         string  `json:"sale_id"`
	ProductID      string  `json:"product_id"`
	ProductName    string  `json:"product_name"`
	Quantity       int     `json:"quantity"`
	UnitPrice      float64 `json:"unit_price"`
	Discount       float64 `json:"discount"`
	LineTotal      float64 `json:"line_total"`
}

type Payment struct {
	ID             string  `json:"id"`
	OrganisationID string  `json:"organisation_id"`
	SaleID         string  `json:"sale_id"`
	Method         string  `json:"method"`
	Amount         float64 `json:"amount"`
	Reference      string  `json:"reference"`
}

type DailySummary struct {
	ID             string    `json:"id"`
	OrganisationID string    `json:"organisation_id"`
	LocationID     string    `json:"location_id"`
	Date           string    `json:"date"`
	TotalSales     int       `json:"total_sales"`
	TotalRevenue   float64   `json:"total_revenue"`
	TotalTax       float64   `json:"total_tax"`
	TotalDiscounts float64   `json:"total_discounts"`
	IsClosed       bool      `json:"is_closed"`
	ClosedBy       string    `json:"closed_by"`
	ClosedAt       string    `json:"closed_at"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
