package models

import "database/sql/driver"

//OrderStatus enumeration
type OrderStatus int64

const (
	Order_Status_Init OrderStatus = iota
	Order_Status_Pending
	Order_Status_Complete
	Order_Status_Failed
)

var types = [...]string{
	"Started",
	"Pending",
	"Complete",
	"Failed",
}

func (orderStatus OrderStatus) String() string {
	return types[orderStatus]
}

func (u *OrderStatus) Scan(value interface{}) error { *u = OrderStatus(value.(int64)); return nil }
func (u OrderStatus) Value() (driver.Value, error)  { return int64(u), nil }
