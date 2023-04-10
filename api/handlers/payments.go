package handlers

import (
	"github.com/dilmurodov/online_banking/api/http"
	"github.com/dilmurodov/online_banking/pkg/models"
	"github.com/gin-gonic/gin"
)

// @Security BearerAuth
// WithDrawalHandler godoc
// @ID withdrawal
// @Summary Withdrawal
// @Description Withdrawal
// @Tags Payment
// @Accept json
// @Produce json
// @Param body body models.WithDrawalRequest true "Withdrawal"
// @Router /api/v1/payments/withdrawal [POST]
// @Success 201 {object} http.Response{data=models.WithDrawalResponse} "Created"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) WithDrawalHandler(c *gin.Context) {

	var (
		verified bool = false
	)

	auth, ok := c.Get("auth")
	if !ok {
		h.handleResponse(c, http.Unauthorized, "unauthorized")
		return
	}
	authObj := auth.(*models.HasAccessModel)

	var req models.WithDrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	accounts, err := h.services.AccountService().GetAccountsByUserID(c.Request.Context(), &models.GetAccountsByUserIDRequest{UserID: authObj.UserId})
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	for _, account := range accounts.Accounts {
		if req.AccountID == account.ID {
			verified = true
			break
		}
	}
	if !verified {
		h.handleResponse(c, http.BadRequest, "account not found")
		return
	}

	resp, err := h.services.PaymentService().WithDrawal(c.Request.Context(), &req)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	h.handleResponse(c, http.Created, resp)
}

// @Security BearerAuth
// DepositHandler godoc
// @ID deposit
// @Summary Deposit
// @Description Deposit
// @Tags Payment
// @Accept json
// @Produce json
// @Param body body models.DepositRequest true "Deposit"
// @Router /api/v1/payments/deposit [POST]
// @Success 201 {object} http.Response{data=models.DepositResponse} "Created"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) DepositHandler(c *gin.Context) {

	var (
		verified bool = false
	)

	// Get auth from context
	auth, ok := c.Get("auth")
	if !ok {
		h.handleResponse(c, http.Unauthorized, "unauthorized")
		return
	}
	authObj := auth.(*models.HasAccessModel)

	// Get request body
	var req models.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	// Check if account is owned by user
	accounts, err := h.services.AccountService().GetAccountsByUserID(c.Request.Context(), &models.GetAccountsByUserIDRequest{UserID: authObj.UserId})
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	for _, account := range accounts.Accounts {
		if req.AccountID == account.ID {
			verified = true
			break
		}
	}
	if !verified {
		h.handleResponse(c, http.BadRequest, "account not found")
		return
	}

	// Call service
	resp, err := h.services.PaymentService().Deposit(c.Request.Context(), &req)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	h.handleResponse(c, http.Created, resp)
}

// @Security BearerAuth
// ConfirmPaymentHandler godoc in TODO:
// @ID confirm
// @Summary Confirm Payment (OTP)
// @Description Confirm Payment (OTP)
// @Tags Payment
// @Accept json
// @Produce json
// @Param body body models.ConfirmPaymentRequest true "Confirm Payment"
// @Router /api/v1/payments/confirm [POST]
// @Success 201 {object} http.Response{data=string} "OK"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) ConfirmPaymentHandler(c *gin.Context) {

	// TODO: check user phone number
}

// @Security BearerAuth
// CapturePaymentHandler godoc
// @ID capture
// @Summary Capture Payment
// @Description Capture Payment
// @Tags Payment
// @Accept json
// @Produce json
// @Param body body models.CaptureTransactionsRequest true "Capture Payment"
// @Router /api/v1/payments/capture [POST]
// @Success 201 {object} http.Response{data=string} "OK"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) CaptureTransactionsHandler(c *gin.Context) {

	var (
		verified bool = false
	)

	// Get auth from context
	auth, ok := c.Get("auth")
	if !ok {
		h.handleResponse(c, http.Unauthorized, "unauthorized")
		return
	}
	authObj := auth.(*models.HasAccessModel)

	// Get request body
	var req models.CaptureTransactionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	// Check if account is owned by user
	accounts, err := h.services.AccountService().GetAccountsByUserID(c.Request.Context(), &models.GetAccountsByUserIDRequest{UserID: authObj.UserId})
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	for _, account := range accounts.Accounts {
		if req.AccountID == account.ID {
			verified = true
			break
		}
	}
	if !verified {
		h.handleResponse(c, http.BadRequest, "account not found")
		return
	}

	// Call service
	err = h.services.PaymentService().CaptureTransactions(c.Request.Context(), &req)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.Created, "OK")
}

// @Security BearerAuth
// TransferHandler godoc
// @ID transfer
// @Summary Transfer
// @Description Transfer
// @Tags Payment
// @Accept json
// @Produce json
// @Param body body models.TransferRequest true "Transfer"
// @Router /api/v1/payments/transfer [POST]
// @Success 201 {object} http.Response{data=models.TransferResponse} "Created"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) TransferHandler(c *gin.Context) {

	var (
		verified bool = false
	)

	// Get auth from context
	auth, ok := c.Get("auth")
	if !ok {
		h.handleResponse(c, http.Unauthorized, "unauthorized")
		return
	}
	authObj := auth.(*models.HasAccessModel)

	// Get request body
	var req models.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	// Check if account is owned by user
	accounts, err := h.services.AccountService().GetAccountsByUserID(c.Request.Context(), &models.GetAccountsByUserIDRequest{UserID: authObj.UserId})
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	for _, account := range accounts.Accounts {
		if req.FromAccountID == account.ID {
			verified = true
			break
		}
	}
	if !verified {
		h.handleResponse(c, http.BadRequest, "account not found")
		return
	}

	// Call service
	resp, err := h.services.PaymentService().Transfer(c.Request.Context(), &req)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	h.handleResponse(c, http.Created, resp)
}
