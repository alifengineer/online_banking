package payment

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dilmurodov/online_banking/config"
	"github.com/dilmurodov/online_banking/pkg/models"
	mock_storage "github.com/dilmurodov/online_banking/storage/mock"
	"github.com/dilmurodov/online_banking/storage/postgres"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPayment_Transfer(t *testing.T) {

	r := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	r.NoError(err)

	s := NewService(
		config.Config{},
		zap.NewNop(),
		postgres.NewStore(db),
	)

	repo := mock_storage.NewMockAccountRepoI(ctrl)
	repoTx := mock_storage.NewMockTxRepoI(ctrl)

	mock.ExpectBegin()
	row1 := sqlmock.NewRows([]string{"guid", "user_id", "balance", "created_at", "updated_at"}).AddRow("TestAccountID1", "TestUserID", 200, "2021-01-01", "2021-01-01")

	row2 := sqlmock.NewRows([]string{"guid", "user_id", "balance", "created_at", "updated_at"}).AddRow("TestAccountID2", "TestUserID", 200, "2021-01-01", "2021-01-01")

	mock.ExpectQuery(`^SELECT (.+?) FROM accounts * `).WithArgs("TestAccountID1").WillReturnRows(row1)

	mock.ExpectQuery(`^SELECT (.+?) FROM accounts * `).WithArgs("TestAccountID2").WillReturnRows(row2)

	mock.ExpectPrepare("INSERT INTO transactions").ExpectExec().WithArgs("TestAccountID1", 100.0, "TestAccountID2", "debit").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectPrepare("INSERT INTO transactions").ExpectExec().WithArgs("TestAccountID2", 100.0, "TestAccountID1", "credit").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	t.Run("SUCCESS", func(t *testing.T) {

		ctx := context.Background()

		var tx *sql.Tx

		acc1 := &models.Account{
			ID: "TestAccountID1",
		}

		acc2 := &models.Account{
			ID: "TestAccountID2",
		}

		inTx := &models.Transaction{
			ID:          "TestTransactionID",
			AccountID:   "TestAccountID1",
			RecipientID: "TestAccountID2",
		}

		req := &models.TransferRequest{
			FromAccountID: "TestAccountID1",
			ToAccountID:   "TestAccountID2",
			Amount:        100.0,
		}

		repoTx.EXPECT().BeginTx(ctx).Return(tx, nil).Times(1).AnyTimes()

		repo.EXPECT().GetAccountByID(ctx, &models.GetAccountByIDRequest{
			ID: "TestAccountID1",
		}).Return(acc1, nil).Times(1).AnyTimes()

		repo.EXPECT().GetAccountByID(ctx, &models.GetAccountByIDRequest{
			ID: "TestAccountID2",
		}).Return(acc2, nil).Times(1).AnyTimes()

		repoTx.EXPECT().CreateTransaction(ctx, tx, inTx).Return(nil).Times(1).AnyTimes()

		repoTx.EXPECT().CreateTransaction(ctx, tx, inTx).Return(nil).Times(1).AnyTimes()

		err = s.Transfer(ctx, req)
		r.NoError(err)
		r.NoError(mock.ExpectationsWereMet())
	})
}

func TestPayment_WithDrawal(t *testing.T) {

	r := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	r.NoError(err)

	s := NewService(
		config.Config{},
		zap.NewNop(),
		postgres.NewStore(db),
	)

	repo := mock_storage.NewMockAccountRepoI(ctrl)
	repoTx := mock_storage.NewMockTxRepoI(ctrl)

	mock.ExpectBegin()
	row1 := sqlmock.NewRows([]string{"guid", "user_id", "balance", "created_at", "updated_at"}).AddRow("TestAccountID1", "TestUserID", 200, "2021-01-01", "2021-01-01")

	mock.ExpectQuery(`^SELECT (.+?) FROM accounts * `).WithArgs("TestAccountID1").WillReturnRows(row1)

	mock.ExpectPrepare("INSERT INTO transactions").ExpectExec().WithArgs("TestAccountID1", 100.0, "TestAccountID1", "debit").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	t.Run("SUCCESS", func(t *testing.T) {

		ctx := context.Background()

		var tx *sql.Tx

		acc1 := &models.Account{
			ID: "TestAccountID1",
		}

		inTx := &models.Transaction{
			ID:          "TestTransactionID",
			AccountID:   "TestAccountID1",
			RecipientID: "TestAccountID1",
		}

		req := &models.WithDrawalRequest{
			AccountID: "TestAccountID1",
			Amount:    100.0,
		}

		repoTx.EXPECT().BeginTx(ctx).Return(tx, nil).Times(1).AnyTimes()

		repo.EXPECT().GetAccountByID(ctx, &models.GetAccountByIDRequest{
			ID: "TestAccountID1",
		}).Return(acc1, nil).Times(1).AnyTimes()

		repoTx.EXPECT().CreateTransaction(ctx, tx, inTx).Return(nil).Times(1).AnyTimes()

		err = s.WithDrawal(ctx, req)
		r.NoError(err)
		r.NoError(mock.ExpectationsWereMet())
	})
}

