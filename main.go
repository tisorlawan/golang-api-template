package main

import (
	"fmt"
	"library/app/utils"
	"library/app/utils/log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Book struct {
	Title         string    `json:"title" validate:"required"`
	Author        string    `json:"author" validate:"required"`
	PublishedDate time.Time `json:"published_date"`
}

func main() {
	e := echo.New()
	fmt.Println("\n⇨ PID:", os.Getpid())

	defer log.Sync()

	// Global middleware
	e.Use(middleware.Recover())
	e.Use(log.GetMiddleware())

	// Register handlers
	e.POST("/books", func(c echo.Context) (err error) {
		book := new(Book)
		if err = c.Bind(book); err != nil {
			return
		}
		if err = c.Validate(book); err != nil {
			return
		}
		log.Info("Book", zap.String("Title", book.Title))
		book.PublishedDate = book.PublishedDate.Add(1000 * time.Hour)
		return c.JSON(http.StatusOK, book)
	})

	e.Validator = utils.NewValidator()
	e.Debug = true
	e.Logger.Fatal(e.Start(":8000"))
}
