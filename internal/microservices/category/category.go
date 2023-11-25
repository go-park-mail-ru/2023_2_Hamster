package category

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase interface {
	CreateTag(ctx context.Context, tag TagInput) (uuid.UUID, error)
	UpdateTag(ctx context.Context, tag *models.Category) error
	DeleteTag(ctx context.Context, tagId uuid.UUID, userId uuid.UUID) error
	GetTags(ctx context.Context, userId uuid.UUID) ([]models.Category, error)
}

type Repository interface {
	CreateTag(ctx context.Context, category models.Category) (uuid.UUID, error)
	UpdateTag(ctx context.Context, tag *models.Category) error
	DeleteTag(ctx context.Context, tagId uuid.UUID) error
	GetTags(ctx context.Context, userId uuid.UUID) ([]models.Category, error)

	CheckNameUniq(ctx context.Context, userId uuid.UUID, parentId uuid.UUID, name string) (bool, error)
	CheckExist(ctx context.Context, userId uuid.UUID, tagId uuid.UUID) (bool, error)
}
