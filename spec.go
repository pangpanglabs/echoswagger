package echoswagger

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo"
)

const SwaggerVersion = "2.0"

func (r *Root) Spec(c echo.Context) error {
	err := r.genSpec(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
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
			fmt.Println("Group:", a)
			// TODO
		}
	}

	return nil
}

func (r *Root) docHandler(swaggerPath string) echo.HandlerFunc {
	t, err := template.New("swagger").Parse(SwaggerUIContent)
	if err != nil {
		panic(err)
	}

	return func(c echo.Context) error {
		buf := new(bytes.Buffer)
		t.Execute(buf, map[string]interface{}{
			"title": r.spec.Info.Title,
			"url":   c.Scheme() + "://" + c.Request().Host + swaggerPath,
		})
		return c.HTMLBlob(http.StatusOK, buf.Bytes())
	}
}
