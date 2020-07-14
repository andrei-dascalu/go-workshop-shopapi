package api

import (
	"net/http"
	"strconv"

	"github.com/andrei-dascalu/go-workshop-shopapi/src/configuration"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/dblayer"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/models"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

//GetMainPage get page
func GetMainPage(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Main Page",
	})
}

//GetProducts get
func GetProducts(c echo.Context) error {

	products, err := dblayer.ShopDB.GetAllProducts()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, products)
}

//GetPromos promos
func GetPromos(c echo.Context) error {

	promos, err := dblayer.ShopDB.GetPromos()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, promos)
}

//AddUser error
func AddUser(c echo.Context) error {

	var customer models.Customer
	err := c.Bind(&customer)
	if err != nil {
		return err
	}
	customer, err = dblayer.ShopDB.AddUser(customer)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, customer)
}

//SignIn signin
func SignIn(c echo.Context) error {

	var customer models.Customer
	err := c.Bind(&customer)
	if err != nil {
		return err
	}
	customer, err = dblayer.ShopDB.SignInUser(customer.Email, customer.Pass)
	if err != nil {
		if err == dblayer.ErrINVALIDPASSWORD {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": err.Error(),
			})
		}
		return err
	}
	return c.JSON(http.StatusOK, customer)
}

//SignOut signout
func SignOut(c echo.Context) error {

	p := c.Param("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	err = dblayer.ShopDB.SignOutUserById(id)
	if err != nil {
		return err
	}

	return nil
}

//GetOrders get orders
func GetOrders(c echo.Context) error {

	p := c.Param("id")
	id, err := strconv.Atoi(p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	orders, err := dblayer.ShopDB.GetCustomerOrdersByID(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, orders)
}

//Charge charge
func Charge(c echo.Context) error {

	request := struct {
		models.Order
		Remember    bool   `json:"rememberCard"`
		UseExisting bool   `json:"useExisting"`
		Token       string `json:"token"`
	}{}

	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, request)
	}

	stripe.Key = configuration.Config.StripeSecretKey
	//test cards available at:	https://stripe.com/docs/testing#cards
	//setting charge parameters
	chargeP := &stripe.ChargeParams{
		Amount:      stripe.Int64(int64(request.Price)),
		Currency:    stripe.String("usd"),
		Description: stripe.String("GoMusic charge..."),
	}
	stripeCustomerID := ""
	//Either remembercard or use exeisting should be enabled but not both
	if request.UseExisting {
		//use existing
		stripeCustomerID, err = dblayer.ShopDB.GetCreditCardCID(request.CustomerID)
		if err != nil {
			return err
		}
	} else {
		cp := &stripe.CustomerParams{}
		cp.SetSource(request.Token)
		customer, err := customer.New(cp)
		if err != nil {
			return err
		}
		stripeCustomerID = customer.ID
		if request.Remember {
			//save card!!
			err = dblayer.ShopDB.SaveCreditCardForCustomer(request.CustomerID, stripeCustomerID)
			if err != nil {
				return err
			}
		}
	}
	//we should check if the customer already ordered the same item or not but for simplicity, let's assume it's a new order
	chargeP.Customer = stripe.String(stripeCustomerID)
	_, err = charge.New(chargeP)
	if err != nil {
		return err
	}

	err = dblayer.ShopDB.AddOrder(request.Order)
	if err != nil {
		return err
	}

	return nil
}
