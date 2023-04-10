package payment

import (
	"context"
	"fmt"
	"log"

	"github.com/dilmurodov/online_banking/pkg/logger"
	"github.com/dilmurodov/online_banking/pkg/models"
)

// Transfer transfers the specified amount from one account to another
func (s *Service) Transfer(ctx context.Context, req *models.TransferRequest) (resp *models.TransferResponse, err error) {
	s.log.Info("---Transfer--->", logger.Any("req", req))
	resp = &models.TransferResponse{}

	// Begin a database transaction for the transfer
	tx, err := s.strg.TxRepo().BeginTx(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", logger.Error(err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the account objects from the account repository
	fromAccount, err := s.strg.Account().GetAccountByID(ctx, &models.GetAccountByIDRequest{
		ID: req.FromAccountID,
	})
	if err != nil {
		s.log.Error("failed to get from account", logger.Error(err))
		return nil, fmt.Errorf("failed to get from account: %w", err)
	}

	toAccount, err := s.strg.Account().GetAccountByID(ctx, &models.GetAccountByIDRequest{
		ID: req.ToAccountID,
	})
	if err != nil {
		s.log.Error("failed to get to account", logger.Error(err))
		return nil, fmt.Errorf("failed to get to account: %w", err)
	}

	// Ensure that the from account has enough funds to transfer
	if fromAccount.Balance < req.Amount {
		s.log.Error("insufficient funds in from account", logger.Error(err))
		return nil, fmt.Errorf("insufficient funds in from account")
	}

	// Debit the amount from the from account and credit it to the to account
	fromAccount.Balance -= req.Amount
	toAccount.Balance += req.Amount

	// Create the debit and credit transactions for the transfer
	debitTx := &models.Transaction{
		AccountID:   fromAccount.ID,
		Amount:      req.Amount,
		Type:        "debit",
		RecipientID: toAccount.ID,
	}

	creditTx := &models.Transaction{
		AccountID:   toAccount.ID,
		Amount:      req.Amount,
		Type:        "credit",
		RecipientID: fromAccount.ID,
	}

	// Save the debit and credit transactions to the database
	createTx1, err := s.strg.TxRepo().CreateTransaction(ctx, tx, debitTx)
	if err != nil {
		s.log.Error("failed to create debit transaction", logger.Error(err))
		return nil, fmt.Errorf("failed to create debit transaction: %w", err)
	}
	resp.Transactions = append(resp.Transactions, createTx1)

	createTx2, err := s.strg.TxRepo().CreateTransaction(ctx, tx, creditTx)
	if err != nil {
		s.log.Error("failed to create credit transaction", logger.Error(err))
		return nil, fmt.Errorf("failed to create credit transaction: %w", err)
	}
	resp.Transactions = append(resp.Transactions, createTx2)

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", logger.Error(err))
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return resp, nil
}

// WithDrawal the specified amount from one account to another
func (s *Service) WithDrawal(ctx context.Context, req *models.WithDrawalRequest) (resp *models.WithDrawalResponse, err error) {
	resp = &models.WithDrawalResponse{}
	s.log.Info("---WithDrawal---", logger.Any("req", req))
	// Begin a database transaction for the transfer
	tx, err := s.strg.TxRepo().BeginTx(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", logger.Any("err", err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the account objects from the account repository
	account, err := s.strg.Account().GetAccountByID(ctx, &models.GetAccountByIDRequest{
		ID: req.AccountID,
	})
	if err != nil {
		s.log.Error("failed to get from account", logger.Any("err", err))
		return nil, fmt.Errorf("failed to get from account: %w", err)
	}

	// Create credit transactions for the transfer
	debitTx := &models.Transaction{
		AccountID:   account.ID,
		Amount:      req.Amount,
		Type:        "debit",
		RecipientID: account.ID,
	}

	// Save the debit and credit transactions to the database
	createTx, err := s.strg.TxRepo().CreateTransaction(ctx, tx, debitTx)
	if err != nil {
		s.log.Error("failed to create debit transaction", logger.Any("err", err))
		return nil, fmt.Errorf("failed to create debit transaction: %w", err)
	}
	resp.Transaction = createTx

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", logger.Any("err", err))
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return resp, nil
}

func (s *Service) CaptureTransactions(ctx context.Context, req *models.CaptureTransactionsRequest) error {
	s.log.Info("---CaptureTransactions--->", logger.Any("req", req))
	// Begin a database transaction for the transfer
	tx, err := s.strg.TxRepo().BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	var (
		accnt *models.Account
		cntr  int
	)

	// Get transactions
	transactions, err := s.strg.TxRepo().GetTransactionsByIDS(ctx, &models.GetTransactionsByIDSRequest{
		IDS: req.TransactionIDS,
	})
	if err != nil {
		_ = tx.Rollback()
		s.log.Error("failed to get transactions", logger.Error(err))
		return fmt.Errorf("failed to get transactions: %w", err)
	}

	for _, v := range transactions.Transactions {
		if err != nil {
			_ = tx.Rollback()
			s.log.Error("failed to get from account", logger.Error(err))
			return fmt.Errorf("failed to get from account: %w", err)
		}

		if v.Type == "credit" {
			cntr++
			log.Println("credit---->", v.Amount)
			accnt, err = s.strg.Account().GetAccountByID(ctx, &models.GetAccountByIDRequest{
				ID: v.AccountID,
			})
			if err != nil {
				_ = tx.Rollback()
				s.log.Error("failed to get from account", logger.Error(err))
				return fmt.Errorf("failed to get from account: %w", err)
			}
			accnt.Balance += v.Amount
		}

		if v.Type == "debit" {
			cntr++
			log.Println("debit---->", v.Amount)
			accnt, err = s.strg.Account().GetAccountByID(ctx, &models.GetAccountByIDRequest{
				ID: v.AccountID,
			})
			if err != nil {
				_ = tx.Rollback()
				s.log.Error("failed to get from account", logger.Error(err))
				return fmt.Errorf("failed to get from account: %w", err)
			}
			accnt.Balance -= v.Amount
		}

		// Update the account balances in the database
		err = s.strg.Account().UpdateAccountBalance(ctx, tx, accnt)
		if err != nil {
			_ = tx.Rollback()
			s.log.Error("failed to update from account", logger.Error(err))
			return fmt.Errorf("failed to update from account: %w", err)
		}
	}

	if cntr != len(transactions.Transactions) {
		_ = tx.Rollback()
		s.log.Error("failed to update from account", logger.Error(err))
		return fmt.Errorf("transactions number mismatch")
	}

	err = s.strg.TxRepo().ApproveTransactions(ctx, tx, &models.ApproveTransactionsRequest{
		TransactionIDS: req.TransactionIDS,
		AccountID:      req.AccountID,
	})
	if err != nil {
		_ = tx.Rollback()
		s.log.Error("failed to approve transactions", logger.Error(err))
		return fmt.Errorf("failed to approve transactions: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		s.log.Error("failed to commit transaction", logger.Error(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Deposit the specified amount to one account
func (s *Service) Deposit(ctx context.Context, req *models.DepositRequest) (resp *models.DepositResponse, err error) {
	s.log.Info("---Deposit--->", logger.Any("req", req))
	resp = &models.DepositResponse{}
	// Begin a database transaction for the transfer
	tx, err := s.strg.TxRepo().BeginTx(ctx)
	if err != nil {
		s.log.Error("---Deposit->BeginTx--->", logger.Error(err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the account objects from the account repository
	account, err := s.strg.Account().GetAccountByID(ctx, &models.GetAccountByIDRequest{
		ID: req.AccountID,
	})
	if err != nil {
		s.log.Error("---Deposit->GetAccountByID--->", logger.Error(err))
		return nil, fmt.Errorf("failed to get from account: %w", err)
	}

	// Create credit transactions for the transfer
	creditTx := &models.Transaction{
		AccountID:   account.ID,
		Amount:      req.Amount,
		Type:        "credit",
		RecipientID: account.ID,
	}

	// Save the debit and credit transactions to the database
	createTxresp, err := s.strg.TxRepo().CreateTransaction(ctx, tx, creditTx)
	if err != nil {
		s.log.Error("---Deposit->CreateTransaction--->", logger.Error(err))
		return nil, fmt.Errorf("failed to create debit transaction: %w", err)
	}
	resp.Transaction = createTxresp

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		s.log.Error("---Deposit->Commit--->", logger.Error(err))
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return resp, nil
}
