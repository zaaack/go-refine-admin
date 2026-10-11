package controllers

import (
	"net/http"
	"strconv"
	"time"

	"go-refine-admin/backend/core/auth"
	"go-refine-admin/backend/core/response"
	"go-refine-admin/backend/genquery"
	"go-refine-admin/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// OrderController is a regular-mode (hand-written) resource. It shows how the
// type-safe query builder composes with custom aggregation/business logic that
// Actions mode does not cover.
type OrderController struct{ DB *gorm.DB }

// List returns orders with type-safe filters AND a revenue summary — something
// the zero-code Actions handler does not do.
//
// The filters use gorm/gen's generated, compile-time type-safe field accessors
// (genquery.Order.Amount.Gte(...)), so a typo or wrong value type fails to build.
//
//	@Summary	List orders (regular mode) with summary
//	@Tags		Order
//	@Produce	json
//	@Param		minAmount	query		number	false	"min amount"
//	@Param		status		query		string	false	"status, repeatable"
//	@Param		channel		query		string	false	"channel"
//	@Param		start		query		string	false	"created >= (RFC3339)"
//	@Param		end			query		string	false	"created <= (RFC3339)"
//	@Param		page		query		int		false	"page"
//	@Param		pageSize	query		int		false	"page size"
//	@Success	200			{object}	response.Envelope
//	@Router		/order [get]
//	@Security	Bearer
func (ctl OrderController) List(c *gin.Context) {
	// Type-safe conditions built with gorm/gen generated accessors.
	var conds []gen.Condition
	if v := c.Query("minAmount"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			conds = append(conds, genquery.Order.Amount.Gte(f))
		}
	}
	if statuses := c.QueryArray("status"); len(statuses) > 0 {
		conds = append(conds, genquery.Order.Status.In(statuses...))
	}
	if ch := c.Query("channel"); ch != "" {
		conds = append(conds, genquery.Order.Channel.Eq(ch))
	}
	if s := c.Query("start"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			conds = append(conds, genquery.Order.CreatedAt.Gte(t))
		}
	}
	if e := c.Query("end"); e != "" {
		if t, err := time.Parse(time.RFC3339, e); err == nil {
			conds = append(conds, genquery.Order.CreatedAt.Lte(t))
		}
	}
	// Role data-scope: 3 = only rows created by the current user.
	if auth.GetDataScope(c) == 3 {
		conds = append(conds, genquery.Order.CreateBy.Eq(auth.GetUserID(c)))
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	// genquery.Order uses the default DB bound by genquery.SetDefault(db).
	list, err := genquery.Order.Where(conds...).
		Order(genquery.Order.ID.Desc()).
		Offset((page - 1) * size).
		Limit(size).
		Find()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err, "")
		return
	}
	total, _ := genquery.Order.Where(conds...).Count()

	// Custom aggregation: total amount + paid amount for the filtered set.
	var summary struct {
		TotalAmount float64 `json:"totalAmount"`
		PaidAmount  float64 `json:"paidAmount"`
		Count       int64   `json:"count"`
	}
	ctl.DB.Model(&models.Order{}).
		Select("COALESCE(SUM(amount),0) as total_amount, "+
			"COALESCE(SUM(CASE WHEN status='paid' OR status='shipped' OR status='done' THEN amount ELSE 0 END),0) as paid_amount, "+
			"COUNT(*) as count").
		Scan(&summary)

	response.OK(c, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": size,
		"summary":  summary,
	}, "")
}

// Stats groups orders by status (regular-mode aggregation).
//
//	@Summary	Order stats by status
//	@Tags		Order
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/order/stats [get]
//	@Security	Bearer
func (ctl OrderController) Stats(c *gin.Context) {
	type row struct {
		Status string  `json:"status"`
		Count  int64   `json:"count"`
		Amount float64 `json:"amount"`
	}
	var rows []row
	ctl.DB.Model(&models.Order{}).
		Select("status, COUNT(*) as count, COALESCE(SUM(amount),0) as amount").
		Group("status").Order("count desc").
		Scan(&rows)
	response.OK(c, rows, "")
}

// MarkPaid transitions an order to paid with hand-written validation.
//
//	@Summary	Mark order paid
//	@Tags		Order
//	@Produce	json
//	@Param		id	path		string	true	"order id"
//	@Success	200	{object}	response.Envelope
//	@Router		/order/{id}/paid [put]
//	@Security	Bearer
func (ctl OrderController) MarkPaid(c *gin.Context) {
	id := c.Param("id")
	var o models.Order
	if err := ctl.DB.First(&o, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, 404, err, "order not found")
		return
	}
	if o.Status != models.OrderPending && o.Status != models.OrderCanceled {
		response.Error(c, http.StatusBadRequest, 400, nil, "only pending/canceled orders can be paid")
		return
	}
	if err := ctl.DB.Model(&o).Update("status", models.OrderPaid).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err, "")
		return
	}
	response.OK(c, o, "")
}

// Create persists a new order with business defaults (regular mode).
//
//	@Summary	Create order
//	@Tags		Order
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.Order	true	"order"
//	@Success	200		{object}	response.Envelope
//	@Router		/order [post]
//	@Security	Bearer
func (ctl OrderController) Create(c *gin.Context) {
	var o models.Order
	if err := c.ShouldBindJSON(&o); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err, "")
		return
	}
	if o.OrderNo == "" {
		o.OrderNo = "NO" + time.Now().Format("20060102150405")
	}
	if o.Status == "" {
		o.Status = models.OrderPending
	}
	if err := ctl.DB.Create(&o).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err, "")
		return
	}
	response.OK(c, o, "")
}

// Delete removes an order (regular mode).
//
//	@Summary	Delete order
//	@Tags		Order
//	@Param		id	path		string	true	"order id"
//	@Success	200	{object}	response.Envelope
//	@Router		/order/{id} [delete]
//	@Security	Bearer
func (ctl OrderController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ctl.DB.Delete(&models.Order{}, id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err, "")
		return
	}
	response.OK(c, nil, "")
}
