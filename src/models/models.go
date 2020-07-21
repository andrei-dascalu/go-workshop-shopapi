package models

import (
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

//TableName table name for Product
func (Product) TableName() string {
	return "products"
}

//Customer customer data
type Customer struct {
	gorm.Model
	Name      string  `json:"name" sql:"-"`
	FirstName string  `gorm:"column:firstname" json:"firstname"  validate:"required"`
	LastName  string  `gorm:"column:lastname" json:"lastname"  validate:"required"`
	Email     string  `gorm:"column:email" json:"email"  validate:"required,email"`
	Password  string  `gorm:"column:password" json:"password"  validate:"required"`
	LoggedIn  bool    `gorm:"column:loggedin" json:"loggedin"`
	Orders    []Order `json:"orders"`
}

//TableName table name for Customer
func (Customer) TableName() string {
	return "customers"
}
