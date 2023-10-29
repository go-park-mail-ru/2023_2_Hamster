package http

import (
	"encoding/json"
	"errors"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/transaction"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/delivery/http/transfer_models"
)

type Handler struct {
	transactionService transaction.Usecase
	logger             logger.CustomLogger
}

const (
	userIdUrlParam    = "userID"
	userloginUrlParam = "login"
)

func NewHandler(uu transaction.Usecase, l logger.CustomLogger) *Handler {
	return &Handler{
		transactionService: uu,
		logger:             l,
	}
}

// @Summary		Get all transaction
// @Tags			Transaction
// @Description	Get User all transaction
// @Produce		json
// @Success		200		{object}	Response[MasTransaction]	"Show transaction"
// @Success		204		{object}	Response[string]	     	"Show actual accounts"
// @Failure		400		{object}	ResponseError			"Client error"
// @Failure     401    	{object}  	ResponseError  		"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  		"Forbidden user"
// @Failure		500		{object}	ResponseError			"Server error"
// @Router		/api/transaction/{userID}/all [get]
func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userID, err := commonHttp.GetIDFromRequest(userIdUrlParam, r)

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidURLParameter, h.logger)
		return
	}

	dataFeed, err := h.transactionService.GetFeed(r.Context(), userID)

	var errNoSuchTransaction *models.NoSuchTransactionError
	if errors.As(err, &errNoSuchTransaction) {
		commonHttp.SuccessResponse(w, http.StatusNoContent, "")
		return
	}

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserFeedServerError, h.logger)
		return
	}

	var dataResponse []models.TransactionTransfer

	for _, transaction := range dataFeed {
		dataResponse = append(dataResponse, models.InitTransactionTransfer(transaction))
	}

	response := MasTransaction{Transactions: dataResponse}
	commonHttp.SuccessResponse(w, http.StatusOK, response)

}

// @Summary		Create transaction
// @Tags			Transaction
// @Description	Create transaction
// @Produce		json
// @Param			transaction		body		CreateTransaction		true		"Input transactin create"
// @Success		200		{object}	Response[TransactionCreateResponse]	"Create transaction"
// @Failure		400		{object}	ResponseError			"Client error"
// @Failure     401    	{object}  	ResponseError  		"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  		"Forbidden user"
// @Failure		500		{object}	ResponseError			"Server error"
// @Router		/api/transaction/create [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	var transactionInput CreateTransaction

	if err := json.NewDecoder(r.Body).Decode(&transactionInput); err != nil {
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
