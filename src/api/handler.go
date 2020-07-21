package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andrei-dascalu/go-workshop-shopapi/src/configuration"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/dblayer"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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

	customer.LoggedIn = false

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

	customer, err = dblayer.ShopDB.SignInUser(customer.Email, customer.Password)
	if err != nil {
		if err == dblayer.ErrINVALIDPASSWORD {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Unknown user",
		})
	}

	customer.Password = ""

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

	err = dblayer.ShopDB.SignOutUserByID(id)
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

	var expectedRequest models.ChargeRequest

	err := c.Bind(&expectedRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, expectedRequest)
	}

	order, err := dblayer.ShopDB.FindOrderByID(expectedRequest.OrderID)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, expectedRequest)
	}

	log.Errorf("%v", order)

	stripe.Key = configuration.Config.StripeSecretKey
	//test cards available at:	https://stripe.com/docs/testing#cards
	//setting charge parameters
	chargeP := &stripe.ChargeParams{
		Amount:      stripe.Int64(int64(16)),
		Currency:    stripe.String("usd"),
		Description: stripe.String("GoMusic charges..."),
	}
	stripeCustomerID := ""
	//Either remembercard or use exeisting should be enabled but not both
	if expectedRequest.UseExisting {
		//use existing
		stripeCustomerID, err = dblayer.ShopDB.GetCreditCardCID(expectedRequest.OrderID)
		if err != nil {
			return err
		}
	} else {
		cp := &stripe.CustomerParams{}
		cp.SetSource(expectedRequest.Token)
		customer, err := customer.New(cp)
		if err != nil {
			return err
		}
		stripeCustomerID = customer.ID
		if expectedRequest.SaveCard {
			//save card!!
			err = dblayer.ShopDB.SaveCreditCardForCustomer(expectedRequest.OrderID, stripeCustomerID)
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

	return nil
}

func CreateOrder(c echo.Context) error {
	var expectedRequest models.CreateOrderRequest

	p := c.Param("id")

	userID, err := strconv.Atoi(p)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	err = c.Bind(&expectedRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	customer, err := dblayer.ShopDB.GetCustomerByID(userID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	cart, err := dblayer.ShopDB.GetActiveCartForUser(customer)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	address, err := dblayer.ShopDB.GetMainAddressForCustomer(customer)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	order, err := dblayer.ShopDB.CreateOrder(cart, address)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	log.Error(order)

	return c.JSON(http.StatusOK, map[string]string{
		"order": fmt.Sprintf("%d", order.ID),
	})
}
