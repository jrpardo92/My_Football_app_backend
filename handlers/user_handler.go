package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"my-football-app/config"
	"my-football-app/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetUsersHandler para la lectura de usuarios
func GetUsersHandler(c echo.Context) error {
	// Obtener la conexión a la base de datos
	db, err := config.GetDBConnection()
	if err != nil {

		return c.String(http.StatusInternalServerError, "Error al conectar con la base de datos")
	}
	defer db.Close()

	// Consulta SQL para obtener todos los usuarios
	query := `SELECT UsuarioID, Nombre, Edad, CiudadResidencia, CorreoElectronico FROM Usuarios`
	rows, err := db.Query(query)
	if err != nil {
		// Utilizar el log de Echo para registrar el error
		c.Echo().Logger.Error("Error al realizar la consulta: ", err)
		return c.String(http.StatusInternalServerError, "Error al ejecutar la consulta")
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UsuarioID, &user.Nombre, &user.Edad, &user.CiudadResidencia, &user.CorreoElectronico); err != nil {
			log.Printf("Error al escanear usuario: %v", err)
			continue
		}
		users = append(users, user)
	}

	// Comprobar si hubo errores durante el iterar las filas
	if err = rows.Err(); err != nil {
		c.Echo().Logger.Error("Error al iterar sobre las filas: ", err)
		return c.String(http.StatusInternalServerError, "Error al procesar los resultados")
	}
	c.Echo().Logger.Info("Conexión a la base de datos establecida")

	if rows.Next() {
		c.Echo().Logger.Info("Encontrado un usuario")
	} else {
		c.Echo().Logger.Info("No se encontraron usuarios")
	}

	// Enviar los usuarios como respuesta JSON
	return c.JSON(http.StatusOK, users)
}

// RegisterUserHandler maneja el registro de nuevos usuarios
func RegisterUserHandler(c echo.Context) error {
	// Instancia para los datos del usuario
	newUser := new(models.User)
	if err := c.Bind(newUser); err != nil {
		return c.String(http.StatusBadRequest, "Datos inválidos")
	}
	// Log para depuración
	c.Echo().Logger.Info(fmt.Sprintf("Datos recibidos: %+v", newUser))

	// Validar datos del usuario (aquí puedes agregar más validaciones según sea necesario)
	if newUser.Nombre == "" || newUser.CorreoElectronico == "" {
		return c.String(http.StatusBadRequest, "Faltan datos requeridos")
	}

	// Conexión a la base de datos
	db, err := config.GetDBConnection()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al conectar con la base de datos")
	}
	defer db.Close()

	// Verificar si el correo electrónico ya está registrado
	var id int
	err = db.QueryRow("SELECT UsuarioID FROM Usuarios WHERE CorreoElectronico = $1", newUser.CorreoElectronico).Scan(&id)
	if err == nil {
		return c.String(http.StatusConflict, "El correo electrónico ya está registrado")
	} else if err != sql.ErrNoRows {
		c.Echo().Logger.Error("Error al verificar el correo electrónico: ", err)
		return c.String(http.StatusInternalServerError, "Error al verificar el correo electrónico")
	}

	// Insertar usuario en la base de datos
	query := `INSERT INTO Usuarios (Nombre, Edad, CiudadResidencia, CorreoElectronico, Contraseña) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(query, newUser.Nombre, newUser.Edad, newUser.CiudadResidencia, newUser.CorreoElectronico, newUser.Contraseña)
	if err != nil {
		//return c.String(http.StatusInternalServerError, "Error al registrar el usuario")
		c.Echo().Logger.Error("Error al insertar usuario en la base de datos: ", err)
		return c.String(http.StatusInternalServerError, "Error al registrar el usuario")
	}

	return c.String(http.StatusOK, "Usuario registrado con éxito")
}

// UpdateUserHandler actualiza los detalles de un usuario
func UpdateUserHandler(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID de usuario inválido")
	}

	var updatedUser models.User
	if err := c.Bind(&updatedUser); err != nil {
		return c.JSON(http.StatusBadRequest, "Datos inválidos")
	}

	db, err := config.GetDBConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error de conexión a la base de datos")
	}
	defer db.Close()

	query := `UPDATE Usuarios SET Nombre = $1, Edad = $2, CiudadResidencia = $3 WHERE UsuarioID = $4`
	_, err = db.Exec(query, updatedUser.Nombre, updatedUser.Edad, updatedUser.CiudadResidencia, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error al actualizar el usuario")
	}

	return c.JSON(http.StatusOK, "Usuario actualizado con éxito")
}

// DeleteUserHandler maneja la eliminación de un usuario
func DeleteUserHandler(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID de usuario inválido")
	}

	db, err := config.GetDBConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error de conexión a la base de datos")
	}
	defer db.Close()

	query := `DELETE FROM Usuarios WHERE UsuarioID = $1`
	result, err := db.Exec(query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error al eliminar el usuario")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error al obtener filas afectadas")
	}

	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, "Usuario no encontrado")
	}

	return c.JSON(http.StatusOK, "Usuario eliminado con éxito")
}

// GetUserHandler maneja las solicitudes para obtener la información de un usuario específico
func GetUserHandler(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID de usuario inválido")
	}

	db, err := config.GetDBConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error de conexión a la base de datos")
	}
	defer db.Close()

	var user models.User
	query := `SELECT UsuarioID, Nombre, Edad, CiudadResidencia, CorreoElectronico FROM Usuarios WHERE UsuarioID = $1`
	row := db.QueryRow(query, userID)
	err = row.Scan(&user.UsuarioID, &user.Nombre, &user.Edad, &user.CiudadResidencia, &user.CorreoElectronico)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, "Usuario no encontrado")
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error al consultar el usuario")
	}

	return c.JSON(http.StatusOK, user)
}
