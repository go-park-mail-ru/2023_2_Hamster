package models

// import (
// 	"time"

// 	valid "github.com/asaskevich/govalidator"
// 	"github.com/google/uuid"
// )

// type Debt struct {
// 	ID          uuid.UUID `json:"id" valid:"-"`
// 	UserID      uint      `json:"user_id" valid:"-"`
// 	Total       float64   `json:"total" valid:"required,greaterzero"`
// 	Date        time.Time `json:"date" valid:"isdate"`
// 	Creditor    string    `json:"creditor"`
// 	Description string    `json:"description" valid:"-"`
// }

// func (d *Debt) DebtValidate() error {
// 	_, err := valid.ValidateStruct(d)
// 	return err
// }
