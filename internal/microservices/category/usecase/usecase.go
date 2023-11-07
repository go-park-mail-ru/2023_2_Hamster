package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase struct {
	categoryRepo category.Repository
	log          logger.Logger
}

func NewUsecase(cr category.Repository, log logger.Logger) *Usecase {
	return &Usecase{
		categoryRepo: cr,
		log:          log,
	}
}

func (u *Usecase) CreateTag(ctx context.Context, tag category.TagInput) (uuid.UUID, error) {
	ok, err := u.categoryRepo.CheckNameUniq(ctx, tag.UserId, tag.ParentId, tag.Name)
	if err != nil {
		return uuid.Nil, fmt.Errorf("[usecase] Error check uniq uniqueness: %w", err)
	}
	if !ok {
		return uuid.Nil, fmt.Errorf("[usecase] Error category alredy exist")
	}

	var newTag models.Category

	newTag.Name = tag.Name
	newTag.ParentID = tag.ParentId
	newTag.ShowIncome = tag.ShowIncome
	newTag.ShowIncome = tag.ShowIncome
	newTag.UserID = tag.UserId

	id, err := u.categoryRepo.CreateTag(ctx, newTag)
	if err != nil {
		return uuid.Nil, fmt.Errorf("[usecase] Error tag creation: %v", err)
	}

	return id, nil
}

func (u *Usecase) UpdateTag(ctx context.Context, tag *models.Category) error {
	ok, err := u.categoryRepo.CheckExist(ctx, tag.UserID, tag.ID)
	if err != nil {
		return fmt.Errorf("[usecase] Error in check tag existens: %v", err)
	}
	if !ok {
		return fmt.Errorf("[usecase] Error tag doesn't exist can't update")
	}

	if err := u.categoryRepo.UpdateTag(ctx, tag); err != nil {
		return fmt.Errorf("[usecase] update tag Error: %v", err)
	}
	return nil
}

func (u *Usecase) DeleteTag(ctx context.Context, tagId uuid.UUID, userId uuid.UUID) error {
	ok, err := u.categoryRepo.CheckExist(ctx, userId, tagId)
	if err != nil {
		return fmt.Errorf("[usecase] Error in check tag existens: %v", err)
	}
	if !ok {
		return fmt.Errorf("[usecase] Error tag doesn't exist can't delete")
	}

	if err := u.categoryRepo.DeleteTag(ctx, tagId); err != nil {
		return fmt.Errorf("[usecase] Error in tag deletion: %v", err)
	}
	return nil
}

func (u *Usecase) GetTags(ctx context.Context, userId uuid.UUID) ([]models.Category, error) {
	tags, err := u.categoryRepo.GetTags(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("[usecase] Error getting user tags: %v", err)
	}

	return tags, nil
}
