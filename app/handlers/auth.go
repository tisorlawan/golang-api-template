package handlers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTUser struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	jwt.StandardClaims
}

func Login(c echo.Context) error {
	var l LoginInput

	if err := c.Bind(&l); err != nil {
		return echo.ErrBadRequest
	}

	// Authenticate
	if l.Username != "admin" || l.Password != "admin" {
		return echo.ErrUnauthorized
	}

	// Set claims
	claims := &JWTUser{
		Name: "Humprey Bogart",
		ID:   "123123",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
