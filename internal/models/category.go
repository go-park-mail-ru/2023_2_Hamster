package models

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `json:"id" valid:"-"`
	UserID      uuid.UUID `json:"user_id" valid:"required"`
	ParentID    uuid.UUID `json:"parent_id" valid:"-"`
	Image       int       `json:"image_id" valid:"-"`
	Name        string    `json:"name" valid:"required"`
	ShowIncome  bool      `json:"show_income" valid:"-"`
	ShowOutcome bool      `json:"show_outcome" valid:"-"`
	Regular     bool      `json:"regular" valid:"-"`
}

func (c *Category) CategoryValidate() error {
	_, err := valid.ValidateStruct(c)
	return err
}
