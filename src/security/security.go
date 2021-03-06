package security

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/andrei-dascalu/go-workshop-shopapi/src/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

//CreateJWTForUser perform login, return token
func CreateJWTForUser(c models.Customer) (string, error) {
	//using plain text password => should come from env
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(ReadPrivateKey(), "testing")

	if err != nil {
		log.Error().Msg("RSA Key Parse Error")

		return "", &echo.HTTPError{
			Code:    http.StatusExpectationFailed,
			Message: "Failed to read key",
		}
	}

	token := jwt.New(jwt.SigningMethodRS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = c.Email
	claims["role"] = "customer"
	claims["exp"] = time.Now().Add(time.Hour * 72)

	t, err := token.SignedString(privateKey)

	if err != nil {
		log.Error().Msg("Failed to sign with private key")

		return "", err
	}

	return t, nil
}

func verifyToken(token string) error {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(ReadPublicKey())

	if err != nil {
		log.Error().Msg("Failed to Parse public key")

		return err
	}

	tokenObject, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		log.Error().Msg("Error parsing token")

		return err
	}

	if !tokenObject.Valid {
		return fmt.Errorf("Invalid Token")
	}

	return nil
}

//ReadPrivateKey get contents of private key
func ReadPrivateKey() []byte {
	path := "config/jwt/private.pem"

	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.Error().Msg("Error reading private key")

		return nil
	}

	return data
}

//ReadPublicKey get contents of public key
func ReadPublicKey() []byte {
	path := "config/jwt/public.pem"

	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.Error().Msg("Error reading public key")

		return nil
	}

	return data
}
