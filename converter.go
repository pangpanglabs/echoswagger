package echoswagger

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

// toSwaggerType returns type、format for a reflect.Type in swagger format
func toSwaggerType(t reflect.Type) (string, string) {
	if t == reflect.TypeOf(time.Time{}) {
		return "string", "date-time"
	}
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer", "int32"
	case reflect.Int64, reflect.Uint64:
		return "integer", "int64"
	case reflect.Float32:
		return "number", "float"
	case reflect.Float64:
		return "number", "double"
	case reflect.String:
		return "string", "string"
	case reflect.Bool:
		return "boolean", "boolean"
	case reflect.Struct:
		return "object", "object"
	case reflect.Map:
		return "object", "map"
	case reflect.Array, reflect.Slice:
		return "array", "array"
	case reflect.Ptr:
		return toSwaggerType(t.Elem())
	default:
		return "string", "string"
	}
}

// toSwaggerPath returns path in swagger format
func toSwaggerPath(path string) string {
	var params []string
	for i := 0; i < len(path); i++ {
		if path[i] == ':' {
			j := i + 1
			for ; i < len(path) && path[i] != '/'; i++ {
			}
			params = append(params, path[j:i])
		}
	}

	for _, name := range params {
		path = strings.Replace(path, ":"+name, "{"+name+"}", 1)
	}
	return proccessPath(path)
}

func proccessPath(path string) string {
	if len(path) == 0 || path[0] != '/' {
		path = "/" + path
	}
	return path
}

// converter returns string to target type converter for a reflect.StructField
func converter(f reflect.StructField) func(s string) (interface{}, error) {
	switch f.Type.Kind() {
	case reflect.Bool:
		return func(s string) (interface{}, error) {
			v, err := strconv.ParseBool(s)
			return v, err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return func(s string) (interface{}, error) {
			v, err := strconv.Atoi(s)
			return v, err
		}
	case reflect.Int64, reflect.Uint64:
		return func(s string) (interface{}, error) {
			v, err := strconv.ParseInt(s, 10, 64)
			return v, err
		}
	case reflect.Float32:
		return func(s string) (interface{}, error) {
			v, err := strconv.ParseFloat(s, 32)
			return float32(v), err
		}
	case reflect.Float64:
		return func(s string) (interface{}, error) {
			v, err := strconv.ParseFloat(s, 64)
			return v, err
		}
	default:
		return func(s string) (interface{}, error) {
			return s, nil
		}
	}
}

func asString(rv reflect.Value) (string, bool) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		av := rv.Int()
		if av != 0 {
			return strconv.FormatInt(av, 10), true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		av := rv.Uint()
		if av != 0 {
			return strconv.FormatUint(av, 10), true
		}
	case reflect.Float64:
		av := rv.Float()
		if av != 0 {
			return strconv.FormatFloat(av, 'g', -1, 64), true
		}
	case reflect.Float32:
		av := rv.Float()
		if av != 0 {
			return strconv.FormatFloat(av, 'g', -1, 32), true
		}
	case reflect.Bool:
		av := rv.Bool()
		if av {
			return strconv.FormatBool(av), true
		}
	case reflect.String:
		av := rv.String()
		if av != "" {
			return av, true
		}
	}
	return "", false
}
