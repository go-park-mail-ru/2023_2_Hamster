package postgresql

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
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

func initPostgresConfigFromEnv() (PostgresConfig, error) {
	var cfg = PostgresConfig{}
	if err := godotenv.Load(); err != nil {
		return cfg, err
	}

	host, existHost := os.LookupEnv("DB_HOST")
	port, existPort := os.LookupEnv("DB_PORT")
	user, existUser := os.LookupEnv("DB_USER")
	pass, existPass := os.LookupEnv("DB_PASSWORD")
	dbname, existName := os.LookupEnv("DB_NAME")
	dbsslmode, existSSL := os.LookupEnv("DB_SSLMODE")

	if !existHost || !existUser || !existPass || !existName || !existSSL || !existPort {
		return cfg, errors.New("existHost or existPort or existUser or existPass or existName is Empty")
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

func InitPostgresDB(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := initPostgresConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.DBSSLMode)

	conn, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		conn.Close()
		return conn, err
	}
	return conn, nil
}
