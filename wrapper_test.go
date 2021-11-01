package echoswagger

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func prepareApiRoot() ApiRoot {
	r := New(echo.New(), "doc/", nil)
	r.SetUI(UISetting{DetachSpec: true})
	return r
}

func prepareApiGroup() ApiGroup {
	r := prepareApiRoot()
	return r.Group("G", "/g")
}

func prepareApi() Api {
	g := prepareApiGroup()
	var h func(e echo.Context) error
	return g.POST("", h)
}

func TestNew(t *testing.T) {
	tests := []struct {
		echo        *echo.Echo
		docPath     string
		info        *Info
		expectPaths []string
		panic       bool
		name        string
	}{
		{
			echo:        echo.New(),
			docPath:     "doc/",
			info:        nil,
			expectPaths: []string{"/doc/", "/doc/swagger.json"},
			panic:       false,
			name:        "Normal",
		},
		{
			echo:    echo.New(),
			docPath: "doc",
			info: &Info{
				Title: "Test project",
				Contact: &Contact{
					URL: "https://github.com/pangpanglabs/echoswagger",
				},
			},
			expectPaths: []string{"/doc", "/doc/swagger.json"},
			panic:       false,
			name:        "Path slash suffix",
		},
		{
			echo:  nil,
			panic: true,
			name:  "Panic",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				assert.Panics(t, func() {
					New(tt.echo, tt.docPath, tt.info)
				})
			} else {
				apiRoot := New(tt.echo, tt.docPath, tt.info)
				assert.NotNil(t, apiRoot.(*Root))

				r := apiRoot.(*Root)
				assert.NotNil(t, r.spec)

				if tt.info == nil {
					assert.Equal(t, r.spec.Info.Title, "Project APIs")
				} else {
					assert.Equal(t, r.spec.Info, tt.info)
				}

				assert.NotNil(t, r.echo)
				assert.Len(t, r.echo.Routes(), 2)
				res := r.echo.Routes()
				paths := []string{res[0].Path, res[1].Path}
				assert.ElementsMatch(t, paths, tt.expectPaths)
			}
		})
	}
}

func TestPath(t *testing.T) {
	tests := []struct {
		docInput              string
		docOutput, specOutput string
		name                  string
	}{
		{
			docInput:   "doc/",
			docOutput:  "/doc/",
			specOutput: "/doc/swagger.json",
			name:       "A",
		}, {
			docInput:   "",
			docOutput:  "/",
			specOutput: "/swagger.json",
			name:       "B",
		}, {
			docInput:   "/doc",
			docOutput:  "/doc",
			specOutput: "/doc/swagger.json",
			name:       "C",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiRoot := New(echo.New(), tt.docInput, nil)
			r := apiRoot.(*Root)
			assert.NotNil(t, r.echo)
			assert.Len(t, r.echo.Routes(), 2)
			res := r.echo.Routes()
			paths := []string{res[0].Path, res[1].Path}
			assert.ElementsMatch(t, paths, []string{tt.docOutput, tt.specOutput})
		})
	}
}

func TestGroup(t *testing.T) {
	r := prepareApiRoot()
	t.Run("Normal", func(t *testing.T) {
		g := r.Group("Users", "users")
		assert.Equal(t, g.(*group).defs, r.(*Root).defs)
	})

	t.Run("Invalid name", func(t *testing.T) {
		assert.Panics(t, func() {
			r.Group("", "")
		})
	})

	t.Run("Repeat name", func(t *testing.T) {
		ga := r.Group("Users", "users")
		assert.Equal(t, ga.(*group).tag.Name, "Users")

		gb := r.Group("Users", "users")
		assert.Equal(t, gb.(*group).tag.Name, "Users")
	})
}

func TestBindGroup(t *testing.T) {
	r := prepareApiRoot()
	e := r.Echo()
	apiGroup := e.Group("/api")

	var h echo.HandlerFunc
	t.Run("Include", func(t *testing.T) {
		v1Group := apiGroup.Group("/v1")
		g := r.BindGroup("APIv1", v1Group)
		assert.Equal(t, g.(*group).tag.Name, "APIv1")

		g.GET("/in", h)
		assert.Len(t, g.(*group).apis, 1)
		assert.Equal(t, g.(*group).apis[0].route.Path, "/api/v1/in")
	})

	t.Run("Exclude", func(t *testing.T) {
		v2Group := apiGroup.Group("/v2")
		g := r.BindGroup("APIv2", v2Group)
		assert.Equal(t, g.(*group).tag.Name, "APIv2")

		v2Group.GET("/ex", h)
		assert.Len(t, g.(*group).apis, 0)
	})
}

