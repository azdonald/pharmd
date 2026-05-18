package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/azdonald/pharmd/backend/api/auth"
	"github.com/azdonald/pharmd/backend/api/locations"
	"github.com/azdonald/pharmd/backend/api/patients"
	"github.com/azdonald/pharmd/backend/api/permissions"
	"github.com/azdonald/pharmd/backend/api/roles"
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
	patients    *patients.ServerInterfaceWrapper
	userRoles   service.UserRoleServiceManager
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
	permRepo := repository.NewPermissionRepositoryImpl(database)

	authSvc := service.NewAuthService(authRepo, userRepo, locationRepo, userRoleRepo)
	userSvc := service.NewUserService(userRepo)
	roleSvc := service.NewRoleService(roleRepo)
	permSvc := service.NewPermissionService(permRepo)
	locationSvc := service.NewLocationService(locationRepo)
	patientSvc := service.NewPatientService(patientRepo)
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

	return &app{
		auth:        authWrapper,
		users:       usersWrapper,
		roles:       rolesWrapper,
		permissions: permWrapper,
		locations:   locationWrapper,
		patients:    patientWrapper,
		userRoles:   userRoleSvc,
	}
}
