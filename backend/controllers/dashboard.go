package controllers

import (
	"enterprise-admin/backend/core/response"
	"enterprise-admin/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Dashboard returns headline counts and order revenue for the home page.
//
//	@Summary	Dashboard stats
//	@Tags		Dashboard
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/dashboard/stats [get]
//	@Security	Bearer
func Dashboard(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userCount, roleCount, productCount, orderCount, dictCount int64
		gdb.Model(&models.SysUser{}).Count(&userCount)
		gdb.Model(&models.SysRole{}).Count(&roleCount)
		gdb.Model(&models.Product{}).Count(&productCount)
		gdb.Model(&models.Order{}).Count(&orderCount)
		gdb.Model(&models.SysDict{}).Count(&dictCount)

		var revenue float64
		gdb.Model(&models.Order{}).
			Where("status IN ?", []string{models.OrderPaid, models.OrderShipped, models.OrderDone}).
			Select("COALESCE(SUM(amount),0)").Scan(&revenue)

		response.OK(c, gin.H{
			"counts": gin.H{
				"user":    userCount,
				"role":    roleCount,
				"product": productCount,
				"order":   orderCount,
				"dict":    dictCount,
			},
			"revenue": revenue,
		}, "")
	}
}
