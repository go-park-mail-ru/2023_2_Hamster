package models

// import (
// 	"time"

// 	valid "github.com/asaskevich/govalidator"
// 	"github.com/google/uuid"
// )

// type Credit struct {
// 	ID          uuid.UUID `json:"id" valid:"-"`
// 	AccountID   uint      `json:"account_id" valid:"-"`
// 	Total       float64   `json:"total" valid:"required,greaterzero"`
// 	DateStart   time.Time `json:"date_start" valid:"isdate"`
// 	DateEnd     time.Time `json:"date_end" valid:"isdate"`
// 	IsAnnuity   bool      `json:"is_annuity" valid:"required"`
// 	Creditor    string    `json:"creditor" valid:"-"`
// 	Description string    `json:"description" valid:"-"`
// 	Payments    int       `json:"payments" valid:"-"`
// 	Bank        string    `json:"bank" valid:"-"`
// }

// func (c *Credit) CreditValidate() error {
// 	_, err := valid.ValidateStruct(c)
// 	return err
// }
