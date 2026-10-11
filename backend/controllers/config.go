package controllers

import (
	"net/http"

	"go-refine-admin/backend/core/response"
	"go-refine-admin/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AppConfig returns frontend-readable configs as {key: value} (public, 参考 go-admin Get2SysApp).
//
//	@Summary	Frontend config map
//	@Tags		Config
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/app-config [get]
func AppConfig(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var list []models.SysConfig
		if err := gdb.Where("is_frontend = ?", 1).Find(&list).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		mp := make(map[string]string, len(list))
		for _, it := range list {
			mp[it.ConfigKey] = it.ConfigValue
		}
		response.OK(c, mp, "")
	}
}

// SetConfigRequest is a batch key-value update payload.
type SetConfigRequest map[string]string

// GetSetConfig returns all configs as {key: value} (authenticated).
//
//	@Summary	All config map
//	@Tags		Config
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/set-config [get]
//	@Security	Bearer
func GetSetConfig(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var list []models.SysConfig
		if err := gdb.Find(&list).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		mp := make(map[string]string, len(list))
		for _, it := range list {
			mp[it.ConfigKey] = it.ConfigValue
		}
		response.OK(c, mp, "")
	}
}

// UpdateSetConfig batch-updates config values by key.
//
//	@Summary	Update config map
//	@Tags		Config
//	@Accept		json
//	@Produce	json
//	@Param		body	body		SetConfigRequest	true	"key-value map"
//	@Success	200	{object}	response.Envelope
//	@Router		/set-config [put]
//	@Security	Bearer
func UpdateSetConfig(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SetConfigRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, 400, err, "")
			return
		}
		for k, v := range req {
			gdb.Model(&models.SysConfig{}).Where("config_key = ?", k).Update("config_value", v)
		}
		response.OK(c, nil, "更新成功")
	}
}
