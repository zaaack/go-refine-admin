package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"enterprise-admin/backend/config"
	"enterprise-admin/backend/core/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims is the JWT payload.
type Claims struct {
	UserID    uint     `json:"userId"`
	Username  string   `json:"username"`
	RoleKey   string   `json:"roleKey"`
	RoleID    uint     `json:"roleId"`
	DataScope int      `json:"dataScope"` // 1 all, 2 custom, 3 self
	Perms     []string `json:"perms"`     // F-node permission strings, "*" = all
	jwt.RegisteredClaims
}

type ctxKey string

const userCtxKey ctxKey = "auth_user"

// GenerateToken issues a signed JWT for the given claims.
func GenerateToken(c Claims) (string, error) {
	cfg := config.Get().JWT
	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		timeout = 24 * time.Hour
	}
	c.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(timeout)),
		Issuer:    cfg.Issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(cfg.Secret))
}

// Middleware validates the Bearer token and stores the claims in the context.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Abort(c, http.StatusUnauthorized, 401, errors.New("missing token"), "unauthorized")
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		if tokenStr == header {
			response.Abort(c, http.StatusUnauthorized, 401, errors.New("invalid token format"), "unauthorized")
			return
		}
		claims, err := parseToken(tokenStr)
		if err != nil {
			response.Abort(c, http.StatusUnauthorized, 401, err, "unauthorized")
			return
		}
		c.Set(string(userCtxKey), claims)
		c.Next()
	}
}

func parseToken(tokenStr string) (*Claims, error) {
	cfg := config.Get().JWT
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// GetClaims returns the authenticated user's claims (nil if absent).
func GetClaims(c *gin.Context) *Claims {
	if v, ok := c.Get(string(userCtxKey)); ok {
		if claims, ok := v.(*Claims); ok {
			return claims
		}
	}
	return nil
}

// GetUserID is a convenience accessor for the current user id.
func GetUserID(c *gin.Context) uint {
	if cl := GetClaims(c); cl != nil {
		return cl.UserID
	}
	return 0
}

// GetDataScope returns the current role's data scope: 1 all, 2 custom, 3 self.
func GetDataScope(c *gin.Context) int {
	if cl := GetClaims(c); cl != nil {
		return cl.DataScope
	}
	return 1
}
