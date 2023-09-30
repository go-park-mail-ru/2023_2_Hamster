package models

import (
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type deposit struct {
	ID           uuid.UUID `json:"id" valid:"-"`
	AccountID    uint      `json:"account_id" valid:"-"`
	Total        float64   `json:"total" valid:"required,greaterzero"`
	DateStart    time.Time `json:"date_start" valid:"isdate"`
	DateEnd      time.Time `json:"date_end" valid:"isdate"`
	InterestRate float64   `json:"interest_rate" valid:"required"`
	Bank         string    `json:"bank"`
}

func (g *Goal) GoalValidate() error {
	_, err := valid.ValidateStruct(g)
	return err
}
