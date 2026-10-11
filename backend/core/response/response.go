package response

import (
	"net/http"

	ii18n "go-refine-admin/backend/core/i18n"

	"github.com/gin-gonic/gin"
)

// Envelope is the unified API response shape:
// { code, data, msg, msgKey, msgParams }.
// code 200 == success; any other code is an error.
//
// i18n contract (frontend renders, backend only returns keys):
//   - msgKey is a stable key from core/i18n (e.g. "auth.invalidCredentials").
//   - msgParams carries interpolation params for the key.
//   - msg is always filled as a human-readable fallback (default language)
//     so old clients / logs / swagger stay usable without i18n.
//   - msgKey empty means "not a translatable business error": msg carries the
//     raw detail (validation / driver errors) and the frontend shows msg as-is.
type Envelope struct {
	Code      int            `json:"code"`
	Data      interface{}    `json:"data"`
	Msg       string         `json:"msg"`
	MsgKey    string         `json:"msgKey,omitempty"`
	MsgParams map[string]any `json:"msgParams,omitempty"`
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
//
// i18n: when err carries a *i18n.KeyError (see i18n.E, used by hooks), its
// msgKey/msgParams are filled automatically and its fallback becomes msg
// unless an explicit msg is given.
func Error(c *gin.Context, httpStatus, code int, err error, msg string) {
	if msg == "" {
		msg = http.StatusText(httpStatus)
	}
	var key string
	var params map[string]any
	if err != nil {
		if k, p, fb, ok := ii18n.KeyOf(err); ok {
			key, params = k, p
			if msg == http.StatusText(httpStatus) {
				msg = fb
			}
		} else if msg == http.StatusText(httpStatus) {
			msg = err.Error()
		}
	}
	c.JSON(httpStatus, Envelope{Code: code, Data: nil, Msg: msg, MsgKey: key, MsgParams: params})
}

// Abort writes an error and stops the chain.
func Abort(c *gin.Context, httpStatus, code int, err error, msg string) {
	Error(c, httpStatus, code, err, msg)
	c.Abort()
}

// ErrorKey writes a translatable business error: msgKey/msgParams are rendered
// by the frontend via react-i18next (namespace "backend"); fallbackMsg keeps
// the envelope readable without i18n. When err != nil and fallbackMsg == "",
// err.Error() is used as the fallback.
func ErrorKey(c *gin.Context, httpStatus, code int, msgKey string, params map[string]any, fallbackMsg string) {
	if fallbackMsg == "" {
		fallbackMsg = http.StatusText(httpStatus)
	}
	c.JSON(httpStatus, Envelope{Code: code, Data: nil, Msg: fallbackMsg, MsgKey: msgKey, MsgParams: params})
}

// AbortKey writes a translatable error and stops the chain.
func AbortKey(c *gin.Context, httpStatus, code int, msgKey string, params map[string]any, fallbackMsg string) {
	ErrorKey(c, httpStatus, code, msgKey, params, fallbackMsg)
	c.Abort()
}
