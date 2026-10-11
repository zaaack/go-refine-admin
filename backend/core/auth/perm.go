package auth

import (
	"net/http"
	"strings"

	"go-refine-admin/backend/core/response"

	"github.com/gin-gonic/gin"
)

// Permissions come from the role's granted F-type menu nodes (Permission field,
// e.g. "sys:user:create") and are carried in the JWT at login time. Changing a
// role's permissions therefore takes effect on the next login.
//
// Hand-written (regular mode) checks:
//
//	auth.RequirePerm("demo:order:paid")        // middleware on a route
//	if !auth.HasPerm(c, "demo:order:export") { // inline check
//	    ...
//	}
//
// The admin role (and the "*" wildcard) always passes.

// HasPerm reports whether the current user holds the permission.
func HasPerm(c *gin.Context, perm string) bool {
	claims := GetClaims(c)
	if claims == nil {
		return false
	}
	if claims.RoleKey == "admin" {
		return true
	}
	for _, p := range claims.Perms {
		if p == "*" || p == perm {
			return true
		}
	}
	return false
}

// RequirePerm aborts with 403 unless the user holds one of the permissions.
func RequirePerm(perms ...string) gin.HandlerFunc {
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
		for _, held := range claims.Perms {
			if held == "*" {
				c.Next()
				return
			}
			for _, need := range perms {
				if held == need {
					c.Next()
					return
				}
			}
		}
		response.Abort(c, http.StatusForbidden, 403, nil,
			"forbidden: missing permission "+strings.Join(perms, ","))
	}
}
