package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/JinDamanee2544/assessment/auth"
	"github.com/JinDamanee2544/assessment/expense"
)

func main() {

	expense.InitDB()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(auth.ValidateToken)

	// Routes
	e.POST("/expenses", expense.CreateExpense)
	e.GET("/expenses", expense.GetAllExpense)
	e.GET("/expenses/:id", expense.GetExpenseByID)
	e.PUT("/expenses/:id", expense.UpdateExpenseByID)
	e.DELETE("/expenses/:id", expense.DeleteExpenseByID)

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
