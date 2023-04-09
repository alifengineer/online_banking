package models

// google uuid

type Account struct {
	ID        string  `json:"guid"`
	UserID    string  `json:"user_id"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type CreateAccountRequest struct {
	UserID  string  `json:"user_id"`
	Balance float64 `json:"balance"`
}

type GetAccountByIDRequest struct {
	ID string `json:"id"`
}

type GetAccountByIDResponse struct {
	Account *Account `json:"account"`
}

type GetAccountsByUserIDRequest struct {
	UserID  string `json:"user_id"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	OrderBy string `json:"order_by"`
	Desc    bool   `json:"desc"`
}

type GetAccountsByUserIDResponse struct {
	Accounts []*Account `json:"accounts"`
	Count    int        `json:"count"`
}
