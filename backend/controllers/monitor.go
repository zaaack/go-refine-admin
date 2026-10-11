package controllers

import (
	"net/http"
	"net/url"
	"time"

	"go-refine-admin/backend/core/auth"
	"go-refine-admin/backend/core/query"
	"go-refine-admin/backend/core/response"
	"go-refine-admin/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// listByQuery is a tiny generic list helper with type-safe filters (Actions-style).
func listByQuery[T any](gdb *gorm.DB, c *gin.Context, defaultSort string) {
	q := c.Request.URL.Query()
	filters, err := query.ParseFilters[T](q)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err, "")
		return
	}
	sorts := query.ParseSort[T](q)
	page, pageSize := query.ParsePage(q)
	db := gdb.Model(new(T))
	for _, sc := range filters {
		db = sc(db)
	}
	if len(sorts) == 0 && defaultSort != "" {
		db = db.Order(defaultSort)
	} else {
		for _, sc := range sorts {
			db = sc(db)
		}
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err, "")
		return
	}
	var list []T
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err, "")
		return
	}
	response.PageOK(c, list, total, page, pageSize, "")
}

// ---- Operation audit log ----

// normalizeOperLogQuery translates go-admin style params to this service's
// field names so both conventions work:
// operUrl->path__like, operIp->ip, operName->username__like,
// title->title__like, requestMethod->method. beginTime/endTime are consumed
// separately as a created_at range (go-admin uses beginTime/endTime).
func normalizeOperLogQuery(q url.Values) (begin, end string) {
	if v := q.Get("operUrl"); v != "" {
		q.Set("path__like", v)
		q.Del("operUrl")
	}
	if v := q.Get("operIp"); v != "" {
		q.Set("ip", v)
		q.Del("operIp")
	}
	if v := q.Get("operName"); v != "" {
		q.Set("username__like", v)
		q.Del("operName")
	}
	if v := q.Get("title"); v != "" && q.Get("title__like") == "" {
		q.Set("title__like", v)
	}
	if v := q.Get("requestMethod"); v != "" {
		q.Set("method", v)
		q.Del("requestMethod")
	}
	begin, end = q.Get("beginTime"), q.Get("endTime")
	q.Del("beginTime")
	q.Del("endTime")
	return begin, end
}

func applyTimeRange(db *gorm.DB, begin, end string) *gorm.DB {
	for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02", time.RFC3339} {
		if begin != "" {
			if t, err := time.ParseInLocation(layout, begin, time.Local); err == nil {
				db = db.Where("created_at >= ?", t)
				break
			}
		}
	}
	for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02", time.RFC3339} {
		if end != "" {
			if t, err := time.ParseInLocation(layout, end, time.Local); err == nil {
				db = db.Where("created_at <= ?", t)
				break
			}
		}
	}
	return db
}

// ListOperLog queries sys_oper_log (write-operation audit).
//
//	@Summary	Operation log list
//	@Tags		Monitor
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/sys/oper-log [get]
//	@Security	Bearer
func ListOperLog(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Request.URL.Query()
		begin, end := normalizeOperLogQuery(q)
		filters, err := query.ParseFilters[models.SysOperLog](q)
		if err != nil {
			response.Error(c, http.StatusBadRequest, 400, err, "")
			return
		}
		sorts := query.ParseSort[models.SysOperLog](q)
		page, pageSize := query.ParsePage(q)
		db := gdb.Model(&models.SysOperLog{})
		for _, sc := range filters {
			db = sc(db)
		}
		db = applyTimeRange(db, begin, end)
		if len(sorts) == 0 {
			db = db.Order("id DESC")
		} else {
			for _, sc := range sorts {
				db = sc(db)
			}
		}
		var total int64
		if err := db.Count(&total).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		var list []models.SysOperLog
		if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		response.PageOK(c, list, total, page, pageSize, "")
	}
}

// GetOperLog returns one audit entry by id.
//
//	@Summary	Operation log detail
//	@Tags		Monitor
//	@Produce	json
//	@Param		id	path		string	true	"id"
//	@Success	200	{object}	response.Envelope
//	@Router		/sys/oper-log/{id} [get]
//	@Security	Bearer
func GetOperLog(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m models.SysOperLog
		if err := gdb.First(&m, c.Param("id")).Error; err != nil {
			response.Error(c, http.StatusNotFound, 404, err, "记录不存在")
			return
		}
		response.OK(c, m, "")
	}
}

// DeleteOperLog removes audit entries by comma-separated ids.
//
//	@Summary	Delete operation log
//	@Tags		Monitor
//	@Produce	json
//	@Param		id	path		string	true	"ids"
//	@Success	200	{object}	response.Envelope
//	@Router		/sys/oper-log/{id} [delete]
//	@Security	Bearer
func DeleteOperLog(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ids := splitIDs(c.Param("id"))
		if len(ids) == 0 {
			response.Error(c, http.StatusBadRequest, 400, errNoID(), "")
			return
		}
		if err := gdb.Where("id IN ?", ids).Delete(&models.SysOperLog{}).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		response.OK(c, nil, "")
	}
}

// CleanOperLog truncates the audit table.
//
//	@Summary	Clean operation log
//	@Tags		Monitor
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/sys/oper-log/clean [delete]
//	@Security	Bearer
func CleanOperLog(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Delete without WHERE requires AllowGlobalUpdate.
		if err := gdb.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.SysOperLog{}).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		response.OK(c, nil, "")
	}
}

// ---- Login log ----

