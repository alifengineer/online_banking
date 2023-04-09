package models

type Transaction struct {
	ID          string  `json:"id"`
	AccountID   string  `json:"account_id"`
	RecipientID string  `json:"recipient_id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	CreatedAt   string  `json:"created_at"`
	Approved    bool    `json:"approved"`
	Done        bool    `json:"done"`
}

type GetTransactionsByAccountIDRequest struct {
	Limit       int    `json:"limit"`
	Offset      int    `json:"offset"`
	OrderBy     string `json:"order_by"`
	Desc        bool   `json:"desc"`
	From        string `json:"from"`
	To          string `json:"to"`
	RecipientID string `json:"recipient_id"`
	AccountID   string `json:"account_id"`
}

type GetTransactionsByAccountIDResponse struct {
	Transactions []*Transaction `json:"transactions"`
	Count        int            `json:"count"`
}

type GetTransactionByIDRequest struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
}

type ApproveTransactionsRequest struct {
	AccountID      string   `json:"account_id"`
	TransactionIDS []string `json:"transaction_ids"`
}

type GetTransactionsByIDSRequest struct {
	IDS       []string `json:"ids"`
	Approved  bool     `json:"approved"`
	Done      bool     `json:"done"`
	AccountID string   `json:"account_id"`
}

type GetTransactionsByIDSResponse struct {
	Transactions []*Transaction `json:"transactions"`
}
