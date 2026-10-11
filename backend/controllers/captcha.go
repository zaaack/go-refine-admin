package controllers

import (
	"time"

	"go-refine-admin/backend/core/response"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

// Mirrors go-admin's captcha: a digit captcha stored in memory, served as a
// base64 PNG together with its id, verified once at login.

var captchaStore = base64Captcha.NewMemoryStore(10240, 5*time.Minute)

func captchaDriver() base64Captcha.Driver {
	return base64Captcha.NewDriverDigit(80, 240, 4, 0.6, 80)
}

// GenerateCaptcha issues a new captcha image.
//
//	@Summary	获取验证码
//	@Tags		Auth
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/captcha [get]
func GenerateCaptcha(c *gin.Context) {
	cap := base64Captcha.NewCaptcha(captchaDriver(), captchaStore)
	// The answer is deliberately discarded so it never reaches logs.
	id, b64, _, err := cap.Generate()
	if err != nil {
		response.Error(c, 500, 500, err, "验证码获取失败")
		return
	}
	// go-admin-compatible shape: image in data, captcha id in id.
	c.JSON(200, gin.H{"code": 200, "data": b64, "id": id, "msg": "success"})
}

// VerifyCaptcha checks and consumes a captcha. Skipped in dev mode, like go-admin.
func VerifyCaptcha(id, code string) bool {
	if id == "" || code == "" {
		return false
	}
	return captchaStore.Verify(id, code, true)
}
