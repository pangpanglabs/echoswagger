package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	e := initServer()
	req := httptest.NewRequest(echo.GET, "/doc/swagger.json", nil)
	rec := httptest.NewRecorder()
	c := e.Echo().NewContext(req, rec)
	b, err := ioutil.ReadFile("./swagger.json")
	assert.Nil(t, err)
	if assert.NoError(t, e.(*echoswagger.Root).Spec(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, string(b), rec.Body.String())
	}
}
