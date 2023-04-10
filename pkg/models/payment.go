package models

type TransferRequest struct {
	FromAccountID string  `json:"from_account_id"`
	ToAccountID   string  `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type TransferResponse struct {
	Transactions []*Transaction `json:"transaction"`
}

type WithDrawalRequest struct {
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type WithDrawalResponse struct {
	Transaction *Transaction `json:"transaction"`
}

type DepositRequest struct {
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type DepositResponse struct {
	Transaction *Transaction `json:"transaction"`
}

type CaptureTransactionsRequest struct {
	TransactionIDS []string `json:"transaction_ids"`
	AccountID      string   `json:"account_id"`
}

type ConfirmPaymentRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
