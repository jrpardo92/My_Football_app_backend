package models

// Match representa un partido de fútbol en la base de datos
type Match struct {
	PartidoID     int    `json:"partidoId"`
	Fecha         string `json:"fecha"`
	Hora          string `json:"hora"`
	Lugar         string `json:"lugar"`
	Descripcion   string `json:"descripcion"`
	CreadorID     int    `json:"creadorId"`
	NombreCreador string `json:"nombreCreador"` // Nombre del creador
	// Agrega aquí otros campos relevantes del creador
}

/* En este modelo:
PartidoID es el identificador único del partido.
Fecha y Hora representan cuándo se jugará el partido. Nota que time.Time se utiliza tanto para la fecha como para la hora. Dependiendo de cómo quieras manejar la fecha y la hora, podrías necesitar ajustar esto (por ejemplo, usando solo la fecha y omitiendo la hora, o combinando fecha y hora en un solo campo).
Lugar es una cadena que describe dónde se jugará el partido.
Descripcion proporciona detalles adicionales sobre el partido.
CreadorID es el ID del usuario que crea el partido, lo que puede ser útil para la gestión de permisos o para mostrar información relacionada con el creador del partido.
*/
