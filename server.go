package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/JinDamanee2544/assessment/expense"
)

func main() {

	expense.InitDB()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/expenses", expense.PostExpense)

	// Start server
	if os.Getenv("PORT") == "" {
		log.Fatal("$PORT must be set")
	}

	e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}
