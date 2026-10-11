package actions

import (
	"enterprise-admin/backend/core/auth"
	"enterprise-admin/backend/core/query"
	"enterprise-admin/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Hooks let you inject business logic around the zero-code CRUD lifecycle
// without leaving Actions mode. Any hook may return an error to abort.
type Hooks[T models.Entity] struct {
	// BeforeList receives the scope slice; append extra gorm scopes (e.g. a
	// tenant filter) before execution.
	BeforeList func(c *gin.Context, scopes *[]query.Scope) error
	// BeforeCreate runs after binding the payload, before insert.
	BeforeCreate func(c *gin.Context, m *T) error
	// BeforeUpdate runs after binding the payload, before update.
	BeforeUpdate func(c *gin.Context, m *T) error
	// BeforeDelete runs before delete; id is the parsed primary key.
	BeforeDelete func(c *gin.Context, id interface{}) error
	// AfterCreate runs after a successful insert.
	AfterCreate func(c *gin.Context, m *T) error
}

// Resource is a zero-code CRUD definition for model T.
//
//	actions.NewResource[models.Product](db).
//	    WithName("product").
//	    WithHooks(Hooks[models.Product]{ BeforeCreate: setCreator }).
//	    Register(v1)
//
// This yields GET/POST /:name, GET/PUT/DELETE /:name/:id with type-safe
// filtering, pagination and unified responses — no handler code required.
type Resource[T models.Entity] struct {
	DB   *gorm.DB
	Name string
	// PermPrefix, when set (e.g. "sys:user"), guards write routes with
	// RequirePerm(prefix+":create"/":edit"/":delete"). Read routes stay open to
	// any authenticated user — menu visibility controls read access, so a role
	// granted only the menu is read-only. Empty prefix = no permission guard.
	PermPrefix string
	Hooks      Hooks[T]
}

// NewResource constructs a Resource bound to a database.
func NewResource[T models.Entity](db *gorm.DB) *Resource[T] {
	return &Resource[T]{DB: db}
}

// WithName overrides the route segment (defaults to the model's table name).
func (r *Resource[T]) WithName(name string) *Resource[T] {
	r.Name = name
	return r
}

// WithPermPrefix enables write-route permission guards with the given prefix.
func (r *Resource[T]) WithPermPrefix(prefix string) *Resource[T] {
	r.PermPrefix = prefix
	return r
}

// WithHooks attaches lifecycle hooks.
func (r *Resource[T]) WithHooks(h Hooks[T]) *Resource[T] {
	r.Hooks = h
	return r
}

// Register wires the five CRUD routes under group. The group is expected to
// already carry auth (and optional RBAC) middleware. Write routes get an
// additional permission guard when WithPermPrefix was set.
func (r *Resource[T]) Register(group *gin.RouterGroup) {
	g := group.Group("/" + r.Name)
	g.GET("", r.list)
	g.GET("/:id", r.get)
	g.POST("", r.writeGuard(":create"), r.create)
	g.PUT("/:id", r.writeGuard(":edit"), r.update)
	g.DELETE("/:id", r.writeGuard(":delete"), r.delete)
}

// writeGuard returns the permission middleware for a write action, or a no-op
// pass-through when no prefix was configured.
func (r *Resource[T]) writeGuard(action string) gin.HandlerFunc {
	if r.PermPrefix == "" {
		return func(c *gin.Context) { c.Next() }
	}
	return auth.RequirePerm(r.PermPrefix + action)
}

// pkDBName returns the primary-key column for T.
func pkDBName(db *gorm.DB, m any) string {
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(m)
	if stmt.Schema != nil && len(stmt.Schema.PrimaryFields) > 0 {
		return stmt.Schema.PrimaryFields[0].DBName
	}
	return "id"
}
