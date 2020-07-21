package models

import (
	"github.com/jinzhu/gorm"
)

//Address customer address
type Address struct {
	gorm.Model
	Customer   Customer
	Address    string `json:"address" gorm:"column:address"`
	CustomerID uint   `json:"customer_id" gorm:"column:customer_id"`
	IsMain     bool   `json:"main_address" gorm:"column:is_main"`
}

//TableName table name for Customer
func (Address) TableName() string {
	return "addresses"
}
