English | [简体中文](./README_zh-CN.md)

# Echoswagger
Swagger UI generator for Echo framework

[![Go Report Card](https://goreportcard.com/badge/github.com/pangpanglabs/echoswagger)](https://goreportcard.com/report/github.com/pangpanglabs/echoswagger)
[![Build Status](https://travis-ci.org/pangpanglabs/echoswagger.svg?branch=master)](https://travis-ci.org/pangpanglabs/echoswagger)
[![codecov](https://codecov.io/gh/pangpanglabs/echoswagger/branch/master/graph/badge.svg)](https://codecov.io/gh/pangpanglabs/echoswagger)

## Feature
- No SwaggerUI HTML/CSS dependency
- Highly integrated with Echo, low intrusive design
- Take advantage of the strong typing language and chain programming to make it easy to use
- Recycle garbage in time, low memory usage

## Installation
```
go get github.com/pangpanglabs/echoswagger
```

## Example
```go
package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

func main() {
	// ApiRoot with Echo instance
	r := echoswagger.New(echo.New(), "", "doc/", nil)

	// Routes with parameters & responses
	r.POST("/", createUser).
		AddParamBody(User{}, "body", "User input struct", true).
		AddResponse(http.StatusCreated, "successful", nil, nil)

	// Start server
	r.Echo().Logger.Fatal(r.Echo().Start(":1323"))
}

type User struct {
	Name string
}

// Handler
func createUser(c echo.Context) error {
	return c.JSON(http.StatusCreated, nil)
}

```

## Usage
#### Create a `ApiRoot` with `New()`, which is a wrapper of `echo.New()`
```
r := echoswagger.New(echo.New(), "/v1", "doc/", nil)
```
> Note: The parameter `basePath` is generally used when the access root path is not the root directory of the website after application is deployed. For example, the URL of an API in the program running locally is: `http://localhost:1323/users`, the actual URL after deployed to server is: `https://www.xxx.com/legacy-api/users`, then, when running locally, `basePath` should be `/`, when running on server, `basePath` should be `/legacy-api`.

You can use the result `ApiRoot` instance to:
- Setup Security definitions, request/response Content-Types, UI options, Scheme, etc.
```
r.AddSecurityAPIKey("JWT", "JWT Token", echoswagger.SecurityInHeader).
	SetRequestContentType("application/x-www-form-urlencoded", "multipart/form-data").
	SetUI(UISetting{HideTop: true}).
	SetScheme("https", "http")
```
- Get `echo.Echo` instance.
```
r.Echo()
```
- Registers a new GET, POST, PUT, DELETE, OPTIONS, HEAD or PATCH route in default group, these are wrappers of Echo's create route methods.
It returns a new `Api` instance.
```
r.GET("/:id", handler)
```
- And: ↓

#### Create a `ApiGroup` with `Group()`, which is a wrapper of `echo.Group()`
```
g := r.Group("Users", "/users")
```
You can use the result `ApiGroup` instance to:
- Set description, etc.
```
g.SetDescription("The desc of group")
```
- Set security for all routes in this group.
```
g.SetSecurity("JWT")
```
- Get `echo.Group` instance.
```
g.EchoGroup()
```
- And: ↓

#### Registers a new route in `ApiGroup`
GET, POST, PUT, DELETE, OPTIONS, HEAD or PATCH methods are supported by Echoswagger, these are wrappers of Echo's create route methods.
```
a := g.GET("/:id", handler)
```
You can use the result `Api` instance to:
- Add parameter with these methods:
```
AddParamPath(p interface{}, name, desc string)

AddParamPathNested(p interface{})

AddParamQuery(p interface{}, name, desc string, required bool)

AddParamQueryNested(p interface{})

AddParamForm(p interface{}, name, desc string, required bool)

AddParamFormNested(p interface{})

AddParamHeader(p interface{}, name, desc string, required bool)

AddParamHeaderNested(p interface{})

AddParamBody(p interface{}, name, desc string, required bool)

AddParamFile(name, desc string, required bool)
```

The methods which name's suffix are `Nested` means these methods treat parameter `p` 's fields as paramters, so it must be a struct type.

e.g.
```
type SearchInput struct {
	Q         string `query:"q" swagger:"desc("Keywords"),required"`
	SkipCount int    `query:"skipCount"`
}
a.AddParamQueryNested(SearchInput{})
```
Is equivalent to:
```
a.AddParamQuery("", "q", "", true).
	AddParamQuery(0, "skipCount", "", false)
```
- Add responses.
```
a.AddResponse(http.StatusOK, "response desc", body{}, nil)
```
- Set Security, request/response Content-Types, summary, description, etc.
```
a.SetSecurity("JWT").
	SetResponseContentType("application/xml").
	SetSummary("The summary of API").
	SetDescription("The desc of API")
```
- Get `echo.Route` instance.
```
a.Route()
```

#### With swagger tag, you can set more info with `AddParam...` methods.
e.g.
```
type User struct {
	Age    int       `swagger:"min(0),max(99)"`
	Gender string    `swagger:"enum(male|female|other),required"`
	Money  []float64 `swagger:"default(0),readOnly"`
}
a.AddParamBody(&User{}, "Body", "", true)
```
The definition is equivalent to:
```
{
    "definitions": {
        "User": {
            "type": "object",
            "properties": {
                "Age": {
                    "type": "integer",
                    "format": "int32",
                    "minimum": 0,
                    "maximum": 99
                },
                "Gender": {
                    "type": "string",
                    "enum": [
                        "male",
                        "female",
                        "other"
                    ],
                    "format": "string"
                },
                "Money": {
                    "type": "array",
                    "items": {
                        "type": "number",
                        "default": 0,
                        "format": "double"
                    },
                    "readOnly": true
                }
            },
            "required": [
                "Gender"
            ]
        }
    }
}
```

**Supported swagger tags:**

Tag | Type | Description
---|:---:|---
collect | `string` | Determines the format of the array if type array is used. Possible values are: <ul><li>`csv` - comma separated values `foo,bar`. <li>`ssv` - space separated values `foo bar`. <li>`tsv` - tab separated values `foo\tbar`. <li>`pipes` - pipe separated values `foo\|bar`. </ul> Default value is `csv`.
desc | `string` | Description.
maximum | `number` | -
minimum | `number` | -
maxLength | `integer` | -
minLength | `integer` | -
allowEmpty | `boolean` | Sets the ability to pass empty-valued parameters. This is valid only for either `query` or `formData` parameters and allows you to send a parameter with a name only or  an empty value. Default value is `false`.
required | `boolean` | Determines whether this parameter is mandatory. If the parameter is `in` "path", this property is `true` without setting. Otherwise, the property MAY be included and its default value is `false`.
readOnly | `boolean` | Relevant only for Schema `"properties"` definitions. Declares the property as "read only". This means that it MAY be sent as part of a response but MUST NOT be sent as part of the request. Properties marked as `readOnly` being `true` SHOULD NOT be in the `required` list of the defined schema. Default value is `false`.
enum | [*] | Enumerate value, multiple values should be separated by "\|"
default | * | Default value, which type is same as the field's type.

## Reference
[OpenAPI Specification 2.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md)

## License

[MIT](https://github.com/pangpanglabs/echoswagger/blob/master/LICENSE)
