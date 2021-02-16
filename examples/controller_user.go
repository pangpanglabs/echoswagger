package main

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/lexholden/echoswagger"
)

type UserController struct{}

func (c UserController) Init(g echoswagger.ApiGroup) {
	g.SetDescription("Operations about user").
		SetExternalDocs("Find out more about our store", "http://swagger.io")

	type User struct {
		X          xml.Name `xml:"Users"`
		Id         int64    `json:"id"`
		UserName   string   `json:"username"`
		FirstName  string   `json:"firstname"`
		LastName   string   `json:"lastname"`
		Email      string   `json:"email"`
		Password   string   `json:"password"`
		Phone      string   `json:"phone"`
		UserStatus int32    `json:"userStatus" swagger:"desc(User Status)"`
	}
	g.POST("", c.Create).
		AddParamBody(&User{}, "body", "Created user object", true).
		SetOperationId("createUser").
		SetDescription("This can only be done by the logged in user.").
		SetSummary("Create user")

	g.POST("/createWithArray", c.CreateWithArray).
		AddParamBody(&[]User{}, "body", "List of user object", true).
		SetOperationId("createUsersWithArrayInput").
		SetRequestContentType("application/json", "application/xml").
		SetSummary("Creates list of users with given input array")

	g.POST("/createWithList", c.CreateWithList).
		AddParamBody(&[]User{}, "body", "List of user object", true).
		SetOperationId("createUsersWithListInput").
		SetSummary("Creates list of users with given input array")

	type ResponseHeader struct {
		XRateLimit    int32     `json:"X-Rate-Limit" swagger:"desc(calls per hour allowed by the user)"`
		XExpiresAfter time.Time `json:"X-Expires-After" swagger:"desc(date in UTC when token expires)"`
	}
	g.GET("/login", c.Login).
		AddParamQuery("", "username", "The user name for login", true).
		AddParamQuery("", "password", "The password for login in clear text", true).
		AddResponse(http.StatusOK, "successful operation", "", ResponseHeader{}).
		AddResponse(http.StatusBadRequest, "Invalid username/password supplied", nil, nil).
		SetOperationId("loginUser").
		SetSummary("Logs user into the system")

	g.GET("/logout", c.Logout).
		SetOperationId("logoutUser").
		SetSummary("Logs out current logged in user session")

	g.GET("/{username}", c.GetByUsername).
		AddParamPath("", "username", "The name that needs to be fetched. Use user1 for testing. ").
		AddResponse(http.StatusOK, "successful operation", &User{}, nil).
		AddResponse(http.StatusBadRequest, "Invalid username supplied", nil, nil).
		AddResponse(http.StatusNotFound, "User not found", nil, nil).
		SetOperationId("getUserByName").
		SetSummary("Get user by user name")

	g.PUT("/{username}", c.UpdateByUsername).
		AddParamPath("", "username", "name that need to be updated").
		AddParamBody(&User{}, "body", "Updated user object", true).
		AddResponse(http.StatusBadRequest, "Invalid user supplied", nil, nil).
		AddResponse(http.StatusNotFound, "User not found", nil, nil).
		SetOperationId("updateUser").
		SetDescription("This can only be done by the logged in user.").
		SetSummary("Updated user")

	g.DELETE("/{username}", c.DeleteByUsername).
		AddParamPath("", "username", "The name that needs to be deleted").
		AddResponse(http.StatusBadRequest, "Invalid username supplied", nil, nil).
		AddResponse(http.StatusNotFound, "User not found", nil, nil).
		SetOperationId("deleteUser").
		SetDescription("This can only be done by the logged in user.").
		SetSummary("Delete user")
}

func (UserController) Create(c echo.Context) error {
	return nil
}

func (UserController) CreateWithArray(c echo.Context) error {
	return nil
}

func (UserController) CreateWithList(c echo.Context) error {
	return nil
}

func (UserController) Login(c echo.Context) error {
	return nil
}

func (UserController) Logout(c echo.Context) error {
	return nil
}

func (UserController) GetByUsername(c echo.Context) error {
	return nil
}

func (UserController) UpdateByUsername(c echo.Context) error {
	return nil
}

func (UserController) DeleteByUsername(c echo.Context) error {
	return nil
}
