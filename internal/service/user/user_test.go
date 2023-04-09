package user

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dilmurodov/online_banking/config"
	"github.com/dilmurodov/online_banking/pkg/models"
	"github.com/dilmurodov/online_banking/pkg/security"
	mock_storage "github.com/dilmurodov/online_banking/storage/mock"
	"github.com/dilmurodov/online_banking/storage/postgres"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestUser_GetUserByID(t *testing.T) {

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

	rows := mock.NewRows([]string{"guid", "first_name", "last_name", "phone", "created_at", "updated_at"}).AddRow("TestUserID", "TestFirstName", "TestLastName", "TestPhone", "2021-01-01 00:00:00", "2021-01-01 00:00:00")

	mock.ExpectQuery(`^SELECT (.+?) FROM "users" * `).WithArgs("TestUserID").WillReturnRows(rows)

	repo := mock_storage.NewMockUserRepoI(ctrl)

	t.Run("SUCCESS", func(t *testing.T) {

		ctx := context.Background()

		in := &models.GetUserByIDRequest{
			UserId: "TestUserID",
		}
		resp := &models.GetUserByIDResponse{
			User: &models.User{
				Guid:      "TestUserID",
				FirstName: "TestFirstName",
				LastName:  "TestLastName",
				Phone:     "TestPhone",
			},
		}

		repo.EXPECT().GetUserByID(ctx, in).Return(resp, nil).Times(1).AnyTimes()

		resp, err := s.GetUserByID(ctx, in)
		r.NoError(err)
		r.NotNil(resp)
	})
}

func TestUser_CreateUser(t *testing.T) {

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

	rows := mock.NewRows([]string{"guid", "first_name", "last_name", "phone", "created_at", "updated_at"}).AddRow("TestUserId", "TestFirstName", "TestLastName", "TestPhone", "2021-01-01 00:00:00", "2021-01-01 00:00:00")

	mock.ExpectQuery(`^INSERT INTO "users" (.+?) * `).WithArgs("TestFirstName", "TestLastName", "TestPhone", "TestPassword").
		WillReturnRows(rows)

	repo := mock_storage.NewMockUserRepoI(ctrl)

	t.Run("SUCCESS", func(t *testing.T) {

		ctx := context.Background()

		in := &models.CreateUserRequest{
			User: &models.User{
				FirstName: "TestFirstName",
				LastName:  "TestLastName",
				Phone:     "TestPhone",
				Password:  "TestPassword",
			},
		}
		resp := &models.User{
			Guid:      "TestUserId",
			FirstName: "TestFirstName",
			LastName:  "TestLastName",
			Phone:     "TestPhone",
			Password:  "TestPassword",
			CreatedAt: "2021-01-01 00:00:00",
			UpdatedAt: "2021-01-01 00:00:00",
		}

		repo.EXPECT().CreateUser(ctx, in).Return(resp, nil).Times(1).AnyTimes()

		_, err := s.CreateUser(ctx, in)
		r.NoError(err)
	})
}

func TestUser_GetUserByCredentials(t *testing.T) {

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

	hpass, err := security.HashPassword("TestPassword")
	r.NoError(err)

	rows := mock.NewRows([]string{"guid", "first_name", "last_name", "phone", "password", "created_at", "updated_at"}).AddRow("TestUserId", "TestFirstName", "TestLastName", "TestPhone", hpass, "2021-01-01 00:00:00", "2021-01-01 00:00:00")

	mock.ExpectQuery(`^SELECT (.+?) FROM "users" * `).WithArgs("TestPhone").WillReturnRows(rows)

	repo := mock_storage.NewMockUserRepoI(ctrl)

	t.Run("SUCCESS", func(t *testing.T) {

		ctx := context.Background()

		in := &models.GetByCredentialsRequest{
			Phone:    "TestPhone",
			Password: "TestPassword",
		}

		mockresp := &models.User{
			Guid:      "TestUserId",
			FirstName: "TestFirstName",
			LastName:  "TestLastName",
			Phone:     "TestPhone",
			Password:  hpass,
			CreatedAt: "2021-01-01 00:00:00",
			UpdatedAt: "2021-01-01 00:00:00",
		}

		repo.EXPECT().GetUserPasswordByPhone(ctx, "TestPhone").Return(mockresp, nil).Times(1).AnyTimes()

		expected, err := s.GetUserByCredentials(ctx, in)
		r.NoError(err)
		r.Equal(expected, mockresp)

	})
}
