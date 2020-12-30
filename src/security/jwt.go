package security

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func RandomLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Warn("test")

		return next(c)
	}
}

//CustomJWTMiddleware custom JWT validation middleware
func CustomJWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		extractor := jwtFromHeader("Bearer")

		data, err := extractor(c)

		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  err.Error(),
				Internal: err,
			}
		}

		err = verifyToken(data)

		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  err.Error(),
				Internal: err,
			}
		}

		c.Logger().Warnf("Whatever: %s", data)

		return next(c)
	}
}

func jwtFromHeader(authScheme string) func(echo.Context) (string, error) {
	header := "Authorization"
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", fmt.Errorf("Missing Token")
	}
}
