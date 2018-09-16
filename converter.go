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
	return connectPath(path)
}

func converter(t reflect.Type) func(s string) (interface{}, error) {
	st, sf := toSwaggerType(t)
	if st == "integer" && sf == "int32" {
		return func(s string) (interface{}, error) {
			v, err := strconv.Atoi(s)
			return v, err
		}
	} else if st == "integer" && sf == "int64" {
		return func(s string) (interface{}, error) {
			v, err := strconv.ParseInt(s, 10, 64)
			return v, err
		}
	} else if st == "number" && sf == "float" {
		return func(s string) (interface{}, error) {
			v, err := strconv.ParseFloat(s, 32)
			return float32(v), err
		}
	} else if st == "number" && sf == "double" {
		return func(s string) (interface{}, error) {
			v, err := strconv.ParseFloat(s, 64)
			return v, err
		}
	} else if st == "boolean" && sf == "boolean" {
		return func(s string) (interface{}, error) {
			v, err := strconv.ParseBool(s)
			return v, err
		}
	} else if st == "array" && sf == "array" {
		return converter(t.Elem())
	} else {
		return func(s string) (interface{}, error) {
			return s, nil
		}
	}
}
