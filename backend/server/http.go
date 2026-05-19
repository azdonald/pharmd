package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/azdonald/pharmd/backend/api/auth"
	"github.com/azdonald/pharmd/backend/api/dispensing"
	"github.com/azdonald/pharmd/backend/api/inventory"
	"github.com/azdonald/pharmd/backend/api/locations"
	"github.com/azdonald/pharmd/backend/api/patients"
	"github.com/azdonald/pharmd/backend/api/permissions"
	"github.com/azdonald/pharmd/backend/api/pos"
	"github.com/azdonald/pharmd/backend/api/prescribers"
	"github.com/azdonald/pharmd/backend/api/prescriptions"
	"github.com/azdonald/pharmd/backend/api/pricing"
	"github.com/azdonald/pharmd/backend/api/product_categories"
	"github.com/azdonald/pharmd/backend/api/products"
	"github.com/azdonald/pharmd/backend/api/purchases"
	"github.com/azdonald/pharmd/backend/api/roles"
	"github.com/azdonald/pharmd/backend/api/suppliers"
	"github.com/azdonald/pharmd/backend/api/users"
	"github.com/azdonald/pharmd/backend/db"
	"github.com/azdonald/pharmd/backend/middleware"
	_ "github.com/azdonald/pharmd/backend/migrations"
	"github.com/azdonald/pharmd/backend/repository"
	"github.com/azdonald/pharmd/backend/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/ichtrojan/olympian"
)

type app struct {
	auth        *auth.ServerInterfaceWrapper
	users       *users.ServerInterfaceWrapper
	roles       *roles.ServerInterfaceWrapper
	permissions *permissions.ServerInterfaceWrapper
	locations   *locations.ServerInterfaceWrapper
	patients            *patients.ServerInterfaceWrapper
	inventory           *inventory.ServerInterfaceWrapper
	suppliers           *suppliers.ServerInterfaceWrapper
	purchases           *purchases.ServerInterfaceWrapper
	pricing             *pricing.ServerInterfaceWrapper
	prescribers         *prescribers.ServerInterfaceWrapper
	prescriptions       *prescriptions.ServerInterfaceWrapper
	dispensing          *dispensing.ServerInterfaceWrapper
	pos                 *pos.ServerInterfaceWrapper
	productCategories   *product_categories.ServerInterfaceWrapper
	products            *products.ServerInterfaceWrapper
	userRoles           service.UserRoleServiceManager
}

func (a *app) start(serverPort string) {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.AuthMiddleware(a.userRoles))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	a.auth.RegisterAuthRoutes(r)
	a.users.RegisterUsersRoutes(r)
	a.roles.RegisterRolesRoutes(r)
	a.permissions.RegisterPermissionsRoutes(r)
	a.locations.RegisterLocationsRoutes(r)
	a.patients.RegisterPatientsRoutes(r)
	a.productCategories.RegisterProductCategoriesRoutes(r)
	a.products.RegisterProductsRoutes(r)
	a.inventory.RegisterInventoryRoutes(r)
	a.suppliers.RegisterSuppliersRoutes(r)
	a.purchases.RegisterPurchasesRoutes(r)
	a.pricing.RegisterPricingRoutes(r)
	a.prescribers.RegisterPrescribersRoutes(r)
	a.prescriptions.RegisterPrescriptionsRoutes(r)
	a.dispensing.RegisterDispensingRoutes(r)
	a.pos.RegisterPOSRoutes(r)

	log.Println("Starting server on port", serverPort)
	if err := http.ListenAndServe(":"+serverPort, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func Run() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	migrator := olympian.NewMigrator(database, olympian.MySQL())
	if err := migrator.Init(); err != nil {
		log.Fatal("Failed to initialize migrations:", err)
	}
	if err := migrator.Migrate(olympian.GetMigrations()); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	appInstance := initDependencies(database)

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8980"
	}

	appInstance.start(serverPort)
}

