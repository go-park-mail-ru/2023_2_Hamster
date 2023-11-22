package postgresql

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
)

type AccountRep struct {
	db     postgresql.DbConn
	logger logger.Logger
}

func NewRepository(db postgresql.DbConn, log logger.Logger) *AccountRep {
	return &AccountRep{
		db:     db,
		logger: log,
	}
}
