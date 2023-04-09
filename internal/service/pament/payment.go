package payment

import (
	"context"
	"fmt"

	"github.com/dilmurodov/online_banking/pkg/logger"
	"github.com/dilmurodov/online_banking/pkg/models"
)

// Transfer transfers the specified amount from one account to another
func (s *Service) Transfer(ctx context.Context, req *models.TransferRequest) error {

	// Begin a database transaction for the transfer
	tx, err := s.strg.TxRepo().BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Get the account objects from the account repository
	fromAccount, err := s.strg.Account().GetAccountByID(ctx, &models.GetAccountByIDRequest{
		ID: req.FromAccountID,
	})
	if err != nil {
		return fmt.Errorf("failed to get from account: %w", err)
	}

	toAccount, err := s.strg.Account().GetAccountByID(ctx, &models.GetAccountByIDRequest{
		ID: req.ToAccountID,
	})
	if err != nil {
		return fmt.Errorf("failed to get to account: %w", err)
	}

	// Ensure that the from account has enough funds to transfer
	if fromAccount.Balance < req.Amount {
		return fmt.Errorf("insufficient funds in from account")
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
	err = s.strg.TxRepo().CreateTransaction(ctx, tx, debitTx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to create debit transaction: %w", err)
	}

	err = s.strg.TxRepo().CreateTransaction(ctx, tx, creditTx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to create credit transaction: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithDrawal the specified amount from one account to another
func (s *Service) WithDrawal(ctx context.Context, req *models.WithDrawalRequest) error {

	// Begin a database transaction for the transfer
	tx, err := s.strg.TxRepo().BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Get the account objects from the account repository
	account, err := s.strg.Account().GetAccountByID(ctx, &models.GetAccountByIDRequest{
		ID: req.AccountID,
	})
	if err != nil {
		return fmt.Errorf("failed to get from account: %w", err)
	}

	// Create credit transactions for the transfer
	debitTx := &models.Transaction{
		AccountID:   account.ID,
		Amount:      req.Amount,
		Type:        "debit",
		RecipientID: account.ID,
	}

	// Save the debit and credit transactions to the database
	err = s.strg.TxRepo().CreateTransaction(ctx, tx, debitTx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to create debit transaction: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) CaptureTransactions(ctx context.Context, req *models.CaptureTransactionsRequest) error {
	s.log.Info("---CaptureTransactions--->", logger.Any("req", req))
	// Begin a database transaction for the transfer
	tx, err := s.strg.TxRepo().BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Get transactions
	transactions, err := s.strg.TxRepo().GetTransactionsByIDS(ctx, &models.GetTransactionsByIDSRequest{
		IDS:       req.TransactionIDS,
		AccountID: req.AccountID,
	})
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to get transactions: %w", err)
	}

	for _, v := range transactions.Transactions {

		accnt, err := s.strg.Account().GetAccountByID(ctx, &models.GetAccountByIDRequest{
			ID: v.AccountID,
		})
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to get from account: %w", err)
		}

		if v.Type == "credit" {
			accnt.Balance += v.Amount
		}

		if v.Type == "debit" {
			accnt.Balance -= v.Amount
		}

		// Update the account balances in the database
		err = s.strg.Account().UpdateAccountBalance(ctx, tx, accnt)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to update from account: %w", err)
		}
	}

	err = s.strg.TxRepo().ApproveTransactions(ctx, tx, &models.ApproveTransactionsRequest{
		TransactionIDS: req.TransactionIDS,
		AccountID:      req.AccountID,
	})
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to approve transactions: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Deposit the specified amount to one account
func (s *Service) Deposit(ctx context.Context, req *models.DepositRequest) error {

	// Begin a database transaction for the transfer
	tx, err := s.strg.TxRepo().BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Get the account objects from the account repository
	account, err := s.strg.Account().GetAccountByID(ctx, &models.GetAccountByIDRequest{
		ID: req.AccountID,
	})
	if err != nil {
		return fmt.Errorf("failed to get from account: %w", err)
	}

	// Create credit transactions for the transfer
	creditTx := &models.Transaction{
		AccountID:   account.ID,
		Amount:      req.Amount,
		Type:        "credit",
		RecipientID: account.ID,
	}

	// Save the debit and credit transactions to the database
	err = s.strg.TxRepo().CreateTransaction(ctx, tx, creditTx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to create debit transaction: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