func TestRouters(t *testing.T) {
	r := prepareApiRoot()
	var h echo.HandlerFunc
	r.GET("/:id", h)
	r.POST("/:id", h)
	r.PUT("/:id", h)
	r.DELETE("/:id", h)
	r.OPTIONS("/:id", h)
	r.HEAD("/:id", h)
	r.PATCH("/:id", h)
	assert.Len(t, r.(*Root).apis, 7)

	g := prepareApiGroup()
	g.GET("/:id", h)
	g.POST("/:id", h)
	g.PUT("/:id", h)
	g.DELETE("/:id", h)
	g.OPTIONS("/:id", h)
	g.HEAD("/:id", h)
	g.PATCH("/:id", h)
	assert.Len(t, g.(*group).apis, 7)
}

func TestAddParam(t *testing.T) {
	name := "name"
	desc := "Param desc"
	type nested struct {
		Name   string `json:"name" form:"name" query:"name"`
		Enable bool   `json:"-" form:"-" query:"-"`
	}

	t.Run("File", func(t *testing.T) {
		a := prepareApi()
		a.AddParamFile(name, desc, true)
		assert.Len(t, a.(*api).operation.Parameters, 1)
		assert.Equal(t, a.(*api).operation.Parameters[0].Name, name)
		assert.Equal(t, a.(*api).operation.Parameters[0].In, string(ParamInFormData))
		assert.Equal(t, a.(*api).operation.Parameters[0].Description, desc)
		assert.Equal(t, a.(*api).operation.Parameters[0].Required, true)
		assert.Equal(t, a.(*api).operation.Parameters[0].Type, "file")
	})

	t.Run("Path", func(t *testing.T) {
		a := prepareApi()
		a.AddParamPath(time.Now(), name, desc)
		assert.Len(t, a.(*api).operation.Parameters, 1)
		assert.Equal(t, a.(*api).operation.Parameters[0].Name, name)
		assert.Equal(t, a.(*api).operation.Parameters[0].In, string(ParamInPath))
		assert.Equal(t, a.(*api).operation.Parameters[0].Description, desc)
		assert.Equal(t, a.(*api).operation.Parameters[0].Required, true)
		assert.Equal(t, a.(*api).operation.Parameters[0].Type, "string")

		a.AddParamPathNested(&nested{})
		assert.Len(t, a.(*api).operation.Parameters, 2)
		assert.Equal(t, a.(*api).operation.Parameters[1].Name, "name_")
		assert.Equal(t, a.(*api).operation.Parameters[1].In, string(ParamInPath))
		assert.Equal(t, a.(*api).operation.Parameters[1].Type, "string")
	})

	t.Run("Query", func(t *testing.T) {
		a := prepareApi()
		a.AddParamQuery(time.Now(), name, desc, true)
		assert.Len(t, a.(*api).operation.Parameters, 1)
		assert.Equal(t, a.(*api).operation.Parameters[0].Name, name)
		assert.Equal(t, a.(*api).operation.Parameters[0].In, string(ParamInQuery))
		assert.Equal(t, a.(*api).operation.Parameters[0].Description, desc)
		assert.Equal(t, a.(*api).operation.Parameters[0].Required, true)
		assert.Equal(t, a.(*api).operation.Parameters[0].Type, "string")

		a.AddParamQueryNested(&nested{})
		assert.Len(t, a.(*api).operation.Parameters, 2)
		assert.Equal(t, a.(*api).operation.Parameters[1].Name, "name_")
		assert.Equal(t, a.(*api).operation.Parameters[1].In, string(ParamInQuery))
		assert.Equal(t, a.(*api).operation.Parameters[1].Type, "string")
	})

	t.Run("FormData", func(t *testing.T) {
		a := prepareApi()
		a.AddParamForm(time.Now(), name, desc, true)
		assert.Len(t, a.(*api).operation.Parameters, 1)
		assert.Equal(t, a.(*api).operation.Parameters[0].Name, name)
		assert.Equal(t, a.(*api).operation.Parameters[0].In, string(ParamInFormData))
		assert.Equal(t, a.(*api).operation.Parameters[0].Description, desc)
		assert.Equal(t, a.(*api).operation.Parameters[0].Required, true)
		assert.Equal(t, a.(*api).operation.Parameters[0].Type, "string")

		a.AddParamFormNested(&nested{})
		assert.Len(t, a.(*api).operation.Parameters, 2)
		assert.Equal(t, a.(*api).operation.Parameters[1].Name, "name_")
		assert.Equal(t, a.(*api).operation.Parameters[1].In, string(ParamInFormData))
		assert.Equal(t, a.(*api).operation.Parameters[1].Type, "string")
	})

	t.Run("Header", func(t *testing.T) {
		a := prepareApi()
		a.AddParamHeader(time.Now(), name, desc, true)
		assert.Len(t, a.(*api).operation.Parameters, 1)
		assert.Equal(t, a.(*api).operation.Parameters[0].Name, name)
		assert.Equal(t, a.(*api).operation.Parameters[0].In, string(ParamInHeader))
		assert.Equal(t, a.(*api).operation.Parameters[0].Description, desc)
		assert.Equal(t, a.(*api).operation.Parameters[0].Required, true)
		assert.Equal(t, a.(*api).operation.Parameters[0].Type, "string")

		a.AddParamHeaderNested(&nested{})
		assert.Len(t, a.(*api).operation.Parameters, 2)
		assert.Equal(t, a.(*api).operation.Parameters[1].Name, "name_")
		assert.Equal(t, a.(*api).operation.Parameters[1].In, string(ParamInHeader))
		assert.Equal(t, a.(*api).operation.Parameters[1].Type, "string")
	})
}

