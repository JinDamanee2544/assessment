//go:build integration

package expense

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func validateToken(next echo.HandlerFunc) echo.HandlerFunc {
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

func initServer() func() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(validateToken)

	e.POST("/expenses", CreateExpense)
	e.GET("/expenses", GetAllExpense)
	e.GET("/expenses/:id", GetExpenseByID)
	e.PUT("/expenses/:id", UpdateExpenseByID)
	e.DELETE("/expenses/:id", DeleteExpenseByID)

	e.Logger.Fatal(e.Start(os.Getenv("PORT")))

	shutdown := func() {
		e.Logger.Fatal(e.Shutdown(context.Background()))
	}

	return shutdown
}

func TestPostExpense(t *testing.T) {
	shutdown := initServer()
	defer shutdown()

	e := Expense{}
	body := bytes.NewBufferString(`{
		"id": "1",
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
		}`)
	res := request("POST", "http://localhost:2565/expenses", body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "strawberry smoothie", e.Title)
	assert.Equal(t, 89, e.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", e.Note)
	assert.Equal(t, []string{"food", "beverage"}, e.Tags)
}
