package handlers

import (
	// Otras importaciones necesarias
	"database/sql"
	"my-football-app/config"
	"my-football-app/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
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

// generateJWT genera un nuevo token JWT para un usuario
func generateJWT(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Crear un mapa para almacenar los claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.UsuarioID
	claims["correoElectronico"] = user.CorreoElectronico
	// Puedes agregar más claims aquí si lo deseas

	// Firmar y obtener el token en formato de cadena
	tokenString, err := token.SignedString([]byte("tuSecretoAquí")) // Usa una clave secreta adecuada
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
