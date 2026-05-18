package roles

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/azdonald/pharmd/backend/middleware"
	"github.com/azdonald/pharmd/backend/models"
	"github.com/azdonald/pharmd/backend/service"
	"github.com/azdonald/pharmd/backend/utils"
	"github.com/go-chi/chi/v5"
)

type ServerImpl struct {
	roleManager service.RoleServiceManager
}

func NewServer(roleManager service.RoleServiceManager) ServerInterface {
	return &ServerImpl{roleManager: roleManager}
}

func (s *ServerImpl) GetRoles(w http.ResponseWriter, r *http.Request, params GetRolesParams) {
	ctx := r.Context()

	page := 1
	limit := 20
	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	roles, err := s.roleManager.ListRoles(ctx, page, limit)
	if err != nil {
		http.Error(w, "Failed to list roles: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseRoles := make([]Role, len(roles))
	for i, role := range roles {
		createdAt := role.CreatedAt.Format("2006-01-02T15:04:05Z")
		updatedAt := role.UpdatedAt.Format("2006-01-02T15:04:05Z")
		desc := ""
		if role.Description != nil {
			desc = *role.Description
		}
		responseRoles[i] = Role{
			Id:          &role.ID,
			Name:        &role.Name,
			Slug:        &role.Slug,
			Description: &desc,
			IsSystem:    &role.IsSystem,
			CreatedAt:   &createdAt,
			UpdatedAt:   &updatedAt,
		}
	}

	total := len(responseRoles)
	response := RoleListResponse{
		Data:  &responseRoles,
		Page:  &page,
		Limit: &limit,
		Total: &total,
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PostRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	role := models.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	created, err := s.roleManager.CreateRole(ctx, role)
	if err != nil {
		http.Error(w, "Failed to create role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if req.PermissionIds != nil {
		if err := s.roleManager.SetRolePermissions(ctx, created.ID, *req.PermissionIds); err != nil {
			http.Error(w, "Failed to set permissions: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	createdAt := created.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := created.UpdatedAt.Format("2006-01-02T15:04:05Z")
	desc := ""
	if created.Description != nil {
		desc = *created.Description
	}

	response := Role{
		Id:          &created.ID,
		Name:        &created.Name,
		Slug:        &created.Slug,
		Description: &desc,
		IsSystem:    &created.IsSystem,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}

	utils.WriteResponse(ctx, w, response, http.StatusCreated)
}

func (s *ServerImpl) GetRolesId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	role, err := s.roleManager.GetRoleByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Role not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := role.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := role.UpdatedAt.Format("2006-01-02T15:04:05Z")
	desc := ""
	if role.Description != nil {
		desc = *role.Description
	}

	response := Role{
		Id:          &role.ID,
		Name:        &role.Name,
		Slug:        &role.Slug,
		Description: &desc,
		IsSystem:    &role.IsSystem,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) PutRolesId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	role := models.Role{
		Name:        utils.PtrOrZero(req.Name),
		Description: req.Description,
	}

	updated, err := s.roleManager.UpdateRole(ctx, id, role)
	if err != nil {
		http.Error(w, "Failed to update role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdAt := updated.CreatedAt.Format("2006-01-02T15:04:05Z")
	updatedAt := updated.UpdatedAt.Format("2006-01-02T15:04:05Z")
	desc := ""
	if updated.Description != nil {
		desc = *updated.Description
	}

	response := Role{
		Id:          &updated.ID,
		Name:        &updated.Name,
		Slug:        &updated.Slug,
		Description: &desc,
		IsSystem:    &updated.IsSystem,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}

	utils.WriteResponse(ctx, w, response, http.StatusOK)
}

func (s *ServerImpl) DeleteRolesId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := s.roleManager.DeleteRole(ctx, id); err != nil {
		http.Error(w, "Failed to delete role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerImpl) GetRolesIdPermissions(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	permIDs, err := s.roleManager.GetRolePermissions(ctx, id)
	if err != nil {
		http.Error(w, "Failed to get role permissions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(ctx, w, map[string]interface{}{"permission_ids": permIDs}, http.StatusOK)
}

func (s *ServerImpl) PutRolesIdPermissions(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req PutRolesIdPermissionsJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.roleManager.SetRolePermissions(ctx, id, req.PermissionIds); err != nil {
		http.Error(w, "Failed to set permissions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(ctx, w, nil, http.StatusOK)
}

func (wrapper *ServerInterfaceWrapper) RegisterRolesRoutes(r *chi.Mux) http.Handler {
	r.With(middleware.RequirePermission(utils.PermRolesRead)).Get("/roles", wrapper.GetRoles)
	r.With(middleware.RequirePermission(utils.PermRolesCreate)).Post("/roles", wrapper.PostRoles)
	r.With(middleware.RequirePermission(utils.PermRolesRead)).Get("/roles/{id}", wrapper.GetRolesId)
	r.With(middleware.RequirePermission(utils.PermRolesUpdate)).Put("/roles/{id}", wrapper.PutRolesId)
	r.With(middleware.RequirePermission(utils.PermRolesDelete)).Delete("/roles/{id}", wrapper.DeleteRolesId)
	r.With(middleware.RequirePermission(utils.PermRolesRead)).Get("/roles/{id}/permissions", wrapper.GetRolesIdPermissions)
	r.With(middleware.RequirePermission(utils.PermRolesUpdate)).Put("/roles/{id}/permissions", wrapper.PutRolesIdPermissions)
	return r
}
