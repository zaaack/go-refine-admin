package query

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Op enumerates supported comparison operators (used by the Actions-mode HTTP
// filter parser). The type-safe query for hand-written controllers is provided
// by gorm/gen (see the generated genquery package), not by this file.
type Op string

const (
	OpEq      Op = "eq"
	OpNeq     Op = "neq"
	OpGt      Op = "gt"
	OpGte     Op = "gte"
	OpLt      Op = "lt"
	OpLte     Op = "lte"
	OpLike    Op = "like"
	OpILike   Op = "ilike"
	OpIn      Op = "in"
	OpNotIn   Op = "notin"
	OpBetween Op = "between"
	OpIsNull  Op = "isnull"
	OpNotNull Op = "notnull"
)

// Scope is a gorm query scope.
type Scope func(*gorm.DB) *gorm.DB

var schemaCache sync.Map // reflect.Type -> *schema.Schema

func getSchema[T any]() (*schema.Schema, error) {
	t := reflect.TypeOf(new(T)).Elem()
	if cached, ok := schemaCache.Load(t); ok {
		return cached.(*schema.Schema), nil
	}
	sch, err := schema.Parse(new(T), &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return nil, err
	}
	schemaCache.Store(t, sch)
	return sch, nil
}

// scopeFor builds a gorm scope for an operator + value. It is used internally by
// the HTTP filter parser to translate validated filter parameters into SQL.
func scopeFor(dbName string, op Op, val interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		switch op {
		case OpEq:
			return db.Where(dbName+" = ?", val)
		case OpNeq:
			return db.Where(dbName+" <> ?", val)
		case OpGt:
			return db.Where(dbName+" > ?", val)
		case OpGte:
			return db.Where(dbName+" >= ?", val)
		case OpLt:
			return db.Where(dbName+" < ?", val)
		case OpLte:
			return db.Where(dbName+" <= ?", val)
		case OpLike:
			return db.Where(dbName+" LIKE ?", "%"+toString(val)+"%")
		case OpILike:
			return db.Where(dbName+" LIKE ?", "%"+strings.ToLower(toString(val))+"%")
		case OpIn:
			if s := toSlice(val); s != nil {
				return db.Where(dbName+" IN ?", s)
			}
			return db
		case OpNotIn:
			if s := toSlice(val); s != nil {
				return db.Where(dbName+" NOT IN ?", s)
			}
			return db
		case OpBetween:
			lo, hi := toPair(val)
			return db.Where(dbName+" BETWEEN ? AND ?", lo, hi)
		case OpIsNull:
			return db.Where(dbName + " IS NULL")
		case OpNotNull:
			return db.Where(dbName + " IS NOT NULL")
		}
		return db
	}
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func toSlice(v interface{}) interface{} {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		return v
	}
	if s, ok := v.(string); ok {
		return strings.Split(s, ",")
	}
	return []interface{}{v}
}

func toPair(v interface{}) (interface{}, interface{}) {
	rv := reflect.ValueOf(v)
	if (rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array) && rv.Len() == 2 {
		return rv.Index(0).Interface(), rv.Index(1).Interface()
	}
	return v, v
}
