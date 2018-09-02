package echoswagger

import (
	"reflect"
	"strconv"
	"sync"

	"github.com/labstack/echo"
)

/*
TODO:
1.pattern
2.opreationId 重复判断

Notice:
1.不会对Email和URL进行验证，因为不影响页面的正常显示
2.只支持对应于SwaggerUI页面的Schema，不支持sw、sww等协议
3.SetSecurity/SetSecurityWithScope 传多个参数表示Security之间是AND关系；多次调用SetSecurity/SetSecurityWithScope Security之间是OR关系
4.只支持基本类型的Map Key
*/

type ApiRouter interface {
	// GET overrides `Echo#GET()` for sub-routes within the ApiGroup.
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api

	// POST overrides `Echo#POST()` for sub-routes within the ApiGroup.
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api

	// PUT overrides `Echo#PUT()` for sub-routes within the ApiGroup.
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api

	// DELETE overrides `Echo#DELETE()` for sub-routes within the ApiGroup.
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api
}

type ApiRoot interface {
	ApiRouter

	// ApiGroup creates ApiGroup. Use this instead of Echo#ApiGroup.
	Group(name, prefix string, m ...echo.MiddlewareFunc) ApiGroup

	// SetRequestContentType sets request content types.
	SetRequestContentType(types ...string) ApiRoot

	// SetResponseContentType sets response content types.
	SetResponseContentType(types ...string) ApiRoot

	// SetExternalDocs sets external docs.
	SetExternalDocs(desc, url string) ApiRoot

	// AddSecurityBasic adds `SecurityDefinition` with type basic.
	AddSecurityBasic(name, desc string) ApiRoot

	// AddSecurityAPIKey adds `SecurityDefinition` with type apikey.
	AddSecurityAPIKey(name, desc string, in SecurityInType) ApiRoot

	// AddSecurityOAuth2 adds `SecurityDefinition` with type oauth2.
	AddSecurityOAuth2(name, desc string, flow OAuth2FlowType, authorizationUrl, tokenUrl string, scopes map[string]string) ApiRoot

	// GetRaw returns raw `Swagger`. Only special case should use.
	GetRaw() *Swagger

	// SetRaw sets raw `Swagger` to ApiRoot. Only special case should use.
	SetRaw(s *Swagger) ApiRoot

	// Echo returns the embeded echo instance
	Echo() *echo.Echo
}

type ApiGroup interface {
	ApiRouter

	// SetDescription sets description for ApiGroup.
	SetDescription(desc string) ApiGroup

	// SetExternalDocs sets external docs for ApiGroup.
	SetExternalDocs(desc, url string) ApiGroup

	// SetSecurity sets Security for all operations within the ApiGroup
	// which names are reigisters by AddSecurity... functions.
	SetSecurity(names ...string) ApiGroup

	// SetSecurity sets Security with scopes for all operations
	// within the ApiGroup which names are reigisters
	// by AddSecurity... functions.
	// Should only use when Security type is oauth2.
	SetSecurityWithScope(s map[string][]string) ApiGroup

	// EchoGroup returns `*echo.Group` within the ApiGroup.
	EchoGroup() *echo.Group
}

type Api interface {
	// AddParamPath adds path parameter.
	AddParamPath(p interface{}, name, desc string) Api

	// AddParamPathNested adds path parameters nested in p.
	AddParamPathNested(p interface{}) Api

	// AddParamQuery adds query parameter.
	AddParamQuery(p interface{}, name, desc string, required bool) Api

	// AddParamQueryNested adds query parameters nested in p.
	AddParamQueryNested(p interface{}) Api

	// AddParamForm adds formData parameter.
	AddParamForm(p interface{}, name, desc string, required bool) Api

	// AddParamFormNested adds formData parameters nested in p.
	AddParamFormNested(p interface{}) Api

	// AddParamHeader adds header parameter.
	AddParamHeader(p interface{}, name, desc string, required bool) Api

	// AddParamHeaderNested adds header parameters nested in p.
	AddParamHeaderNested(p interface{}) Api

	// AddParamBody adds body parameter.
	AddParamBody(p interface{}, name, desc string, required bool) Api

	// AddParamBody adds file parameter.
	AddParamFile(name, desc string, required bool) Api

	// AddResponse adds response for Api.
	AddResponse(code int, desc string, schema interface{}, header interface{}) Api

	// SetResponseContentType sets request content types.
	SetRequestContentType(types ...string) Api

	// SetResponseContentType sets response content types.
	SetResponseContentType(types ...string) Api

	// SetOperationId sets operationId
	SetOperationId(id string) Api

	// SetDescription marks Api as deprecated.
	SetDeprecated() Api

	// SetDescription sets description.
	SetDescription(desc string) Api

	// SetExternalDocs sets external docs.
	SetExternalDocs(desc, url string) Api

	// SetExternalDocs sets summary.
	SetSummary(summary string) Api

	// SetSecurity sets Security which names are reigisters
	// by AddSecurity... functions.
	SetSecurity(names ...string) Api

	// SetSecurity sets Security for Api which names are
	// reigisters by AddSecurity... functions.
	// Should only use when Security type is oauth2.
	SetSecurityWithScope(s map[string][]string) Api

	// TODO return echo.Router
}

