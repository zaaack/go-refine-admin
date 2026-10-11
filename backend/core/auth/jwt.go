package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"go-refine-admin/backend/config"
	ii18n "go-refine-admin/backend/core/i18n"
	"go-refine-admin/backend/core/response"
	"go-refine-admin/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
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

// NewJTI generates a random JWT ID (session identifier stored in sys_user_token).
func NewJTI() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return strings.ReplaceAll(time.Now().Format("20060102150405.000000000"), ".", "")
	}
	return hex.EncodeToString(b[:])
}

// GenerateToken issues a signed JWT for the given claims.
// It assigns a fresh JTI (RegisteredClaims.ID) and returns token + jti.
func GenerateToken(c Claims) (tokenStr, jti string, err error) {
	cfg := config.Get().JWT
	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		timeout = 24 * time.Hour
	}
	jti = NewJTI()
	c.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(timeout)),
		Issuer:    cfg.Issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        jti,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenStr, err = token.SignedString([]byte(cfg.Secret))
	return tokenStr, jti, err
}

// TokenTimeout returns the configured JWT lifetime (fallback 24h).
func TokenTimeout() time.Duration {
	cfg := config.Get().JWT
	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		timeout = 24 * time.Hour
	}
	return timeout
}

// Middleware validates the Bearer token, checks the server-side token table
// (sys_user_token) for revocation, and stores the claims in the context.
// Tokens issued before the token-table upgrade carry no JTI and are accepted
// as-is for backward compatibility.
func Middleware(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.AbortKey(c, http.StatusUnauthorized, 401, ii18n.AuthMissingToken, nil, "unauthorized")
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		if tokenStr == header {
			response.AbortKey(c, http.StatusUnauthorized, 401, ii18n.AuthInvalidToken, nil, "unauthorized")
			return
		}
		claims, err := parseToken(tokenStr)
		if err != nil {
			response.AbortKey(c, http.StatusUnauthorized, 401, ii18n.AuthInvalidToken, nil, "unauthorized")
			return
		}
		if claims.ID != "" && gdb != nil {
			var sess models.SysUserToken
			if err := gdb.Where("jti = ?", claims.ID).First(&sess).Error; err != nil {
				response.AbortKey(c, http.StatusUnauthorized, 401, ii18n.AuthInvalidToken, nil, "token已下线或已注销，请重新登录")
				return
			}
			if sess.Status != 0 || time.Now().After(sess.ExpireAt) {
				response.AbortKey(c, http.StatusUnauthorized, 401, ii18n.AuthInvalidToken, nil, "token已下线或已过期，请重新登录")
				return
			}
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

// GetJTI returns the current token's JTI (session id, "" if absent).
func GetJTI(c *gin.Context) string {
	if cl := GetClaims(c); cl != nil {
		return cl.ID
	}
	return ""
}

// GetDataScope returns the current role's data scope: 1 all, 2 custom, 3 self.
func GetDataScope(c *gin.Context) int {
	if cl := GetClaims(c); cl != nil {
		return cl.DataScope
	}
	return 1
}
