package service

import (
	"github.com/dilmurodov/online_banking/config"
	"github.com/dilmurodov/online_banking/internal/service/account"
	payment "github.com/dilmurodov/online_banking/internal/service/pament"
	"github.com/dilmurodov/online_banking/internal/service/user"
	"github.com/dilmurodov/online_banking/pkg/logger"
	"github.com/dilmurodov/online_banking/storage"
)

type ServiceManagerI interface {
	UserService() user.ServiceI
	AccountService() account.ServiceI
	PaymentService() payment.ServiceI
}

type serviceManager struct {
	userService    user.ServiceI
	accountService account.ServiceI
	paymentService payment.ServiceI
}

func NewServiceManager(cfg config.Config, log logger.LoggerI, strg storage.StorageI) ServiceManagerI {

	userService := user.NewService(cfg, log, strg)
	accountService := account.NewService(cfg, log, strg)
	paymentService := payment.NewService(cfg, log, strg)

	return &serviceManager{
		userService:    userService,
		accountService: accountService,
		paymentService: paymentService,
	}
}

func (s *serviceManager) UserService() user.ServiceI {
	return s.userService
}

func (s *serviceManager) AccountService() account.ServiceI {
	return s.accountService
}

func (s *serviceManager) PaymentService() payment.ServiceI {
	return s.paymentService
}
