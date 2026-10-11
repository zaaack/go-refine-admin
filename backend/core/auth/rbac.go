package auth

import (
	"net/http"

	ii18n "go-refine-admin/backend/core/i18n"
	"go-refine-admin/backend/core/response"

	"github.com/gin-gonic/gin"
)

// RequireRoles is a lightweight RBAC gate. The "admin" role always passes.
// Attach it after auth.Middleware() to protect a route by role.
func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		claims := GetClaims(c)
		if claims == nil {
			response.AbortKey(c, http.StatusUnauthorized, 401, ii18n.AuthUnauthorized, nil, "unauthorized")
			return
		}
		if claims.RoleKey == "admin" {
			c.Next()
			return
		}
		if _, ok := allowed[claims.RoleKey]; ok {
			c.Next()
			return
		}
		response.AbortKey(c, http.StatusForbidden, 403, ii18n.AuthForbiddenRole, nil, "forbidden: insufficient role")
	}
}
