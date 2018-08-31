package echoswagger

import (
	"reflect"

	"github.com/labstack/echo"
)

type RawDefine struct {
	Value  reflect.Value
	Schema *JSONSchema
}

type Root struct {
	apis   []api
	spec   *Swagger
	echo   *echo.Echo
	groups []group
}

type group struct {
	apis      []api
	echoGroup *echo.Group
	security  []map[string][]string
	tag       Tag
}

type api struct {
	route     *echo.Route
	security  []map[string][]string
	method    string
	operation Operation
}

func New(e *echo.Echo, basePath, docPath string, i *Info) *Root {
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

	r := &Root{
		echo: e,
		spec: &Swagger{
			Info:                i,
			SecurityDefinitions: make(map[string]*SecurityDefinition),
			BasePath:            basePath,
			Definitions:         make(map[string]*JSONSchema),
		},
	}

	e.GET(docPath, r.docHandler(specPath))
	e.GET(specPath, r.Spec)
	return r
}
