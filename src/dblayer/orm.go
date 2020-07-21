package dblayer

import (
	"github.com/andrei-dascalu/go-workshop-shopapi/src/configuration"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//ShopDB global ORM
var ShopDB *DBORM

//DBORM ORM type wrapper
type DBORM struct {
	*gorm.DB
}

//InitORM init ORM
func InitORM() error {
	db, err := gorm.Open("mysql", configuration.Config.GetConnectionString())
	if err != nil {
		return err
	}

	ShopDB = &DBORM{
		DB: db,
	}

	return nil
}

func (db *DBORM) GetAllProducts() (products []models.Product, err error) {
	return products, db.Find(&products).Error
}

func (db *DBORM) GetPromos() (products []models.Product, err error) {
	return products, db.Where("promotion IS NOT NULL").Find(&products).Error

}

func (db *DBORM) GetCustomerByName(firstname string, lastname string) (customer models.Customer, err error) {
	return customer, db.Where(&models.Customer{FirstName: firstname, LastName: lastname}).Find(&customer).Error
}

func (db *DBORM) GetCustomerByID(id int) (customer models.Customer, err error) {
	return customer, db.First(&customer, id).Error
}

func (db *DBORM) GetProduct(id int) (product models.Product, error error) {
	return product, db.First(&product, id).Error
}

//AddUser add user
func (db *DBORM) AddUser(customer models.Customer) (models.Customer, error) {
	//pass received password by reference so that we can change it to it's hashed version
	pass, err := hashPassword(customer.Password)

	if err != nil {
		return customer, err
	}

	customer.Password = pass

	err = db.Create(&customer).Error

	customer.Password = ""

	return customer, err
}

func hashPassword(s string) (string, error) {
	//converd password string to byte slice
	sBytes := []byte(s)
	//Obtain hashed password
	hashedBytes, err := bcrypt.GenerateFromPassword(sBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	//update password string with the hashed version
	return string(hashedBytes[:]), nil
}

//SignInUser sign in
func (db *DBORM) SignInUser(email, pass string) (customer models.Customer, err error) {
	//Obtain a *gorm.DB object representing our customer's row
	result := db.Table("customers").Where(&models.Customer{Email: email})
	err = result.First(&customer).Error

	if err != nil {
		return customer, err
	}

	if !checkPassword(customer.Password, pass) {
		return customer, ErrINVALIDPASSWORD
	}

	customer.Password = ""
	//update the loggedin field
	err = result.Update("loggedin", 1).Error
	if err != nil {
		return customer, err
	}
	//return the new customer row
	return customer, result.Find(&customer).Error
}

func checkPassword(existingHash, incomingPass string) bool {
	//this method will return an error if the hash does not match the provided password string
	return bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(incomingPass)) == nil
}

//SignOutUserByID signout
func (db *DBORM) SignOutUserByID(id int) error {
	customer := models.Customer{
		Model: gorm.Model{
			ID: uint(id),
		},
	}
	return db.Table("Customers").Where(&customer).Update("loggedin", 0).Error
}

//GetCustomerOrdersByID get customers
func (db *DBORM) GetCustomerOrdersByID(id int) (orders []models.Order, err error) {
	return orders, db.Table("orders").Select("*").Joins("join customers on customers.id = customer_id").Joins("join products on products.id = product_id").Where("customer_id=?", id).Scan(&orders).Error //db.Find(&orders, models.Order{CustomerID: id}).Error
}

//GetCreditCardCID get cc
func (db *DBORM) GetCreditCardCID(id int) (string, error) {

	cusomterWithCCID := struct {
		models.Customer
		CCID string `gorm:"column:cc_customerid"`
	}{}

	return cusomterWithCCID.CCID, db.First(&cusomterWithCCID, id).Error
}

//SaveCreditCardForCustomer save cc
func (db *DBORM) SaveCreditCardForCustomer(id int, ccid string) error {
	result := db.Table("customers").Where("id=?", id)
	return result.Update("cc_customerid", ccid).Error
}
