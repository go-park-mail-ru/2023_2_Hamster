package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/app"
	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/postgresql"
	redisDB "github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/redis"
	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/server"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
)

// @title		Hamster API
// @version		1.0.1
// @description	Server API for Hamster Money Service Application

// @contact.name   Hamster API Support
// @contact.email  dimka.komarov@bk.ru
// @contact.email  grigorikovalenko@gmail.com
// @contact.url    https://t.me/CodeMaster482

// @host		localhost:8080
// @BasePath	/user/{userID}/account/feed

// @securityDefinitions	AuthKey
// @in				header
// @name				session_id

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 99999*time.Second)
	defer cancel()

	log := logger.NewLogger(ctx)

	// Postgre Connection
	db, err := postgresql.InitPostgresDB(ctx)
	if err != nil {
		log.Errorf("Error Initializing PostgreSQL database: %v", err)
		return
	}
	defer func() {
		db.Close()

		log.Info("Db closed without errors")
	}()

	log.Info("Db connection successfully")

	// redis-cli init
	redisCli, err := redisDB.InitRedisCli(ctx)
	if err != nil {
		log.Errorf("Error Initializing Redis-cli: %v", err)
		return
	}
	defer func() {
		if err := redisCli.Close(); err != nil {
			log.Errorf("Error Closing Redis connection: %v", err)
		}
		log.Info("Redis closed without errors")
	}()

	_, pingErr := redisCli.Ping(context.Background()).Result()
	if pingErr != nil {
		log.Errorf("Failed to ping Redis server: %v", pingErr)
	}

	log.Info("Redis connection successfully")

	router := app.Init(db, redisCli, log)

	var srv server.Server
	if err := srv.Init(router); err != nil {
		log.Fatalf("error while launching server: %v", err)
	}

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("can't launch server: %v", err)
		}
	}()
	log.Infof("server launcher at %s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exiting")
}
