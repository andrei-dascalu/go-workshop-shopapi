package models

import "github.com/jinzhu/gorm"

//Payment payment data
type Payment struct {
	gorm.Model
	Order    Order
	Status   string `gorm:"column:status;type:varchar(128)"`
	StripeID string `gorm:"column:stripe_id;type:varchar(250)"`
}

//TableName table name for Payment
func (Payment) TableName() string {
	return "payments"
}
