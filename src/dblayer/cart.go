package dblayer

import (
	"time"

	"github.com/andrei-dascalu/go-workshop-shopapi/src/models"
)

//CreateNewCart create new cart
func (db *DBORM) CreateNewCart(user models.Customer) (cart models.Cart, err error) {
	now := time.Now()

	cart = models.Cart{
		Customer:       user,
		CustomerID:     user.ID,
		Status:         models.Cart_Status_Active,
		ExpirationDate: now.Add(time.Hour * 2),
	}

	return cart, db.Create(&cart).Error
}

//GetActiveCartForUser find active cart
func (db *DBORM) GetActiveCartForUser(user models.Customer) (cart models.Cart, err error) {
	return cart, db.Where("customer_id = ? AND expiration >= ?", user.ID, time.Now()).Find(&cart).Error
}

//AddProductToCart add product to cart
func (db *DBORM) AddProductToCart(product models.Product, quantity int, cart models.Cart) error {
	//pass received password by reference so that we can change it to it's hashed version

	item, err := db.GetCartItemForProduct(product, cart)

	if err != nil {
		item, err = db.CreateCartItemForProduct(product, quantity, cart)

		if err != nil {
			return err
		}

		return nil
	}

	item.Quantity += quantity
	_, err = db.UpdateCartItemForProduct(item)

	if err != nil {
		return err
	}

	return nil
}

//GetCartItemForProduct get cart item
func (db *DBORM) GetCartItemForProduct(product models.Product, cart models.Cart) (item models.CartItem, err error) {
	return item, db.Where("product_id = ?", product.ID).Find(&item).Error
}

//CreateCartItemForProduct create cart item
func (db *DBORM) CreateCartItemForProduct(product models.Product, quantity int, cart models.Cart) (item models.CartItem, err error) {
	sellingUnitPrice := product.Price

	if product.Price > product.Promotion {
		sellingUnitPrice = product.Promotion
	}

	item = models.CartItem{
		Product:   product,
		ProductID: product.ID,
		Cart:      cart,
		CartID:    cart.ID,
		UnitPrice: sellingUnitPrice,
		Quantity:  quantity,
	}

	return item, db.Create(&item).Error
}

//UpdateCartItemForProduct update a cart item
func (db *DBORM) UpdateCartItemForProduct(item models.CartItem) (models.CartItem, error) {
	return item, db.Save(&item).Error
}