func initDependencies(database *sql.DB) *app {
	authRepo := repository.NewAuthRepositoryImpl(database)
	userRepo := repository.NewUserRepositoryImpl(database)
	userRoleRepo := repository.NewUserRoleRepositoryImpl(database)
	roleRepo := repository.NewRoleRepositoryImpl(database)
	locationRepo := repository.NewLocationRepositoryImpl(database)
	patientRepo := repository.NewPatientRepositoryImpl(database)
	inventoryRepo := repository.NewInventoryRepositoryImpl(database)
	purchaseRepo := repository.NewPurchaseOrderRepositoryImpl(database)
	pricingRepo := repository.NewPricingRepositoryImpl(database)
	prescriberRepo := repository.NewPrescriberRepositoryImpl(database)
	prescriptionRepo := repository.NewPrescriptionRepositoryImpl(database)
	dispensingRepo := repository.NewDispensingRepositoryImpl(database)
	posRepo := repository.NewPOSRepositoryImpl(database)
	supplierRepo := repository.NewSupplierRepositoryImpl(database)
	categoryRepo := repository.NewProductCategoryRepositoryImpl(database)
	productRepo := repository.NewProductRepositoryImpl(database)
	permRepo := repository.NewPermissionRepositoryImpl(database)

	authSvc := service.NewAuthService(authRepo, userRepo, locationRepo, userRoleRepo)
	userSvc := service.NewUserService(userRepo)
	roleSvc := service.NewRoleService(roleRepo)
	permSvc := service.NewPermissionService(permRepo)
	locationSvc := service.NewLocationService(locationRepo)
	patientSvc := service.NewPatientService(patientRepo)
	inventorySvc := service.NewInventoryService(inventoryRepo)
	purchaseSvc := service.NewPurchaseOrderService(purchaseRepo)
	pricingSvc := service.NewPricingService(pricingRepo)
	prescriberSvc := service.NewPrescriberService(prescriberRepo)
	prescriptionSvc := service.NewPrescriptionService(prescriptionRepo)
	dispensingSvc := service.NewDispensingService(dispensingRepo)
	posSvc := service.NewPOSService(posRepo)
	supplierSvc := service.NewSupplierService(supplierRepo)
	productCategorySvc := service.NewProductCategoryService(categoryRepo)
	productSvc := service.NewProductService(productRepo)
	userRoleSvc := service.NewUserRoleService(userRoleRepo, userRepo, roleRepo)

	authServer := auth.NewServer(authSvc)
	authWrapper := &auth.ServerInterfaceWrapper{Handler: authServer}

	usersServer := users.NewServer(userSvc, userRoleSvc)
	usersWrapper := &users.ServerInterfaceWrapper{Handler: usersServer}

	rolesServer := roles.NewServer(roleSvc)
	rolesWrapper := &roles.ServerInterfaceWrapper{Handler: rolesServer}

	permServer := permissions.NewServer(permSvc)
	permWrapper := &permissions.ServerInterfaceWrapper{Handler: permServer}

	locationServer := locations.NewServer(locationSvc)
	locationWrapper := &locations.ServerInterfaceWrapper{Handler: locationServer}

	patientServer := patients.NewServer(patientSvc)
	patientWrapper := &patients.ServerInterfaceWrapper{Handler: patientServer}

	inventoryServer := inventory.NewServer(inventorySvc)
	inventoryWrapper := &inventory.ServerInterfaceWrapper{Handler: inventoryServer}

	purchaseServer := purchases.NewServer(purchaseSvc)
	purchaseWrapper := &purchases.ServerInterfaceWrapper{Handler: purchaseServer}

	pricingServer := pricing.NewServer(pricingSvc)
	pricingWrapper := &pricing.ServerInterfaceWrapper{Handler: pricingServer}

	prescriberServer := prescribers.NewServer(prescriberSvc)
	prescriberWrapper := &prescribers.ServerInterfaceWrapper{Handler: prescriberServer}

	prescriptionServer := prescriptions.NewServer(prescriptionSvc, prescriberSvc)
	prescriptionWrapper := &prescriptions.ServerInterfaceWrapper{Handler: prescriptionServer}

	dispensingServer := dispensing.NewServer(dispensingSvc)
	dispensingWrapper := &dispensing.ServerInterfaceWrapper{Handler: dispensingServer}

	posServer := pos.NewServer(posSvc)
	posWrapper := &pos.ServerInterfaceWrapper{Handler: posServer}

	supplierServer := suppliers.NewServer(supplierSvc)
	supplierWrapper := &suppliers.ServerInterfaceWrapper{Handler: supplierServer}

	categoryServer := product_categories.NewServer(productCategorySvc)
	categoryWrapper := &product_categories.ServerInterfaceWrapper{Handler: categoryServer}

	productServer := products.NewServer(productSvc)
	productWrapper := &products.ServerInterfaceWrapper{Handler: productServer}

	return &app{
		auth:        authWrapper,
		users:       usersWrapper,
		roles:       rolesWrapper,
		permissions: permWrapper,
		locations:   locationWrapper,
		patients:            patientWrapper,
		inventory:           inventoryWrapper,
		suppliers:           supplierWrapper,
		purchases:           purchaseWrapper,
		pricing:             pricingWrapper,
		prescribers:         prescriberWrapper,
		prescriptions:       prescriptionWrapper,
		dispensing:          dispensingWrapper,
		pos:                 posWrapper,
		productCategories:   categoryWrapper,
		products:            productWrapper,
		userRoles:           userRoleSvc,
	}
}
