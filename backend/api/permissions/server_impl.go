package permissions

import (
	"net/http"

	"github.com/azdonald/pharmd/backend/service"
	"github.com/azdonald/pharmd/backend/utils"
	"github.com/go-chi/chi/v5"
)

type serverImpl struct {
	permManager service.PermissionServiceManager
}

func NewServer(permManager service.PermissionServiceManager) ServerInterface {
	return &serverImpl{permManager: permManager}
}

func (s *serverImpl) GetPermissions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	perms, err := s.permManager.ListPermissions(ctx)
	if err != nil {
		http.Error(w, "Failed to list permissions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responsePerms := make([]Permission, len(perms))
	for i, p := range perms {
		desc := ""
		if p.Description != nil {
			desc = *p.Description
		}
		responsePerms[i] = Permission{
			Id:          &p.ID,
			Name:        &p.Name,
			Slug:        &p.Slug,
			Description: &desc,
		}
	}

	utils.WriteResponse(ctx, w, map[string]interface{}{"data": responsePerms}, http.StatusOK)
}

func (s *serverImpl) GetPermissionsId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	perms, err := s.permManager.ListPermissions(ctx)
	if err != nil {
		http.Error(w, "Failed to list permissions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, p := range perms {
		if p.ID == id {
			desc := ""
			if p.Description != nil {
				desc = *p.Description
			}
			response := Permission{
				Id:          &p.ID,
				Name:        &p.Name,
				Slug:        &p.Slug,
				Description: &desc,
			}
			utils.WriteResponse(ctx, w, response, http.StatusOK)
			return
		}
	}

	http.Error(w, "Permission not found", http.StatusNotFound)
}

func (wrapper *ServerInterfaceWrapper) RegisterPermissionsRoutes(r *chi.Mux) http.Handler {
	r.Get("/permissions", wrapper.GetPermissions)
	r.Get("/permissions/{id}", wrapper.GetPermissionsId)
	return r
}