func TestAddSchema(t *testing.T) {
	type body struct {
		Name   string `json:"name"`
		Enable bool   `json:"-"`
	}

	t.Run("Multiple", func(t *testing.T) {
		a := prepareApi()
		a.AddParamBody(&body{}, "body", "body desc", true)
		assert.Len(t, a.(*api).operation.Parameters, 1)
		assert.Equal(t, a.(*api).operation.Parameters[0].Name, "body")
		assert.Equal(t, a.(*api).operation.Parameters[0].In, string(ParamInBody))
		assert.Equal(t, a.(*api).operation.Parameters[0].Description, "body desc")
		assert.Equal(t, a.(*api).operation.Parameters[0].Required, true)

		assert.Panics(t, func() {
			a.AddParamBody(body{}, "body", "body desc", true)
		})
	})

	t.Run("GetKey", func(t *testing.T) {
		a := prepareApi()
		a.AddParamBody(&body{}, "body", "body desc", true)
		assert.Len(t, a.(*api).operation.Parameters, 1)
		assert.Equal(t, a.(*api).operation.Parameters[0].Name, "body")
		assert.Equal(t, a.(*api).operation.Parameters[0].In, string(ParamInBody))
		assert.Equal(t, a.(*api).operation.Parameters[0].Description, "body desc")
		assert.Equal(t, a.(*api).operation.Parameters[0].Required, true)

		a.AddResponse(http.StatusOK, "response desc", body{}, nil)
		ca := strconv.Itoa(http.StatusOK)
		assert.Len(t, a.(*api).operation.Responses, 1)
		assert.Equal(t, a.(*api).operation.Responses[ca].Description, "response desc")

		assert.NotNil(t, a.(*api).defs)
		da := a.(*api).defs
		assert.Len(t, (*da), 1)
		assert.NotNil(t, (*da)["body"])

		a.AddResponse(http.StatusBadRequest, "response desc", body{Name: "name"}, nil)
		cb := strconv.Itoa(http.StatusBadRequest)
		assert.Len(t, a.(*api).operation.Responses, 2)
		assert.Equal(t, a.(*api).operation.Responses[cb].Description, "response desc")

		assert.NotNil(t, a.(*api).defs)
		db := a.(*api).defs
		assert.Len(t, (*db), 2)
		assert.NotNil(t, (*db)["body"])
	})
}

func TestAddResponse(t *testing.T) {
	a := prepareApi()
	a.AddResponse(http.StatusOK, "successful", nil, nil)
	var f = func() {}
	assert.Panics(t, func() {
		a.AddResponse(http.StatusBadRequest, "bad request", f, nil)
	})
	assert.Panics(t, func() {
		a.AddResponse(http.StatusBadRequest, "bad request", nil, time.Now())
	})
}

func TestUI(t *testing.T) {
	t.Run("DefaultCDN", func(t *testing.T) {
		r := New(echo.New(), "doc/", nil)
		se := r.(*Root)
		req := httptest.NewRequest(echo.GET, "/doc/", nil)
		rec := httptest.NewRecorder()
		c := se.echo.NewContext(req, rec)
		h := se.docHandler("doc/")

		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), DefaultCDN)
		}
	})

	t.Run("SetUI", func(t *testing.T) {
		r := New(echo.New(), "doc/", nil)
		se := r.(*Root)
		req := httptest.NewRequest(echo.GET, "/doc/", nil)
		rec := httptest.NewRecorder()
		c := se.echo.NewContext(req, rec)
		h := se.docHandler("doc/")

		cdn := "https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.18.0"
		r.SetUI(UISetting{
			HideTop: true,
			CDN:     cdn,
		})

		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), cdn)
			assert.NotContains(t, rec.Body.String(), `var specStr = ""`)
			assert.Contains(t, rec.Body.String(), "#swagger-ui>.swagger-container>.topbar")
		}

		r.SetUI(UISetting{
			DetachSpec: true,
		})

		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), `var specStr = ""`)
			assert.Contains(t, rec.Body.String(), "#swagger-ui>.swagger-container>.topbar")
		}
	})
}

