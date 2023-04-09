package account

import (
	"context"

	"github.com/dilmurodov/online_banking/config"
	"github.com/dilmurodov/online_banking/pkg/logger"
	"github.com/dilmurodov/online_banking/pkg/models"
	"github.com/dilmurodov/online_banking/storage"
)

type ServiceI interface {
	CreateAccount(ctx context.Context, req *models.CreateAccountRequest) (resp *models.Account, err error)
	GetAccountsByUserID(ctx context.Context, req *models.GetAccountsByUserIDRequest) (resp *models.GetAccountsByUserIDResponse, err error)
	GetAccountByID(ctx context.Context, req *models.GetAccountByIDRequest) (resp *models.Account, err error)
	GetAccountTransactions(ctx context.Context, req *models.GetTransactionsByAccountIDRequest) (resp *models.GetTransactionsByAccountIDResponse, err error)
	GetAccountTransactionByID(ctx context.Context, req *models.GetTransactionByIDRequest) (resp *models.Transaction, err error)
}

type Service struct {
	cfg  config.Config
	log  logger.LoggerI
	strg storage.StorageI
}

func NewService(cfg config.Config, log logger.LoggerI, strg storage.StorageI) *Service {
	return &Service{
		cfg:  cfg,
		log:  log,
		strg: strg,
	}
}
