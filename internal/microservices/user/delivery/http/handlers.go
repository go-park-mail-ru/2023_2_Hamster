package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mailru/easyjson"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/delivery/http/transfer_models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

const (
	accountID = "account_id"
)

type Handler struct {
	userService user.Usecase
	logger      logger.Logger
}

func NewHandler(uu user.Usecase, l logger.Logger) *Handler {
	return &Handler{
		userService: uu,
		logger:      l,
	}
}

// @Summary		Get User
// @Tags		User
// @Description	Get user with chosen ID
// @Produce		json
// @Success		200		{object}	Response[transfer_models.UserTransfer] "Show user"
// @Failure		400		{object}	ResponseError	"Client error"
// @Failure     401    	{object}  	ResponseError  		"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  		"Forbidden user"
// @Failure		500		{object}	ResponseError	"Server error"
// @Router		/api/user/ [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) { // ?
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	usrTranfer := transfer_models.InitUserTransfer(*user)

	commonHttp.SuccessResponse(w, http.StatusOK, usrTranfer)
}

// @Summary		Get Balance
// @Tags			User
// @Description	Get User balance
// @Produce		json
// @Success		200		{object}	Response[transfer_models.BalanceResponse] "Show balance"
// @Failure		400		{object}	ResponseError	"Client error"
// @Failure     401    	{object}  	ResponseError  		"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  		"Forbidden user"
// @Failure		500		{object}	ResponseError	"Server error"
// @Router		/api/user/balance [get]
func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) { // ?
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}
	balance, err := h.userService.GetUserBalance(r.Context(), user.ID)

	var errNoSuchUserIdBalanceError *models.NoSuchUserIdBalanceError
	if errors.As(err, &errNoSuchUserIdBalanceError) {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.BalanceNotFound, h.logger)
		return
	}

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.BalanceGetServerError, h.logger)
		return
	}

	response := transfer_models.BalanceResponse{Balance: balance}
	commonHttp.SuccessResponse[transfer_models.BalanceResponse](w, http.StatusOK, response)
}

// @Summary		Get Planned Budget
// @Tags			User
// @Description	Get User planned budget
// @Produce		json
// @Success		200		{object} 	Response[transfer_models.BudgetPlannedResponse]	"Show planned budget"
// @Failure		400		{object}	ResponseError			"Client error"
// @Failure     401    	{object}  	ResponseError  			"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  			"Forbidden user"
// @Failure		500		{object}	ResponseError			"Server error"
// @Router		/api/user/plannedBudget [get]
func (h *Handler) GetPlannedBudget(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	budget, err := h.userService.GetPlannedBudget(r.Context(), user.ID)

	var errNoSuchPlannedBudgetError *models.NoSuchPlannedBudgetError
	if errors.As(err, &errNoSuchPlannedBudgetError) {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.PlannedBudgetNotFound, h.logger)
		return
	}

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.PlannedBudgetGetServerError, h.logger)
		return
	}

	response := transfer_models.BudgetPlannedResponse{BudgetPlanned: budget}
	commonHttp.SuccessResponse[transfer_models.BudgetPlannedResponse](w, http.StatusOK, response)
}

// @Summary		Get Actual Budget
// @Tags			User
// @Description	Get User actual budget
// @Produce		json
// @Success		200		{object}	Response[transfer_models.BudgetActualResponse]	"Show actual budget"
// @Failure		400		{object}	ResponseError			"Client error"
// @Failure     401    	{object}  	ResponseError  		"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  		"Forbidden user"
// @Failure		500		{object}	ResponseError			"Server error"
// @Router		/api/user/actualBudget [get]
func (h *Handler) GetCurrentBudget(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	budget, err := h.userService.GetCurrentBudget(r.Context(), user.ID)

	// var errNoSuchCurrentBudget *models.NoSuchCurrentBudget
	// if errors.As(err, &errNoSuchCurrentBudget) {
	// 	commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.CurrentBudgetNotFound, h.logger)
	// 	return
	// }

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.CurrentBudgetGetServerError, h.logger)
		return
	}

	response := transfer_models.BudgetActualResponse{BudgetActual: budget}
	commonHttp.SuccessResponse[transfer_models.BudgetActualResponse](w, http.StatusOK, response)
}

