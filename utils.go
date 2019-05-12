package echoswagger

import (
	"reflect"
	"strings"
)

func contains(list []string, s string) bool {
	for _, t := range list {
		if t == s {
			return true
		}
	}
	return false
}

func containsMap(list []map[string][]string, m map[string][]string) bool {
LoopMaps:
	for _, t := range list {
		if len(t) != len(m) {
			continue
		}
		for k, v := range t {
			if mv, ok := m[k]; !ok || !equals(mv, v) {
				continue LoopMaps
			}
		}
		return true
	}
	return false
}

func equals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, t := range a {
		if !contains(b, t) {
			return false
		}
	}
	return true
}

func indirect(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		ev := v.Elem()
		if !ev.IsValid() {
			ev = reflect.New(v.Type().Elem())
		}
		return indirect(ev)
	}
	return v
}

func indirectValue(p interface{}) reflect.Value {
	v := reflect.ValueOf(p)
	return indirect(v)
}

func indirectType(p interface{}) reflect.Type {
	t := reflect.TypeOf(p)
LoopType:
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		goto LoopType
	}
	return t
}

// "" → "/"
// "/" → "/"
// "a" → "/a"
// "/a" → "/a"
// "/a/" → "/a/"
func connectPath(paths ...string) string {
	var result string
	for i, path := range paths {
		// add prefix slash
		if len(path) == 0 || path[0] != '/' {
			path = "/" + path
		}
		// remove suffix slash, ignore last path
		if i != len(paths)-1 && len(path) != 0 && path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
		result += path
	}
	return result
}

func removeTrailingSlash(path string) string {
	l := len(path) - 1
	if l > 0 && strings.HasSuffix(path, "/") {
		path = path[:l]
	}
	return path
}

func trimSuffixSlash(s, suffix string) string {
	s = connectPath(s)
	suffix = connectPath(suffix)
	s = removeTrailingSlash(s)
	suffix = removeTrailingSlash(suffix)
	return strings.TrimSuffix(s, suffix)
}
