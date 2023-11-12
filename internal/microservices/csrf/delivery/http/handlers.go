package http

import (
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/csrf"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
)

type Handler struct {
	tokenServices csrf.Usecase
	logger        logger.Logger
}

func NewHandler(tu csrf.Usecase, l logger.Logger) *Handler {
	return &Handler{
		tokenServices: tu,
		logger:        l,
	}
}

// @Summary		Get csrf token
// @Tags			Csrf
// @Description	Get csrf token
// @Produce		json
// @Success		200		{object}	Response[getCSRFResponce]	"Csrf Token"
// @Failure     	401    	{object}  	ResponseError  				"Unauthorized user"
// @Failure		500		{object}	ResponseError				"Server error"
// @Router		/api/csrf/ [get]
func (h *Handler) GetCSRF(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	token, err := h.tokenServices.GenerateCSRFToken(user.ID)
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, csrfGetError, h.logger)
		return
	}

	resp := getCSRFResponce{CSRF: token}
	commonHttp.SuccessResponse(w, http.StatusOK, resp)
}