// @Summary		Get User Accounts
// @Tags			User
// @Description	Get User accounts
// @Produce		json
// @Success		200		{object}	Response[transfer_models.Account]	     	"Show actual accounts"
// @Success		204		{object}	Response[string]	     	"Show actual accounts"
// @Failure		400		{object}	ResponseError		"Client error"
// @Failure     401    	{object}  	ResponseError  		"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  		"Forbidden user"
// @Failure		500		{object}	ResponseError		"Server error"
// @Router		/api/user/accounts/all [get]
func (h *Handler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	accountInfo, err := h.userService.GetAccounts(r.Context(), user.ID)

	var errNoSuchAccounts *models.NoSuchAccounts

	if errors.As(err, &errNoSuchAccounts) {
		h.logger.Info(err.Error())
		commonHttp.SuccessResponse(w, http.StatusNoContent, "")
		return
	}

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.AccountNotFound, h.logger)
		return
	}

	response := transfer_models.Account{AccountMas: accountInfo}
	commonHttp.SuccessResponse[transfer_models.Account](w, http.StatusOK, response)
}

// @Summary		Get Feed
// @Tags			User
// @Description	Get Feed user info
// @Produce		json
// @Success		200		{object}	Response[transfer_models.UserFeed]	     	"Show actual accounts"
// @Failure		400		{object}	ResponseError		"Client error"
// @Failure     401    	{object}  	ResponseError  		"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  		"Forbidden user"
// @Failure		500		{object}	ResponseError		"Server error"
// @Router		/api/user/feed [get]
func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	dataFeed, err := h.userService.GetFeed(r.Context(), user.ID)

	var errNoSuchPlannedBudgetError *models.NoSuchPlannedBudgetError
	var errNoSuchUserIdBalanceError *models.NoSuchUserIdBalanceError
	var errNoSuchAccounts *models.NoSuchAccounts

	if errors.As(err, &errNoSuchAccounts) ||
		errors.As(err, &errNoSuchPlannedBudgetError) ||
		errors.As(err, &errNoSuchUserIdBalanceError) {

		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.UserFeedNotFound, h.logger)
		return
	}

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserFeedServerError, h.logger)
		return
	}

	commonHttp.SuccessResponse(w, http.StatusOK, dataFeed)
}

// @Summary		PUT Update
// @Tags			User
// @Description	Update user info
// @Accept      json
// @Produce		json
// @Param			user		body		transfer_models.UserUdate		true		"user info update"
// @Success		200		{object}	Response[NilBody]	     	"Update user info"
// @Failure		400		{object}	ResponseError		"Client error"
// @Failure     401    	{object}  	ResponseError  		"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  		"Forbidden user"
// @Failure		500		{object}	ResponseError		"Server error"
// @Router		/api/user/update [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	var updProfile transfer_models.UserUdate

	if err := easyjson.UnmarshalFromReader(r.Body, &updProfile); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	if err := updProfile.CheckValid(); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	if err := h.userService.UpdateUser(r.Context(), updProfile.ToUser(user)); err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.UserNotFound, h.logger)
			return
		}

		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserServerError, h.logger)
			return
		}
	}

	commonHttp.SuccessResponse(w, http.StatusOK, commonHttp.NilBody{})
}

// @Summary     PUT Update Photo
// @Tags        User
// @Description Update user photo
// @Accept      multipart/form-data
// @Produce     json
// @Param       userID        path  string  true  "User ID"
// @Param       upload        formData file    true  "New photo to upload"
// @Param       path          formData string  true  "Path to old photo"
// @Success     200           {object} Response[transfer_models.PhotoUpdate] "Photo updated successfully"
// @Failure     400           {object} ResponseError   "Client error"
// @Failure     401    	{object}  	ResponseError  		"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  		"Forbidden user"
// @Failure     500           {object} ResponseError   "Server error"
// @Router      /api/user/updatePhoto [put]
func (h *Handler) UpdatePhoto(w http.ResponseWriter, r *http.Request) { // need test
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	err = r.ParseMultipartForm(transfer_models.MaxFileSize)
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserFileServerError, h.logger)
		return
	}

	file, fileTmp, err := r.FormFile("upload")
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.UserFileUnableUpload, h.logger)
		return
	}

	buf, _ := io.ReadAll(file)
	file.Close()
	if file, err = fileTmp.Open(); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.UserFileUnableOpen, h.logger)
		return
	}

	if http.DetectContentType(buf) != "image/jpeg" && http.DetectContentType(buf) != "image/png" && http.DetectContentType(buf) != "image/jpg" {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.UserFileNotCorrectType, h.logger)
		return
	}
	defer file.Close()

	var oldName uuid.UUID
	oldName, err = uuid.Parse(r.PostFormValue("path"))
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.UserFileNotPath, h.logger)
		return
	}

	if oldName != uuid.Nil {
		err = os.Remove(transfer_models.FolderPath + fmt.Sprintf("%s.jpg", oldName.String()))
		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.UserFileNotDelete, h.logger)
			return
		}
	}

	name, err := h.userService.UpdatePhoto(r.Context(), user.ID)
	var errNoSuchUser *models.NoSuchUserError
	if errors.As(err, &errNoSuchUser) {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.UserNotFound, h.logger)
		return
	}

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserFileServerNotUpdateError, h.logger)
		return
	}

	f, err := os.Create(fmt.Sprintf("%s%s.jpg", transfer_models.FolderPath, name.String()))
	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserFileServerNotCreate, h.logger)
		return
	}
	defer f.Close()

	if _, err = io.Copy(f, file); err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserFileServerNotCreate, h.logger)
		return
	}

	commonHttp.SuccessResponse(w, http.StatusOK, transfer_models.PhotoUpdate{Path: name})
}