func TestPayment_Deposit(t *testing.T) {

	r := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	r.NoError(err)

	s := NewService(
		config.Config{},
		zap.NewNop(),
		postgres.NewStore(db),
	)

	repo := mock_storage.NewMockAccountRepoI(ctrl)
	repoTx := mock_storage.NewMockTxRepoI(ctrl)

	mock.ExpectBegin()
	row1 := sqlmock.NewRows([]string{"guid", "user_id", "balance", "created_at", "updated_at"}).AddRow("TestAccountID1", "TestUserID", 200, "2021-01-01", "2021-01-01")

	mock.ExpectQuery(`^SELECT (.+?) FROM accounts * `).WithArgs("TestAccountID1").WillReturnRows(row1)

	mock.ExpectPrepare("INSERT INTO transactions").ExpectExec().WithArgs("TestAccountID1", 100.0, "TestAccountID1", "credit").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	t.Run("SUCCESS", func(t *testing.T) {

		ctx := context.Background()

		var tx *sql.Tx

		acc1 := &models.Account{
			ID: "TestAccountID1",
		}

		inTx := &models.Transaction{
			ID:          "TestTransactionID",
			AccountID:   "TestAccountID1",
			RecipientID: "TestAccountID1",
		}

		req := &models.DepositRequest{
			AccountID: "TestAccountID1",
			Amount:    100.0,
		}

		repoTx.EXPECT().BeginTx(ctx).Return(tx, nil).Times(1).AnyTimes()

		repo.EXPECT().GetAccountByID(ctx, &models.GetAccountByIDRequest{
			ID: "TestAccountID1",
		}).Return(acc1, nil).Times(1).AnyTimes()

		repoTx.EXPECT().CreateTransaction(ctx, tx, inTx).Return(nil).Times(1).AnyTimes()

		err = s.Deposit(ctx, req)
		r.NoError(err)
		r.NoError(mock.ExpectationsWereMet())
	})
}

func TestPayment_CaptureTransactions(t *testing.T) {
	r := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	r.NoError(err)

	s := NewService(
		config.Config{},
		zap.NewNop(),
		postgres.NewStore(db),
	)

	repo := mock_storage.NewMockAccountRepoI(ctrl)
	repoTx := mock_storage.NewMockTxRepoI(ctrl)

	mock.ExpectBegin()
	row1 := sqlmock.NewRows([]string{"guid", "user_id", "balance", "created_at", "updated_at"}).AddRow("TestAccountID1", "TestUserID", 200, "2021-01-01", "2021-01-01")

	txrow := sqlmock.NewRows([]string{"guid", "account_id", "transaction_amount", "transaction_type", "recipient_id", "created_at"}).AddRow("TestTransactionID", "TestAccountID1", 100.0, "debit", "TestAccountID1", "2021-01-01")

	mock.ExpectQuery(`^SELECT (.+?) FROM transactions * `).WithArgs(pq.Array([]string{"TestTransactionID"}), "TestAccountID1").WillReturnRows(txrow)

	mock.ExpectQuery(`^SELECT (.+?) FROM accounts * `).WithArgs("TestAccountID1").WillReturnRows(row1)

	mock.ExpectPrepare(`UPDATE accounts`).ExpectExec().WithArgs(100.0, "TestAccountID1").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	t.Run("SUCCESS", func(t *testing.T) {

		ctx := context.Background()

		var tx *sql.Tx

		acc1 := &models.Account{
			ID: "TestAccountID1",
		}

		inTx := &models.Transaction{
			ID:          "TestTransactionID",
			AccountID:   "TestAccountID1",
			RecipientID: "TestAccountID1",
		}

		req := &models.CaptureTransactionsRequest{
			AccountID: "TestAccountID1",
			TransactionIDS: []string{
				"TestTransactionID",
			},
		}

		txsresp := &models.GetTransactionsByIDSResponse{
			Transactions: []*models.Transaction{
				{
					ID:          "TestTransactionID",
					AccountID:   "TestAccountID1",
					RecipientID: "TestAccountID1",
				},
			},
		}

		repoTx.EXPECT().BeginTx(ctx).Return(tx, nil).Times(1).AnyTimes()

		repoTx.EXPECT().GetTransactionsByIDS(ctx, &models.GetTransactionsByIDSRequest{
			IDS: []string{
				"TestTransactionID",
			},
			AccountID: "TestAccountID1",
		}).Return(txsresp, nil).Times(1).AnyTimes()

		repo.EXPECT().GetAccountByID(ctx, &models.GetAccountByIDRequest{
			ID: "TestAccountID1",
		}).Return(acc1, nil).Times(1).AnyTimes()

		repoTx.EXPECT().CreateTransaction(ctx, tx, inTx).Return(nil).Times(1).AnyTimes()

		err = s.CaptureTransactions(ctx, req)
		r.NoError(err)
		r.NoError(mock.ExpectationsWereMet())
	})
}
