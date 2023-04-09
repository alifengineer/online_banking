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
// @Success 201 {object} http.Response{data=string} "OK"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) WithDrawalHandler(c *gin.Context) {

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
	req.AccountID = authObj.UserId

	err := h.services.PaymentService().WithDrawal(c.Request.Context(), &req)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
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
// @Success 201 {object} http.Response{data=string} "OK"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) DepositHandler(c *gin.Context) {

	auth, ok := c.Get("auth")
	if !ok {
		h.handleResponse(c, http.Unauthorized, "unauthorized")
		return
	}
	authObj := auth.(*models.HasAccessModel)

	var req models.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	req.AccountID = authObj.UserId

	err := h.services.PaymentService().Deposit(c.Request.Context(), &req)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
}

// @Security BearerAuth
// ConfirmPaymentHandler godoc
// @ID confirm
// @Summary Confirm Payment
// @Description Confirm Payment
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

	var req models.CaptureTransactionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	err := h.services.PaymentService().CaptureTransactions(c.Request.Context(), &req)
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
// @Success 201 {object} http.Response{data=string} "OK"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *Handler) TransferHandler(c *gin.Context) {

	auth, ok := c.Get("auth")
	if !ok {
		h.handleResponse(c, http.Unauthorized, "unauthorized")
		return
	}
	authObj := auth.(*models.HasAccessModel)

	var req models.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	req.FromAccountID = authObj.UserId

	err := h.services.PaymentService().Transfer(c.Request.Context(), &req)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
}
