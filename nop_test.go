package echoswagger

import (
	"net/http"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func testHandler(echo.Context) error { return nil }

func TestNop(t *testing.T) {
	e := echo.New()
	r := NewNop(e)

	path := "test"
	expectHandler := "github.com/pangpanglabs/echoswagger.testHandler"
	a := r.Add(http.MethodConnect, path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodConnect)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, path)

	a = r.GET(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodGet)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, path)

	a = r.POST(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodPost)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, path)

	a = r.PUT(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodPut)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, path)

	a = r.DELETE(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodDelete)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, path)

	a = r.OPTIONS(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodOptions)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, path)

	a = r.HEAD(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodHead)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, path)

	a = r.PATCH(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodPatch)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, path)

	g := r.Group("", "g/")
	assert.EqualValues(t, g.EchoGroup(), r.Echo().Group("g/"))

	eg := g.EchoGroup()
	assert.Equal(t, eg, r.BindGroup("", eg).EchoGroup())

	assert.Equal(t, r.SetRequestContentType(), r)
	assert.Equal(t, r.SetResponseContentType(), r)
	assert.Equal(t, r.SetExternalDocs("", ""), r)
	assert.Equal(t, r.AddSecurityBasic("", ""), r)
	assert.Equal(t, r.AddSecurityAPIKey("", "", ""), r)
	assert.Equal(t, r.AddSecurityOAuth2("", "", "", "", "", nil), r)
	assert.Equal(t, r.SetUI(UISetting{}), r)
	assert.Equal(t, r.SetScheme(), r)
	assert.Nil(t, r.GetRaw())
	assert.Equal(t, r.SetRaw(nil), r)
	assert.Equal(t, r.Echo(), e)

	expectPath := "g/" + path
	a = g.Add(http.MethodConnect, path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodConnect)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, expectPath)

	a = g.GET(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodGet)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, expectPath)

	a = g.POST(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodPost)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, expectPath)

	a = g.PUT(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodPut)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, expectPath)

	a = g.DELETE(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodDelete)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, expectPath)

	a = g.OPTIONS(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodOptions)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, expectPath)

	a = g.HEAD(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodHead)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, expectPath)

	a = g.PATCH(path, testHandler)
	assert.Equal(t, a.Route().Method, http.MethodPatch)
	assert.Equal(t, a.Route().Name, expectHandler)
	assert.Equal(t, a.Route().Path, expectPath)

	assert.Equal(t, g.SetDescription(""), g)
	assert.Equal(t, g.SetExternalDocs("", ""), g)
	assert.Equal(t, g.SetSecurity(), g)
	assert.Equal(t, g.SetSecurityWithScope(nil), g)

	assert.Equal(t, a.AddParamPath(nil, "", ""), a)
	assert.Equal(t, a.AddParamPathNested(nil), a)
	assert.Equal(t, a.AddParamQuery(nil, "", "", false), a)
	assert.Equal(t, a.AddParamQueryNested(nil), a)
	assert.Equal(t, a.AddParamForm(nil, "", "", false), a)
	assert.Equal(t, a.AddParamFormNested(nil), a)
	assert.Equal(t, a.AddParamHeader(nil, "", "", false), a)
	assert.Equal(t, a.AddParamHeaderNested(nil), a)
	assert.Equal(t, a.AddParamBody(nil, "", "", false), a)
	assert.Equal(t, a.AddParamFile("", "", false), a)
	assert.Equal(t, a.SetRequestContentType(), a)
	assert.Equal(t, a.SetResponseContentType(), a)
	assert.Equal(t, a.AddResponse(0, "", nil, nil), a)
	assert.Equal(t, a.SetOperationId(""), a)
	assert.Equal(t, a.SetDeprecated(), a)
	assert.Equal(t, a.SetDescription(""), a)
	assert.Equal(t, a.SetExternalDocs("", ""), a)
	assert.Equal(t, a.SetSummary(""), a)
	assert.Equal(t, a.SetSecurity(), a)
	assert.Equal(t, a.SetSecurityWithScope(nil), a)
}
