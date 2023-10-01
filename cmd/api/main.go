package main

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/postgresql"
	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/app"
)

func main() {
	ctx, cancel
	log := logger.CreateCustomLogger()

	db, err := postgresql.InitPostgresDB()
	if err != nil {
		log.Errorf("Error Initializing PostgreSQL database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Errorf("Error Closing database connection: %v", err)
		}
	}()
	
	router := app.Init(db, log)
	router, err := 
}
