package echoswagger

import (
	"reflect"
)

func (Items) generate(t reflect.Type) *Items {
	st, sf := toSwaggerType(t)
	item := &Items{
		Type: st,
	}
	if st == "array" {
		item.Items = Items{}.generate(t.Elem())
		item.CollectionFormat = "multi"
	} else {
		item.Format = sf
	}
	return item
}

func (Parameter) generate(f reflect.StructField, in ParamInType) *Parameter {
	name, _ := getFieldName(f, in)
	if name == "-" {
		return nil
	}
	st, sf := toSwaggerType(f.Type)
	pm := &Parameter{
		Name: name,
		In:   string(in),
		Type: st,
	}
	if st == "array" {
		pm.Items = Items{}.generate(f.Type.Elem())
		pm.CollectionFormat = "multi"
	} else {
		pm.Format = sf
	}

	pm.handleSwaggerTags(f, name, in)
	return pm
}

func (Header) generate(f reflect.StructField) *Header {
	name, _ := getFieldName(f, ParamInHeader)
	if name == "-" {
		return nil
	}
	st, sf := toSwaggerType(f.Type)
	h := &Header{
		Type: st,
	}
	if st == "array" {
		h.Items = Items{}.generate(f.Type.Elem())
		h.CollectionFormat = "multi"
	} else {
		h.Format = sf
	}

	h.handleSwaggerTags(f, name)
	return h
}

func (r *RawDefineDic) genSchema(v reflect.Value) *JSONSchema {
	if !v.IsValid() {
		return nil
	}
	v = indirect(v)
	st, sf := toSwaggerType(v.Type())
	schema := &JSONSchema{}
	if st == "array" {
		schema.Type = JSONType(st)
		if v.Len() == 0 {
			v = reflect.MakeSlice(v.Type(), 1, 1)
		}
		schema.Items = r.genSchema(v.Index(0))
	} else if st == "object" && sf == "map" {
		schema.Type = JSONType(st)
		if v.Len() == 0 {
			v = reflect.New(v.Type().Elem())
		} else {
			v = v.MapIndex(v.MapKeys()[0])
		}
		schema.AdditionalProperties = r.genSchema(v)
	} else if st == "object" {
		key := r.addDefinition(v)
		schema.Ref = DefPrefix + key
	} else {
		schema.Type = JSONType(st)
		schema.Format = sf
		zv := reflect.Zero(v.Type())
		if v.CanInterface() && zv.CanInterface() && v.Interface() != zv.Interface() {
			schema.Example = v.Interface()
		}
	}
	return schema
}

func (api) genHeader(v reflect.Value) map[string]*Header {
	rt := indirect(v).Type()
	if rt.Kind() != reflect.Struct {
		return nil
	}
	mh := make(map[string]*Header)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		h := Header{}.generate(f)
		if h != nil {
			name, _ := getFieldName(f, ParamInHeader)
			mh[name] = h
		}
	}
	return mh
}
