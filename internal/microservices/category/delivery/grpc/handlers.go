package grpc

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category"
	proto "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type categoryGRPC struct {
	CategoryServices category.Usecase
	logger           logger.Logger

	proto.UnimplementedCategoryServiceServer
}

func NewCategoryGRPC(categoryServices category.Usecase, logger logger.Logger) *categoryGRPC {
	return &categoryGRPC{
		CategoryServices: categoryServices,
		logger:           logger,
	}
}

func (c *categoryGRPC) CreateTag(ctx context.Context, in *proto.CreateTagRequest) (*proto.CreateTagResponse, error) {
	userId, _ := uuid.Parse(in.UserId)
	ParentID, _ := uuid.Parse(in.ParentId)
	request := category.TagInput{
		UserId:      userId,
		ParentId:    ParentID,
		Name:        in.Name,
		ShowIncome:  in.ShowIncome,
		ShowOutcome: in.ShowOutcome,
		Regular:     in.Regular,
		Image:       in.Image,
	}
	id, err := c.CategoryServices.CreateTag(ctx, request)

	return &proto.CreateTagResponse{TagId: id.String()}, err
}

func (c *categoryGRPC) GetTags(ctx context.Context, in *proto.UserIdRequest) (*proto.GetTagsResponse, error) {
	userId, _ := uuid.Parse(in.UserId)

	tags, err := c.CategoryServices.GetTags(ctx, userId)
	if err != nil {
		var errNoTags *models.ErrNoTags
		if errors.As(err, &errNoTags) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	var generatedCategories []*proto.Category

	for _, tag := range tags {
		generatedCategory := &proto.Category{
			Id:          tag.ID.String(),
			UserId:      tag.UserID.String(),
			ParentId:    tag.ParentID.String(),
			Name:        tag.Name,
			ShowIncome:  tag.ShowIncome,
			ShowOutcome: tag.ShowOutcome,
			Regular:     tag.Regular,
			Image:       tag.Image,
		}

		generatedCategories = append(generatedCategories, generatedCategory)
	}

	return &proto.GetTagsResponse{Categories: generatedCategories}, nil
}

func (c *categoryGRPC) UpdateTag(ctx context.Context, in *proto.Category) (*proto.Category, error) {
	cId, _ := uuid.Parse(in.Id)
	cUserId, _ := uuid.Parse(in.UserId)
	cParentId, _ := uuid.Parse(in.ParentId)
	tag := &models.Category{
		ID:          cId,
		UserID:      cUserId,
		ParentID:    cParentId,
		Name:        in.Name,
		ShowIncome:  in.ShowIncome,
		ShowOutcome: in.ShowOutcome,
		Regular:     in.Regular,
		Image:       in.Image,
	}
	err := c.CategoryServices.UpdateTag(ctx, tag)

	updatedProtoCategory := &proto.Category{
		Id:          tag.ID.String(),
		UserId:      tag.UserID.String(),
		ParentId:    tag.ParentID.String(),
		Name:        tag.Name,
		ShowIncome:  tag.ShowIncome,
		ShowOutcome: tag.ShowOutcome,
		Regular:     tag.Regular,
		Image:       int32(tag.Image),
	}

	return updatedProtoCategory, err
}

func (c *categoryGRPC) DeleteTag(ctx context.Context, in *proto.DeleteRequest) (*empty.Empty, error) {
	cId, _ := uuid.Parse(in.TagId)
	cUserId, _ := uuid.Parse(in.UserId)

	err := c.CategoryServices.DeleteTag(ctx, cId, cUserId)

	return &empty.Empty{}, err
}
