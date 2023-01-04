package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	// For testing purpose
	// e.DELETE("/expenses/:id", expense.DeleteExpenseByID)

	// Start server
	if os.Getenv("PORT") == "" {
		log.Fatal("$PORT must be set")
	}

	go func() {
		e.Logger.Fatal(e.Start(os.Getenv("PORT")))
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	// Wait for interrupt signal
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := e.Shutdown(ctx)
	if err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("Shutting down server...")
}
