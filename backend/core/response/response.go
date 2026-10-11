package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Envelope is the unified API response shape: { code, data, msg }.
// code 200 == success; any other code is an error.
type Envelope struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// PageData wraps a paginated list.
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

const (
	// CodeOK is the success code (mirrors go-admin convention).
	CodeOK = 200
	// CodeErr is the generic error code.
	CodeErr = 500
)

// OK writes a success response with arbitrary data.
func OK(c *gin.Context, data interface{}, msg string) {
	if msg == "" {
		msg = "ok"
	}
	c.JSON(http.StatusOK, Envelope{Code: CodeOK, Data: data, Msg: msg})
}

// PageOK writes a paginated list response.
func PageOK(c *gin.Context, list interface{}, total int64, page, pageSize int, msg string) {
	if msg == "" {
		msg = "ok"
	}
	c.JSON(http.StatusOK, Envelope{Code: CodeOK, Data: PageData{List: list, Total: total, Page: page, PageSize: pageSize}, Msg: msg})
}

// Error writes an error response. httpStatus is the real HTTP status; code is the
// business code carried in the envelope.
func Error(c *gin.Context, httpStatus, code int, err error, msg string) {
	if msg == "" {
		msg = http.StatusText(httpStatus)
	}
	if err != nil && msg == http.StatusText(httpStatus) {
		msg = err.Error()
	}
	c.JSON(httpStatus, Envelope{Code: code, Data: nil, Msg: msg})
}

// Abort writes an error and stops the chain.
func Abort(c *gin.Context, httpStatus, code int, err error, msg string) {
	Error(c, httpStatus, code, err, msg)
	c.Abort()
}
