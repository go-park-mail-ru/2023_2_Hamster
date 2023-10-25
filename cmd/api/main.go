package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/app"
	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/server"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/postgresql"
)

// @title		Hamster API
// @version		1.0.1
// @description	Server API for Hamster Money Service Application

// @contact.name   Hamster API Support
// @contact.email  dimka.komarov@bk.ru

// @host		localhost:8090
// @BasePath	/user/{userID}/account/feed

// @securityDefinitions	AuthKey
// @in					header
// @name				Authorization

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log := logger.CreateCustomLogger()

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

	router := app.Init(db, log)
	var srv server.Server
	go func() {
		if err := srv.Run(router); err != nil {
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
