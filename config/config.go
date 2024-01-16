package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("dev") // nombre del archivo de configuración
	viper.AddConfigPath("config")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// funcion para conectar la base de datos
func GetDBConnection() (*sql.DB, error) {
	// Aquí deberías cargar tus configuraciones de `dev.yaml` o `prod.yaml`
	// Por ejemplo, usando variables de entorno o un paquete de configuración
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "demo12"
	dbname := "FutbolApp"

	// Construyendo la cadena de conexión
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Conectando a la base de datos
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Comprobar la conexión
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Conectado exitosamente a la base de datos")
	return db, nil
}
