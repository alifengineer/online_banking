package postgres

import (
	"context"
	"fmt"

	"github.com/dilmurodov/online_banking/config"
	"github.com/dilmurodov/online_banking/storage"

	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	db          *sql.DB
	userRepo    *userRepo
	accountRepo *accountRepo
	txRepo      *txRepo
}

func NewPostgres(ctx context.Context, cfg config.Config) (storage.StorageI, error) {

	pgUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	)

	db, err := sql.Open("postgres", pgUrl)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return NewStore(db), nil
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:          db,
		accountRepo: &accountRepo{db: db},
		userRepo:    &userRepo{db: db},
		txRepo:      &txRepo{db: db},
	}
}

func (s *Store) CloseDB() {
	s.db.Close()
}

func (s *Store) User() storage.UserRepoI {
	if s.userRepo != nil {
		return NewUserRepo(s.db)
	}
	return s.userRepo
}

func (s *Store) Account() storage.AccountRepoI {
	if s.accountRepo != nil {
		return NewAccountRepo(s.db)
	}
	return s.accountRepo
}

func (s *Store) TxRepo() storage.TxRepoI {
	if s.txRepo != nil {
		return NewTxRepo(s.db)
	}
	return s.txRepo
}
