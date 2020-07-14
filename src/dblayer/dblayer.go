package dblayer

import (
	"errors"

	"github.com/andrei-dascalu/go-workshop-shopapi/src/models"
)

//DBLayer interface for db
type DBLayer interface {
	GetAllProducts() ([]models.Product, error)
	GetPromos() ([]models.Product, error)
	GetCustomerByName(string, string) (models.Customer, error)
	GetCustomerByID(int) (models.Customer, error)
	GetProduct(int) (models.Product, error)
	AddUser(models.Customer) (models.Customer, error)
	SignInUser(username, password string) (models.Customer, error)
	SignOutUserById(int) error
	GetCustomerOrdersByID(int) ([]models.Order, error)
	AddOrder(models.Order) error
	GetCreditCardCID(int) (string, error)
	SaveCreditCardForCustomer(int, string) error
}

/*
Passwords:
mal:123
john:1111
jayne:123
River:abc
*/

//ErrINVALIDPASSWORD custom error type
var ErrINVALIDPASSWORD = errors.New("Invalid password")
