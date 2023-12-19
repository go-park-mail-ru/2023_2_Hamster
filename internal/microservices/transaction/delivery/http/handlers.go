package http

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/mailru/easyjson"

	genAccount "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/delivery/http/transfer_models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Handler struct {
	transactionService transaction.Usecase
	userService        user.Usecase
	client             genAccount.AccountServiceClient
	logger             logger.Logger
}

const (
	transactionID = "transaction_id"

	// userIdUrlParam    = "userID"
	// userloginUrlParam = "login"
)

func NewHandler(uu transaction.Usecase, userUsecase user.Usecase, cl genAccount.AccountServiceClient, l logger.Logger) *Handler {
	return &Handler{
		transactionService: uu,
		userService:        userUsecase,
		client:             cl,
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

// @Summary		Export .csv Transactions
// @Tags		Transaction
// @Description	Sends a .csv file with transactions based on the specified criteria.
// @Produce		plain
// @Failure		400		{object}	ResponseError	"Bad request - Transaction error"
// @Failure		401		{object}	ResponseError	"Unauthorized - User unauthorized"
// @Failure		403		{object}	ResponseError	"Forbidden - User doesn't have rights"
// @Failure		404		{object}	ResponseError	"Not Found - No transactions found for the specified criteria"
// @Failure		500		{object}	ResponseError	"Internal Server Error - Server error"
// @Router		/api/transaction/export [get]
// @Param		startDate	query	string	true	"Start date (format: 'YYYY-MM-DD')"
// @Param		endDate		query	string	true	"End date (format: 'YYYY-MM-DD')"
// @Param		authorization	header	string	true	"session_id"
// @Success        200     {string}    "Successfully exported transactions"   {example: "TransactionID,Amount,Date\n1,100,2023-01-01\n2,150,2023-01-02\n"}
func (h *Handler) ExportTransactions(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	query, err := commonHttp.GetQueryParam(r)
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidURLParameter, h.logger)
		return
	}

	var errNoSuchTransaction *models.NoSuchTransactionError
	dataFeed, err := h.transactionService.GetTransactionForExport(r.Context(), user.ID, query)
	if errors.As(err, &errNoSuchTransaction) {
		commonHttp.ErrorResponse(w, http.StatusNotFound, err, "no transactions found", h.logger)
		return
	} else if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserFeedServerError, h.logger)
		return
	}

	fileName := "Transactions_" + user.Login + query.StartDate.String() + "-" + query.EndDate.String() + ".csv"

	// 1. Create a CSV file and write the `dataFeed` into it.
	csvFile, err := os.Create(fileName)
	if err != nil {
		h.logger.Error("can't crete datafile")
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, "can't create .csv file", h.logger)
		return
	}
	defer csvFile.Close()

	var csvHeader []string

	t := reflect.TypeOf(dataFeed[0])
	for i := 1; i < t.NumField(); i++ {
		field := t.Field(i)
		csvHeader = append(csvHeader, field.Name)
	}

	csvWriter := csv.NewWriter(csvFile)
	if err = csvWriter.Write(csvHeader); err != nil {
		h.logger.Errorf("Error in csv writing")
	}

	for _, row := range dataFeed {
		record := row.String()
		if err := csvWriter.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		log.Fatal(err)
	}

	// 2. Create a new multipart form writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 3. Create a new form file where the CSV file will be written.
	part, err := writer.CreateFormFile("file", "dataFeed.csv")
	if err != nil {
		log.Fatal(err)
	}

	// 4. Write the CSV file into the form file.
	_, err = csvFile.Seek(0, 0)
	if err != nil {
		h.logger.Errorf("Error in csv set start")
	}

	_, err = io.Copy(part, csvFile)
	if err != nil {
		h.logger.Errorf("Error in csv writing")
	}

	// 5. Write the multipart form's boundary to the response header.
	w.Header().Set("Content-Type", writer.FormDataContentType())

	// 6. Write the multipart form to the response.
	writer.Close()
	_, err = io.Copy(w, body)
	if err != nil {
		h.logger.Errorf("Error in csv writing")
	}
}

func (h *Handler) ImportTransactions(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	// Parse the multipart form in the request
	err = r.ParseMultipartForm(10 << 20) // Max memory 10MB
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, "Error parsing the request", h.logger)
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("csv")
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, "Error getting the file", h.logger)
		return
	}
	defer file.Close()

	if header.Header.Get("Content-Type") != "text/csv" {
		commonHttp.ErrorResponse(w, http.StatusUnsupportedMediaType, nil, "File must be a CSV", h.logger)
		return
	}

	const maxFileSize = 10 << 20 // 10MB
	if header.Size > maxFileSize {
		commonHttp.ErrorResponse(w, http.StatusRequestEntityTooLarge, nil, "File is too large", h.logger)
		return
	}
	// Create a new CSV reader reading from the file
	reader := csv.NewReader(file)

	var errNoSuchAccounts *models.NoSuchAccounts
	accounts, err := h.userService.GetAccounts(r.Context(), user.ID)
	if errors.As(err, &errNoSuchAccounts) {
		h.logger.Println(errNoSuchAccounts)
	} else if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, "Error getting accounts", h.logger)
		return
	}

	accountCache := sync.Map{}
	if len(accounts) != 0 {
		for _, account := range accounts {
			accountCache.Store(account.MeanPayment, account.ID)
		}
	}

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, "Error reading the CSV file", h.logger)
			return
		}

		accountIncome := record[0]
		accountOutcome := record[1]

		income, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, "Error converting the amount to float", h.logger)
			return
		}
		if income == 0 {
			income = 0
		}

		outcome, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, "Error converting the amount to float", h.logger)
			return
		}
		if outcome == 0 {
			outcome = 0
		}

		date, err := time.Parse(record[4], time.RFC3339Nano)
		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, "Error wrong time format", h.logger)
			return
		}

		payer := record[5]
		description := record[6]

		var accountIncomeId uuid.UUID
		if value, ok := accountCache.Load(accountIncome); ok {
			accountIncomeId = value.(uuid.UUID)
		} else {
			account, err := h.client.Create(r.Context(), &genAccount.CreateRequest{
				UserId:         user.ID.String(),
				Balance:        0.0,
				Accumulation:   true,
				BalanceEnabled: true,
				MeanPayment:    accountIncome,
			})
			if err != nil {
				commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, "Import error account add failed", h.logger)
				return
			}
			accountIncomeId, _ = uuid.Parse(account.AccountId)
			accountCache.Store(accountIncome, accountIncomeId)
		}

		var accountOutcomeId uuid.UUID
		if value, ok := accountCache.Load(accountOutcome); ok {
			accountOutcomeId = value.(uuid.UUID)
		} else {
			account, err := h.client.Create(r.Context(), &genAccount.CreateRequest{
				UserId:         user.ID.String(),
				Balance:        0.0,
				Accumulation:   true,
				BalanceEnabled: true,
				MeanPayment:    accountOutcome,
			})
			if err != nil {
				commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, "Import error account add failed", h.logger)
				return
			}
			accountOutcomeId, _ = uuid.Parse(account.AccountId)
			accountCache.Store(accountOutcome, accountOutcomeId)
		}

		// Parse the record to a Transaction struct
		transaction := models.Transaction{
			AccountIncomeID:  accountIncomeId,
			AccountOutcomeID: accountOutcomeId,
			Income:           income,
			Outcome:          outcome,
			Date:             date,
			Payer:            payer,
			Description:      description,
		}

		// Create the transaction in the database
		_, err = h.transactionService.CreateTransaction(r.Context(), &transaction)
		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, "Error creating the transaction", h.logger)
			return
		}
	}
}
