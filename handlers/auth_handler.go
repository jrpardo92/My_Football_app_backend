package handlers

import (
	// Otras importaciones necesarias
	"database/sql"
	"log"
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
		return c.JSON(http.StatusInternalServerError, "Error de conexión a la base de datos")
	}
	defer db.Close()

	var user models.User
	err = db.QueryRow("SELECT UsuarioID, Contraseña FROM Usuarios WHERE CorreoElectronico = $1", loginReq.CorreoElectronico).Scan(&user.UsuarioID, &user.Contraseña)
	//comprobar correo
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, "Correo electrónico")
		}
		return c.JSON(http.StatusInternalServerError, "Error al consultar la base de datos")
	}
	// Comprobar la contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Contraseña), []byte(loginReq.Contraseña)); err != nil {
		return c.JSON(http.StatusUnauthorized, " contraseña incorrectos")
	}

	// Generar JWT
	tokenString, err := generateJWT(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error al generar el token")
	}

	// Devolver el token en la respuesta
	return c.JSON(http.StatusOK, echo.Map{
		"mensaje": "Inicio de sesión exitoso",
		"token":   tokenString,
	})
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
	tokenString, err := token.SignedString([]byte("aB3!f#9ZpQ8rS*7tVx4wC&1dE@5gH+6j")) // Utiliza una clave secreta adecuada y segura
	if err != nil {
		log.Printf("Error al generar el token JWT: %v", err)
		return "", err
	}
	log.Printf("Token generado: %s", tokenString)
	return tokenString, nil
}
