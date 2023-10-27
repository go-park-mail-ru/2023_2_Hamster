package http

import "github.com/go-park-mail-ru/2023_2_Hamster/internal/models"

type MasTransaction struct {
	Transactions []models.TransactionTransfer `json:"transaction"`
}
