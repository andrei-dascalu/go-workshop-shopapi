package models

import "database/sql/driver"

//OrderStatus enumeration
type OrderStatus int64

//CartStatus enumeration
type CartStatus int64

const (
	Order_Status_Init OrderStatus = iota
	Order_Status_Pending
	Order_Status_Complete
	Order_Status_Failed
)

const (
	Cart_Status_Active CartStatus = iota
	Cart_Status_Complete
	Cart_Status_Closed
)

var orderTypes = [...]string{
	"Started",
	"Pending",
	"Complete",
	"Failed",
}

var cartTypes = [...]string{
	"Active",
	"Complete",
	"Closed",
}

func (orderStatus OrderStatus) String() string {
	return orderTypes[orderStatus]
}

func (cartStatus CartStatus) String() string {
	return cartTypes[cartStatus]
}

func (u *OrderStatus) Scan(value interface{}) error { *u = OrderStatus(value.(int64)); return nil }
func (u OrderStatus) Value() (driver.Value, error)  { return int64(u), nil }

func (u *CartStatus) Scan(value interface{}) error { *u = CartStatus(value.(int64)); return nil }
func (u CartStatus) Value() (driver.Value, error)  { return int64(u), nil }
