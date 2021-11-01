package echoswagger

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestSpec(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		r := prepareApiRoot()
		e := r.(*Root).echo
		req := httptest.NewRequest(echo.GET, "/doc/swagger.json", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		j := `{"swagger":"2.0","info":{"title":"Project APIs","version":""},"host":"example.com","paths":{}}`
		if assert.NoError(t, r.(*Root).specHandler("/doc/")(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, j, rec.Body.String())
		}
	})

	t.Run("BasicGenerater", func(t *testing.T) {
		r := prepareApiRoot()
		e := r.(*Root).echo
		req := httptest.NewRequest(echo.GET, "/doc/swagger.json", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		j := `{"swagger":"2.0","info":{"title":"Project APIs","version":""},"paths":{}}`
		s, err := r.(*Root).GetSpec(c, "/doc/")
		assert.Nil(t, err)
		rs, err := json.Marshal(s)
		assert.Nil(t, err)
		assert.JSONEq(t, j, string(rs))
	})

	t.Run("BasicIntegrated", func(t *testing.T) {
		r := prepareApiRoot()
		r.SetUI(UISetting{DetachSpec: false})
		e := r.(*Root).echo
		req := httptest.NewRequest(echo.GET, "/doc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		j := `{"swagger":"2.0","info":{"title":"Project APIs","version":""},"paths":{}}`
		s, err := r.(*Root).GetSpec(c, "/doc/")
		assert.Nil(t, err)
		rs, err := json.Marshal(s)
		assert.Nil(t, err)
		assert.JSONEq(t, j, string(rs))
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
		if assert.NoError(t, r.(*Root).specHandler("/doc")(c)) {
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
		j := `{"swagger":"2.0","info":{"title":"Project APIs","version":""},"host":"example.com","paths":{"/ping":{"get":{"responses":{"default":{"description":"successful operation"}}}},"/users/{id}":{"delete":{"tags":["Users"],"responses":{"default":{"description":"successful operation"}}}}},"tags":[{"name":"Users"}]}`
		if assert.NoError(t, r.(*Root).specHandler("/doc")(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, j, rec.Body.String())
		}

		assert.Nil(t, r.(*Root).echo)
		assert.Nil(t, r.(*Root).defs)
		assert.Len(t, r.(*Root).groups, 0)
		assert.Len(t, r.(*Root).apis, 0)
	})
}

func TestReferer(t *testing.T) {
	tests := []struct {
		name, referer, host, docPath, basePath string
	}{
		{
			referer:  "http://localhost:1323/doc",
			host:     "localhost:1323",
			docPath:  "/doc",
			name:     "A",
			basePath: "",
		},
		{
			referer:  "http://localhost:1323/doc",
			host:     "localhost:1323",
			docPath:  "/doc/",
			name:     "B",
			basePath: "",
		},
		{
			referer:  "http://localhost:1323/doc/",
			host:     "localhost:1323",
			docPath:  "/doc",
			name:     "C",
			basePath: "",
		},
		{
			referer:  "http://localhost:1323/api/v1/doc",
			host:     "localhost:1323",
			docPath:  "/doc",
			name:     "D",
			basePath: "/api/v1",
		},
		{
			referer:  "1/doc",
			host:     "127.0.0.1",
			docPath:  "/doc",
			name:     "E",
			basePath: "",
		},
		{
			referer:  "http://user:pass@github.com",
			host:     "github.com",
			docPath:  "/",
			name:     "F",
			basePath: "",
		},
		{
			referer:  "https://www.github.com/v1/docs/?q=1",
			host:     "www.github.com",
			docPath:  "/docs/",
			name:     "G",
			basePath: "/v1",
		},
		{
			referer:  "https://www.github.com/?q=1#tag=TAG",
			host:     "www.github.com",
			docPath:  "",
			name:     "H",
			basePath: "",
		},
		{
			referer:  "https://www.github.com/",
			host:     "www.github.com",
			docPath:  "/doc",
			name:     "I",
			basePath: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := prepareApiRoot()
			e := r.(*Root).echo
			req := httptest.NewRequest(echo.GET, "http://127.0.0.1/doc/swagger.json", nil)
			req.Header.Add("referer", tt.referer)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if assert.NoError(t, r.(*Root).specHandler(tt.docPath)(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				var v struct {
					Host     string `json:"host"`
					BasePath string `json:"basePath"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &v)
				assert.NoError(t, err)
				assert.Equal(t, tt.host, v.Host)
				assert.Equal(t, tt.basePath, v.BasePath)
			}
		})
	}
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
	if assert.NoError(t, r.(*Root).specHandler("/doc")(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Len(t, r.(*Root).spec.Definitions, 2)
	}
}
