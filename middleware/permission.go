package middleware

import (
	"net/http"
	"slices"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal/permission"
	"github.com/aryaadinulfadlan/go-social-api/internal/shared"
	"github.com/google/uuid"
)

func RequirePermission(permissionRepository permission.Repository, permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := shared.GetUserFromContext(r)
			if !CheckUserPermission(permissionRepository, permission, r, user.Id) {
				helpers.ForbiddenError(w, shared.ErrForbidden.Error())
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CheckUserPermission(permissionRepository permission.Repository, permission string, r *http.Request, userID uuid.UUID) bool {
	user := shared.GetUserFromContext(r)
	permissions, err := permissionRepository.GetPermissionNamesByRoleId(r.Context(), user.RoleId)
	if err != nil {
		return false
	}
	if slices.Contains(permissions, permission) {
		return true
	}
	return false
}
