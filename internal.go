package echoswagger

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type ParamInType string

const (
	ParamInQuery    ParamInType = "query"
	ParamInHeader   ParamInType = "header"
	ParamInPath     ParamInType = "path"
	ParamInFormData ParamInType = "formData"
	ParamInBody     ParamInType = "body"
)

type UISetting struct {
	DetachSpec bool
	HideTop    bool
	CDN        string
}

type RawDefineDic map[string]RawDefine

type RawDefine struct {
	Value  reflect.Value
	Schema *JSONSchema
}

func (r *Root) docHandler(docPath string) echo.HandlerFunc {
	t, err := template.New("swagger").Parse(SwaggerUIContent)
	if err != nil {
		panic(err)
	}
	return func(c echo.Context) error {
		cdn := r.ui.CDN
		if cdn == "" {
			cdn = DefaultCDN
		}
		buf := new(bytes.Buffer)
		params := map[string]interface{}{
			"title":    r.spec.Info.Title,
			"cdn":      cdn,
			"specName": SpecName,
		}
		if !r.ui.DetachSpec {
			spec, err := r.GetSpec(c, docPath)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			b, err := json.Marshal(spec)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			params["spec"] = string(b)
			params["docPath"] = docPath
			params["hideTop"] = true
		} else {
			params["hideTop"] = r.ui.HideTop
		}
		t.Execute(buf, params)
		return c.HTMLBlob(http.StatusOK, buf.Bytes())
	}
}

func (r *RawDefineDic) getKey(v reflect.Value) (bool, string) {
	for k, d := range *r {
		if reflect.DeepEqual(d.Value.Interface(), v.Interface()) {
			return true, k
		}
	}
	name := v.Type().Name()
	for k := range *r {
		if name == k {
			name += "_"
		}
	}
	return false, name
}

func (r *routers) appendRoute(route *echo.Route) *api {
	opr := Operation{
		Responses: make(map[string]*Response),
	}
	a := api{
		route:     route,
		defs:      r.defs,
		operation: opr,
	}
	r.apis = append(r.apis, a)
	return &r.apis[len(r.apis)-1]
}

func (g *api) addParams(p interface{}, in ParamInType, name, desc string, required, nest bool) Api {
	if !isValidParam(reflect.TypeOf(p), nest, false) {
		panic("echoswagger: invalid " + string(in) + " param")
	}
	rt := indirectType(p)
	st, sf := toSwaggerType(rt)
	if st == "object" && sf == "object" {
		g.operation.handleParamStruct(rt, in)
	} else {
		name = g.operation.rename(name)
		pm := &Parameter{
			Name:        name,
			In:          string(in),
			Description: desc,
			Required:    required,
			Type:        st,
		}
		if st == "array" {
			pm.Items = Items{}.generate(rt.Elem())
			pm.CollectionFormat = "multi"
		} else {
			pm.Format = sf
		}
		g.operation.Parameters = append(g.operation.Parameters, pm)
	}
	return g
}

func (g *api) addBodyParams(p interface{}, name, desc string, required bool) Api {
	if !isValidSchema(reflect.TypeOf(p), false) {
		panic("echoswagger: invalid body parameter")
	}
	for _, param := range g.operation.Parameters {
		if param.In == string(ParamInBody) {
			panic("echoswagger: multiple body parameters are not allowed")
		}
	}

	rv := indirectValue(p)
	pm := &Parameter{
		Name:        name,
		In:          string(ParamInBody),
		Description: desc,
		Required:    required,
		Schema:      g.defs.genSchema(rv),
	}
	g.operation.Parameters = append(g.operation.Parameters, pm)
	return g
}

func (o Operation) rename(s string) string {
	for _, p := range o.Parameters {
		if p.Name == s {
			return o.rename(s + "_")
		}
	}
	return s
}

func (o *Operation) handleParamStruct(rt reflect.Type, in ParamInType) {
	for i := 0; i < rt.NumField(); i++ {
		if rt.Field(i).Type.Kind() == reflect.Struct && rt.Field(i).Anonymous {
			o.handleParamStruct(rt.Field(i).Type, in)
		} else {
			pm := Parameter{}.generate(rt.Field(i), in)
			if pm != nil {
				pm.Name = o.rename(pm.Name)
				o.Parameters = append(o.Parameters, pm)
			}
		}
	}
}
