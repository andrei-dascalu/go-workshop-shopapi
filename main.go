package main

import (
	"github.com/andrei-dascalu/go-workshop-shopapi/src/api"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/configuration"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/dblayer"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/models"
	"github.com/labstack/gommon/log"
)

// @title Webshop API Workshop
// @contact.email andrei.dascalu@qubiz.com
// @contact.name Andrei
// @host localhost:8181
// @BasePath /

func main() {
	configuration.InitConfiguration()
	err := dblayer.InitORM()

	if err != nil {
		log.Errorf("Err: %v", err)
	}
	dblayer.ShopDB.DB.AutoMigrate(&models.Order{}, &models.Payment{}, &models.Address{}, &models.Cart{}, &models.Product{}, &models.Customer{})
	api.RunAPIWithHandlers()
}
