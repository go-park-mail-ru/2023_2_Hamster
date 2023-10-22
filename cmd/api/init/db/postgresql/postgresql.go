package postgresql

import (
	"errors"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBName     string
	DBPassword string
	DBSSLMode  string
}

func initPostgresConfigFromEnv() (PostgresConfig, error) {
	var cfg = PostgresConfig{}
	if err := godotenv.Load(); err != nil {
		return cfg, err
	}

	host, existHost := os.LookupEnv("DB_HOST")
	user, existUser := os.LookupEnv("DB_USER")
	pass, existPass := os.LookupEnv("DB_PASSWORD")
	dbname, existName := os.LookupEnv("DB_NAME")
	dbsslmode, existSSL := os.LookupEnv("DB_SSLMODE")

	if !existHost || !existUser || !existPass || !existName || !existSSL {
		return cfg, errors.New("existHost or existPort or existUser or existPass or existName is Empty")
	}

	cfg = PostgresConfig{
		DBHost:     host,
		DBPort:     "5432",
		DBUser:     user,
		DBName:     dbname,
		DBPassword: pass,
		DBSSLMode:  dbsslmode,
	}
	return cfg, nil
}

func InitPostgresDB() (*sqlx.DB, error) {
	cfg, err := initPostgresConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.DBSSLMode)

	db, err := sqlx.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		errClose := db.Close()
		if errClose != nil {
			return nil, fmt.Errorf("can't close postgresql (%w) after failed ping: %w", errClose, err)
		}
		return db, err
	}
	return db, nil
}
