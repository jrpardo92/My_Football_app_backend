package routes

import (
	"my-football-app/handlers"

	"github.com/labstack/echo/v4"
)

// InitializeRoutes configura las rutas para la aplicación
func InitializeRoutes(e *echo.Echo) {
	// Ruta para la página de inicio
	e.GET("/", handlers.HomeHandler)

	// Rutas para obtener usuarios
	e.GET("/users", handlers.GetUsersHandler)

	// Ruta para registrar usuarios
	e.POST("/register", handlers.RegisterUserHandler)

	// Ruta para actualizar un usuario
	e.PUT("/users/:id", handlers.UpdateUserHandler)

	// Ruta para eliminar un usuario
	e.DELETE("/users/:id", handlers.DeleteUserHandler)
}
