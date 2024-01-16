package models

// User representa un usuario en la base de datos
type User struct {
	UsuarioID         int    `json:"usuarioId"`         // Identificador único del usuario
	Nombre            string `json:"nombre"`            // Nombre del usuario
	Edad              int    `json:"edad"`              // Edad del usuario
	CiudadResidencia  string `json:"ciudadResidencia"`  // Ciudad de residencia del usuario
	CorreoElectronico string `json:"correoElectronico"` // Correo electrónico del usuario
	Contraseña        string `json:"contraseña"`        // Contraseña del usuario (no se incluye en JSON)
}
