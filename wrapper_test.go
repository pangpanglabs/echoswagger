package echoswagger

import (
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		echo        *echo.Echo
		basePath    string
		docPath     string
		info        *Info
		expectPaths []string
		panic       bool
		name        string
	}{
		{
			echo:        echo.New(),
			basePath:    "/",
			docPath:     "doc/",
			info:        nil,
			expectPaths: []string{"/doc/", "/doc/swagger.json"},
			panic:       false,
			name:        "Normal",
		},
		{
			echo:     echo.New(),
			basePath: "/",
			docPath:  "doc",
			info: &Info{
				Title: "Test project",
				Contact: &Contact{
					URL: "https://github.com/elvinchan/echoswagger",
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
					New(tt.echo, tt.basePath, tt.docPath, tt.info)
				})
			} else {
				apiRoot := New(tt.echo, tt.basePath, tt.docPath, tt.info)
				assert.NotNil(t, apiRoot.(*Root))

				root := apiRoot.(*Root)
				assert.NotNil(t, root.spec)

				assert.Equal(t, root.spec.BasePath, tt.basePath)
				if tt.info == nil {
					assert.Equal(t, root.spec.Info.Title, "Project APIs")
				} else {
					assert.Equal(t, root.spec.Info, tt.info)
				}

				assert.NotNil(t, root.echo)
				assert.Len(t, root.echo.Routes(), 2)
				res := root.echo.Routes()
				paths := []string{res[0].Path, res[1].Path}
				assert.ElementsMatch(t, paths, tt.expectPaths)
			}
		})
	}
}
