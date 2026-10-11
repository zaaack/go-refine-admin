package middleware

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"time"

	"go-refine-admin/backend/core/auth"
	"go-refine-admin/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OperationLog records each request into sys_oper_log (sampled: skips static/assets).
func OperationLog(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// capture body (best effort, only for small payloads)
		var bodyBuf bytes.Buffer
		if c.Request.Body != nil && c.Request.ContentLength < 1<<16 {
			tee := io.TeeReader(c.Request.Body, &bodyBuf)
			b, _ := io.ReadAll(tee)
			c.Request.Body = io.NopCloser(bytes.NewReader(b))
			if len(b) > 2000 {
				b = b[:2000]
			}
			_ = b
		}

		c.Next()

		path := c.Request.URL.Path
		if path == "/health" || path == "/swagger/*any" || path == "/swagger/index.html" {
			return
		}
		claims := auth.GetClaims(c)
		rec := models.SysOperLog{
			Title:    path,
			Method:   c.Request.Method,
			Path:     path,
			IP:       clientIP(c.Request),
			Duration: time.Since(start).Milliseconds(),
			Body:     bodyBuf.String(),
		}
		if claims != nil {
			rec.UserID = claims.UserID
			rec.Username = claims.Username
		}
		if len(c.Errors) > 0 {
			rec.Status = 1
			rec.ErrMsg = c.Errors.Last().Error()
		}
		// fire-and-forget; ignore errors to avoid breaking the request
		go func() { _ = gdb.Create(&rec).Error }()
	}
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return fwd
	}
	if real := r.Header.Get("X-Real-IP"); real != "" {
		return real
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