type routers struct {
	apis []api
	defs *RawDefineDic
}

type Root struct {
	routers
	spec   *Swagger
	echo   *echo.Echo
	groups []group
	once   sync.Once
	err    error
}

type group struct {
	routers
	echoGroup *echo.Group
	security  []map[string][]string
	tag       Tag
}

type api struct {
	route     *echo.Route
	defs      *RawDefineDic
	security  []map[string][]string
	method    string
	operation Operation
}

// New creates ApiRoot instance.
// Multiple ApiRoot are allowed in one project.
func New(e *echo.Echo, basePath, docPath string, i *Info) ApiRoot {
	if e == nil {
		panic("echoswagger: invalid Echo instance")
	}
	basePath = proccessPath(basePath)
	docPath = proccessPath(docPath)

	var connector string
	if docPath[len(docPath)-1] != '/' {
		connector = "/"
	}
	specPath := docPath + connector + "swagger.json"

	if i == nil {
		i = &Info{
			Title: "Project APIs",
		}
	}
	defs := make(RawDefineDic)
	r := &Root{
		echo: e,
		spec: &Swagger{
			Info:                i,
			SecurityDefinitions: make(map[string]*SecurityDefinition),
			BasePath:            basePath,
			Definitions:         make(map[string]*JSONSchema),
		},
		routers: routers{
			defs: &defs,
		},
	}

	e.GET(docPath, r.docHandler(specPath))
	e.GET(specPath, r.Spec)
	return r
}

func (r *Root) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return r.appendRoute(echo.GET, r.echo.GET(path, h, m...))
}

func (r *Root) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return r.appendRoute(echo.POST, r.echo.POST(path, h, m...))
}

func (r *Root) PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return r.appendRoute(echo.PUT, r.echo.PUT(path, h, m...))
}

func (r *Root) DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return r.appendRoute(echo.DELETE, r.echo.DELETE(path, h, m...))
}

func (r *Root) Group(name, prefix string, m ...echo.MiddlewareFunc) ApiGroup {
	if name == "" {
		panic("echowagger: invalid name of ApiGroup")
	}
	echoGroup := r.echo.Group(prefix, m...)
	group := group{
		echoGroup: echoGroup,
		routers: routers{
			defs: r.defs,
		},
	}
	var counter int
LoopTags:
	for _, t := range r.groups {
		if t.tag.Name == name {
			if counter > 0 {
				name = name[:len(name)-2]
			}
			counter++
			name += "_" + strconv.Itoa(counter)
			goto LoopTags
		}
	}
	group.tag = Tag{Name: name}
	r.groups = append(r.groups, group)
	return &r.groups[len(r.groups)-1]
}

func (r *Root) SetRequestContentType(types ...string) ApiRoot {
	r.spec.Consumes = types
	return r
}

func (r *Root) SetResponseContentType(types ...string) ApiRoot {
	r.spec.Produces = types
	return r
}

func (r *Root) SetExternalDocs(desc, url string) ApiRoot {
	r.spec.ExternalDocs = &ExternalDocs{
		Description: desc,
		URL:         url,
	}
	return r
}

func (r *Root) AddSecurityBasic(name, desc string) ApiRoot {
	if !r.checkSecurity(name) {
		return r
	}
	sd := &SecurityDefinition{
		Type:        string(SecurityBasic),
		Description: desc,
	}
	r.spec.SecurityDefinitions[name] = sd
	return r
}

func (r *Root) AddSecurityAPIKey(name, desc string, in SecurityInType) ApiRoot {
	if !r.checkSecurity(name) {
		return r
	}
	sd := &SecurityDefinition{
		Type:        string(SecurityAPIKey),
		Description: desc,
		Name:        name,
		In:          string(in),
	}
	r.spec.SecurityDefinitions[name] = sd
	return r
}

func (r *Root) AddSecurityOAuth2(name, desc string, flow OAuth2FlowType, authorizationUrl, tokenUrl string, scopes map[string]string) ApiRoot {
	if !r.checkSecurity(name) {
		return r
	}
	sd := &SecurityDefinition{
		Type:             string(SecurityOAuth2),
		Description:      desc,
		Flow:             string(flow),
		AuthorizationURL: authorizationUrl,
		TokenURL:         tokenUrl,
		Scopes:           scopes,
	}
	r.spec.SecurityDefinitions[name] = sd
	return r
}

func (r *Root) GetRaw() *Swagger {
	return r.spec
}

func (r *Root) SetRaw(s *Swagger) ApiRoot {
	r.spec = s
	return r
}

func (r *Root) Echo() *echo.Echo {
	return r.echo
}