// normalizeLoginLogQuery maps go-admin params: ipaddr->ip,
// loginLocation->location__like. beginTime/endTime handled as created_at range.
func normalizeLoginLogQuery(q url.Values) (begin, end string) {
	if v := q.Get("ipaddr"); v != "" {
		q.Set("ip", v)
		q.Del("ipaddr")
	}
	if v := q.Get("loginLocation"); v != "" {
		q.Set("location__like", v)
		q.Del("loginLocation")
	}
	if v := q.Get("username"); v != "" && q.Get("username__like") == "" {
		q.Set("username__like", v)
	}
	begin, end = q.Get("beginTime"), q.Get("endTime")
	q.Del("beginTime")
	q.Del("endTime")
	return begin, end
}

// ListLoginLog queries sys_login_log.
//
//	@Summary	Login log list
//	@Tags		Monitor
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/sys/login-log [get]
//	@Security	Bearer
func ListLoginLog(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Request.URL.Query()
		begin, end := normalizeLoginLogQuery(q)
		filters, err := query.ParseFilters[models.SysLoginLog](q)
		if err != nil {
			response.Error(c, http.StatusBadRequest, 400, err, "")
			return
		}
		sorts := query.ParseSort[models.SysLoginLog](q)
		page, pageSize := query.ParsePage(q)
		db := gdb.Model(&models.SysLoginLog{})
		for _, sc := range filters {
			db = sc(db)
		}
		db = applyTimeRange(db, begin, end)
		if len(sorts) == 0 {
			db = db.Order("id DESC")
		} else {
			for _, sc := range sorts {
				db = sc(db)
			}
		}
		var total int64
		if err := db.Count(&total).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		var list []models.SysLoginLog
		if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		response.PageOK(c, list, total, page, pageSize, "")
	}
}

// GetLoginLog returns one login-log entry by id.
//
//	@Summary	Login log detail
//	@Tags		Monitor
//	@Produce	json
//	@Param		id	path		string	true	"id"
//	@Success	200	{object}	response.Envelope
//	@Router		/sys/login-log/{id} [get]
//	@Security	Bearer
func GetLoginLog(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m models.SysLoginLog
		if err := gdb.First(&m, c.Param("id")).Error; err != nil {
			response.Error(c, http.StatusNotFound, 404, err, "记录不存在")
			return
		}
		response.OK(c, m, "")
	}
}

// DeleteLoginLog removes login-log entries.
//
//	@Summary	Delete login log
//	@Tags		Monitor
//	@Produce	json
//	@Param		id	path		string	true	"ids"
//	@Success	200	{object}	response.Envelope
//	@Router		/sys/login-log/{id} [delete]
//	@Security	Bearer
func DeleteLoginLog(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ids := splitIDs(c.Param("id"))
		if len(ids) == 0 {
			response.Error(c, http.StatusBadRequest, 400, errNoID(), "")
			return
		}
		if err := gdb.Where("id IN ?", ids).Delete(&models.SysLoginLog{}).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		response.OK(c, nil, "")
	}
}

// CleanLoginLog truncates the login-log table.
//
//	@Summary	Clean login log
//	@Tags		Monitor
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/sys/login-log/clean [delete]
//	@Security	Bearer
func CleanLoginLog(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := gdb.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.SysLoginLog{}).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		response.OK(c, nil, "")
	}
}

// ---- Online sessions (token table) ----

// ListOnline lists active tokens (online users).
//
//	@Summary	Online user list
//	@Tags		Monitor
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/auth/online [get]
//	@Security	Bearer
func ListOnline(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Request.URL.Query()
		// Default to online-only; clients can pass status__eq=1 for offline history.
		if q.Get("status") == "" && q.Get("status__eq") == "" {
			q.Set("status", "0")
		}
		filters, err := query.ParseFilters[models.SysUserToken](q)
		if err != nil {
			response.Error(c, http.StatusBadRequest, 400, err, "")
			return
		}
		sorts := query.ParseSort[models.SysUserToken](q)
		page, pageSize := query.ParsePage(q)
		db := gdb.Model(&models.SysUserToken{})
		for _, sc := range filters {
			db = sc(db)
		}
		if len(sorts) == 0 {
			db = db.Order("id DESC")
		} else {
			for _, sc := range sorts {
				db = sc(db)
			}
		}
		var total int64
		if err := db.Count(&total).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		var list []models.SysUserToken
		if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		// Never expose full context; UA is kept for device audit but token JTI stays.
		response.PageOK(c, list, total, page, pageSize, "")
	}
}

// ForceOffline kicks a session offline by token id (admin only).
// It marks sys_user_token offline so the JWT middleware rejects it afterwards.
//
//	@Summary	Force user offline
//	@Tags		Monitor
//	@Produce	json
//	@Param		id	path		string	true	"token id"
//	@Success	200	{object}	response.Envelope
//	@Router		/auth/online/{id} [delete]
//	@Security	Bearer
func ForceOffline(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetClaims(c)
		if claims == nil || claims.RoleKey != "admin" {
			response.Error(c, http.StatusForbidden, 403, errForbidden(), "仅管理员可强制下线")
			return
		}
		ids := splitIDs(c.Param("id"))
		if len(ids) == 0 {
			response.Error(c, http.StatusBadRequest, 400, errNoID(), "")
			return
		}
		if err := gdb.Model(&models.SysUserToken{}).Where("id IN ?", ids).Update("status", 1).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		response.OK(c, nil, "")
	}
}
