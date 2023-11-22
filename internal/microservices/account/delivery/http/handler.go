package http

import (
	"encoding/json"
	"errors"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account"
)

type Handler struct {
	accountService account.Usecase
	logger         logger.Logger
}

func NewHandler(au account.Usecase, l logger.Logger) *Handler {
	return &Handler{
		accountService: au,
		logger:         l,
	}
}

// @Summary			Create account
// @Tags			Account
// @Description		Create account
// @Produce			json
// @Param			account		body		CreateAccount		true		"Input account create"
// @Success			200		{object}	Response[accountCreateResponse]				"Create account"
// @Failure			400		{object}	ResponseError									"Client error"
// @Failure     	401    	{object}  	ResponseError  									"Unauthorized user"
// @Failure     	403    	{object}  	ResponseError  									"Forbidden user"
// @Failure			500		{object}	ResponseError									"Server error"
// @Router		/api/account/create [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	var accountInput CreateAccount

	if err := json.NewDecoder(r.Body).Decode(&accountInput); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	if err := accountInput.CheckValid(); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	accountID, err := h.accountService.CreateAccount(r.Context(), accountInput.ToAccount())
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, accountNotCreate, h.logger)
		return
	}

	accountResponse := accountCreateResponse{accountID: accountID}
	commonHttp.SuccessResponse(w, http.StatusOK, accountResponse)

}

// @Summary		PUT Update
// @Tags			Account
// @Description	Put account
// @Produce		json
// @Param			account		body		UpdateAccount		true		    "Input transactin update"
// @Success		200		{object}	Response[NilBody]				                "Update account"
// @Failure		400		{object}	ResponseError									"Client error"
// @Failure     401    	{object}  	ResponseError  									"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  									"Forbidden user"
// @Failure		500		{object}	ResponseError									"Server error"
// @Router		/api/account/update [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	var updaccountInput Updaccount

	if err := json.NewDecoder(r.Body).Decode(&updaccountInput); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	if err := updaccountInput.CheckValid(); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	if err := h.accountService.Updateaccount(r.Context(), updaccountInput.Toaccount(user)); err != nil {
		var errNoSuchaccount *models.NoSuchaccountError
		if errors.As(err, &errNoSuchaccount) {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, accountNotSuch, h.logger)
			return
		}

		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponse(w, http.StatusForbidden, err, commonHttp.ForbiddenUser, h.logger)
			return
		}

		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, accountCreateServerError, h.logger)
			return
		}
	}

	commonHttp.SuccessResponse(w, http.StatusOK, commonHttp.NilBody{})

}

// @Summary		Delete Account
// @Tags		Account
// @Description	Delete account with chosen ID
// @Produce		json
// @Success		200		{object}	Response[NilBody]	  	    "Account deleted"
// @Failure		400		{object}	ResponseError				"Account error"
// @Failure		401		{object}	ResponseError  			    "User unathorized"
// @Failure		403		{object}	ResponseError				"User hasn't rights"
// @Failure		500		{object}	ResponseError				"Server error"
// @Router		/api/account/{account_id}/delete [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	accountID, err := commonHttp.GetIDFromRequest(accountID, r)

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidURLParameter, h.logger)
		return
	}
	user, err := commonHttp.GetUserFromRequest(r)

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	err = h.accountService.Deleteaccount(r.Context(), accountID, user.ID)

	if err != nil {
		var errNoSuchaccount *models.NoSuchaccountError
		if errors.As(err, &errNoSuchaccount) {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, accountNotSuch, h.logger)
			return
		}

		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponse(w, http.StatusForbidden, err, commonHttp.ForbiddenUser, h.logger)
			return
		}
		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, accountDeleteServerError, h.logger)
			return
		}
	}
	commonHttp.SuccessResponse(w, http.StatusOK, commonHttp.NilBody{})

}
