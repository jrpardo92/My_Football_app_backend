package main

import (
	"my-football-app/routes"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Inicializar las rutas
	routes.InitializeRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// HomeHandler devuelve un mensaje de bienvenida
func HomeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Bienvenido a My Futbol App")
}
