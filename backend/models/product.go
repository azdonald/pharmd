package models

import "time"

type Product struct {
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	CategoryID     string     `json:"category_id"`
	Classification string     `json:"classification"`
	BrandName      string     `json:"brand_name"`
	GenericName    string     `json:"generic_name"`
	Manufacturer   string     `json:"manufacturer"`
	Barcode        string     `json:"barcode"`
	NDC            string     `json:"ndc"`
	UnitOfMeasure  string     `json:"unit_of_measure"`
	Strength       string     `json:"strength"`
	Form           string     `json:"form"`
	ReorderLevel   int        `json:"reorder_level"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type ProductCategory struct {
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	ParentID       string     `json:"parent_id"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type GenericSubstitution struct {
	ID                  string    `json:"id"`
	OrganisationID      string    `json:"organisation_id"`
	ProductID           string    `json:"product_id"`
	SubstituteProductID string    `json:"substitute_product_id"`
	Notes               string    `json:"notes"`
	CreatedAt           time.Time `json:"created_at"`
}
