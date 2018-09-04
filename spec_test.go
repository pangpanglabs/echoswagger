package echoswagger

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var handleWithFilter func(handlerFunc echo.HandlerFunc, c echo.Context) error

func TestSpec(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		r := prepareApiRoot()
		e := r.(*Root).echo
		req := httptest.NewRequest(echo.GET, "/doc/swagger.json", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		j := `{"swagger":"2.0","info":{"title":"Project APIs","version":""},"host":"example.com","basePath":"/","schemes":["http"],"paths":{}}`
		if assert.NoError(t, r.(*Root).Spec(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, j, rec.Body.String())
		}
	})

	t.Run("Methods", func(t *testing.T) {
		r := prepareApiRoot()
		var h echo.HandlerFunc
		r.GET("/", h)
		r.POST("/", h)
		r.PUT("/", h)
		r.DELETE("/", h)
		r.OPTIONS("/", h)
		r.HEAD("/", h)
		r.PATCH("/", h)
		e := r.(*Root).echo
		req := httptest.NewRequest(echo.GET, "/doc/swagger.json", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, r.(*Root).Spec(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			s := r.(*Root).spec
			assert.Len(t, s.Paths, 1)
			assert.NotNil(t, s.Paths["/"].(*Path).Get)
			assert.NotNil(t, s.Paths["/"].(*Path).Post)
			assert.NotNil(t, s.Paths["/"].(*Path).Put)
			assert.NotNil(t, s.Paths["/"].(*Path).Delete)
			assert.NotNil(t, s.Paths["/"].(*Path).Options)
			assert.NotNil(t, s.Paths["/"].(*Path).Head)
			assert.NotNil(t, s.Paths["/"].(*Path).Patch)
		}
	})

	t.Run("ErrorGroupSecurity", func(t *testing.T) {
		r := prepareApiRoot()
		e := r.(*Root).echo
		var h echo.HandlerFunc
		r.Group("G", "/g").SetSecurity("JWT").GET("/", h)
		req := httptest.NewRequest(echo.GET, "/doc/swagger.json", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, r.(*Root).Spec(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("ErrorApiSecurity", func(t *testing.T) {
		r := prepareApiRoot()
		e := r.(*Root).echo
		var h echo.HandlerFunc
		r.GET("/", h).SetSecurity("JWT")
		req := httptest.NewRequest(echo.GET, "/doc/swagger.json", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, r.(*Root).Spec(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("CleanUp", func(t *testing.T) {
		r := prepareApiRoot()
		e := r.(*Root).echo
		g := r.Group("Users", "users")

		var ha echo.HandlerFunc
		g.DELETE("/:id", ha)

		var hb echo.HandlerFunc
		r.GET("/ping", hb)

		req := httptest.NewRequest(echo.GET, "/doc/swagger.json", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		j := `{"swagger":"2.0","info":{"title":"Project APIs","version":""},"host":"example.com","basePath":"/","schemes":["http"],"paths":{"/ping":{"get":{"responses":{"default":{"description":"successful operation"}}}},"/users/{id}":{"delete":{"tags":["Users"],"responses":{"default":{"description":"successful operation"}}}}},"tags":[{"name":"Users"}]}`
		if assert.NoError(t, r.(*Root).Spec(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, j, rec.Body.String())
		}

		assert.Nil(t, r.(*Root).echo)
		assert.Nil(t, r.(*Root).defs)
		assert.Len(t, r.(*Root).groups, 0)
		assert.Len(t, r.(*Root).apis, 0)
	})
}

func TestAddDefinition(t *testing.T) {
	type DA struct {
		Name string
		DB   struct {
			Name string
		}
	}
	var da DA
	r := prepareApiRoot()
	var h echo.HandlerFunc
	a := r.GET("/", h)
	a.AddParamBody(&da, "DA", "DA Struct", false)
	assert.Equal(t, len(a.(*api).operation.Parameters), 1)
	assert.Equal(t, "DA", a.(*api).operation.Parameters[0].Name)
	assert.Equal(t, "DA Struct", a.(*api).operation.Parameters[0].Description)
	assert.Equal(t, "body", a.(*api).operation.Parameters[0].In)
	assert.NotNil(t, a.(*api).operation.Parameters[0].Schema)
	assert.Equal(t, "#/definitions/DA", a.(*api).operation.Parameters[0].Schema.Ref)

	assert.NotNil(t, a.(*api).defs)
	assert.Equal(t, reflect.ValueOf(&da).Elem(), (*a.(*api).defs)["DA"].Value)
	assert.Equal(t, reflect.ValueOf(&da.DB).Elem(), (*a.(*api).defs)[""].Value)

	e := r.(*Root).echo
	req := httptest.NewRequest(echo.GET, "/doc/swagger.json", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, r.(*Root).Spec(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Len(t, r.(*Root).spec.Definitions, 2)
	}
}
