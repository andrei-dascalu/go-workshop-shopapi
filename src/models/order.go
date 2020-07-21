package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Order order data
type Order struct {
	gorm.Model
	Cart              Cart
	DeliveryAddress   Address
	CartID            int         `json:"cart_id" gorm:"column:cart_id"`
	DeliveryAddressID int         `json:"address_id" gorm:"column:address_id"`
	PurchaseDate      time.Time   `gorm:"column:purchase_date" json:"purchase_date"`
	PaymentID         int         `json:"payment_id" gorm:"column:payment_id"`
	Status            OrderStatus `json:"status" gorm:"column:status"`
	TotalPrice        float64     `gorm:"column:total_price" json:"total_price"`
}

//TableName table name for Order
func (Order) TableName() string {
	return "orders"
}
