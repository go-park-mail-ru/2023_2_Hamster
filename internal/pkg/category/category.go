package category

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase interface {
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
	CreateCategory(ctx context.Context, category models.Category) (uuid.UUID, error)
	GetFeed(ctx context.Context, categoryID uuid.UUID) ([]models.Category, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
}

type Repository interface {
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
	CreateCategory(ctx context.Context, category models.Category) (uuid.UUID, error)
	GetFeed(ctx context.Context, userID uuid.UUID) ([]models.Category, error)
	UpdateCategory(ctx context.Context, transaction *models.Category) error
}
