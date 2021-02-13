[English](./README.md) | 简体中文

# Echoswagger
[Echo](https://github.com/labstack/echo) 框架的 [Swagger UI](https://github.com/swagger-api/swagger-ui) 生成器

[![Go Report Card](https://goreportcard.com/badge/github.com/lexholden/echoswagger)](https://goreportcard.com/report/github.com/lexholden/echoswagger)
[![Build Status](https://travis-ci.org/lexholden/echoswagger.svg?branch=master)](https://travis-ci.org/lexholden/echoswagger)
[![codecov](https://codecov.io/gh/lexholden/echoswagger/branch/master/graph/badge.svg)](https://codecov.io/gh/lexholden/echoswagger)

## 特性
- 不依赖任何SwaggerUI的HTML/CSS文件
- 与Echo高度整合，低侵入式设计
- 利用强类型语言和链式编程的优势，简单易用
- 及时的垃圾回收，低内存占用

## 安装
```
go get github.com/lexholden/echoswagger
```

## Go modules 支持
如果你的项目已经使用Go modules，你可以:
- 选择v2版本的Echoswagger搭配Echo v4版本
- 选择v1版本的Echoswagger搭配Echo v3及以下版本

使用v2版本，只需要:
- `go get github.com/lexholden/echoswagger/v2`
- 在你的项目中import `github.com/labstack/echo/v4` 和 `github.com/lexholden/echoswagger/v2`

同时，v1版本将继续更新。关于Go modules的详细内容，请参考 [Go Wiki](https://github.com/golang/go/wiki/Modules)

## 示例
```go
package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/lexholden/echoswagger"
)

func main() {
	// ApiRoot with Echo instance
	r := echoswagger.New(echo.New(), "/doc", nil)

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

## 用法
#### 用`New()`创建`ApiRoot`，此方法是对`echo.New()`方法的封装
```go
r := echoswagger.New(echo.New(), "/doc", nil)
```
你可以用这个`ApiRoot`来：
- 设置Security定义, 请求/响应Content-Type，UI选项，Scheme等。
```go
r.AddSecurityAPIKey("JWT", "JWT Token", echoswagger.SecurityInHeader).
	SetRequestContentType("application/x-www-form-urlencoded", "multipart/form-data").
	SetUI(UISetting{HideTop: true}).
	SetScheme("https", "http")
```
- 获取`echo.Echo`实例。
```go
r.Echo()
```
- 在默认组中注册一个GET、POST、PUT、DELETE、OPTIONS、HEAD或PATCH路由，这些是对Echo的注册路由方法的封装。
此方法返回一个`Api`实例。
```go
r.GET("/:id", handler)
```
- 以及： ↓

#### 用`Group()`创建`ApiGroup`，此方法是对`echo.Group()`方法的封装
```go
g := r.Group("Users", "/users")
```
你可以用这个`ApiGroup`来：
- 设置描述，等。
```go
g.SetDescription("The desc of group")
```
- 为此组中的所有路由设置Security。
```go
g.SetSecurity("JWT")
```
- 获取`echo.Group`实例。
```go
g.EchoGroup()
```
- 以及： ↓

#### 在`ApiGroup`中注册一个新的路由
Echoswagger支持GET、POST、PUT、DELETE、OPTIONS、HEAD或PATCH方法，这些是对Echo的注册路由方法的封装。
```go
a := g.GET("/:id", handler)
```
你可以使用此`Api`实例来：
- 使用以下方法添加参数：
```go
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

后缀带有`Nested`的方法把参数`p`的字段看做多个参数，所以它必须是结构体类型的。

例：
```go
type SearchInput struct {
	Q         string `query:"q" swagger:"desc(Keywords),required"`
	SkipCount int    `query:"skipCount"`
}
a.AddParamQueryNested(SearchInput{})
```
等价于：
```go
a.AddParamQuery("", "q", "Keywords", true).
	AddParamQuery(0, "skipCount", "", false)
```
- 添加响应。
```go
a.AddResponse(http.StatusOK, "response desc", body{}, nil)
```
- 设置Security，请求/响应的Content-Type，概要，描述，等。
```go
a.SetSecurity("JWT").
	SetResponseContentType("application/xml").
	SetSummary("The summary of API").
	SetDescription("The desc of API")
```
- 获取`echo.Route`实例。
```go
a.Route()
```

#### 使用`swagger`标签，你可以在`AddParam...`方法中设置更多信息
例：
```go
type User struct {
	Age    int       `swagger:"min(0),max(99)"`
	Gender string    `swagger:"enum(male|female|other),required"`
	Money  []float64 `swagger:"default(0),readOnly"`
}
a.AddParamBody(&User{}, "Body", "", true)
```
此定义等价于:
```json
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

**支持的`swagger`标签：**

Tag | Type | Description
---|:---:|---
desc | `string` | 描述。
min | `number` | -
max | `number` | -
minLen | `integer` | -
maxLen | `integer` | -
allowEmpty | `boolean` | 设置传递空值参数的功能。 这仅对`query`或`formData`参数有效，并允许你发送仅具有名称或空值的参数。默认值为“false”。
required | `boolean` | 确定此参数是否必需。如果参数是`in`“path”，则此属性默认为“true”。否则，可以设置此属性，其默认值为“false”。
readOnly | `boolean` | 仅与Schema`"properties"`定义相关。将属性声明为“只读”。这意味着它可以作为响应的一部分发送，但绝不能作为请求的一部分发送。标记为“readOnly”的属性为“true”，不应位于已定义模式的“required”列表中。默认值为“false”。
enum | [*] | 枚举值，多个值应以“\|”分隔。
default | * | 默认值，该类型与字段的类型相同。

## 参考
[OpenAPI Specification 2.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md)

## License

[MIT](https://github.com/lexholden/echoswagger/blob/master/LICENSE)
