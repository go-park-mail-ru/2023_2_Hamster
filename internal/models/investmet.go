package models

// import (
// 	"time"

// 	valid "github.com/asaskevich/govalidator"
// 	"github.com/google/uuid"
// )

// type Investment struct {
// 	ID         uuid.UUID `json:"id" valid:"-"`
// 	UserID     uint      `json:"user_id" valid:"-"`
// 	Name       string    `json:"name" valid:"required"`
// 	Total      float64   `json:"total" valid:"required,greaterzero"`
// 	DateStart  time.Time `json:"date_start" valid:"isdate"`
// 	DateEnd    time.Time `json:"date_end" valid:"isdate"`
// 	Price      float64   `json:"price" valid:"required"`
// 	Percentage float64   `json:"percentage" valid:"required"`
// }

// func (i *Investment) InvestmentValidate() error {
// 	_, err := valid.ValidateStruct(i)
// 	return err
// }
