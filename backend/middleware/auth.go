package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/azdonald/pharmd/backend/service"
	"github.com/azdonald/pharmd/backend/utils"
)

func AuthMiddleware(userRoleManager service.UserRoleServiceManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path

			if utils.ValidationExemptRoutes[path] {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			claims, err := utils.ExtractClaimFromToken(token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userClaims", claims)
			ctx = context.WithValue(ctx, "organisation_id", claims.BusinessID)
			ctx = context.WithValue(ctx, "user_id", claims.ID)

			permissions, err := userRoleManager.GetUserPermissions(ctx, claims.ID, claims.BusinessID)
			if err != nil {
				log.Printf("[AuthMiddleware] Failed to load permissions for user %s (org %s): %v", claims.ID, claims.BusinessID, err)
				http.Error(w, "Failed to load user permissions", http.StatusInternalServerError)
				return
			}

			permissionSet := make(map[string]struct{}, len(permissions))
			for _, permission := range permissions {
				permissionSet[permission] = struct{}{}
			}
			ctx = context.WithValue(ctx, "permissions", permissionSet)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			permissionSet, ok := r.Context().Value("permissions").(map[string]struct{})
			if !ok {
				http.Error(w, "Permissions context is missing", http.StatusForbidden)
				return
			}
			if _, allowed := permissionSet[permission]; !allowed {
				http.Error(w, "Forbidden: missing permission "+permission, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
