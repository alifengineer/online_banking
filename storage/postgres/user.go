package postgres

import (
	"context"
	"database/sql"

	"github.com/dilmurodov/online_banking/pkg/customerrors"
	"github.com/dilmurodov/online_banking/pkg/models"
	_ "github.com/lib/pq"
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
		return nil, &customerrors.UserNotFoundError{Guid: req.UserId}
	} else if err != nil {
		return nil, &customerrors.InternalServerError{Message: err.Error()}
	}

	resp.User = user

	return resp, nil
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
		FROM "users"
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
		return nil, &customerrors.UserNotFoundWithPhoneError{Phone: phone}
	} else if err != nil {
		return nil, &customerrors.InternalServerError{Message: err.Error()}
	}

	return resp, nil
}

func (u *userRepo) CreateUser(ctx context.Context, req *models.CreateUserRequest) (resp *models.User, err error) {
	resp = &models.User{}

	query := `
		INSERT INTO "users" (
			first_name,
			last_name,
			phone,
			password
		) VALUES (
			$1, $2, $3, $4
		) RETURNING guid, first_name, last_name, phone, created_at, updated_at
	`

	row := u.db.QueryRowContext(ctx, query,
		req.User.FirstName,
		req.User.LastName,
		req.User.Phone,
		req.User.Password,
	)

	err = row.Scan(
		&resp.Guid,
		&resp.FirstName,
		&resp.LastName,
		&resp.Phone,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)
	if err != nil {
		return nil, &customerrors.InternalServerError{Message: err.Error()}
	}

	return resp, nil
}
