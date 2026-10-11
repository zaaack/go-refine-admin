package middleware

import (
	"enterprise-admin/backend/config"

	"github.com/gin-gonic/gin"
)

// CORS applies cross-origin headers from configuration.
func CORS() gin.HandlerFunc {
	cfg := config.Get().CORS
	allowOrigin := "*"
	if len(cfg.AllowOrigins) == 1 && cfg.AllowOrigins[0] == "*" {
		allowOrigin = "*"
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if allowOrigin == "*" {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				for _, o := range cfg.AllowOrigins {
					if o == origin {
						c.Header("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}
		}
		c.Header("Access-Control-Allow-Methods", join(cfg.AllowMethods))
		c.Header("Access-Control-Allow-Headers", join(cfg.AllowHeaders))
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func join(s []string) string {
	out := ""
	for i, v := range s {
		if i > 0 {
			out += ", "
		}
		out += v
	}
	return out
}
