package models

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type Sex string

const (
	Male   Sex = "M"
	Female Sex = "F"
	Other  Sex = "O" // :-)
)

type User struct {
	ID        uuid.UUID `json:"id" valid:"-"`
	Username  string    `json:"name" valid:"-"`
	Email     string    `json:"email" valid:"required,email"`
	FirstName string    `json:"firstName" valid:"required,runelength(2|20)"`
	LastName  string    `json:"lastName" valid:"required,runelength(2|20)"`
	Password  string    `json:"password" valid:"required,runelength(7|30),passwordcheck"`
	Sex       Sex       `json:"sex" valid:"required,in(F|M|O)"`
	AvatarURL string    `json:"avatar_url" vaild:"-"`
}

func (u *User) UserValidate() error {
	_, err := valid.ValidateStruct(u)
	return err
}