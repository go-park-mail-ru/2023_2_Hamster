package models

import (
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type Goal struct {
	ID          uuid.UUID `json:"id" valid:"-"`
	UserId      uuid.UUID `json:"user_id" valid:"-"`
	Name        string    `json:"name" valid:"required"`
	Description string    `json:"description" valid:"-"`
	Target      float64   `json:"total" valid:"required,greaterzero"`
	Date        time.Time `json:"date" valid:"isdate"`
}

func (g *Goal) GoalValidate() error {
	_, err := valid.ValidateStruct(g)
	return err
}
