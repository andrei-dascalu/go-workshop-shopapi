package dto

type AddProductDto struct {
	ProductID uint `json:"product_id" form:"product_id" query:"product_id"`
	Quantity  int  `json:"quantity" form:"quantity" query:"quantity"`
}
