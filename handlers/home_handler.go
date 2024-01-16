package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HomeHandler devuelve un mensaje de bienvenida
func HomeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Bienvenido a My Futbol App")
}
