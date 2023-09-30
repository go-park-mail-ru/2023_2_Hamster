package dsn

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBName     string
	DBPassword string
	DBSSLMode  string
}

func InitPostgresConfigFromEnv() (PostgresConfig, error) {
	var cfg = PostgresConfig{}
	if err := godotenv.Load(); err != nil {
		return cfg, err
	}

	host, existHost := os.LookupEnv("DB_HOST")
	port, existPort := os.LookupEnv("DB_PORT")
	user, existUser := os.LookupEnv("DB_USER")
	pass, existPass := os.LookupEnv("DB_PASS")
	dbname, existName := os.LookupEnv("DB_NAME")
	dbsslmode, existSSL := os.LookupEnv("DB_SSLMODE")

	if !existHost || !existPort || !existUser || !existPass || !existName || existSSL {
		return cfg, fmt.Errorf("existHost or existPort or existUser or existPass or existName is Empty")
	}
	cfg = PostgresConfig{
		DBHost:     host,
		DBPort:     port,
		DBUser:     user,
		DBName:     dbname,
		DBPassword: pass,
		DBSSLMode:  dbsslmode,
	}
	return cfg, nil
}

// TO DO InitPostgresDB()
