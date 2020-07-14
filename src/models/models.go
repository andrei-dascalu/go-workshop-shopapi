package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Product product data
type Product struct {
	gorm.Model
	Image       string  `json:"img"`
	SmallImage  string  `gorm:"column:smallimg" json:"small_img"`
	ImagAlt     string  `json:"imgalt" gorm:"column:imgalt"`
	Price       float64 `json:"price"`
	Promotion   float64 `json:"promotion"` //sql.NullFloat64
	PoructName  string  `gorm:"column:productname" json:"productname"`
	Description string
}

func (Product) TableName() string {
	return "products"
}

//Customer customer data
type Customer struct {
	gorm.Model
	Name      string  `json:"name" sql:"-"`
	FirstName string  `gorm:"column:firstname" json:"firstname"`
	LastName  string  `gorm:"column:lastname" json:"lastname"`
	Email     string  `gorm:"column:email" json:"email"`
	Pass      string  `json:"password"`
	LoggedIn  bool    `gorm:"column:loggedin" json:"loggedin"`
	Orders    []Order `json:"orders"`
}

func (Customer) TableName() string {
	return "customers"
}

//Order order data
type Order struct {
	gorm.Model
	Product      Product
	Customer     Customer
	CustomerID   uint      `json:"customer_id" gorm:"column:customer_id"`
	ProductID    uint      `json:"product_id" gorm:"column:product_id"`
	Price        float64   `gorm:"column:price" json:"sell_price"`
	PurchaseDate time.Time `gorm:"column:purchase_date" json:"purchase_date"`
	PaymentID    uint      `json:"payment_id" gorm:"column:payment_id"`
}

func (Order) TableName() string {
	return "orders"
}

//Payment payment data
type Payment struct {
	gorm.Model
	Order    Order
	Status   string `gorm:"column:status;type:varchar(128)"`
	StripeID string `gorm:"column:stripe_id;type:varchar(250)"`
}

func (Payment) TableName() string {
	return "payments"
}
