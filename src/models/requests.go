package models

type ChargeRequest struct {
	OrderID     int    `json:"order_id"`
	UseExisting bool   `json:"use_existing"`
	SaveCard    bool   `json:"save_card"`
	Token       string `json:"token"`
}

//CreateOrderRequest create order
type CreateOrderRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
