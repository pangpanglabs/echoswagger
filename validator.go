package echoswagger

import (
	"net/url"
	"reflect"
	"regexp"
	"time"
)

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isValidEmail(str string) bool {
	return emailRegexp.MatchString(str)
}

// See: https://github.com/swagger-api/swagger-js/blob/7414ad062ba9b6d9cc397c72e7561ec775b35a9f/lib/shred/parseUri.js#L28
func isValidURL(str string) bool {
	if _, err := url.ParseRequestURI(str); err != nil {
		return false
	}
	return true
}

func isValidParam(t reflect.Type, nest, inner bool) bool {
	if t == nil {
		return false
	}
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.String:
		if !nest || (nest && inner) {
			return true
		}
	case reflect.Array, reflect.Slice:
		return isValidParam(t.Elem(), nest, true)
	case reflect.Ptr:
		return isValidParam(t.Elem(), nest, inner)
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) && (!nest || nest && inner) {
			return true
		} else if !inner {
			for i := 0; i < t.NumField(); i++ {
				inner := true
				if t.Field(i).Type.Kind() == reflect.Struct && t.Field(i).Anonymous {
					inner = false
				}
				if !isValidParam(t.Field(i).Type, nest, inner) {
					return false
				}
			}
			return true
		}
	}
	return false
}

// isValidSchema reports a type is valid for body param.
// valid case:
// 1. Struct
// 2. Struct{ A int64 }
// 3. *[]Struct
// 4. [][]Struct
// 5. []Struct{ A []Struct }
// 6. []Struct{ A Map[string]string }
// 7. *Struct{ A []Map[int64]Struct }
// 8. Map[string]string
// 9. []int64
// invalid case:
// 1. interface{}
// 2. Map[Struct]string
func isValidSchema(t reflect.Type, inner bool, pres ...reflect.Type) bool {
	if t == nil {
		return false
	}
	for _, pre := range pres {
		if t == pre {
			return true
		}
	}

	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Interface:
		return true
	case reflect.Array, reflect.Slice:
		return isValidSchema(t.Elem(), inner, pres...)
	case reflect.Map:
		return isBasicType(t.Key()) && isValidSchema(t.Elem(), true, pres...)
	case reflect.Ptr:
		return isValidSchema(t.Elem(), inner, pres...)
	case reflect.Struct:
		pres = append(pres, t)
		if t == reflect.TypeOf(time.Time{}) {
			return true
		}
		for i := 0; i < t.NumField(); i++ {
			if !isValidSchema(t.Field(i).Type, true, pres...) {
				return false
			}
		}
		return true
	}
	return false
}

func isBasicType(t reflect.Type) bool {
	if t == nil {
		return false
	}
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.String:
		return true
	case reflect.Ptr:
		return isBasicType(t.Elem())
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			return true
		}
	}
	return false
}

func isValidScheme(s string) bool {
	if s == "http" || s == "https" || s == "ws" || s == "wss" {
		return true
	}
	return false
}