// @Summary		POST 	Add User in Account
// @Tags				User
// @Description	Post 	User
// @Produce		json
// @Param		User	body		models.AddUserAccount		 true		    "Add user in account"
// @Success		200		{object}	Response[NilBody]				    "Update account"
// @Failure		400		{object}	ResponseError						"Client error"
// @Failure     401    	{object}  	ResponseError  						"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  						"Forbidden user"
// @Failure		500		{object}	ResponseError						"Server error"
// @Router		/api/user/addUserInAccount [post]
func (h *Handler) AddUserInAccount(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	var accountInput models.AddUserAccount

	if err := easyjson.UnmarshalFromReader(r.Body, &accountInput); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	if err := h.userService.AddUserInAccount(r.Context(), accountInput, user.ID); err != nil {
		var errNoSuchUser *models.NoSuchUserInLogin
		if errors.As(err, &errNoSuchUser) {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.UserNotFoundLogin, h.logger)
			return
		}

		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponse(w, http.StatusForbidden, err, commonHttp.ForbiddenUser, h.logger)
			return
		}

		var errDuplicate *models.DuplicateError
		if errors.As(err, &errDuplicate) {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.UserDuplicate, h.logger)
			return
		}

		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserServerError, h.logger)
			return
		}
	}

	commonHttp.SuccessResponse(w, http.StatusOK, commonHttp.NilBody{})
}

// @Summary		PUT 	Unsibscribe in Account
// @Tags				User
// @Description	Post 	User
// @Produce		json
// @Success		200		{object}	Response[NilBody]				    "Unsibscribe in Account"
// @Failure		400		{object}	ResponseError						"Client error"
// @Failure     401    	{object}  	ResponseError  						"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  						"Forbidden user"
// @Failure		500		{object}	ResponseError						"Server error"
// @Router		/api/user/unsubscribeAccount/{account_id} [put]
func (h *Handler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
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

	if err := h.userService.Unsubscribe(r.Context(), accountID, user.ID); err != nil {
		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserServerError, h.logger)
			return
		}
	}

	commonHttp.SuccessResponse(w, http.StatusOK, commonHttp.NilBody{})

}

// @Summary		DELETE 	Delete user in Account
// @Tags				User
// @Description	Post 	User
// @Produce		json
// @Param		User	body		models.DeleteInAccount		 true		    "Delete user in Account"
// @Success		200		{object}	Response[NilBody]				    "Unsibscribe in Account"
// @Failure		400		{object}	ResponseError						"Client error"
// @Failure     401    	{object}  	ResponseError  						"Unauthorized user"
// @Failure     403    	{object}  	ResponseError  						"Forbidden user"
// @Failure		500		{object}	ResponseError						"Server error"
// @Router		/api/user/deleteUserInAccount [delete]
func (h *Handler) DeleteUserInAccount(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), h.logger)
		return
	}

	var accountDelete models.DeleteInAccount

	if err := easyjson.UnmarshalFromReader(r.Body, &accountDelete); err != nil {
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err, commonHttp.InvalidBodyRequest, h.logger)
		return
	}

	if err := h.userService.DeleteUserInAccount(r.Context(), accountDelete.UserID, accountDelete.AccountID, user.ID); err != nil {
		var errNoSuchUser *models.NoSuchUserInLogin
		if errors.As(err, &errNoSuchUser) {
			commonHttp.ErrorResponse(w, http.StatusBadRequest, err, transfer_models.UserNotFound, h.logger)
			return
		}

		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponse(w, http.StatusForbidden, err, commonHttp.ForbiddenUser, h.logger)
			return
		}

		if err != nil {
			commonHttp.ErrorResponse(w, http.StatusInternalServerError, err, transfer_models.UserServerError, h.logger)
			return
		}
	}

	commonHttp.SuccessResponse(w, http.StatusOK, commonHttp.NilBody{})
}
