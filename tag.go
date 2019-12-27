package echoswagger

import (
	"reflect"
	"strconv"
	"strings"
)

// getTag reports is a tag exists and it's content
// search tagName in all tags when index = -1
func getTag(field reflect.StructField, tagName string, index int) (bool, string) {
	t := field.Tag.Get(tagName)
	s := strings.Split(t, ",")

	if len(s) < index+1 {
		return false, ""
	}

	return true, strings.TrimSpace(s[index])
}

func getSwaggerTags(field reflect.StructField) map[string]string {
	t := field.Tag.Get("swagger")
	r := make(map[string]string)
	for _, v := range strings.Split(t, ",") {
		leftIndex := strings.Index(v, "(")
		rightIndex := strings.LastIndex(v, ")")
		if leftIndex > 0 && rightIndex > leftIndex {
			r[v[:leftIndex]] = v[leftIndex+1 : rightIndex]
		} else {
			r[v] = ""
		}
	}
	return r
}

func getFieldName(f reflect.StructField, in ParamInType) (string, bool) {
	var name string
	switch in {
	case ParamInQuery:
		name = f.Tag.Get("query")
	case ParamInFormData:
		name = f.Tag.Get("form")
	case ParamInBody, ParamInHeader, ParamInPath:
		_, name = getTag(f, "json", 0)
	}
	if name != "" {
		return name, true
	} else {
		return f.Name, false
	}
}

func (p *Parameter) handleSwaggerTags(field reflect.StructField, name string, in ParamInType) {
	tags := getSwaggerTags(field)

	if t, ok := tags["desc"]; ok {
		p.Description = t
	}
	if t, ok := tags["min"]; ok {
		if m, err := strconv.ParseFloat(t, 64); err == nil {
			p.Minimum = &m
		}
	}
	if t, ok := tags["max"]; ok {
		if m, err := strconv.ParseFloat(t, 64); err == nil {
			p.Maximum = &m
		}
	}
	if t, ok := tags["minLen"]; ok {
		if m, err := strconv.Atoi(t); err == nil {
			p.MinLength = &m
		}
	}
	if t, ok := tags["maxLen"]; ok {
		if m, err := strconv.Atoi(t); err == nil {
			p.MaxLength = &m
		}
	}
	if _, ok := tags["allowEmpty"]; ok {
		p.AllowEmptyValue = true
	}
	if _, ok := tags["required"]; ok || in == ParamInPath {
		p.Required = true
	}

	convert := converter(field.Type)
	if t, ok := tags["enum"]; ok {
		enums := strings.Split(t, "|")
		var es []interface{}
		for _, s := range enums {
			v, err := convert(s)
			if err != nil {
				continue
			}
			es = append(es, v)
		}
		p.Enum = es
	}
	if t, ok := tags["default"]; ok {
		v, err := convert(t)
		if err == nil {
			p.Default = v
		}
	}

	// Move part of tags in Parameter to Items
	if p.Type == "array" {
		items := p.Items.latest()
		items.Minimum = p.Minimum
		items.Maximum = p.Maximum
		items.MinLength = p.MinLength
		items.MaxLength = p.MaxLength
		items.Enum = p.Enum
		items.Default = p.Default
		p.Minimum = nil
		p.Maximum = nil
		p.MinLength = nil
		p.MaxLength = nil
		p.Enum = nil
		p.Default = nil
	}
}

func (s *JSONSchema) handleSwaggerTags(f reflect.StructField, name string) {
	propSchema := s.Properties[name]
	tags := getSwaggerTags(f)

	if t, ok := tags["desc"]; ok {
		propSchema.Description = t
	}
	if t, ok := tags["min"]; ok {
		if m, err := strconv.ParseFloat(t, 64); err == nil {
			propSchema.Minimum = &m
		}
	}
	if t, ok := tags["max"]; ok {
		if m, err := strconv.ParseFloat(t, 64); err == nil {
			propSchema.Maximum = &m
		}
	}
	if t, ok := tags["minLen"]; ok {
		if m, err := strconv.Atoi(t); err == nil {
			propSchema.MinLength = &m
		}
	}
	if t, ok := tags["maxLen"]; ok {
		if m, err := strconv.Atoi(t); err == nil {
			propSchema.MaxLength = &m
		}
	}
	if _, ok := tags["required"]; ok {
		s.Required = append(s.Required, name)
	}
	if _, ok := tags["readOnly"]; ok {
		propSchema.ReadOnly = true
	}

	convert := converter(f.Type)
	if t, ok := tags["enum"]; ok {
		enums := strings.Split(t, "|")
		var es []interface{}
		for _, s := range enums {
			v, err := convert(s)
			if err != nil {
				continue
			}
			es = append(es, v)
		}
		propSchema.Enum = es
	}
	if t, ok := tags["default"]; ok {
		v, err := convert(t)
		if err == nil {
			propSchema.DefaultValue = v
		}
	}

	// Move part of tags in Schema to Items
	if propSchema.Type == "array" {
		items := propSchema.Items.latest()
		items.Minimum = propSchema.Minimum
		items.Maximum = propSchema.Maximum
		items.MinLength = propSchema.MinLength
		items.MaxLength = propSchema.MaxLength
		items.Enum = propSchema.Enum
		items.DefaultValue = propSchema.DefaultValue
		propSchema.Minimum = nil
		propSchema.Maximum = nil
		propSchema.MinLength = nil
		propSchema.MaxLength = nil
		propSchema.Enum = nil
		propSchema.DefaultValue = nil
	}
}

