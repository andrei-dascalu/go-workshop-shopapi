package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Cart shopping cart
type Cart struct {
	gorm.Model
	Customer       Customer
	Status         CartStatus `json:"cart-status" gorm:"column:cart_status"`
	CustomerID     uint       `json:"customer_id" gorm:"column:customer_id"`
	ExpirationDate time.Time  `json:"expiration" gorm:"column:expiration"`
}

//TableName table name for Customer
func (Cart) TableName() string {
	return "carts"
}

//CartItem link between cart and a product
type CartItem struct {
	gorm.Model
	Product   Product
	Cart      Cart
	CartID    uint    `json:"cart_id" gorm:"column:cart_id"`
	ProductID uint    `json:"product_id" gorm:"column:product_id"`
	Quantity  int     `json:"quantity" gorm:"column:quantity"`
	UnitPrice float64 `gorm:"column:price" json:"unit_price"`
}

//TableName table name for Customer
func (CartItem) TableName() string {
	return "cart_items"
}
