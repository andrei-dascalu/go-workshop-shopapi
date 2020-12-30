package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andrei-dascalu/go-workshop-shopapi/src/configuration"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/dblayer"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/dto"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/models"
	"github.com/andrei-dascalu/go-workshop-shopapi/src/security"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

//GetMainPage get page
func GetMainPage(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(map[string]string{
		"message": "Main Page",
	})
}

//GetProducts get
func GetProducts(c *fiber.Ctx) error {
	products, err := dblayer.ShopDB.GetAllProducts()
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(products)
}

//AddProductToCart add prod
func AddProductToCart(c *fiber.Ctx) error {
	var productDto dto.AddProductDto

	err := c.BodyParser(&productDto)

	if err != nil {
		log.Error(err)
		return err
	}

	return c.Status(http.StatusOK).JSON(productDto)
}

//GetPromos promos
func GetPromos(c *fiber.Ctx) error {
	promos, err := dblayer.ShopDB.GetPromos()
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(promos)
}

//AddUser error
func AddUser(c *fiber.Ctx) error {
	var customer models.Customer

	err := c.BodyParser(&customer)
	if err != nil {
		return err
	}

	customer.LoggedIn = false

	customer, err = dblayer.ShopDB.AddUser(customer)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(customer)
}

//SignIn signin
func SignIn(c *fiber.Ctx) error {
	var customer models.Customer
	err := c.BodyParser(&customer)

	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  err.Error(),
			Internal: err,
		}
	}

	customer, err = dblayer.ShopDB.SignInUser(customer.Email, customer.Password)
	if err != nil {
		if err == dblayer.ErrINVALIDPASSWORD {
			return c.Status(http.StatusForbidden).JSON(map[string]string{
				"error": err.Error(),
			})
		}
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": "Unknown user",
		})
	}

	token, err := security.CreateJWTForUser(customer)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{
			"error": err.Error(),
		})
	}

	response := models.LoginResponse{
		Token: token,
	}

	return c.Status(http.StatusOK).JSON(response)
}

//SignOut signout
func SignOut(c *fiber.Ctx) error {
	p := c.Param("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
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
func GetOrders(c *fiber.Ctx) error {
	p := c.Param("id")
	id, err := strconv.Atoi(p)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
	}
	orders, err := dblayer.ShopDB.GetCustomerOrdersByID(id)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(orders)
}

//Charge charge
func Charge(c *fiber.Ctx) error {
	var expectedRequest models.ChargeRequest

	err := c.BodyParser(&expectedRequest)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(expectedRequest)
	}

	order, err := dblayer.ShopDB.FindOrderByID(expectedRequest.OrderID)

	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(expectedRequest)
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

//CreateOrder create order
func CreateOrder(c *fiber.Ctx) error {
	var expectedRequest models.CreateOrderRequest

	p := c.Param("id")

	userID, err := strconv.Atoi(p)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
	}

	err = c.BodyParser(&expectedRequest)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
	}

	customer, err := dblayer.ShopDB.GetCustomerByID(userID)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
	}

	cart, err := dblayer.ShopDB.GetActiveCartForUser(customer)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
	}

	address, err := dblayer.ShopDB.GetMainAddressForCustomer(customer)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
	}

	order, err := dblayer.ShopDB.CreateOrder(cart, address)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
	}

	log.Error(order)

	return c.Status(http.StatusOK).JSON(map[string]string{
		"order": fmt.Sprintf("%d", order.ID),
	})
}
