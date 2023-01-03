package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/JinDamanee2544/assessment/expense"
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

func main() {

	expense.InitDB()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(validateToken)

	// Routes
	e.POST("/expenses", expense.CreateExpense)
	e.GET("/expenses", expense.GetAllExpense)
	e.GET("/expenses/:id", expense.GetExpenseByID)
	e.PUT("/expenses/:id", expense.UpdateExpenseByID)

	// Start server
	if os.Getenv("PORT") == "" {
		log.Fatal("$PORT must be set")
	}

	e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}
