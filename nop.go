package echoswagger

import (
	"github.com/labstack/echo/v4"
)

type NopRoot struct {
	echo *echo.Echo
}

var _ ApiRoot = NewNop(nil)

func NewNop(e *echo.Echo) ApiRoot {
	return &NopRoot{echo: e}
}

type nopGroup struct {
	echoGroup *echo.Group
}

type nopApi struct {
	route *echo.Route
}

func (r *NopRoot) Add(method string, path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: r.echo.Add(method, path, h, m...)}
}

func (r *NopRoot) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: r.echo.GET(path, h, m...)}
}

func (r *NopRoot) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: r.echo.POST(path, h, m...)}
}

func (r *NopRoot) PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: r.echo.PUT(path, h, m...)}
}

func (r *NopRoot) DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: r.echo.DELETE(path, h, m...)}
}

func (r *NopRoot) OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: r.echo.OPTIONS(path, h, m...)}
}

func (r *NopRoot) HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: r.echo.HEAD(path, h, m...)}
}

func (r *NopRoot) PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: r.echo.PATCH(path, h, m...)}
}

func (r *NopRoot) Group(_ string, prefix string, m ...echo.MiddlewareFunc) ApiGroup {
	return &nopGroup{echoGroup: r.echo.Group(prefix, m...)}
}

func (r *NopRoot) BindGroup(_ string, g *echo.Group) ApiGroup {
	return &nopGroup{echoGroup: g}
}

func (r *NopRoot) SetRequestContentType(_ ...string) ApiRoot {
	return r
}

func (r *NopRoot) SetResponseContentType(_ ...string) ApiRoot {
	return r
}

func (r *NopRoot) SetExternalDocs(_, _ string) ApiRoot {
	return r
}

func (r *NopRoot) AddSecurityBasic(_, _ string) ApiRoot {
	return r
}

func (r *NopRoot) AddSecurityAPIKey(_, _ string, _ SecurityInType) ApiRoot {
	return r
}

func (r *NopRoot) AddSecurityOAuth2(_, _ string, _ OAuth2FlowType,
	_, _ string, _ map[string]string) ApiRoot {
	return r
}

func (r *NopRoot) SetUI(_ UISetting) ApiRoot {
	return r
}

func (r *NopRoot) SetScheme(_ ...string) ApiRoot {
	return r
}

func (r *NopRoot) GetRaw() *Swagger {
	return nil
}

func (r *NopRoot) SetRaw(_ *Swagger) ApiRoot {
	return r
}

func (r *NopRoot) Echo() *echo.Echo {
	return r.echo
}

func (g *nopGroup) Add(method, path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: g.echoGroup.Add(method, path, h, m...)}
}

func (g *nopGroup) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: g.echoGroup.GET(path, h, m...)}
}

func (g *nopGroup) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: g.echoGroup.POST(path, h, m...)}
}

func (g *nopGroup) PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: g.echoGroup.PUT(path, h, m...)}
}

func (g *nopGroup) DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: g.echoGroup.DELETE(path, h, m...)}
}

func (g *nopGroup) OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: g.echoGroup.OPTIONS(path, h, m...)}
}

func (g *nopGroup) HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: g.echoGroup.HEAD(path, h, m...)}
}

func (g *nopGroup) PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) Api {
	return &nopApi{route: g.echoGroup.PATCH(path, h, m...)}
}

func (g *nopGroup) SetDescription(_ string) ApiGroup {
	return g
}

func (g *nopGroup) SetExternalDocs(_, _ string) ApiGroup {
	return g
}

func (g *nopGroup) SetSecurity(_ ...string) ApiGroup {
	return g
}

func (g *nopGroup) SetSecurityWithScope(_ map[string][]string) ApiGroup {
	return g
}

func (g *nopGroup) EchoGroup() *echo.Group {
	return g.echoGroup
}

func (a *nopApi) AddParamPath(_ interface{}, _, _ string) Api {
	return a
}

func (a *nopApi) AddParamPathNested(_ interface{}) Api {
	return a
}

func (a *nopApi) AddParamQuery(_ interface{}, _, _ string, _ bool) Api {
	return a
}

func (a *nopApi) AddParamQueryNested(_ interface{}) Api {
	return a
}

func (a *nopApi) AddParamForm(_ interface{}, _, _ string, _ bool) Api {
	return a
}

func (a *nopApi) AddParamFormNested(_ interface{}) Api {
	return a
}

func (a *nopApi) AddParamHeader(_ interface{}, _, _ string, _ bool) Api {
	return a
}

func (a *nopApi) AddParamHeaderNested(_ interface{}) Api {
	return a
}

func (a *nopApi) AddParamBody(_ interface{}, _, _ string, _ bool) Api {
	return a
}

func (a *nopApi) AddParamFile(_, _ string, _ bool) Api {
	return a
}

func (a *nopApi) SetRequestContentType(_ ...string) Api {
	return a
}

func (a *nopApi) SetResponseContentType(_ ...string) Api {
	return a
}

func (a *nopApi) AddResponse(_ int, _ string, _, _ interface{}) Api {
	return a
}

func (a *nopApi) SetOperationId(_ string) Api {
	return a
}

func (a *nopApi) SetDeprecated() Api {
	return a
}

func (a *nopApi) SetDescription(_ string) Api {
	return a
}

func (a *nopApi) SetExternalDocs(_, _ string) Api {
	return a
}

func (a *nopApi) SetSummary(_ string) Api {
	return a
}

func (a *nopApi) SetSecurity(_ ...string) Api {
	return a
}

func (a *nopApi) SetSecurityWithScope(_ map[string][]string) Api {
	return a
}

func (a *nopApi) Route() *echo.Route {
	return a.route
}
