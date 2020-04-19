package main

import (
	"fmt"
	"library/app/handlers"
	"library/app/utils"
	"library/app/utils/log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*handlers.JWTUser)

	return c.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Hi %s - %s, welcome!", claims.Name, claims.ID),
	})
}

func main() {
	e := echo.New()
	fmt.Println("\n⇨ PID:", os.Getpid())

	defer log.Sync()

	// Global middleware
	e.Use(middleware.Recover())
	e.Use(log.GetMiddleware())

	// Register handlers
	e.POST("/login", handlers.Login)
	e.POST("/books", handlers.HandleBook)

	// Restricted

	config := middleware.JWTConfig{
		Claims:     &handlers.JWTUser{},
		SigningKey: []byte("secret"),
	}
	jwt := middleware.JWTWithConfig(config)

	r := e.Group("/admin", jwt)
	r.POST("", Restricted)

	// Profiling
	// profiling.Wrap(e)

	e.Validator = utils.NewValidator()
	e.Logger.Fatal(e.Start(":8001"))
}
