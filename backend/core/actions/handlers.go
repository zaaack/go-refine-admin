package actions

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"go-refine-admin/backend/core/i18n"
	"go-refine-admin/backend/core/query"
	"go-refine-admin/backend/core/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// list handles GET /:name  (paginated, filtered, sorted).
//
//	@Summary		List resources (Actions mode)
//	@Tags		Actions
//	@Produce	json
//	@Param		resource	path		string	true	"resource name, e.g. demo/product"
//	@Param		page		query		int		false	"page"
//	@Param		pageSize	query		int		false	"page size"
//	@Param		sort		query		string	false	"sort, e.g. -createdAt,name"
//	@Success	200	{object}	response.Envelope
//	@Router		/{resource} [get]
// NOTE: dynamic filters (field__op=value, e.g. name__like=x) are accepted as
// arbitrary query params at runtime; they are intentionally not enumerated
// here. Pass them via axios options.params when using the generated client.
func (r *Resource[T]) list(c *gin.Context) {
	q := c.Request.URL.Query()

	filters, err := query.ParseFilters[T](q)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err, "")
		return
	}
	sorts := query.ParseSort[T](q)
	page, pageSize := query.ParsePage(q)

	var scopes []query.Scope
	scopes = append(scopes, filters...)
	scopes = append(scopes, sorts...)
	if r.Hooks.BeforeList != nil {
		if err := r.Hooks.BeforeList(c, &scopes); err != nil {
			response.Error(c, http.StatusBadRequest, 400, err, "")
			return
		}
	}

	db := r.DB.Model(new(T))
	for _, sc := range scopes {
		db = sc(db)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err, "")
		return
	}
	var list []T
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err, "")
		return
	}
	response.PageOK(c, list, total, page, pageSize, "")
}

// get handles GET /:name/:id.
//
//	@Summary	Get one resource (Actions mode)
//	@Tags		Actions
//	@Produce	json
//	@Param		resource	path		string	true	"resource name, e.g. demo/product"
//	@Param		id			path		string	true	"primary key"
//	@Success	200	{object}	response.Envelope
//	@Router		/{resource}/{id} [get]
func (r *Resource[T]) get(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.ErrorKey(c, http.StatusBadRequest, 400, i18n.CommonInvalidID, nil, "invalid id")
		return
	}
	var m T
	pk := pkDBName(r.DB, m)
	if err := r.DB.Model(new(T)).Where(pk+" = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorKey(c, http.StatusNotFound, 404, i18n.CommonNotFound, nil, "not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, 500, err, "")
		return
	}
	response.OK(c, m, "")
}

// create handles POST /:name.
//
//	@Summary	Create resource (Actions mode)
//	@Tags		Actions
//	@Accept		json
//	@Produce	json
//	@Param		resource	path		string	true	"resource name, e.g. demo/product"
//	@Param		body		body		object	true	"record"
//	@Success	200	{object}	response.Envelope
//	@Router		/{resource} [post]
func (r *Resource[T]) create(c *gin.Context) {
	var m T
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err, "")
		return
	}
	if r.Hooks.BeforeCreate != nil {
		if err := r.Hooks.BeforeCreate(c, &m); err != nil {
			response.Error(c, http.StatusBadRequest, 400, err, "")
			return
		}
	}
	if err := r.DB.Create(&m).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err, "")
		return
	}
	if r.Hooks.AfterCreate != nil {
		_ = r.Hooks.AfterCreate(c, &m)
	}
	response.OK(c, m, "")
}

// update handles PUT /:name/:id.
//
//	@Summary	Update resource (Actions mode)
//	@Tags		Actions
//	@Accept		json
//	@Produce	json
//	@Param		resource	path		string	true	"resource name, e.g. demo/product"
//	@Param		id			path		string	true	"primary key"
//	@Param		body		body		object	true	"record"
//	@Success	200	{object}	response.Envelope
//	@Router		/{resource}/{id} [put]
func (r *Resource[T]) update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.ErrorKey(c, http.StatusBadRequest, 400, i18n.CommonInvalidID, nil, "invalid id")
		return
	}
	var m T
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err, "")
		return
	}
	if r.Hooks.BeforeUpdate != nil {
		if err := r.Hooks.BeforeUpdate(c, &m); err != nil {
			response.Error(c, http.StatusBadRequest, 400, err, "")
			return
		}
	}
	pk := pkDBName(r.DB, m)
	if err := r.DB.Model(new(T)).Where(pk+" = ?", id).Updates(&m).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err, "")
		return
	}
	var out T
	r.DB.Model(new(T)).Where(pk+" = ?", id).First(&out)
	response.OK(c, out, "")
}

// delete handles DELETE /:name/:id  (comma-separated ids allowed).
//
//	@Summary	Delete resource (Actions mode)
//	@Tags		Actions
//	@Param		resource	path		string	true	"resource name, e.g. demo/product"
//	@Param		id			path		string	true	"primary key(s), comma-separated"
//	@Success	200	{object}	response.Envelope
//	@Router		/{resource}/{id} [delete]
func (r *Resource[T]) delete(c *gin.Context) {
	ids, err := parseIDs(c.Param("id"))
	if err != nil {
		response.ErrorKey(c, http.StatusBadRequest, 400, i18n.CommonInvalidID, nil, "invalid id")
		return
	}
	if r.Hooks.BeforeDelete != nil {
		for _, id := range ids {
			if err := r.Hooks.BeforeDelete(c, id); err != nil {
				response.Error(c, http.StatusBadRequest, 400, err, "")
				return
			}
		}
	}
	pk := pkDBName(r.DB, new(T))
	if err := r.DB.Model(new(T)).Where(pk+" IN ?", ids).Delete(new(T)).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err, "")
		return
	}
	response.OK(c, nil, "")
}

func parseID(s string) (interface{}, error) {
	if u, err := strconv.ParseUint(s, 10, 64); err == nil {
		return u, nil
	}
	return s, nil
}

func parseIDs(s string) ([]interface{}, error) {
	parts := strings.Split(s, ",")
	out := make([]interface{}, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if u, err := strconv.ParseUint(p, 10, 64); err == nil {
			out = append(out, u)
		} else {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil, i18n.E(i18n.CommonNoID, "no id provided")
	}
	return out, nil
}
