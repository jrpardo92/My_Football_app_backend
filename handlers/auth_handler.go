package handlers

import (
	// Otras importaciones necesarias
	"database/sql"
	"my-football-app/config"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest define la estructura de la solicitud de inicio de sesión
type LoginRequest struct {
	CorreoElectronico string `json:"correoElectronico"`
	Contraseña        string `json:"contraseña"`
}

// LoginHandler maneja el inicio de sesión de los usuarios
func LoginHandler(c echo.Context) error {
	var loginReq LoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, "Datos inválidos")
	}

	db, err := config.GetDBConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			"Error de conexión a la base de datos")
	}
	defer db.Close()

	// Buscar usuario por correo electrónico
	var hashedPassword string
	err = db.QueryRow("SELECT Contraseña FROM Usuarios WHERE CorreoElectronico = $1", loginReq.CorreoElectronico).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, "Correo electrónico incorrecto")
		}
		return c.JSON(http.StatusInternalServerError, "Error al consultar la base de datos")
	}

	// Comprobar la contraseña hasheada
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginReq.Contraseña)); err != nil {
		return c.JSON(http.StatusUnauthorized, "contraseña incorrecta")
	}

	// Lógica para el inicio de sesión exitoso (puede incluir generación de token JWT, etc.)
	// ...
	return c.JSON(http.StatusOK, "Inicio de sesión exitoso")
}
