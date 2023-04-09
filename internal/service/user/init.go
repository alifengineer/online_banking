package user

import (
	"context"

	"github.com/dilmurodov/online_banking/config"
	"github.com/dilmurodov/online_banking/pkg/logger"
	"github.com/dilmurodov/online_banking/pkg/models"
	"github.com/dilmurodov/online_banking/storage"
)

type ServiceI interface {
	GetUserByID(context.Context, *models.GetUserByIDRequest) (*models.GetUserByIDResponse, error)
	CreateUser(context.Context, *models.CreateUserRequest) (*models.User, error)
	GetUserByCredentials(ctx context.Context, req *models.GetByCredentialsRequest) (*models.User, error)
	GetUserPasswordByPhone(ctx context.Context, phone string) (resp *models.User, err error)
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