func TestScheme(t *testing.T) {
	r := prepareApiRoot()
	schemes := []string{"http", "https"}
	r.SetScheme(schemes...)
	assert.ElementsMatch(t, r.(*Root).spec.Schemes, schemes)

	assert.Panics(t, func() {
		r.SetScheme("grpc")
	})
}

func TestRaw(t *testing.T) {
	r := prepareApiRoot()
	s := r.GetRaw()
	assert.NotNil(t, s)
	assert.NotNil(t, s.Info)
	assert.Equal(t, s.Info.Version, "")

	s.Info.Version = "1.0"
	r.SetRaw(s)
	assert.Equal(t, s.Info.Version, "1.0")
}

func TestContentType(t *testing.T) {
	c := []string{"application/x-www-form-urlencoded", "multipart/form-data"}
	p := []string{"application/vnd.github.v3+json", "application/vnd.github.v3.raw+json", "application/vnd.github.v3.text+json"}

	t.Run("In Root", func(t *testing.T) {
		r := prepareApiRoot()
		r.SetRequestContentType(c...)
		r.SetResponseContentType(p...)
		assert.NotNil(t, r.(*Root))
		assert.NotNil(t, r.(*Root).spec)
		assert.Len(t, r.(*Root).spec.Consumes, 2)
		assert.ElementsMatch(t, r.(*Root).spec.Consumes, c)
		assert.Len(t, r.(*Root).spec.Produces, 3)
		assert.ElementsMatch(t, r.(*Root).spec.Produces, p)
	})

	t.Run("In Api", func(t *testing.T) {
		a := prepareApi()
		a.SetRequestContentType(c...)
		a.SetResponseContentType(p...)
		assert.NotNil(t, a.(*api))
		assert.Len(t, a.(*api).operation.Consumes, 2)
		assert.ElementsMatch(t, a.(*api).operation.Consumes, c)
		assert.Len(t, a.(*api).operation.Produces, 3)
		assert.ElementsMatch(t, a.(*api).operation.Produces, p)
	})
}

func TestOperationId(t *testing.T) {
	id := "TestOperation"

	a := prepareApi()
	a.SetOperationId(id)
	assert.Equal(t, a.(*api).operation.OperationID, id)
}

func TestDeprecated(t *testing.T) {
	a := prepareApi()
	a.SetDeprecated()
	assert.Equal(t, a.(*api).operation.Deprecated, true)
}

func TestDescription(t *testing.T) {
	d := "Test desc"

	g := prepareApiGroup()
	g.SetDescription(d)
	assert.Equal(t, g.(*group).tag.Description, d)

	a := prepareApi()
	a.SetDescription(d)
	assert.Equal(t, a.(*api).operation.Description, d)
}

func TestExternalDocs(t *testing.T) {
	e := ExternalDocs{
		Description: "Test desc",
		URL:         "http://127.0.0.1/",
	}

	r := prepareApiRoot()
	r.SetExternalDocs(e.Description, e.URL)
	assert.Equal(t, r.(*Root).spec.ExternalDocs, &e)

	g := prepareApiGroup()
	g.SetExternalDocs(e.Description, e.URL)
	assert.Equal(t, g.(*group).tag.ExternalDocs, &e)

	a := prepareApi()
	a.SetExternalDocs(e.Description, e.URL)
	assert.Equal(t, a.(*api).operation.ExternalDocs, &e)
}

func TestSummary(t *testing.T) {
	s := "Test summary"

	a := prepareApi()
	a.SetSummary(s)
	assert.Equal(t, a.(*api).operation.Summary, s)
}

func TestEcho(t *testing.T) {
	r := prepareApiRoot()
	assert.NotNil(t, r.Echo())

	g := prepareApiGroup()
	assert.NotNil(t, g.EchoGroup())

	a := prepareApi()
	assert.NotNil(t, a.Route())
}

func TestHandlers(t *testing.T) {
	t.Run("ErrorGroupSecurity", func(t *testing.T) {
		r := prepareApiRoot()
		e := r.(*Root).echo
		var h echo.HandlerFunc
		r.Group("G", "/g").SetSecurity("JWT").GET("/", h)
		req := httptest.NewRequest(echo.GET, "/doc/swagger.json", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, r.(*Root).specHandler("/doc")(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
		if assert.NoError(t, r.(*Root).docHandler("/doc")(c)) {
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
		if assert.NoError(t, r.(*Root).specHandler("/doc")(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
		if assert.NoError(t, r.(*Root).docHandler("/doc")(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}
