package main

import (
	"github.com/andrei-dascalu/go-workshop-shopapi/src/api"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/configuration"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/dblayer"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/models"
	"github.com/rs/zerolog/log"
)

func main() {
	configuration.InitConfiguration()
	err := dblayer.InitORM()

	if err != nil {
		log.Error().Err(err).Msg("Error")
	}
	dblayer.ShopDB.DB.AutoMigrate(&models.Order{}, &models.Payment{}, &models.Address{}, &models.Cart{}, &models.Product{}, &models.Customer{})
	api.RunAPIWithHandlers()
}
