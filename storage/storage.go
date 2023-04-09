package storage

import (
	"context"
	"database/sql"

	"github.com/dilmurodov/online_banking/pkg/models"
)

type StorageI interface {
	CloseDB()
	User() UserRepoI
	Account() AccountRepoI
	TxRepo() TxRepoI
}

type UserRepoI interface {
	GetUserByID(context.Context, *models.GetUserByIDRequest) (*models.GetUserByIDResponse, error)
	CreateUser(context.Context, *models.CreateUserRequest) (*models.User, error)
	GetUserPasswordByPhone(ctx context.Context, phone string) (resp *models.User, err error)
}

type AccountRepoI interface {
	GetAccountByID(context.Context, *models.GetAccountByIDRequest) (*models.Account, error)
	CreateAccount(context.Context, *models.CreateAccountRequest) (*models.Account, error)
	GetAccountsByUserID(context.Context, *models.GetAccountsByUserIDRequest) (resp *models.GetAccountsByUserIDResponse, err error)
	UpdateAccountBalance(context.Context, *sql.Tx, *models.Account) error
}

type TxRepoI interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateTransaction(ctx context.Context, tx *sql.Tx, transaction *models.Transaction) error
	GetTransactionsByAccountID(ctx context.Context, req *models.GetTransactionsByAccountIDRequest) (resp *models.GetTransactionsByAccountIDResponse, err error)
	GetTransactionByID(ctx context.Context, req *models.GetTransactionByIDRequest) (resp *models.Transaction, err error)
	GetTransactionsByIDS(ctx context.Context, req *models.GetTransactionsByIDSRequest) (resp *models.GetTransactionsByIDSResponse, err error)
}
