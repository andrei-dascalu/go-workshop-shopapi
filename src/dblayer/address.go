package dblayer

import "github.com/andrei-dascalu/go-workshop-shopapi/src/models"

//GetMainAddressForCustomer get main address
func (db *DBORM) GetMainAddressForCustomer(c models.Customer) (a models.Address, err error) {
	return a, db.Where("is_main = ?", true).Find(&a).Error
}
