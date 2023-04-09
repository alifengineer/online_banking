package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/dilmurodov/online_banking/pkg/models"
	_ "github.com/lib/pq"

	"github.com/google/uuid"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (u *userRepo) GetUserByID(ctx context.Context, req *models.GetUserByIDRequest) (resp *models.GetUserByIDResponse, err error) {

	resp = &models.GetUserByIDResponse{}

	params := []interface{}{}

	query := `
		SELECT
			guid,
			first_name, 
			last_name, 
			phone, 
			created_at, 
			updated_at
		FROM "users"
		`

	filter := `WHERE guid = $1 AND deleted_at = 0`
	params = append(params, req.UserId)

	if req.Phone != "" {
		filter += " AND phone = $2"
		params = append(params, req.Phone)
	}

	query += filter

	row := u.db.QueryRowContext(ctx, query, params...)

	user := &models.User{}
	err = row.Scan(
		&user.Guid,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}
	resp.User = user

	return resp, err
}

// Get User By Phone

func (u *userRepo) GetUserPasswordByPhone(ctx context.Context, phone string) (resp *models.User, err error) {

	query := `
		SELECT
			guid,
			first_name,
			last_name,
			phone,
			password,
			created_at,
			updated_at
		FROM users
		WHERE phone = $1 AND deleted_at = 0
	`

	row := u.db.QueryRowContext(ctx, query, phone)

	resp = &models.User{}
	err = row.Scan(
		&resp.Guid,
		&resp.FirstName,
		&resp.LastName,
		&resp.Phone,
		&resp.Password,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)

	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}

	return resp, err
}

func (u *userRepo) CreateUser(ctx context.Context, req *models.CreateUserRequest) (resp *models.User, err error) {
	resp = &models.User{}

	id := uuid.New()

	query := `
		INSERT INTO users (
			guid,
			first_name,
			last_name,
			phone,
			password,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	_, err = u.db.ExecContext(ctx, query,
		id.String(),
		req.User.FirstName,
		req.User.LastName,
		req.User.Phone,
		req.User.Password,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
	)

	u.db.QueryRowContext(ctx,
		`SELECT guid, first_name, last_name, phone, created_at, updated_at FROM users WHERE guid = $1`,
		id,
	).Scan(
		&resp.Guid,
		&resp.FirstName,
		&resp.LastName,
		&resp.Phone,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)

	return resp, err
}
