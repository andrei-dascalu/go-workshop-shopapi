package api

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"

	//docs folder
	_ "github.com/andrei-dascalu/go-workshop-shopapi/docs"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/security"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

//RunAPIWithHandlers start  API
func RunAPIWithHandlers() {
	//default echo router
	router := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	router.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	//get products
	router.Get("/products", GetProducts)
	//get promos
	router.Get("/promos", GetPromos)

	router.Get("/", GetMainPage)

	usersGroup := router.Group("/users")
	{
		usersGroup.Post("/charge", Charge)
		usersGroup.Post("/signin", SignIn)
		usersGroup.Post("", AddUser)
	}

	userGroup := router.Group("/user")
	{
		userGroup.Use(security.CustomJWTMiddleware)
		userGroup.Post("/:id/signout", SignOut)
		userGroup.Post("/:id/addProduct", AddProductToCart)
		userGroup.Post("/:id/orders", GetOrders)
		userGroup.Post("/:id/createOrder", CreateOrder)
	}

	router.Listen(":8080")
}
