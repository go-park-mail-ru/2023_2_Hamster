package http

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	response "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	genAccount "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/delivery/http/transfer_models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	transactionService transaction.Usecase
	userService        user.Usecase
	tagService         generated.CategoryServiceClient
	client             genAccount.AccountServiceClient
	logger             logger.Logger
}

func NewHandler(uu transaction.Usecase,
	userUsecase user.Usecase,
	tagCli generated.CategoryServiceClient,
	cl genAccount.AccountServiceClient,
	l logger.Logger,
) *Handler {
	return &Handler{
		transactionService: uu,
		tagService:         tagCli,
		userService:        userUsecase,
		client:             cl,
		logger:             l,
	}
}

// @Summary		Export .csv Transactions
// @Tags		Transaction
// @Description	Sends a .csv file with transactions based on the specified criteria.
// @Produce		plain
// @Success     200     {string}    "Successfully exported transactions"   {example: "TransactionID,Amount,Date\n1,100,2023-01-01\n2,150,2023-01-02\n"}
// @Failure		400		{object}	ResponseError	"Bad request - Transaction error"
// @Failure		401		{object}	ResponseError	"Unauthorized - User unauthorized"
// @Failure		403		{object}	ResponseError	"Forbidden - User doesn't have rights"
// @Failure		404		{object}	ResponseError	"Not Found - No transactions found for the specified criteria"
// @Failure		500		{object}	ResponseError	"Internal Server Error - Server error"
// @Router		/api/transaction/export [get]
// @Param		startDate	query	string	true	"Start date (format: 'YYYY-MM-DD')"
// @Param		endDate		query	string	true	"End date (format: 'YYYY-MM-DD')"
// @Param		authorization	header	string	true	"session_id"
func (h *Handler) ExportTransactions(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.logger)
		return
	}

	query, err := response.GetQueryParam(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, response.InvalidURLParameter, h.logger)
		return
	}

	var errNoSuchTransaction *models.NoSuchTransactionError
	dataFeed, err := h.transactionService.GetTransactionForExport(r.Context(), user.ID, query)
	if errors.As(err, &errNoSuchTransaction) {
		response.ErrorResponse(w, http.StatusNotFound, err, "no transactions found", h.logger)
		return
	} else if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserFeedServerError, h.logger)
		return
	}

	fileName := "Transactions_" + user.Login + query.StartDate.String() + "-" + query.EndDate.String() + ".csv"

	// 1. Create a CSV file and write the `dataFeed` into it.
	csvFile, err := os.Create(fileName)
	if err != nil {
		h.logger.Error("can't crete datafile")
		response.ErrorResponse(w, http.StatusInternalServerError, err, "can't create .csv file", h.logger)
		return
	}
	defer func() {
		// 7. Close the CSV file.
		if err := csvFile.Close(); err != nil {
			h.logger.Errorf("Error closing CSV file: %v", err)
		}

		// 8. Delete the CSV file after serving it.
		if err := os.Remove(fileName); err != nil {
			h.logger.Errorf("Error deleting CSV file: %v", err)
		}
	}()

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

	// Set the appropriate headers for a CSV file.
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))

	// Seek to the start of the CSV file.
	_, err = csvFile.Seek(0, 0)
	if err != nil {
		h.logger.Errorf("Error in csv set start")
	}

	// Write the CSV file directly to the response.
	_, err = io.Copy(w, csvFile)
	if err != nil {
		h.logger.Errorf("Error in csv writing")
	}
}

