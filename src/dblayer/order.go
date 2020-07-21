package dblayer

import (
	"time"

	"github.com/andrei-dascalu/go-workshop-shopapi/src/models"
)

//AddOrder save order
func (db *DBORM) AddOrder(order models.Order) (models.Order, error) {
	return order, db.Create(&order).Error
}

//FindOrderByID find order
func (db *DBORM) FindOrderByID(id int) (order models.Order, err error) {
	return order, db.Table("orders").Where("id=?", id).Error
}

//CreateOrder crate order
func (db *DBORM) CreateOrder(c models.Cart, a models.Address) (models.Order, error) {
	order := models.Order{
		Cart:              c,
		CartID:            int(c.ID),
		DeliveryAddress:   a,
		DeliveryAddressID: int(a.ID),
		PurchaseDate:      time.Now(),
	}

	item, err := db.AddOrder(order)

	if err != nil {
		return models.Order{}, err
	}

	return item, nil
}
