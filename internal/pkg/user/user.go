package user

import "github.com/go-park-mail-ru/2023_2_Hamster/internal/models"

// Bussiness logic methods to work with user
type Usecase interface {
	GetByID(userID uint32) (*models.User, error)
	ChangeInfo(user *models.User) error
}

type Repository interface {
	GetByID(userID uint32) (*models.User, error)

	CreateUser(user models.User) error
}
