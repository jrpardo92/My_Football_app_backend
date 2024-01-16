package handlers

import (
	"fmt"
	"my-football-app/config"
	"my-football-app/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CreateMatchHandler maneja la creación de nuevos partidos
func CreateMatchHandler(c echo.Context) error {
	newMatch := new(models.Match)
	if err := c.Bind(newMatch); err != nil {
		c.Echo().Logger.Error("Error al vincular datos: ", err)
		return c.JSON(http.StatusBadRequest, "Datos inválidos")
	}
	c.Echo().Logger.Info(fmt.Sprintf("Datos recibidos: %+v", newMatch))

	db, err := config.GetDBConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error de conexión a la base de datos")
	}
	defer db.Close()

	query := `INSERT INTO Partidos (Fecha, Hora, Lugar, Descripcion, CreadorID) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(query, newMatch.Fecha, newMatch.Hora, newMatch.Lugar, newMatch.Descripcion, newMatch.CreadorID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error al crear el partido")
	}

	return c.JSON(http.StatusCreated, "Partido creado con éxito")
}

// GetMatchesHandler maneja las solicitudes para obtener todos los partidos
func GetMatchesHandler(c echo.Context) error {
	db, err := config.GetDBConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error al conectar con la base de datos")
	}
	defer db.Close()

	var matches []models.Match
	query := `SELECT Partidos.PartidoID, Partidos.Fecha, Partidos.Hora, Partidos.Lugar, Partidos.Descripcion, Partidos.CreadorID, Usuarios.Nombre 
              FROM Partidos 
              JOIN Usuarios ON Partidos.CreadorID = Usuarios.UsuarioID`
	rows, err := db.Query(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error al consultar partidos")
	}
	defer rows.Close()

	for rows.Next() {
		var match models.Match
		if err := rows.Scan(&match.PartidoID, &match.Fecha, &match.Hora, &match.Lugar, &match.Descripcion, &match.CreadorID, &match.NombreCreador); err != nil {
			return c.JSON(http.StatusInternalServerError, "Error al leer datos de partido")
		}
		matches = append(matches, match)
	}

	if err = rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, "Error al procesar los resultados")
	}

	return c.JSON(http.StatusOK, matches)
}
