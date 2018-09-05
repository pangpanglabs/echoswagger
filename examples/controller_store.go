package main

import (
	"net/http"
	"time"

	"github.com/elvinchan/echoswagger"
	"github.com/labstack/echo"
)

type StoreController struct{}

func (c StoreController) Init(g echoswagger.ApiGroup) {
	g.SetDescription("Access to Petstore orders")

	g.GET("/inventory", c.GetInventory).
		AddResponse(http.StatusOK, "successful operation", map[string]int32{}, nil).
		SetResponseContentType("application/json").
		SetOperationId("getInventory").
		SetDescription("Returns a map of status codes to quantities").
		SetSummary("Returns pet inventories by status").
		SetSecurity("api_key")

	type Order struct {
		Id       int64     `json:"id"`
		PetId    int64     `json:"petId"`
		Quantity int64     `json:"quantity"`
		ShipDate time.Time `json:"shipDate"`
		Status   string    `json:"status" swagger:"desc(Order Status),enum(placed|approved|delivered)"`
		Complete bool      `json:"complete" swagger:"default(false)"`
	}
	g.POST("/order", c.CreateOrder).
		AddParamBody(Order{}, "body", "order placed for purchasing the pet", true).
		AddResponse(http.StatusOK, "successful operation", Order{}, nil).
		AddResponse(http.StatusBadRequest, "Invalid Order", nil, nil).
		SetOperationId("placeOrder").
		SetSummary("Place an order for a pet")

	type GetOrderId struct {
		orderId int64 `swagger:"max(10.0),min(1.0),desc(ID of pet that needs to be fetched)"`
	}
	g.GET("/order/{orderId}", c.GetOrderById).
		AddParamPathNested(&GetOrderId{}).
		AddResponse(http.StatusOK, "successful operation", Order{}, nil).
		AddResponse(http.StatusBadRequest, "Invalid ID supplied", nil, nil).
		AddResponse(http.StatusNotFound, "Order not found", nil, nil).
		SetOperationId("getOrderById").
		SetDescription("For valid response try integer IDs with value >= 1 and <= 10.         Other values will generated exceptions").
		SetSummary("Find purchase order by ID")

	type DeleteOrderId struct {
		orderId int64 `swagger:"min(1.0),desc(ID of the order that needs to be deleted)"`
	}
	g.DELETE("/order/{orderId}", c.DeleteOrderById).
		AddParamPathNested(&DeleteOrderId{}).
		AddResponse(http.StatusBadRequest, "Invalid ID supplied", nil, nil).
		AddResponse(http.StatusNotFound, "Order not found", nil, nil).
		SetOperationId("deleteOrder").
		SetDescription("For valid response try integer IDs with positive integer value.         Negative or non-integer values will generate API errors").
		SetSummary("Delete purchase order by ID")
}

func (StoreController) GetInventory(c echo.Context) error {
	return nil
}

func (StoreController) CreateOrder(c echo.Context) error {
	return nil
}

func (StoreController) GetOrderById(c echo.Context) error {
	return nil
}

func (StoreController) DeleteOrderById(c echo.Context) error {
	return nil
}
