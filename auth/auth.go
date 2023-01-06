package auth

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		correctToken := os.Getenv("TOKEN")
		if correctToken == "" {
			log.Fatal("TOKEN is not set")
		}

		if token != os.Getenv("TOKEN") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}
		return next(c)
	}
}
