package models

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type Category struct {
	ID     uuid.UUID `json:"id" valid:"-"`
	UserID uint      `json:"user_id" valid:"-"`
	Name   string    `json:"name" valid:"required"`
}

func (c *Category) CategoryValidate() error {
	_, err := valid.ValidateStruct(c)
	return err
}
