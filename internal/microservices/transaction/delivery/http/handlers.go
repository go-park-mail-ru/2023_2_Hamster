package http

import (
	"errors"
	"net/http"

	"github.com/mailru/easyjson"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/delivery/http/transfer_models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction"
)

type Handler struct {
	transactionService transaction.Usecase
	userService        user.Usecase
	logger             logger.Logger
}

const (
	transactionID = "transaction_id"

	// userIdUrlParam    = "userID"
	// userloginUrlParam = "login"
)

func NewHandler(uu transaction.Usecase, user user.Usecase, l logger.Logger) *Handler {
	return &Handler{
		transactionService: uu,
		userService:        user,
		logger:             l,
	}
}

// @Summary		Get count transaction
// @Tags		Transaction
// @Description	Get User count transaction
// @Produce		json
// @Success		200		{object}	Response[TransactionCount] "Show transaction count"
// @Failure		400		{object}	ResponseError			 "Client error"
// @Failure     401    	{object}    ResponseError  			 "Unauthorized user"
// @Failure     403    	{object}    ResponseError  			 "Forbidden user"
// @Failure		500		{object}	ResponseError			 "Server error"
// @Router		/api/transaction/count [get]
func (h *Handler) GetCount(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	transactionCount, err := h.transactionService.GetCount(r.Context(), user.ID)

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, "can't get count transaction info", h.logger)
		return
	}

	response := TransactionCount{Count: transactionCount}
	commonHttp.SuccessResponse(w, http.StatusOK, response)

}

// @Summary		Get all transaction
// @Tags		Transaction
// @Description	Get User all transaction
// @Produce		json
// @Param       request query       models.QueryListOptions false   "Query Params"
// @Success		200		{object}	Response[MasTransaction] "Show transaction"
// @Success		204		{object}	Response[string]	     "Show actual accounts"
// @Failure		400		{object}	ResponseError			 "Client error"
// @Failure     401    	{object}    ResponseError  			 "Unauthorized user"
// @Failure     403    	{object}    ResponseError  			 "Forbidden user"
// @Failure		500		{object}	ResponseError			 "Server error"
// @Router		/api/transaction/feed [get]
func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	query, err := commonHttp.GetQueryParam(r)
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidURLParameter, h.logger)
		return
	}
	dataFeed, err := h.transactionService.GetFeed(r.Context(), user.ID, query)

	var errNoSuchTransaction *models.NoSuchTransactionError
	if errors.As(err, &errNoSuchTransaction) {
		commonHttp.SuccessResponse(w, http.StatusNoContent, commonHttp.NilBody{})
		return
	}

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserFeedServerError, h.logger)
		return
	}

	var dataResponse []models.TransactionTransfer

	var userT *models.User
	for _, transaction := range dataFeed {
		userT, err = h.userService.GetUser(r.Context(), transaction.UserID)
		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserFeedServerError, h.logger)
		}
		dataResponse = append(dataResponse, models.InitTransactionTransfer(transaction, userT.Login))
	}

	response := MasTransaction{Transactions: dataResponse}
	commonHttp.SuccessResponse(w, http.StatusOK, response)

}

// @Summary		Create transaction
// @Tags			Transaction
// @Description	Create transaction
// @Produce		json
// @Param			transaction		body		CreateTransaction		true		"Input transactin create"
// @Success		200		{object}	Response[TransactionCreateResponse]				"Create transaction"
// @Failure		400		{object}	ResponseError									"Client error"
// @Failure     401    	{object}  	ResponseError  									"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  									"Forbidden user"
// @Failure		500		{object}	ResponseError									"Server error"
// @Router		/api/transaction/create [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	var transactionInput CreateTransaction

	if err := easyjson.UnmarshalFromReader(r.Body, &transactionInput); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	if err := transactionInput.CheckValid(); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	transactionID, err := h.transactionService.CreateTransaction(r.Context(), transactionInput.ToTransaction(user))
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, TransactionNotCreate, h.logger)
		return
	}

	transactionResponse := TransactionCreateResponse{TransactionID: transactionID}
	commonHttp.SuccessResponse(w, http.StatusOK, transactionResponse)
}

// @Summary		PUT Update
// @Tags			Transaction
// @Description	Put transaction
// @Produce		json
// @Param			transaction		body		UpdTransaction		true		    "Input transactin update"
// @Success		200		{object}	Response[NilBody]				                "Update transaction"
// @Failure		400		{object}	ResponseError									"Client error"
// @Failure     401    	{object}  	ResponseError  									"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  									"Forbidden user"
// @Failure		500		{object}	ResponseError									"Server error"
// @Router		/api/transaction/update [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	var updTransactionInput UpdTransaction

	if err := easyjson.UnmarshalFromReader(r.Body, &updTransactionInput); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	if err := updTransactionInput.CheckValid(); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	if err := h.transactionService.UpdateTransaction(r.Context(), updTransactionInput.ToTransaction(user)); err != nil {
		var errNoSuchTransaction *models.NoSuchTransactionError
		if errors.As(err, &errNoSuchTransaction) {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, TransactionNotSuch, h.logger)
			return
		}

		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponse(w, http.StatusForbidden, err, commonHttp.ForbiddenUser, h.logger)
			return
		}

		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, TransactionCreateServerError, h.logger)
			return
		}
	}

	commonHttp.SuccessResponse(w, http.StatusOK, commonHttp.NilBody{})

}

// @Summary		Delete Transaction
// @Tags		Transaction
// @Description	Delete transaction with chosen ID
// @Produce		json
// @Success		200		{object}	Response[NilBody]	  	    "Transaction deleted"
// @Failure		400		{object}	ResponseError				"Transaction error"
// @Failure		401		{object}	ResponseError  			    "User unathorized"
// @Failure		403		{object}	ResponseError				"User hasn't rights"
// @Failure		500		{object}	ResponseError				"Server error"
// @Router		/api/transaction/{transaction_id}/delete [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	transactionID, err := commonHttp.GetIDFromRequest(transactionID, r)

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidURLParameter, h.logger)
		return
	}
	user, err := commonHttp.GetUserFromRequest(r)

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	err = h.transactionService.DeleteTransaction(r.Context(), transactionID, user.ID)

	if err != nil {
		var errNoSuchTransaction *models.NoSuchTransactionError
		if errors.As(err, &errNoSuchTransaction) {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, TransactionNotSuch, h.logger)
			return
		}

		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponse(w, http.StatusForbidden, err, commonHttp.ForbiddenUser, h.logger)
			return
		}
		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, TransactionDeleteServerError, h.logger)
			return
		}
	}
	commonHttp.SuccessResponse(w, http.StatusOK, commonHttp.NilBody{})
}
