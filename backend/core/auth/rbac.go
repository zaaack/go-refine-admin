package auth

import (
	"net/http"

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
			response.Abort(c, http.StatusUnauthorized, 401, nil, "unauthorized")
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
		response.Abort(c, http.StatusForbidden, 403, nil, "forbidden: insufficient role")
	}
}
