package http

import (
	"encoding/json"
	"net/http"

	contextutils "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/context_utils"
	response "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Handler struct {
	cu  category.Usecase
	log logger.Logger
}

func NewHandler(cu category.Usecase, log logger.Logger) *Handler {
	return &Handler{
		cu:  cu,
		log: log,
	}
}

// @Summary		Create Tag
// @Tags			Category
// @Description	Creates tag
// @Accept 		json
// @Produce		json
// @Param			tag		body		category.TagInput			true			"tag info"
// @Success		200		{object}	Response[uuid.UUID]					"tag with id creted"
// @Failure		400		{object}	ResponseError					"Incorrect Input"
// @Failure		401		{object}	ResponseError					"auth error relogin"
// @Failure		429		{object}	ResponseError					"Server error"
// @Router		/api/addTag	[post]
func (h *Handler) CreateTag(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	var tag category.TagInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tag); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error Corupted request body: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	tag.UserId = user.ID

	id, err := h.cu.CreateTag(r.Context(), tag)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't crate tag", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusOK, id)
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
// @Router		/api/getTags	[post]
func (h *Handler) GetTags(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error: %v", err)
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	tags, err := h.cu.GetTags(r.Context(), user.ID)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] get tags Error: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't get tags", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusOK, tags)
}

// @Summary		Update Tag
// @Tags			Category
// @Description	Get all tags for user
// @Accept 		json
// @Produce		json
// @Param			tag		body		models.Category		true		"tag info"
// @Success		200		{object}	Response[models.Category]		"tag slice"
// @Failure		400		{object}	ResponseError					"Incorrect Input"
// @Failure		401		{object}	ResponseError					"auth error relogin"
// @Failure		429		{object}	ResponseError					"Server error"
// @Router		/api/getTags	[post]
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

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tag); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error Corupted request body: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	tag.UserID = user.ID

	if err := h.cu.UpdateTag(r.Context(), &tag); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Update Error: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Can't Update tag", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusOK, tag)
}

func (h *Handler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error: %v", err)
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	var tagId string

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tagId); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error Corupted request body: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	tagUuid, err := uuid.Parse(tagId)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] uuid format Error: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corupted data uuid format mismatch", h.log)
		return
	}

	if err := h.cu.DeleteTag(r.Context(), tagUuid, user.ID); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Delete Error: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Can't Delete tag", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusBadRequest, tagId)
}
