package main

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/postgresql"
)

func main() {
	log := logger.CreateCustomLogger()

	db, err := postgresql.InitPostgresDB()
	if err != nil {
		log.Errorf("error when trying connect database")
	}
}
