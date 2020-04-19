package main

import (
	"fmt"
	"library/app/handlers"
	"library/app/utils"
	"library/app/utils/log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	fmt.Println("\n⇨ PID:", os.Getpid())

	defer log.Sync()

	// Global middleware
	e.Use(middleware.Recover())
	e.Use(log.GetMiddleware())

	// Register handlers
	e.POST("/books", handlers.HandleBook)

	e.Validator = utils.NewValidator()
	e.Debug = true // To enable validation error message
	e.Logger.Fatal(e.Start(":8000"))
}
