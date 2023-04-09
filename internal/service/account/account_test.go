package account

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dilmurodov/online_banking/config"
	"github.com/dilmurodov/online_banking/pkg/models"
	mock_storage "github.com/dilmurodov/online_banking/storage/mock"
	"github.com/dilmurodov/online_banking/storage/postgres"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAccount_CreateAccount(t *testing.T) {

	r := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	r.NoError(err)

	rows := mock.NewRows([]string{"guid"}).AddRow("TestUserID")
	mock.ExpectPrepare("INSERT INTO accounts").ExpectQuery().WithArgs("TestUserID", 0.0).WillReturnRows(rows)

	repo := mock_storage.NewMockAccountRepoI(ctrl)
	s := NewService(config.Config{}, zap.NewNop(), postgres.NewStore(db))

	t.Run("SUCCESS", func(t *testing.T) {
		in := &models.CreateAccountRequest{
			UserID: "TestUserID",
		}
		mockresp := &models.Account{
			ID: "TestUserID",
		}

		repo.EXPECT().CreateAccount(context.Background(), in).Return(mockresp, nil).Times(1).AnyTimes()
		_, err := s.CreateAccount(context.Background(), in)
		r.NoError(err)
		r.NoError(mock.ExpectationsWereMet())
	})

	r.NoError(err)
}

func TestAccount_GetAccountByID(t *testing.T) {

	r := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	r.NoError(err)

	rows := mock.NewRows([]string{"guid", "user_id", "balance", "created_at", "updated_at"}).AddRow("TestUserID", "TestUserID", 0.0, "2021-01-01", "2021-01-01")
	mock.ExpectQuery(`^SELECT (.+?) FROM accounts * `).WithArgs("TestUserID").WillReturnRows(rows)

	repo := mock_storage.NewMockAccountRepoI(ctrl)
	s := NewService(config.Config{}, zap.NewNop(), postgres.NewStore(db))

	in := &models.GetAccountByIDRequest{
		ID: "TestUserID",
	}
	mockresp := &models.Account{
		ID: "TestUserID",
	}

	t.Run("SUCCESS", func(t *testing.T) {

		repo.EXPECT().GetAccountByID(context.Background(), in).Return(mockresp, nil).Times(1).AnyTimes()
		_, err := s.GetAccountByID(context.Background(), in)
		r.NoError(err)
		r.NoError(mock.ExpectationsWereMet())
	})

	r.NoError(err)
}

func TestAccount_GetAccountsByUserID(t *testing.T) {

	r := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	r.NoError(err)

	rows := mock.NewRows([]string{"guid", "user_id", "balance", "created_at", "updated_at", "count"}).AddRow("TestUserID", "TestUserID", 0.0, "2021-01-01", "2021-01-01", 1)
	mock.ExpectQuery(`^SELECT (.+?) FROM accounts * `).WithArgs("TestUserID").WillReturnRows(rows)

	repo := mock_storage.NewMockAccountRepoI(ctrl)
	s := NewService(config.Config{}, zap.NewNop(), postgres.NewStore(db))

	in := &models.GetAccountsByUserIDRequest{
		UserID: "TestUserID",
	}
	mockresp := &models.GetAccountsByUserIDResponse{
		Accounts: []*models.Account{
			{
				ID: "TestUserID",
			},
		},
	}

	t.Run("SUCCESS", func(t *testing.T) {

		repo.EXPECT().GetAccountsByUserID(context.Background(), in).Return(mockresp, nil).Times(1).AnyTimes()
		_, err := s.GetAccountsByUserID(context.Background(), in)
		r.NoError(err)
		r.NoError(mock.ExpectationsWereMet())
	})

	r.NoError(err)
}

func TestAccount_GetAccountTransactions(t *testing.T) {

	r := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	r.NoError(err)

	rows := mock.NewRows([]string{"guid", "account_id", "transaction_amount", "transaction_type", "recipient_id", "created_at", "count"}).AddRow("TestTransactionID", "TestUserID", 0.0, "TestType", "TestUserID", "2021-01-01", 1)
	mock.ExpectQuery(`^SELECT (.+?) FROM transactions * `).WithArgs("TestUserID").WillReturnRows(rows)

	repo := mock_storage.NewMockTxRepoI(ctrl)

	s := NewService(config.Config{}, zap.NewNop(), postgres.NewStore(db))

	t.Run("SUCCESS", func(t *testing.T) {
		ctx := context.Background()

		in := &models.GetTransactionsByAccountIDRequest{
			AccountID: "TestUserID",
		}
		mockresp := &models.GetTransactionsByAccountIDResponse{
			Transactions: []*models.Transaction{
				{
					ID:        "TestTransactionID",
					AccountID: "TestUserID",
					Amount:    0.0,
					Type:      "TestType",
				},
			},
		}

		repo.EXPECT().GetTransactionsByAccountID(ctx, in).Return(mockresp, nil).Times(1).AnyTimes()
		_, err := s.GetAccountTransactions(ctx, in)

		r.NoError(err)
		r.NoError(mock.ExpectationsWereMet())
	})

	r.NoError(err)
}

func TestAccount_GetAccountTransactionByID(t *testing.T) {

	r := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	r.NoError(err)

	rows := mock.NewRows([]string{"guid", "account_id", "transaction_amount", "transaction_type", "recipient_id", "created_at"}).AddRow("TestTransactionID", "TestUserID", 0.0, "TestType", "TestUserID", "2021-01-01")
	mock.ExpectQuery(`^SELECT (.+?) FROM transactions * `).WithArgs("TestTransactionID", "TestUserID").WillReturnRows(rows)

	repo := mock_storage.NewMockTxRepoI(ctrl)

	s := NewService(config.Config{}, zap.NewNop(), postgres.NewStore(db))

	t.Run("SUCCESS", func(t *testing.T) {
		ctx := context.Background()

		in := &models.GetTransactionByIDRequest{
			ID:        "TestTransactionID",
			AccountID: "TestUserID",
		}
		mockresp := &models.Transaction{
			ID:        "TestTransactionID",
			AccountID: "TestUserID",
			Amount:    0.0,
			Type:      "TestType",
		}

		repo.EXPECT().GetTransactionByID(ctx, in).Return(mockresp, nil).Times(1).AnyTimes()
		_, err := s.GetAccountTransactionByID(ctx, in)

		r.NoError(err)
		r.NoError(mock.ExpectationsWereMet())
	})

	r.NoError(err)
}
