# echoswagger
Swagger UI generator for Echo framework

[![Go Report Card](https://goreportcard.com/badge/github.com/elvinchan/echoswagger)](https://goreportcard.com/report/github.com/elvinchan/echoswagger)
[![Build Status](https://travis-ci.org/elvinchan/echoswagger.svg?branch=master)](https://travis-ci.org/elvinchan/echoswagger)
[![codecov](https://codecov.io/gh/elvinchan/echoswagger/branch/master/graph/badge.svg)](https://codecov.io/gh/elvinchan/echoswagger)

## Feature
- No SwaggerUI html/css dependency
- Highly integrated with echo, low intrusive design
- Take advantage of the strong typing language to make it easy to use
- Recycle garbage in time, low memory usage

## Example
```go
package main

import (
	"net/http"

	"github.com/elvinchan/echoswagger"
	"github.com/labstack/echo"
)

func main() {
	// ApiRoot with Echo instance
	e := echoswagger.New(echo.New(), "/v2", "doc/", nil)

	// Routes with parameters & responses
	e.POST("/", createUser).
		AddParamBody(User{}, "body", "User input struct", true).
		AddResponse(http.StatusCreated, "successful", nil, nil)

	// Start server
	e.Echo().Logger.Fatal(e.Echo().Start(":1323"))
}

type User struct {
	Name string
}

// Handler
func createUser(c echo.Context) error {
	return c.JSON(http.StatusCreated, nil)
}

```
