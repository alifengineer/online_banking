package account

import (
	"context"

	"github.com/dilmurodov/online_banking/pkg/logger"
	"github.com/dilmurodov/online_banking/pkg/models"
)

func (self *Service) CreateAccount(ctx context.Context, req *models.CreateAccountRequest) (resp *models.Account, err error) {
	self.log.Info("---CreateAccount--->", logger.Any("req", req))

	resp, err = self.strg.Account().CreateAccount(ctx, req)
	if err != nil {
		self.log.Error("---CreateAccount--->", logger.Any("err", err))
		return nil, err
	}

	return
}

func (self *Service) GetAccountsByUserID(ctx context.Context, req *models.GetAccountsByUserIDRequest) (resp *models.GetAccountsByUserIDResponse, err error) {
	self.log.Info("---GetAllAccounts--->", logger.Any("req", req))

	resp, err = self.strg.Account().GetAccountsByUserID(ctx, req)

	return
}

func (s *Service) GetAccountByID(ctx context.Context, req *models.GetAccountByIDRequest) (resp *models.Account, err error) {
	s.log.Info("---GetAccountByID--->", logger.Any("req", req))

	resp, err = s.strg.Account().GetAccountByID(ctx, req)

	return
}

func (s *Service) GetAccountTransactions(ctx context.Context, req *models.GetTransactionsByAccountIDRequest) (resp *models.GetTransactionsByAccountIDResponse, err error) {
	s.log.Info("---GetAccountTransactions--->", logger.Any("req", req))

	resp, err = s.strg.TxRepo().GetTransactionsByAccountID(ctx, req)
	if err != nil {
		s.log.Error("---GetAccountTransactions--->", logger.Any("err", err))
		return nil, err
	}

	return
}

func (s *Service) GetAccountTransactionByID(ctx context.Context, req *models.GetTransactionByIDRequest) (resp *models.Transaction, err error) {
	s.log.Info("---GetAccountTransactionByID--->", logger.Any("req", req))

	resp, err = s.strg.TxRepo().GetTransactionByID(ctx, req)
	if err != nil {
		s.log.Error("---GetAccountTransactionByID--->", logger.Any("err", err))
		return nil, err
	}

	return resp, nil
}