func (h *Header) handleSwaggerTags(f reflect.StructField, name string) {
	tags := getSwaggerTags(f)

	if t, ok := tags["desc"]; ok {
		h.Description = t
	}
	if t, ok := tags["min"]; ok {
		if m, err := strconv.ParseFloat(t, 64); err == nil {
			h.Minimum = &m
		}
	}
	if t, ok := tags["max"]; ok {
		if m, err := strconv.ParseFloat(t, 64); err == nil {
			h.Maximum = &m
		}
	}
	if t, ok := tags["minLen"]; ok {
		if m, err := strconv.Atoi(t); err == nil {
			h.MinLength = &m
		}
	}
	if t, ok := tags["maxLen"]; ok {
		if m, err := strconv.Atoi(t); err == nil {
			h.MaxLength = &m
		}
	}

	convert := converter(f.Type)
	if t, ok := tags["enum"]; ok {
		enums := strings.Split(t, "|")
		var es []interface{}
		for _, s := range enums {
			v, err := convert(s)
			if err != nil {
				continue
			}
			es = append(es, v)
		}
		h.Enum = es
	}
	if t, ok := tags["default"]; ok {
		v, err := convert(t)
		if err == nil {
			h.Default = v
		}
	}

	// Move part of tags in Header to Items
	if h.Type == "array" {
		items := h.Items.latest()
		items.Minimum = h.Minimum
		items.Maximum = h.Maximum
		items.MinLength = h.MinLength
		items.MaxLength = h.MaxLength
		items.Enum = h.Enum
		items.Default = h.Default
		h.Minimum = nil
		h.Maximum = nil
		h.MinLength = nil
		h.MaxLength = nil
		h.Enum = nil
		h.Default = nil
	}
}

func (t *Items) latest() *Items {
	if t.Items != nil {
		return t.Items.latest()
	}
	return t
}

func (s *JSONSchema) latest() *JSONSchema {
	if s.Items != nil {
		return s.Items.latest()
	}
	return s
}

// Not support nested elements tag eg:"a>b>c"
// Not support tags: ",chardata", ",cdata", ",comment"
// Not support embedded structure with tag ",innerxml"
// Only support nested elements tag in array type eg:"Name []string `xml:"names>name"`"
func (s *JSONSchema) handleXMLTags(f reflect.StructField) {
	b, a := getTag(f, "xml", 1)
	if b && contains([]string{"chardata", "cdata", "comment"}, a) {
		return
	}

	if b, t := getTag(f, "xml", 0); b {
		if t == "-" || s.Ref != "" {
			return
		} else if t == "" {
			t = f.Name
		}

		if s.XML == nil {
			s.XML = &XMLSchema{}
		}
		if a == "attr" {
			s.XML.Attribute = t
		} else {
			s.XML.Name = t
		}
	}
}

func (s *JSONSchema) handleChildXMLTags(rest string, r *RawDefineDic) {
	if rest == "" {
		return
	}

	if s.Items == nil && s.Ref == "" {
		if s.XML == nil {
			s.XML = &XMLSchema{}
		}
		s.XML.Name = rest
	} else if s.Ref != "" {
		key := s.Ref[len(DefPrefix):]
		if sc, ok := (*r)[key]; ok && sc.Schema != nil {
			if sc.Schema.XML == nil {
				sc.Schema.XML = &XMLSchema{}
			}
			sc.Schema.XML.Name = rest
		}
	} else {
		if s.XML == nil {
			s.XML = &XMLSchema{}
		}
		s.XML.Wrapped = true
		i := strings.Index(rest, ">")
		if i <= 0 {
			s.XML.Name = rest
		} else {
			s.XML.Name = rest[:i]
			rest = rest[i+1:]
			s.Items.handleChildXMLTags(rest, r)
		}
	}
}