func (g *group) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	a := g.appendRoute(echo.GET, g.echoGroup.GET(path, h, m...))
	a.operation.Tags = []string{g.tag.Name}
	return a
}

func (g *group) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	a := g.appendRoute(echo.POST, g.echoGroup.POST(path, h, m...))
	a.operation.Tags = []string{g.tag.Name}
	return a
}

func (g *group) PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	a := g.appendRoute(echo.PUT, g.echoGroup.PUT(path, h, m...))
	a.operation.Tags = []string{g.tag.Name}
	return a
}

func (g *group) DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	a := g.appendRoute(echo.DELETE, g.echoGroup.DELETE(path, h, m...))
	a.operation.Tags = []string{g.tag.Name}
	return a
}

func (g *group) SetDescription(desc string) ApiGroup {
	g.tag.Description = desc
	return g
}

func (g *group) SetExternalDocs(desc, url string) ApiGroup {
	g.tag.ExternalDocs = &ExternalDocs{
		Description: desc,
		URL:         url,
	}
	return g
}

func (g *group) SetSecurity(names ...string) ApiGroup {
	if len(names) == 0 {
		return g
	}
	g.security = setSecurity(g.security, names...)
	return g
}

func (g *group) SetSecurityWithScope(s map[string][]string) ApiGroup {
	g.security = setSecurityWithScope(g.security, s)
	return g
}

func (g *group) EchoGroup() *echo.Group {
	return g.echoGroup
}

func (a *api) AddParamPath(p interface{}, name, desc string) Api {
	return a.addParams(p, ParamInPath, name, desc, true, false)
}

func (a *api) AddParamPathNested(p interface{}) Api {
	return a.addParams(p, ParamInPath, "", "", true, true)
}

func (a *api) AddParamQuery(p interface{}, name, desc string, required bool) Api {
	return a.addParams(p, ParamInQuery, name, desc, required, false)
}

func (a *api) AddParamQueryNested(p interface{}) Api {
	return a.addParams(p, ParamInQuery, "", "", false, true)
}

func (a *api) AddParamForm(p interface{}, name, desc string, required bool) Api {
	return a.addParams(p, ParamInFormData, name, desc, required, false)
}

func (a *api) AddParamFormNested(p interface{}) Api {
	return a.addParams(p, ParamInFormData, "", "", false, true)
}

func (a *api) AddParamHeader(p interface{}, name, desc string, required bool) Api {
	return a.addParams(p, ParamInHeader, name, desc, required, false)
}

func (a *api) AddParamHeaderNested(p interface{}) Api {
	return a.addParams(p, ParamInHeader, "", "", false, true)
}

func (a *api) AddParamBody(p interface{}, name, desc string, required bool) Api {
	return a.addBodyParams(p, name, desc, required)
}

func (a *api) AddParamFile(name, desc string, required bool) Api {
	name = a.operation.rename(name)
	a.operation.Parameters = append(a.operation.Parameters, &Parameter{
		Name:        name,
		In:          string(ParamInFormData),
		Description: desc,
		Required:    required,
		Type:        "file",
	})
	return a
}

// Notice: header must be nested in a struct.
func (a *api) AddResponse(code int, desc string, schema interface{}, header interface{}) Api {
	r := &Response{
		Description: desc,
	}

	st := reflect.TypeOf(schema)
	if st != nil {
		if !isValidSchema(st, false) {
			panic("echoswagger: invalid response schema")
		}
		r.Schema = a.defs.genSchema(reflect.ValueOf(schema))
	}

	ht := reflect.TypeOf(header)
	if ht != nil {
		if !isValidParam(reflect.TypeOf(header), true, false) {
			panic("echoswagger: invalid response header")
		}
		r.Headers = a.genHeader(reflect.ValueOf(header))
	}

	cstr := strconv.Itoa(code)
	a.operation.Responses[cstr] = r
	return a
}

func (a *api) SetRequestContentType(types ...string) Api {
	a.operation.Consumes = types
	return a
}

func (a *api) SetResponseContentType(types ...string) Api {
	a.operation.Produces = types
	return a
}

func (a *api) SetOperationId(id string) Api {
	a.operation.OperationID = id
	return a
}

func (a *api) SetDeprecated() Api {
	a.operation.Deprecated = true
	return a
}

func (a *api) SetDescription(desc string) Api {
	a.operation.Description = desc
	return a
}

func (a *api) SetExternalDocs(desc, url string) Api {
	a.operation.ExternalDocs = &ExternalDocs{
		Description: desc,
		URL:         url,
	}
	return a
}

func (a *api) SetSummary(summary string) Api {
	a.operation.Summary = summary
	return a
}

func (a *api) SetSecurity(names ...string) Api {
	if len(names) == 0 {
		return a
	}
	a.security = setSecurity(a.security, names...)
	return a
}

func (a *api) SetSecurityWithScope(s map[string][]string) Api {
	a.security = setSecurityWithScope(a.security, s)
	return a
}
