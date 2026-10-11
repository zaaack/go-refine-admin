package middleware

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"time"

	"go-refine-admin/backend/core/auth"
	"go-refine-admin/backend/core/netutil"
	"go-refine-admin/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OperationLog records write operations (POST/PUT/DELETE) into sys_oper_log.
// Read operations (GET/HEAD/OPTIONS) are skipped to keep the audit table focused.
func OperationLog(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		// Only audit write operations.
		auditable := method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH"
		var bodyBuf bytes.Buffer
		if auditable && c.Request.Body != nil && c.Request.ContentLength < 1<<16 {
			tee := io.TeeReader(c.Request.Body, &bodyBuf)
			b, _ := io.ReadAll(tee)
			c.Request.Body = io.NopCloser(bytes.NewReader(b))
		}

		c.Next()

		if !auditable {
			return
		}
		path := c.Request.URL.Path
		// Skip health/swagger/captcha noise even on write paths.
		if path == "/health" || strings.HasPrefix(path, "/swagger") ||
			strings.HasPrefix(path, "/captcha") || strings.HasPrefix(path, "/api/v1/captcha") {
			return
		}
		claims := auth.GetClaims(c)
		ip := netutil.ClientIP(c.Request)
		ua := c.Request.UserAgent()
		if len(ua) > 500 {
			ua = ua[:500]
		}
		body := maskSecrets(bodyBuf.String())
		if len(body) > 2000 {
			body = body[:2000]
		}
		statusCode := c.Writer.Status()
		rec := models.SysOperLog{
			Title:      method + " " + path,
			Method:     method,
			Path:       path,
			IP:         ip,
			Location:   netutil.Location(ip),
			UserAgent:  ua,
			Duration:   time.Since(start).Milliseconds(),
			Body:       body,
			StatusCode: statusCode,
		}
		if claims != nil {
			rec.UserID = claims.UserID
			rec.Username = claims.Username
		}
		if statusCode >= 400 || len(c.Errors) > 0 {
			rec.Status = 1
			if len(c.Errors) > 0 {
				rec.ErrMsg = c.Errors.Last().Error()
			} else {
				rec.ErrMsg = httpStatusText(statusCode)
			}
		}
		// fire-and-forget; ignore errors to avoid breaking the request
		go func() { _ = gdb.Create(&rec).Error }()
	}
}

// maskSecrets replaces password-like JSON values to avoid storing credentials.
var secretPattern = regexp.MustCompile(`(?i)("(oldPassword|newPassword|password)"\s*:\s*")[^"]*(")`)

func maskSecrets(s string) string {
	return secretPattern.ReplaceAllString(s, `$1******$3`)
}

func httpStatusText(code int) string {
	switch code {
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return ""
	}
}
