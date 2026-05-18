package users

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/azdonald/pharmd/backend/middleware"
	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/service"
	"github.com/azdonald/pharmd/backend/utils"
	"github.com/go-chi/chi/v5"
)

type ServerImpl struct {
	userManager     service.UserServiceManager
	userRoleManager service.UserRoleServiceManager
}

func NewServer(userManager service.UserServiceManager, userRoleManager service.UserRoleServiceManager) ServerInterface {
	return &ServerImpl{
		userManager:     userManager,
		userRoleManager: userRoleManager,
	}
}

func (s *ServerImpl) GetUsers(w http.ResponseWriter, r *http.Request, params GetUsersParams) {
	ctx := r.Context()

	page := 1
	limit := 20
	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	users, err := s.userManager.ListUsers(ctx, page, limit)
	if err != nil {
		http.Error(w, "Failed to list users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseUsers := make([]User, len(users))
	for i, u := range users {
		createdAt := u.CreatedAt.Format("2006-01-02T15:04:05Z")
		updatedAt := u.UpdatedAt.Format("2006-01-02T15:04:05Z")
		responseUsers[i] = User{
			Id:        &u.ID,
			FirstName: &u.FirstName,
			LastName:  &u.LastName,
			Email:     &u.Email,
			IsActive:  &u.IsActive,
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		}
	}

	total := len(responseUsers)
	response := UserListResponse{
		Data:  &responseUsers,
		Page:  &page,
		Limit: &limit,
		Total: &total,
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PostUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     string(req.Email),
		Password:  []byte(randomString(8)),
	}

	created, err := s.userManager.CreateUser(ctx, user)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if req.RoleId != nil && *req.RoleId != "" {
		organisationID, ok := ctx.Value("organisation_id").(string)
		if !ok || organisationID == "" {
			http.Error(w, "Failed to assign role: organisation context is required", http.StatusInternalServerError)
			return
		}
		if err := s.userRoleManager.AssignRoleToUser(ctx, created.ID, *req.RoleId, organisationID); err != nil {
			http.Error(w, "Failed to assign role: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := created.UpdatedAt.Format("2006-01-02T15:04:05Z")

	response := User{
		Id:        &created.ID,
		FirstName: &created.FirstName,
		LastName:  &created.LastName,
		Email:     &created.Email,
		IsActive:  &created.IsActive,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	utils.WriteResponse(ctx, w, response, http.StatusCreated)
}

func (s *ServerImpl) GetUsersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	user, err := s.userManager.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := user.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := user.UpdatedAt.Format("2006-01-02T15:04:05Z")

	response := User{
		Id:        &user.ID,
		FirstName: &user.FirstName,
		LastName:  &user.LastName,
		Email:     &user.Email,
		IsActive:  &user.IsActive,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PutUsersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existing, err := s.userManager.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if req.FirstName != nil {
		existing.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		existing.LastName = *req.LastName
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	updated, err := s.userManager.UpdateUser(ctx, id, *existing)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")

	response := User{
		Id:        &updated.ID,
		FirstName: &updated.FirstName,
		LastName:  &updated.LastName,
		Email:     &updated.Email,
		IsActive:  &updated.IsActive,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) DeleteUsersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := s.userManager.DeleteUser(ctx, id); err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerImpl) PutUsersIdRoles(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req PutUsersIdRolesJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	organisationID, ok := ctx.Value("organisation_id").(string)
	if !ok || organisationID == "" {
		http.Error(w, "Organisation context is required", http.StatusInternalServerError)
		return
	}

	if err := s.userRoleManager.AssignRoleToUser(ctx, id, req.RoleId, organisationID); err != nil {
		http.Error(w, "Failed to assign role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(ctx, w, nil, http.StatusOK)
}

func (wrapper *ServerInterfaceWrapper) RegisterUsersRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermUsersRead)).Get("/users", wrapper.GetUsers)
	r.With(middleware.RequirePermission(utils.PermUsersCreate)).Post("/users", wrapper.PostUsers)
	r.With(middleware.RequirePermission(utils.PermUsersRead)).Get("/users/{id}", wrapper.GetUsersId)
	r.With(middleware.RequirePermission(utils.PermUsersUpdate)).Put("/users/{id}", wrapper.PutUsersId)
	r.With(middleware.RequirePermission(utils.PermUsersDelete)).Delete("/users/{id}", wrapper.DeleteUsersId)
	r.With(middleware.RequirePermission(utils.PermUsersUpdate)).Put("/users/{id}/roles", wrapper.PutUsersIdRoles)
	return r
}

func randomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rng.Intn(len(chars))]
	}
	return string(result)
}
