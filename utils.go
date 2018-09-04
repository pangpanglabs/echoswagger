package echoswagger

import "reflect"

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
