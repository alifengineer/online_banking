package handlers

import (
	"github.com/dilmurodov/online_banking/api/http"
	"github.com/dilmurodov/online_banking/config"
	"github.com/dilmurodov/online_banking/pkg/models"
	"github.com/dilmurodov/online_banking/pkg/util"
	"github.com/gin-gonic/gin"
)

// CreateAccount godoc
// @Security BearerAuth
// @ID create_account
// @Router /api/v1/user/accounts [POST]
// @Summary Create Account
// @Description Create Account
// @Tags Account
// @Accept json
// @Produce json
// @Success 201 {object} http.Response{data=string} "Created"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) AccountCreateHandler(c *gin.Context) {

	authObj, ok := c.Get("auth")
	if !ok {
		h.handleResponse(c, http.Unauthorized, "unauthorized")
		return
	}

	auth := authObj.(*models.HasAccessModel)

	req := &models.CreateAccountRequest{
		UserID:  auth.UserId,
		Balance: 0,
	}

	resp, err := h.services.AccountService().CreateAccount(c.Request.Context(), req)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.Created, resp)
}

// GetAccounts godoc
// @Security BearerAuth
// @ID get_accounts
// @Router /api/v1/user/accounts [GET]
// @Summary Get Accounts
// @Description Get Accounts
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {object} http.Response{data=string} "OK"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) AccountsGetHandler(c *gin.Context) {

	authObj, ok := c.Get("auth")
	if !ok {
		h.handleResponse(c, http.Unauthorized, "unauthorized")
		return
	}

	auth := authObj.(*models.HasAccessModel)

	req := &models.GetAccountsByUserIDRequest{
		UserID: auth.UserId,
	}

	resp, err := h.services.AccountService().GetAccountsByUserID(c.Request.Context(), req)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, resp)
}

// GetAccount godoc
// @Security BearerAuth
// @ID get_account
// @Router /api/v1/user/accounts/{id} [GET]
// @Summary Get Account
// @Description Get Account
// @Tags Account
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} http.Response{data=string} "OK"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) AccountGetHandler(c *gin.Context) {

	accountID := c.Param("id")
	if !util.IsValidUUID(accountID) {
		h.handleResponse(c, http.BadRequest, "Invalid account ID")
		return
	}

	req := &models.GetAccountByIDRequest{
		ID: accountID,
	}

	resp, err := h.services.AccountService().GetAccountByID(c.Request.Context(), req)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, resp)
}

// GetAccountTransactions godoc
// @Security BearerAuth
// @ID get_account_transactions
// @Router /api/v1/user/accounts/{id}/transactions [GET]
// @Summary Get Account Transactions
// @Description Get Account Transactions
// @Tags Account
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} http.Response{data=string} "OK"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) AccountTransactionsHandler(c *gin.Context) {

	accountID := c.Param("id")
	if !util.IsValidUUID(accountID) {
		h.handleResponse(c, http.BadRequest, "Invalid account ID")
		return
	}

	req := &models.GetTransactionsByAccountIDRequest{
		AccountID: accountID,
	}

	resp, err := h.services.AccountService().GetAccountTransactions(c.Request.Context(), req)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, resp)
}

// GetAccountTransaction godoc
// @Security BearerAuth
// @ID get_account_transaction
// @Router /api/v1/user/accounts/{id}/transactions/{transaction_id} [GET]
// @Summary Get Account Transaction
// @Description Get Account Transaction
// @Tags Account
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param transaction_id path string true "Transaction ID"
// @Success 200 {object} http.Response{data=string} "OK"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) AccountTransactionByIDHandler(c *gin.Context) {

	accountID := c.Param("id")
	if !util.IsValidUUID(accountID) {
		h.handleResponse(c, http.BadRequest, "Invalid account ID")
		return
	}

	transactionID := c.Param("transaction_id")
	if !util.IsValidUUID(transactionID) {
		h.handleResponse(c, http.BadRequest, "Invalid transaction ID")
		return
	}

	req := &models.GetTransactionByIDRequest{
		ID:        transactionID,
		AccountID: accountID,
	}

	resp, err := h.services.AccountService().GetAccountTransactionByID(c.Request.Context(), req)
	if err != nil && err.Error() == config.RECORD_NOT_FOUND {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	} else if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, resp)
}
