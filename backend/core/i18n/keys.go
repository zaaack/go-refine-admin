// Package i18n is the frontend-backend i18n contract.
//
// Architecture: the backend NEVER translates. It returns stable message keys
// (msgKey) plus interpolation params (msgParams) inside the response envelope,
// and the frontend renders them with react-i18next:
//
//	backend  ->  { code: 401, msgKey: "auth.invalidCredentials", msg: "<fallback>" }
//	frontend ->  t("backend:auth.invalidCredentials")
//
// Rules:
//   - Keys are stable identifiers; never rename a key, only add new ones.
//   - msg (fallback) must always be filled so logs / swagger / old clients work.
//   - Dynamic technical details (DB / driver / validation errors) are NOT keys:
//     leave msgKey empty and put err.Error() in msg; the frontend shows msg as-is.
//   - Frontend dictionaries live in frontend/src/i18n/locales/{zh-CN,en}.json
//     under the "backend" namespace and MUST cover every key below.
package i18n

// Backend message keys returned in Envelope.MsgKey.
const (
	// Auth / login.
	AuthCaptchaInvalid     = "auth.captchaInvalid"
	AuthInvalidCredentials = "auth.invalidCredentials"
	AuthAccountDisabled    = "auth.accountDisabled"
	AuthUnauthorized       = "auth.unauthorized"
	AuthMissingToken       = "auth.missingToken"
	AuthInvalidToken       = "auth.invalidToken"
	AuthForbiddenRole      = "auth.forbiddenRole"
	AuthMissingPermission  = "auth.missingPermission" // params: { perm }
	// Captcha.
	CaptchaGenerateFailed = "captcha.generateFailed"
	// Roles.
	RoleBuiltinImmutable = "role.builtinImmutable"
	RoleAdminKeyReserved = "role.adminKeyReserved"
	// Orders (regular mode).
	OrderNotFound       = "order.notFound"
	OrderOnlyPendingPay = "order.onlyPendingCanPay"
	// Generic Actions-mode errors.
	CommonInvalidID = "common.invalidId"
	CommonNotFound  = "common.notFound"
	CommonNoID      = "common.noId"
)

// Params builds the msgParams map: i18n.Params("perm", "a,b") -> {"perm":"a,b"}.
func Params(kv ...any) map[string]any {
	if len(kv) == 0 {
		return nil
	}
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		if k, ok := kv[i].(string); ok {
			m[k] = kv[i+1]
		}
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

// KeyError is a translatable business error. Hooks (which can only return
// error) use E() to carry a msgKey through the generic response.Error path:
// response.Error detects KeyError and fills msgKey/msgParams automatically,
// so call sites need no signature changes.
type KeyError struct {
	Key      string
	Params   map[string]any
	Fallback string
}

func (e *KeyError) Error() string { return e.Fallback }

// E builds a translatable error: E(RoleBuiltinImmutable, "内置...").
// The fallback keeps logs and non-i18n clients readable.
func E(key, fallback string, kv ...any) error {
	return &KeyError{Key: key, Params: Params(kv...), Fallback: fallback}
}

// KeyOf unwraps err to its msgKey/msgParams/fallback. ok == false means err
// is not translatable: the caller should use err.Error() as plain msg.
func KeyOf(err error) (key string, params map[string]any, fallback string, ok bool) {
	type causer interface{ Unwrap() error }
	for err != nil {
		if ke, ok := err.(*KeyError); ok {
			return ke.Key, ke.Params, ke.Fallback, true
		}
		u, ok := err.(causer)
		if !ok {
			break
		}
		err = u.Unwrap()
	}
	return "", nil, "", false
}
