package service

import (
	"context"

	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/repository"
)

type AuthServiceManager interface {
	Register(ctx context.Context, org models.Organisation, user models.User) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetOrganisationByID(ctx context.Context, id string) (*models.Organisation, error)
	CompleteOnboarding(ctx context.Context, orgID string) error
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
}

type UserServiceManager interface {
	ListUsers(ctx context.Context, page, limit int) ([]models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	UpdateUser(ctx context.Context, id string, user models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserRoleServiceManager interface {
	GetUserPermissions(ctx context.Context, userID, orgID string) ([]string, error)
	AssignRoleToUser(ctx context.Context, userID, roleID, orgID string) error
}

type RoleServiceManager interface {
	ListRoles(ctx context.Context, page, limit int) ([]models.Role, error)
	GetRoleByID(ctx context.Context, id string) (*models.Role, error)
	CreateRole(ctx context.Context, role models.Role) (*models.Role, error)
	UpdateRole(ctx context.Context, id string, role models.Role) (*models.Role, error)
	DeleteRole(ctx context.Context, id string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]string, error)
	SetRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error
}

type PermissionServiceManager interface {
	ListPermissions(ctx context.Context) ([]models.Permission, error)
}

type LocationServiceManager interface {
	ListLocations(ctx context.Context, page, limit int) ([]models.Location, error)
	GetLocationByID(ctx context.Context, id string) (*models.Location, error)
	CreateLocation(ctx context.Context, location models.Location) (*models.Location, error)
	UpdateLocation(ctx context.Context, id string, location models.Location) (*models.Location, error)
	DeleteLocation(ctx context.Context, id string) error
}

type ProductCategoryServiceManager interface {
	ListCategories(ctx context.Context) ([]models.ProductCategory, error)
	GetCategoryByID(ctx context.Context, id string) (*models.ProductCategory, error)
	CreateCategory(ctx context.Context, category models.ProductCategory) (*models.ProductCategory, error)
	UpdateCategory(ctx context.Context, id string, category models.ProductCategory) (*models.ProductCategory, error)
	DeleteCategory(ctx context.Context, id string) error
}

type ProductServiceManager interface {
	ListProducts(ctx context.Context, page, limit int, query, categoryID string) ([]models.Product, int, error)
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	GetProductByBarcode(ctx context.Context, barcode string) (*models.Product, error)
	CreateProduct(ctx context.Context, product models.Product) (*models.Product, error)
	ImportProductsCSV(ctx context.Context, orgID string, records [][]string) (int, int, []string)
	UpdateProduct(ctx context.Context, id string, product models.Product) (*models.Product, error)
	DeleteProduct(ctx context.Context, id string) error
	ListSubstitutes(ctx context.Context, productID string) ([]models.GenericSubstitution, error)
	AddSubstitute(ctx context.Context, productID string, sub models.GenericSubstitution) (*models.GenericSubstitution, error)
	RemoveSubstitute(ctx context.Context, productID, substituteID string) error
}

type PurchaseOrderServiceManager interface {
	ListPurchaseOrders(ctx context.Context, page, limit int, status string) ([]models.PurchaseOrder, int, error)
	GetPurchaseOrderByID(ctx context.Context, id string) (*models.PurchaseOrder, error)
	CreatePurchaseOrder(ctx context.Context, po models.PurchaseOrder, items []models.PurchaseOrderItem) (*models.PurchaseOrder, error)
	ApprovePurchaseOrder(ctx context.Context, id string) (*models.PurchaseOrder, error)
	RejectPurchaseOrder(ctx context.Context, id string) (*models.PurchaseOrder, error)
	ReceiveGoods(ctx context.Context, id string, items []models.PurchaseOrderItem, notes string) (*models.PurchaseOrder, error)
}

type SupplierServiceManager interface {
	ListSuppliers(ctx context.Context, page, limit int, query string) ([]models.Supplier, int, error)
	GetSupplierByID(ctx context.Context, id string) (*models.Supplier, error)
	CreateSupplier(ctx context.Context, supplier models.Supplier) (*models.Supplier, error)
	UpdateSupplier(ctx context.Context, id string, supplier models.Supplier) (*models.Supplier, error)
	DeleteSupplier(ctx context.Context, id string) error
	ListSupplierProducts(ctx context.Context, supplierID string) ([]models.SupplierProduct, error)
	SetSupplierProducts(ctx context.Context, supplierID string, products []models.SupplierProduct) ([]models.SupplierProduct, error)
}

type InventoryServiceManager interface {
	CreateBatch(ctx context.Context, batch models.StockBatch) (*models.StockBatch, error)
	ListStock(ctx context.Context, locationID string, page, limit int, query string) ([]repository.InventoryBatchView, int, error)
	CreateAdjustment(ctx context.Context, movement models.StockMovement) (*models.StockMovement, error)
	ListAlerts(ctx context.Context, locationID string) ([]repository.InventoryAlertView, error)
	ListExpiring(ctx context.Context, locationID string, days int) ([]repository.InventoryExpiringView, error)
	StockCount(ctx context.Context, items []models.StockCountItem) (int, error)
}

type OrganisationServiceManager interface {
	GetOrganisationByID(ctx context.Context, id string) (*models.Organisation, error)
	UpdateOrganisation(ctx context.Context, id string, org models.Organisation) error
}

type PatientServiceManager interface {
	ListPatients(ctx context.Context, page, limit int, query string) ([]models.Patient, int, error)
	GetPatientByID(ctx context.Context, id string) (*models.Patient, error)
	CreatePatient(ctx context.Context, patient models.Patient) (*models.Patient, error)
	UpdatePatient(ctx context.Context, id string, patient models.Patient) (*models.Patient, error)
	DeletePatient(ctx context.Context, id string) error
	ListPatientAllergies(ctx context.Context, patientID string) ([]models.PatientAllergy, error)
	AddPatientAllergy(ctx context.Context, patientID string, allergy models.PatientAllergy) (*models.PatientAllergy, error)
	RemovePatientAllergy(ctx context.Context, patientID, allergyID string) error
	ListPatientConditions(ctx context.Context, patientID string) ([]models.PatientCondition, error)
	AddPatientCondition(ctx context.Context, patientID string, condition models.PatientCondition) (*models.PatientCondition, error)
}