// @Summary 	Export Transactions from CSV
// @Tags 		Transaction
// @Description `Uploads a CSV file containing transactions and processes them to be stored in the system.
// @Accept  mult`ipart/form-data
// @Produce json
// @Param 	csv formData file true "CSV file containing transactions data"
// @Success 200 {string} string "Successfully imported transactions"
// @Failure 400 {object} ResponseError "Bad request - Transaction error"
// @Failure 401 {object} ResponseError "Unauthorized - User unauthorized"
// @Failure 403 {object} ResponseError "Forbidden - User doesn't have rights"
// @Failure 404 {object} ResponseError "Not Found - No transactions found for the specified criteria"
// @Failure 413 {object} ResponseError "Request Entity Too Large - File is too large"
// @Failure 500 {object} ResponseError "Internal Server Error - Server error"
// @Router /api/transaction/import [post]
func (h *Handler) ImportTransactions(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.logger)
		return
	}

	// Parse the multipart form in the request
	err = r.ParseMultipartForm(15 << 20) // Max memory 10MB
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, "File too big or contains errors", h.logger)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("csvFile")
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, "Error getting the file", h.logger)
		return
	}
	defer file.Close()

	// Create a new CSV reader reading from the file
	reader := csv.NewReader(file)

	var errNoSuchAccounts *models.NoSuchAccounts

	accounts, err := h.userService.GetAccounts(r.Context(), user.ID)
	if err != nil && !errors.As(err, &errNoSuchAccounts) {
		response.ErrorResponse(w, http.StatusInternalServerError, err, "Error getting accounts", h.logger)
		return
	}

	accountCache := sync.Map{}
	if len(accounts) != 0 {
		for _, account := range accounts {
			accountCache.Store(account.MeanPayment, account.ID)
		}
	}

	gtags, err := h.tagService.GetTags(r.Context(), &generated.UserIdRequest{
		UserId: user.ID.String(),
	})
	if err != nil && status.Code(err) != codes.NotFound {
		response.ErrorResponse(w, http.StatusInternalServerError, err, "Error getting accounts", h.logger)
		return
	}

	var tags []models.Category

	if len(gtags.Categories) != 0 {
		tags = make([]models.Category, len(gtags.Categories))

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
				Image:       gtag.Image,
			}
			tags[i] = tag
		}
	}

	tagCache := sync.Map{}
	if len(tags) != 0 {
		for _, tag := range tags {
			tagCache.Store(tag.Name, tag.ID)
		}
	}

	var i int
	// Iterate through the records
	for i = 0; true; i++ {
		// Read each record from csv
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			response.ErrorResponse(w, http.StatusBadRequest, err, "Error reading the CSV file", h.logger)
			return
		}

		if len(record) < 8 {
			continue
			// response.ErrorResponse(w, http.StatusBadRequest, err, "wrong format not enough args for transaction", h.logger)
			// return
		}

		if len(record) > 9 {
			continue
			// response.ErrorResponse(w, http.StatusBadRequest, err, "wrong format too much args for transaction", h.logger)
			// return
		}

		accountIncome := record[0]
		accountOutcome := record[1]

		income, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			response.ErrorResponse(w, http.StatusBadRequest, err, "Error converting the income amount to float", h.logger)
			return
		}

		outcome, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			response.ErrorResponse(w, http.StatusBadRequest, err, "Error converting the outcome amount to float", h.logger)
			return
		}

		if (income == 0) == (outcome == 0) {
			response.ErrorResponse(w, http.StatusBadRequest, err, "transaction can be consumables or profitable", h.logger)
			return
		}

		date, err := time.Parse(time.RFC3339, record[4])
		if err != nil {
			response.ErrorResponse(w, http.StatusBadRequest, err, "Error wrong time format", h.logger)
			return
		}

		payer := record[5]
		description := record[6]

		tagName := record[7]

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
				response.ErrorResponse(w, http.StatusInternalServerError, err, "Import error account add failed", h.logger)
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
				response.ErrorResponse(w, http.StatusInternalServerError, err, "Import error account add failed", h.logger)
				return
			}
			accountOutcomeId, _ = uuid.Parse(account.AccountId)
			accountCache.Store(accountOutcome, accountOutcomeId)
		}

		var tagId uuid.UUID
		if value, ok := tagCache.Load(tagName); ok {
			tagId = value.(uuid.UUID)
		} else {
			tag, err := h.tagService.CreateTag(r.Context(), &generated.CreateTagRequest{
				UserId:      user.ID.String(),
				Name:        tagName,
				Image:       0,
				ShowIncome:  false,
				Regular:     false,
				ShowOutcome: true,
			})
			if err != nil {
				response.ErrorResponse(w, http.StatusInternalServerError, err, "Import error tag add failed", h.logger)
				return
			}

			tagId, _ = uuid.Parse(tag.TagId)
			tagCache.Store(tagName, tagId)
		}

		// Parse the record to a Transaction struct
		transaction := models.Transaction{
			UserID:           user.ID,
			AccountIncomeID:  accountIncomeId,
			AccountOutcomeID: accountOutcomeId,
			Income:           income,
			Outcome:          outcome,
			Date:             date,
			Payer:            payer,
			Description:      description,
		}

		transaction.Categories = append(transaction.Categories, models.CategoryName{
			ID:   tagId,
			Name: tagName,
		})

		// Create the transaction in the database
		_, err = h.transactionService.CreateTransaction(r.Context(), &transaction)
		if err != nil {
			response.ErrorResponse(w, http.StatusInternalServerError, err, "Error creating the transaction", h.logger)
			return
		}
	}

	if i == 0 {
		response.ErrorResponse(w, http.StatusBadRequest, err, "empty file", h.logger)
		return
	}

	response.SuccessResponse(w, http.StatusOK, "Successfully imported transactions")
}
