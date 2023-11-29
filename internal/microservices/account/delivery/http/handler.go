package http

import (
	"encoding/json"
	"errors"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	genAccount "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/delivery/grpc/generated"
	"github.com/google/uuid"
)

const (
	accountID = "account_id"
)

type Handler struct {
	client genAccount.AccountServiceClient
	logger logger.Logger
}

func NewHandler(
	client genAccount.AccountServiceClient,
	log logger.Logger) *Handler {
	return &Handler{
		client: client,
		logger: log,
	}
}

// @Summary			Create account
// @Tags			Account
// @Description		Create account
// @Produce			json
// @Param			account		body		CreateAccount		true			"Input account create"
// @Success			200			{object}	Response[AccountCreateResponse]		"Create account"
// @Failure			400			{object}	ResponseError						"Client error"
// @Failure     	401    		{object}  	ResponseError  						"Unauthorized user"
// @Failure     	403    		{object}  	ResponseError  						"Forbidden user"
// @Failure			500			{object}	ResponseError						"Server error"
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

	account, err := h.client.Create(r.Context(), &genAccount.CreateRequest{
		UserId:         user.ID.String(),
		Balance:        float32(accountInput.Balance),
		Accumulation:   accountInput.Accumulation,
		BalanceEnabled: accountInput.BalanceEnabled,
		MeanPayment:    accountInput.MeanPayment,
	})
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, AccountNotCreate, h.logger)
		return
	}

	accountID, _ := uuid.Parse(account.AccountId)
	accountResponse := AccountCreateResponse{AccountID: accountID}
	commonHttp.SuccessResponse(w, http.StatusOK, accountResponse)

}

// @Summary		PUT 	Update
// @Tags				Account
// @Description	Put 	account
// @Produce		json
// @Param		account	body		UpdateAccount		true		    "Input transactin update"
// @Success		200		{object}	Response[NilBody]				    "Update account"
// @Failure		400		{object}	ResponseError						"Client error"
// @Failure     401    	{object}  	ResponseError  						"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  						"Forbidden user"
// @Failure		500		{object}	ResponseError						"Server error"
// @Router		/api/account/update [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	var updateAccountInput UpdateAccount

	if err := json.NewDecoder(r.Body).Decode(&updateAccountInput); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	if err := updateAccountInput.CheckValid(); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	if _, err := h.client.Update(r.Context(), &genAccount.UpdasteRequest{
		Id:             updateAccountInput.ID.String(),
		UserId:         user.ID.String(),
		Balance:        float32(updateAccountInput.Balance),
		Accumulation:   updateAccountInput.Accumulation,
		BalanceEnabled: updateAccountInput.BalanceEnabled,
		MeanPayment:    updateAccountInput.MeanPayment,
	}); err != nil {
		var errNoSuchaccount *models.NoSuchAccounts
		if errors.As(err, &errNoSuchaccount) {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, AccountNotSuch, h.logger)
			return
		}

		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponse(w, http.StatusForbidden, err, commonHttp.ForbiddenUser, h.logger)
			return
		}

		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, AccountCreateServerError, h.logger)
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

	_, err = h.client.Delete(r.Context(), &genAccount.DeleteRequest{
		AccountId: accountID.String(),
		UserId:    user.ID.String(),
	})

	if err != nil {
		var errNoSuchaccount *models.NoSuchAccounts
		if errors.As(err, &errNoSuchaccount) {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, AccountNotSuch, h.logger)
			return
		}

		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponse(w, http.StatusForbidden, err, commonHttp.ForbiddenUser, h.logger)
			return
		}
		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, AccountCreateServerError, h.logger)
			return
		}
	}
	commonHttp.SuccessResponse(w, http.StatusOK, commonHttp.NilBody{})
}
