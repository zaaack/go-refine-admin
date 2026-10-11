package query

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// fieldIndex builds a lookup from API field names (Go name, db name, or json
// tag) to the gorm schema field. This lets clients filter/sort by the same name
// they see in JSON while we still validate against the real schema.
func fieldIndex[T any]() (map[string]*schema.Field, error) {
	sch, err := getSchema[T]()
	if err != nil {
		return nil, err
	}
	t := reflect.TypeOf(new(T)).Elem()
	idx := make(map[string]*schema.Field, len(sch.Fields)*4)
	for _, f := range sch.Fields {
		idx[f.Name] = f
		idx[f.DBName] = f
		idx[strings.ToLower(f.DBName)] = f
		idx[strings.ToLower(f.Name)] = f
		if sf, ok := t.FieldByName(f.Name); ok {
			if jt := sf.Tag.Get("json"); jt != "" {
				jn := strings.Split(jt, ",")[0]
				if jn != "" && jn != "-" {
					idx[jn] = f
					idx[strings.ToLower(jn)] = f
				}
			}
		}
	}
	return idx, nil
}

// ParseFilters builds gorm scopes from HTTP query parameters for model T.
//
// Parameter convention:  field__op=value
//   - field  : a struct field of T addressed by Go name, db column or json tag
//              (validated against the gorm schema)
//   - op     : eq|neq|gt|gte|lt|lte|like|ilike|in|notin|between|isnull|notnull
//              (omitted => eq)
// Unknown fields or type-incompatible values return an error, making the query
// conditions type-safe at the API boundary.
func ParseFilters[T any](values url.Values) ([]Scope, error) {
	idx, err := fieldIndex[T]()
	if err != nil {
		return nil, err
	}
	var scopes []Scope
	for rawKey, rawVals := range values {
		if len(rawVals) == 0 {
			continue
		}
		field, op := splitKey(rawKey)
		// skip pagination / sorting / reserved keys
		if field == "page" || field == "pageSize" || field == "sort" || field == "_" {
			continue
		}
		f, ok := idx[field]
		if !ok {
			f, ok = idx[strings.ToLower(field)]
		}
		if !ok {
			return nil, fmt.Errorf("unknown filter field %q", field)
		}
		raw := rawVals[0]
		switch op {
		case OpIsNull:
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB { return db.Where(f.DBName + " IS NULL") })
		case OpNotNull:
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB { return db.Where(f.DBName + " IS NOT NULL") })
		case OpBetween:
			parts := strings.SplitN(raw, ",", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("field %q between requires two comma-separated values", field)
			}
			lo, err := convertScalar(f.FieldType, parts[0])
			if err != nil {
				return nil, err
			}
			hi, err := convertScalar(f.FieldType, parts[1])
			if err != nil {
				return nil, err
			}
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
				return db.Where(f.DBName+" BETWEEN ? AND ?", lo, hi)
			})
		case OpIn, OpNotIn:
			parts := strings.Split(raw, ",")
			slice := reflect.MakeSlice(reflect.SliceOf(f.FieldType), 0, len(parts))
			for _, p := range parts {
				v, err := convertScalar(f.FieldType, p)
				if err != nil {
					return nil, err
				}
				slice = reflect.Append(slice, reflect.ValueOf(v))
			}
			col := f.DBName
			neg := op == OpNotIn
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
				if neg {
					return db.Where(col+" NOT IN ?", slice.Interface())
				}
				return db.Where(col+" IN ?", slice.Interface())
			})
		case OpLike, OpILike:
			like := "%" + raw + "%"
			col := f.DBName
			if op == OpILike {
				like = "%" + strings.ToLower(raw) + "%"
			}
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB { return db.Where(col+" LIKE ?", like) })
		default:
			v, err := convertScalar(f.FieldType, raw)
			if err != nil {
				return nil, err
			}
			col, oper := f.DBName, sqlOp(op)
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB { return db.Where(col+" "+oper+" ?", v) })
		}
	}
	return scopes, nil
}

func splitKey(key string) (field string, op Op) {
	if i := strings.Index(key, "__"); i >= 0 {
		return key[:i], Op(key[i+2:])
	}
	return key, OpEq
}

// sqlOp maps an Op to its SQL comparison operator.
func sqlOp(op Op) string {
	switch op {
	case OpNeq:
		return "<>"
	case OpGt:
		return ">"
	case OpGte:
		return ">="
	case OpLt:
		return "<"
	case OpLte:
		return "<="
	default:
		return "="
	}
}

// ParseSort parses a "sort" param: "name,-created_at" (leading '-' = desc).
// Unknown fields are ignored.
func ParseSort[T any](values url.Values) []Scope {
	raw := values.Get("sort")
	if raw == "" {
		return nil
	}
	idx, err := fieldIndex[T]()
	if err != nil {
		return nil
	}
	var scopes []Scope
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		desc := false
		if strings.HasPrefix(part, "-") {
			desc = true
			part = part[1:]
		}
		f, ok := idx[part]
		if !ok {
			f, ok = idx[strings.ToLower(part)]
		}
		if !ok {
			continue
		}
		dbName, order := f.DBName, "ASC"
		if desc {
			order = "DESC"
		}
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB { return db.Order(dbName + " " + order) })
	}
	return scopes
}

// ParsePage extracts page/pageSize (defaults 1 / 10, pageSize capped at 200).
func ParsePage(values url.Values) (page, pageSize int) {
	page, _ = strconv.Atoi(values.Get("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ = strconv.Atoi(values.Get("pageSize"))
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 200 {
		pageSize = 200
	}
	return
}

// convertScalar converts a raw string into the given Go type.
func convertScalar(t reflect.Type, raw string) (interface{}, error) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.String:
		return raw, nil
	case reflect.Bool:
		return strconv.ParseBool(raw)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.ParseInt(raw, 10, 64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.ParseUint(raw, 10, 64)
	case reflect.Float32, reflect.Float64:
		return strconv.ParseFloat(raw, 64)
	default:
		if t == reflect.TypeOf(time.Time{}) {
			for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"} {
				if tm, err := time.Parse(layout, raw); err == nil {
					return tm, nil
				}
			}
		}
		return raw, nil
	}
}
