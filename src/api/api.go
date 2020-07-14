package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"

	_ "github.com/andrei-dascalu/go-workshop-shopapi/docs"
)

//RunAPIWithHandlers start  API
func RunAPIWithHandlers() {
	//default echo router
	r := echo.New()

	r.Use(middleware.Logger())
	//r.Use(middleware.Recover())

	//get products
	r.GET("/products", GetProducts)
	//get promos
	r.GET("/promos", GetPromos)

	r.GET("/", GetMainPage)

	r.GET("/swagger/*", echoSwagger.WrapHandler)

	userGroup := r.Group("/user")
	{
		userGroup.POST("/:id/signout", SignOut)
		userGroup.GET("/:id/orders", GetOrders)
	}

	usersGroup := r.Group("/users")
	{
		usersGroup.POST("/charge", Charge)
		usersGroup.POST("/signin", SignIn)
		usersGroup.POST("", AddUser)
	}

	r.Logger.Fatal(r.Start(":8080"))
}
