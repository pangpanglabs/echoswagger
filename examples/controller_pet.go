package main

import (
	"net/http"

	"github.com/elvinchan/echoswagger"
	"github.com/labstack/echo"
)

type PetController struct{}

func (c PetController) Init(g echoswagger.ApiGroup) {
	g.SetDescription("Everything about your Pets").
		SetExternalDocs("Find out more", "http://swagger.io")

	security := map[string][]string{
		"petstore_auth": []string{"write:pets", "read:pets"},
	}

	type Category struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	type Tag struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	type Pet struct {
		Id        int64    `json:"id"`
		Category  Category `json:"category"`
		Name      string   `json:"name" swagger:"required"`
		PhotoUrls []string `json:"photoUrls" xml:"photoUrl" swagger:"required"`
		Tags      []Tag    `json:"tags" xml:"tag"`
		Status    string   `json:"status" swagger:"enum(available|pending|sold),desc(pet status in the store)"`
	}
	pet := Pet{Name: "doggie"}
	g.POST("", c.Create).
		AddParamBody(&pet, "body", "Pet object that needs to be added to the store", true).
		AddResponse(http.StatusMethodNotAllowed, "Invalid input", nil, nil).
		SetRequestContentType("application/json", "application/xml").
		SetOperationId("addPet").
		SetSummary("Add a new pet to the store").
		SetSecurityWithScope(security)

	g.PUT("", c.Update).
		AddParamBody(&pet, "body", "Pet object that needs to be added to the store", true).
		AddResponse(http.StatusBadRequest, "Invalid ID supplied", nil, nil).
		AddResponse(http.StatusNotFound, "Pet not found", nil, nil).
		AddResponse(http.StatusMethodNotAllowed, "Validation exception", nil, nil).
		SetRequestContentType("application/json", "application/xml").
		SetOperationId("updatePet").
		SetSummary("Update an existing pet").
		SetSecurityWithScope(security)

	type StatusParam struct {
		Status []string `query:"status" swagger:"required,desc(Status values that need to be considered for filter),default(available),enum(available|pending|sold)"`
	}
	g.GET("/findByStatus", c.FindByStatus).
		AddParamQueryNested(&StatusParam{}).
		AddResponse(http.StatusOK, "successful operation", &[]Pet{pet}, nil).
		AddResponse(http.StatusBadRequest, "Invalid status value", nil, nil).
		SetOperationId("findPetsByStatus").
		SetDescription("Multiple status values can be provided with comma separated strings").
		SetSummary("Finds Pets by status").
		SetSecurityWithScope(security)

	g.GET("/findByTags", c.FindByTags).
		AddParamQuery([]string{}, "tags", "Tags to filter by", true).
		AddResponse(http.StatusOK, "successful operation", &[]Pet{pet}, nil).
		AddResponse(http.StatusBadRequest, "Invalid tag value", nil, nil).
		SetOperationId("findPetsByTags").
		SetDeprecated().
		SetDescription("Muliple tags can be provided with comma separated strings. Use         tag1, tag2, tag3 for testing.").
		SetSummary("Finds Pets by tags").
		SetSecurityWithScope(security)

	g.GET("/{petId}", c.GetById).
		AddParamPath(0, "petId", "ID of pet to return").
		AddResponse(http.StatusOK, "successful operation", &pet, nil).
		AddResponse(http.StatusBadRequest, "Invalid ID supplied", nil, nil).
		AddResponse(http.StatusNotFound, "Pet not found", nil, nil).
		SetOperationId("getPetById").
		SetDescription("Returns a single pet").
		SetSummary("Find pet by ID").
		SetSecurity("api_key")

	g.POST("/{petId}", c.CreateById).
		AddParamPath(0, "petId", "ID of pet that needs to be updated").
		AddParamForm("", "name", "Updated name of the pet", false).
		AddParamForm("", "status", "Updated status of the pet", false).
		AddResponse(http.StatusMethodNotAllowed, "Invalid input", nil, nil).
		SetRequestContentType("application/x-www-form-urlencoded").
		SetOperationId("updatePetWithForm").
		SetSummary("Updates a pet in the store with form data").
		SetSecurityWithScope(security)

	g.DELETE("/{petId}", c.DeleteById).
		AddParamHeader("", "api_key", "", false).
		AddParamPath(int64(0), "petId", "Pet id to delete").
		AddResponse(http.StatusBadRequest, "Invalid ID supplied", nil, nil).
		AddResponse(http.StatusNotFound, "Pet not found", nil, nil).
		SetOperationId("deletePet").
		SetSummary("Deletes a pet").
		SetSecurityWithScope(security)

	type ApiResponse struct {
		Code    int32  `json:"code"`
		Type    string `json:"type"`
		Message string `json:"message"`
	}
	g.POST("/{petId}/uploadImage", c.UploadImageById).
		AddParamPath("", "petId", "ID of pet to update").
		AddParamForm("", "additionalMetadata", "Additional data to pass to server", false).
		AddParamFile("file", "file to upload", false).
		AddResponse(http.StatusOK, "successful operation", &ApiResponse{}, nil).
		SetRequestContentType("multipart/form-data").
		SetResponseContentType("application/json").
		SetOperationId("uploadFile").
		SetSummary("uploads an image").
		SetSecurityWithScope(security)
}

func (PetController) Create(c echo.Context) error {
	return nil
}

func (PetController) Update(c echo.Context) error {
	return nil
}

func (PetController) FindByStatus(c echo.Context) error {
	return nil
}

func (PetController) FindByTags(c echo.Context) error {
	return nil
}

func (PetController) GetById(c echo.Context) error {
	return nil
}

func (PetController) CreateById(c echo.Context) error {
	return nil
}

func (PetController) DeleteById(c echo.Context) error {
	return nil
}

func (PetController) UploadImageById(c echo.Context) error {
	return nil
}
