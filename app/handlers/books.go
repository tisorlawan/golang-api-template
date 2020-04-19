package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Book struct {
	Title         string    `json:"title" validate:"required"`
	Author        string    `json:"author" validate:"required"`
	PublishedDate time.Time `json:"published_date"`
}

func HandleBook(c echo.Context) (err error) {
	book := new(Book)
	if err = c.Bind(book); err != nil {
		return
	}
	if err = c.Validate(book); err != nil {
		return
	}
	// log.Info("Book", zap.String("Title", book.Title))
	// book.PublishedDate = book.PublishedDate.Add(1000 * time.Hour)
	return c.JSON(http.StatusOK, book)
}
