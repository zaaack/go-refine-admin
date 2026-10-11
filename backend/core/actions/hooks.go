package actions

import (
	"go-refine-admin/backend/core/auth"
	"go-refine-admin/backend/core/query"
	"go-refine-admin/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DataScopeHook enforces the role's data_scope on list queries, mirroring
// go-admin's data permission: scope 3 ("self") restricts rows to those created
// by the current user. Attach it as a BeforeList hook:
//
//	actions.NewResource[models.Product](db).WithHooks(actions.Hooks[models.Product]{
//	    BeforeList: actions.DataScopeHook[models.Product](),
//	})
//
// (scope 1 = all data, no filter; scope 2 = custom, extend this hook to fit.)
func DataScopeHook[T models.Entity]() func(c *gin.Context, scopes *[]query.Scope) error {
	return func(c *gin.Context, scopes *[]query.Scope) error {
		if auth.GetDataScope(c) == 3 {
			uid := auth.GetUserID(c)
			*scopes = append(*scopes, func(db *gorm.DB) *gorm.DB {
				return db.Where("create_by = ?", uid)
			})
		}
		return nil
	}
}
