package main

import (
	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

func main() {
	e := initServer().Echo()
	e.Logger.Fatal(e.Start(":1323"))
}

func initServer() echoswagger.ApiRoot {
	e := echo.New()

	se := echoswagger.New(e, "/v2", "doc/", &echoswagger.Info{
		Title:          "Swagger Petstore",
		Description:    "This is a sample server Petstore server.  You can find out more about     Swagger at [http://swagger.io](http://swagger.io) or on [irc.freenode.net, #swagger](http://swagger.io/irc/).      For this sample, you can use the api key `special-key` to test the authorization     filters.",
		Version:        "1.0.0",
		TermsOfService: "http://swagger.io/terms/",
		Contact: &echoswagger.Contact{
			Email: "apiteam@swagger.io",
		},
		License: &echoswagger.License{
			Name: "Apache 2.0",
			URL:  "http://www.apache.org/licenses/LICENSE-2.0.html",
		},
	})

	se.AddSecurityOAuth2("petstore_auth", "", echoswagger.OAuth2FlowImplicit,
		"http://petstore.swagger.io/oauth/dialog", "", map[string]string{
			"write:pets": "modify pets in your account",
			"read:pets":  "read your pets",
		},
	).AddSecurityAPIKey("api_key", "", echoswagger.SecurityInHeader)

	se.SetExternalDocs("Find out more about Swagger", "http://swagger.io").
		SetResponseContentType("application/xml", "application/json").
		SetUI(echoswagger.UISetting{HideTop: true})

	PetController{}.Init(se.Group("pet", "/pet"))
	StoreController{}.Init(se.Group("store", "/store"))
	UserController{}.Init(se.Group("user", "/user"))

	return se
}
