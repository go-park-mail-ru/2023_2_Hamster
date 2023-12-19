package http

import (
	"github.com/mailru/easyjson"
	"net/http"

	contextutils "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/context_utils"
	response "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category"
	genCategory "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Handler struct {
	client genCategory.CategoryServiceClient
	log    logger.Logger
}

func NewHandler(client genCategory.CategoryServiceClient, log logger.Logger) *Handler {
	return &Handler{
		client: client,
		log:    log,
	}
}

// @Summary		Create Tag
// @Tags			Category
// @Description	Creates tag
// @Accept 		json
// @Produce		json
// @Param			tag		body		category.TagInput		true		"tag info"
// @Success		200		{object}	Response[category.CategoryCreateResponse]				"tag with id creted"
// @Failure		400		{object}	ResponseError					"Incorrect Input"
// @Failure		401		{object}	ResponseError					"auth error relogin"
// @Failure		429		{object}	ResponseError					"Server error"
// @Router		/api/tag/create	[post]
func (h *Handler) CreateTag(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	var tag category.TagInput

	if err := easyjson.UnmarshalFromReader(r.Body, &tag); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error Corupted request body: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	tag.UserId = user.ID

	id, err := h.client.CreateTag(r.Context(), &genCategory.CreateTagRequest{
		UserId:      tag.UserId.String(),
		ParentId:    tag.ParentId.String(),
		Name:        tag.Name,
		ShowIncome:  tag.ShowIncome,
		ShowOutcome: tag.ShowOutcome,
		Regular:     tag.Regular,
	})
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't crate tag", h.log)
		return
	}
	tagID, _ := uuid.Parse(id.TagId)
	categoryResponse := category.CategoryCreateResponse{CategoryID: tagID}
	response.SuccessResponse(w, http.StatusOK, categoryResponse)
}

// @Summary		Get Tags
// @Tags			Category
// @Description	Get all tags for user
// @Accept 		json
// @Produce		json
// @Success		200		{object}	Response[[]models.Category]		"tag slice"
// @Failure		400		{object}	ResponseError					"Incorrect Input"
// @Failure		401		{object}	ResponseError					"auth error relogin"
// @Failure		429		{object}	ResponseError					"Server error"
// @Router		/api/tag/all	[get]
func (h *Handler) GetTags(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error: %v", err)
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	gtags, err := h.client.GetTags(r.Context(), &genCategory.UserIdRequest{UserId: user.ID.String()})
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] get tags Error: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't get tags", h.log)
		return
	}

	tags := make([]models.Category, len(gtags.Categories))

	for i, gtag := range gtags.Categories {
		id, _ := uuid.Parse(gtag.Id)
		userID, _ := uuid.Parse(gtag.UserId)
		parentID, _ := uuid.Parse(gtag.ParentId)

		tag := models.Category{
			ID:          id,
			UserID:      userID,
			ParentID:    parentID,
			Name:        gtag.Name,
			ShowIncome:  gtag.ShowIncome,
			ShowOutcome: gtag.ShowOutcome,
			Regular:     gtag.Regular,
		}
		tags[i] = tag
	}

	response.SuccessResponse(w, http.StatusOK, tags)
}

// @Summary		Update Tag
// @Tags			Category
// @Description	Update Tag
// @Accept 		json
// @Produce		json
// @Param			tag		body		models.Category		true		"tag info"
// @Success		200		{object}	Response[models.Category]		"tag to update"
// @Failure		400		{object}	ResponseError					"Incorrect Input"
// @Failure		401		{object}	ResponseError					"auth error relogin"
// @Failure		429		{object}	ResponseError					"Server error"
// @Router		/api/tag/{tagId}/update	[put]
func (h *Handler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error: %v", err)
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	var tag models.Category

	if err := easyjson.UnmarshalFromReader(r.Body, &tag); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error Corupted request body: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	tag.UserID = user.ID

	upd, err := h.client.UpdateTag(r.Context(), &genCategory.Category{
		Id:          tag.ID.String(),
		UserId:      tag.UserID.String(),
		ParentId:    tag.ParentID.String(),
		Name:        tag.Name,
		ShowIncome:  tag.ShowIncome,
		ShowOutcome: tag.ShowOutcome,
		Regular:     tag.Regular,
	})

	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Update Error: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Can't Update tag", h.log)
		return
	}

	tagID, _ := uuid.Parse(upd.Id)
	tagUserID, _ := uuid.Parse(upd.UserId)
	tagParentID, _ := uuid.Parse(upd.ParentId)

	tag.ID = tagID
	tag.UserID = tagUserID
	tag.ParentID = tagParentID
	tag.Name = upd.Name
	tag.ShowIncome = upd.ShowIncome
	tag.ShowOutcome = upd.ShowOutcome
	tag.Regular = upd.Regular

	response.SuccessResponse(w, http.StatusOK, tag)
}

// @Summary		Delete Tag
// @Tags			Category
// @Description	delete tag
// @Accept 		json
// @Produce		json
// @Param			tag		body		string				true		"tag id"
// @Success		200		{object}	Response[models.Category]		"tag slice"
// @Failure		400		{object}	ResponseError					"Incorrect Input"
// @Failure		401		{object}	ResponseError					"auth error relogin"
// @Failure		429		{object}	ResponseError					"Server error"
// @Router		/api/tag/delete	[delete]
func (h *Handler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error: %v", err)
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	var tagId category.TagDeleteInput

	if err := easyjson.UnmarshalFromReader(r.Body, &tagId); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error Corupted request body: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	_, err = h.client.DeleteTag(r.Context(), &genCategory.DeleteRequest{TagId: tagId.ID.String(), UserId: user.ID.String()})
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Delete Error: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Can't Delete tag", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusOK, tagId)
}
