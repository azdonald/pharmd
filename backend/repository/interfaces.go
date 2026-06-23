package repository

import (
	"context"

	"github.com/azdonald/pharmd/backend/models"
)

type AuthRepository interface {
	CreateOrganisation(ctx context.Context, org models.Organisation) error
	CreateUser(ctx context.Context, user models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetOrganisationByID(ctx context.Context, id string) (*models.Organisation, error)
	UpdateOrganisationOnboarding(ctx context.Context, orgID string) error
	UpdatePassword(ctx context.Context, userID, hashedPassword string) error
}

type UserRepository interface {
	ListUsers(ctx context.Context, page, limit int) ([]models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, user models.User) error
	UpdateUser(ctx context.Context, id string, user models.User) error
	DeleteUser(ctx context.Context, id string) error
}

type UserRoleRepository interface {
	GetUserPermissions(ctx context.Context, userID, orgID string) ([]string, error)
	AssignRoleToUser(ctx context.Context, userID, roleID, orgID string) error
}

type RoleRepository interface {
	ListRoles(ctx context.Context, page, limit int) ([]models.Role, error)
	GetRoleByID(ctx context.Context, id string) (*models.Role, error)
	CreateRole(ctx context.Context, role models.Role) error
	UpdateRole(ctx context.Context, id string, role models.Role) error
	DeleteRole(ctx context.Context, id string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]string, error)
	SetRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error
	CloneSystemRolesForOrg(ctx context.Context, orgID string) (map[string]string, error)
}

type PermissionRepository interface {
	ListPermissions(ctx context.Context) ([]models.Permission, error)
}

type LocationRepository interface {
	ListLocations(ctx context.Context, page, limit int) ([]models.Location, error)
	GetLocationByID(ctx context.Context, id string) (*models.Location, error)
	CreateLocation(ctx context.Context, location models.Location) error
	UpdateLocation(ctx context.Context, id string, location models.Location) error
	DeleteLocation(ctx context.Context, id string) error
}

type ProductCategoryRepository interface {
	ListCategories(ctx context.Context) ([]models.ProductCategory, error)
	GetCategoryByID(ctx context.Context, id string) (*models.ProductCategory, error)
	CreateCategory(ctx context.Context, category models.ProductCategory) error
	UpdateCategory(ctx context.Context, id string, category models.ProductCategory) error
	DeleteCategory(ctx context.Context, id string) error
}

type ProductRepository interface {
	ListProducts(ctx context.Context, page, limit int, query, categoryID string) ([]models.Product, int, error)
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	GetProductByBarcode(ctx context.Context, barcode string) (*models.Product, error)
	CreateProduct(ctx context.Context, product models.Product) error
	BulkCreateProducts(ctx context.Context, products []models.Product) error
	UpdateProduct(ctx context.Context, id string, product models.Product) error
	DeleteProduct(ctx context.Context, id string) error
	ListSubstitutes(ctx context.Context, productID string) ([]models.GenericSubstitution, error)
	AddSubstitute(ctx context.Context, sub models.GenericSubstitution) error
	RemoveSubstitute(ctx context.Context, productID, substituteID string) error
}

type InventoryBatchView struct {
	StockBatch models.StockBatch
	ProductName    string
	BrandName      string
	GenericName    string
	Classification string
	ReorderLevel   int
}

type InventoryAlertView struct {
	ProductID     string
	ProductName   string
	BrandName     string
	TotalQuantity int
	ReorderLevel  int
	LocationID    string
}

type InventoryExpiringView struct {
	StockBatch     models.StockBatch
	ProductName    string
	DaysUntilExpiry int
}

type PurchaseOrderRepository interface {
	ListPurchaseOrders(ctx context.Context, page, limit int, status string) ([]models.PurchaseOrder, int, error)
	GetPurchaseOrderByID(ctx context.Context, id string) (*models.PurchaseOrder, error)
	CreatePurchaseOrder(ctx context.Context, po models.PurchaseOrder, items []models.PurchaseOrderItem) error
	UpdatePOStatus(ctx context.Context, id, status, approvedBy string) error
	ReceiveGoods(ctx context.Context, poID string, items []models.PurchaseOrderItem, notes, userID string) error
}

type PricingRepository interface {
	ListPrices(ctx context.Context, productID, locationID string, page, limit int) ([]models.ProductPrice, int, error)
	GetPriceByID(ctx context.Context, id string) (*models.ProductPrice, error)
	UpsertPrice(ctx context.Context, price models.ProductPrice) error
	DeletePrice(ctx context.Context, id string) error
	ListDiscountRules(ctx context.Context, page, limit int) ([]models.DiscountRule, int, error)
	GetDiscountRuleByID(ctx context.Context, id string) (*models.DiscountRule, error)
	CreateDiscountRule(ctx context.Context, rule models.DiscountRule) error
	UpdateDiscountRule(ctx context.Context, id string, rule models.DiscountRule) error
	DeleteDiscountRule(ctx context.Context, id string) error
}

type SupplierRepository interface {
	ListSuppliers(ctx context.Context, page, limit int, query string) ([]models.Supplier, int, error)
	GetSupplierByID(ctx context.Context, id string) (*models.Supplier, error)
	CreateSupplier(ctx context.Context, supplier models.Supplier) error
	UpdateSupplier(ctx context.Context, id string, supplier models.Supplier) error
	DeleteSupplier(ctx context.Context, id string) error
	ListSupplierProducts(ctx context.Context, supplierID string) ([]models.SupplierProduct, error)
	SetSupplierProducts(ctx context.Context, supplierID string, products []models.SupplierProduct) error
}

type InventoryRepository interface {
	CreateBatch(ctx context.Context, batch models.StockBatch) error
	GetBatchByID(ctx context.Context, id string) (*models.StockBatch, error)
	UpdateBatchQty(ctx context.Context, id string, remainingQty int) error
	ListStock(ctx context.Context, locationID string, page, limit int, query string) ([]InventoryBatchView, int, error)
	ListAlerts(ctx context.Context, locationID string) ([]InventoryAlertView, error)
	ListExpiring(ctx context.Context, locationID string, days int) ([]InventoryExpiringView, error)
	CreateMovement(ctx context.Context, movement models.StockMovement) error
	StockCount(ctx context.Context, items []models.StockCountItem, userID string) (int, error)
}

type PatientRepository interface {
	ListPatients(ctx context.Context, page, limit int, query string) ([]models.Patient, int, error)
	GetPatientByID(ctx context.Context, id string) (*models.Patient, error)
	CreatePatient(ctx context.Context, patient models.Patient) error
	UpdatePatient(ctx context.Context, id string, patient models.Patient) error
	DeletePatient(ctx context.Context, id string) error
	ListPatientAllergies(ctx context.Context, patientID string) ([]models.PatientAllergy, error)
	AddPatientAllergy(ctx context.Context, allergy models.PatientAllergy) error
	RemovePatientAllergy(ctx context.Context, patientID, allergyID string) error
	ListPatientConditions(ctx context.Context, patientID string) ([]models.PatientCondition, error)
	AddPatientCondition(ctx context.Context, condition models.PatientCondition) error
}
