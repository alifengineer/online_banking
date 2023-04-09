package models

type TransferRequest struct {
	FromAccountID string  `json:"from_account_id"`
	ToAccountID   string  `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type WithDrawalRequest struct {
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type DepositRequest struct {
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type CaptureTransactionsRequest struct {
	TransactionIDS []string `json:"transaction_ids"`
	AccountID      string   `json:"account_id"`
}

type ConfirmPaymentRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
