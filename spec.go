package echoswagger

import (
	"encoding/xml"
	"net/http"
	"reflect"

	"github.com/labstack/echo"
)

const (
	DefPrefix      = "#/definitions/"
	SwaggerVersion = "2.0"
)

func (r *Root) Spec(c echo.Context) error {
	r.once.Do(func() {
		r.err = r.genSpec(c)
		r.cleanUp()
	})
	if r.err != nil {
		return c.String(http.StatusInternalServerError, r.err.Error())
	}
	return c.JSON(http.StatusOK, r.spec)
}

func (r *Root) genSpec(c echo.Context) error {
	r.spec.Swagger = SwaggerVersion
	r.spec.Paths = make(map[string]interface{})
	r.spec.Host = c.Request().Host
	r.spec.Schemes = []string{c.Scheme()}

	for i := range r.groups {
		group := &r.groups[i]
		r.spec.Tags = append(r.spec.Tags, &group.tag)
		for j := range group.apis {
			a := &group.apis[j]
			if err := a.operation.addSecurity(r.spec.SecurityDefinitions, group.security); err != nil {
				return err
			}
			if err := r.transfer(a); err != nil {
				return err
			}
		}
	}

	for i := range r.apis {
		if err := r.transfer(&r.apis[i]); err != nil {
			return err
		}
	}

	for k, v := range *r.defs {
		r.spec.Definitions[k] = v.Schema
	}
	return nil
}

func (r *Root) transfer(a *api) error {
	if err := a.operation.addSecurity(r.spec.SecurityDefinitions, a.security); err != nil {
		return err
	}

	path := toSwaggerPath(a.route.Path)
	if len(a.operation.Responses) == 0 {
		a.operation.Responses["default"] = &Response{
			Description: "successful operation",
		}
	}

	if p, ok := r.spec.Paths[path]; ok {
		p.(*Path).oprationAssign(a.method, &a.operation)
	} else {
		p := &Path{}
		p.oprationAssign(a.method, &a.operation)
		r.spec.Paths[path] = p
	}
	return nil
}

func (p *Path) oprationAssign(method string, operation *Operation) {
	switch method {
	case echo.GET:
		p.Get = operation
	case echo.POST:
		p.Post = operation
	case echo.PUT:
		p.Put = operation
	case echo.PATCH:
		p.Patch = operation
	case echo.DELETE:
		p.Delete = operation
	}
}

func (r *Root) cleanUp() {
	r.echo = nil
	r.groups = nil
	r.apis = nil
	r.defs = nil
}

// addDefinition adds definition specification and returns
// key of RawDefineDic
func (r *RawDefineDic) addDefinition(v reflect.Value) string {
	exist, key := r.getKey(v)
	if exist {
		return key
	}

	schema := &JSONSchema{
		Type:       "object",
		Properties: make(map[string]*JSONSchema),
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		name := getFieldName(f, ParamInBody)
		if name == "-" {
			continue
		}
		if f.Type == reflect.TypeOf(xml.Name{}) {
			schema.handleXMLTags(f)
			continue
		}
		sp := r.genSchema(v.Field(i))
		sp.handleXMLTags(f)
		if sp.XML != nil {
			sp.handleChildXMLTags(sp.XML.Name, r)
		}
		schema.Properties[name] = sp

		schema.handleSwaggerTags(f, name)
	}

	(*r)[key] = RawDefine{
		Value:  v,
		Schema: schema,
	}

	if schema.XML == nil {
		schema.XML = &XMLSchema{}
	}
	if schema.XML.Name == "" {
		schema.XML.Name = v.Type().Name()
	}
	return key
}
