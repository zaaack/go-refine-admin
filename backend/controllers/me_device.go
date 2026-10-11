package controllers

import (
	"net/http"
	"time"

	"go-refine-admin/backend/core/auth"
	ii18n "go-refine-admin/backend/core/i18n"
	"go-refine-admin/backend/core/response"
	"go-refine-admin/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// deviceRow is a sanitized view of SysUserToken for the self-service
// "my devices" list: no JTI is exposed, and the current session is flagged.
type deviceRow struct {
	ID        uint      `json:"id"`
	IP        string    `json:"ip"`
	Location  string    `json:"location"`
	Browser   string    `json:"browser"`
	OS        string    `json:"os"`
	Device    string    `json:"device"`
	LoginTime time.Time `json:"loginTime"`
	ExpireAt  time.Time `json:"expireAt"`
	Current   bool      `json:"current"`
}

// MyDevices lists the current user's active sessions (online devices).
//
//	@Summary	My online devices
//	@Tags		Auth
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/me/devices [get]
//	@Security	Bearer
func MyDevices(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetClaims(c)
		if claims == nil {
			response.ErrorKey(c, http.StatusUnauthorized, 401, ii18n.AuthUnauthorized, nil, "unauthorized")
			return
		}
		curJTI := auth.GetJTI(c)
		var tokens []models.SysUserToken
		if err := gdb.Where("user_id = ? AND status = 0 AND expire_at > ?", claims.UserID, time.Now()).
			Order("login_time DESC").Find(&tokens).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		list := make([]deviceRow, 0, len(tokens))
		for _, t := range tokens {
			list = append(list, deviceRow{
				ID: t.ID, IP: t.IP, Location: t.Location,
				Browser: t.Browser, OS: t.OS, Device: t.Device,
				LoginTime: t.LoginTime, ExpireAt: t.ExpireAt,
				Current: t.JTI == curJTI,
			})
		}
		response.OK(c, list, "")
	}
}

// OfflineDevice revokes one of the current user's own sessions by token id.
// The current session cannot be kicked here (use logout instead).
//
//	@Summary	Offline one of my devices
//	@Tags		Auth
//	@Produce	json
//	@Param		id	path		string	true	"token id"
//	@Success	200	{object}	response.Envelope
//	@Router		/me/devices/{id} [delete]
//	@Security	Bearer
func OfflineDevice(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetClaims(c)
		if claims == nil {
			response.ErrorKey(c, http.StatusUnauthorized, 401, ii18n.AuthUnauthorized, nil, "unauthorized")
			return
		}
		id := c.Param("id")
		var t models.SysUserToken
		if err := gdb.Where("user_id = ?", claims.UserID).First(&t, id).Error; err != nil {
			response.ErrorKey(c, http.StatusNotFound, 404, ii18n.CommonNotFound, nil, "记录不存在")
			return
		}
		if t.JTI == auth.GetJTI(c) {
			response.Error(c, http.StatusBadRequest, 400, nil, "不能下线当前设备，请使用退出登录")
			return
		}
		if err := gdb.Model(&t).Update("status", 1).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		response.OK(c, nil, "")
	}
}

// OfflineOtherDevices revokes every active session of the current user except
// the caller's own. Returns the number of sessions kicked.
//
//	@Summary	Offline all my other devices
//	@Tags		Auth
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/me/devices/offline-others [post]
//	@Security	Bearer
func OfflineOtherDevices(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetClaims(c)
		if claims == nil {
			response.ErrorKey(c, http.StatusUnauthorized, 401, ii18n.AuthUnauthorized, nil, "unauthorized")
			return
		}
		res := gdb.Model(&models.SysUserToken{}).
			Where("user_id = ? AND status = 0 AND jti <> ?", claims.UserID, auth.GetJTI(c)).
			Update("status", 1)
		if res.Error != nil {
			response.Error(c, http.StatusInternalServerError, 500, res.Error, "")
			return
		}
		response.OK(c, gin.H{"count": res.RowsAffected}, "")
	}
}
