package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dilmurodov/online_banking/pkg/customerrors"
	"github.com/dilmurodov/online_banking/pkg/models"
	"github.com/dilmurodov/online_banking/pkg/util"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type txRepo struct {
	db *sql.DB
}

func NewTxRepo(db *sql.DB) *txRepo {
	return &txRepo{db: db}
}

func (s *txRepo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

// CreateTransaction creates a new transaction in a given transaction object
func (r *txRepo) CreateTransaction(ctx context.Context, tx *sql.Tx, transaction *models.Transaction) (resp *models.Transaction, err error) {
	resp = &models.Transaction{}
	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO transactions 
			(account_id, 
			transaction_amount,
			recipient_id,
			transaction_type
		) 
		VALUES ($1, $2, $3, $4) 
		RETURNING 
			guid, 
			transaction_amount, 
			transaction_type, 
			recipient_id, 
			created_at`)
	if err != nil {
		return nil, &customerrors.InternalServerError{Message: err.Error()}
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx,
		transaction.AccountID,
		transaction.Amount,
		transaction.RecipientID,
		transaction.Type,
	)
	err = row.Scan(
		&resp.ID,
		&resp.Amount,
		&resp.Type,
		&resp.RecipientID,
		&resp.CreatedAt,
	)
	if err != nil {
		return nil, &customerrors.InternalServerError{Message: err.Error()}
	}
	return resp, nil
}

// GetTransactionsByAccountID returns all transactions associated with the given account ID
func (r *txRepo) GetTransactionsByAccountID(ctx context.Context, req *models.GetTransactionsByAccountIDRequest) (*models.GetTransactionsByAccountIDResponse, error) {
	var (
		count         int
		params        = []interface{}{}
		doneTimestamp sql.NullString
	)
	transactions := make([]*models.Transaction, 0)

	query :=
		`SELECT 
			guid, 
			account_id, 
			transaction_amount, 
			transaction_type,
			recipient_id, 
			created_at,
			count(1) filter (where deleted_at IS NULL) OVER() AS count,
			approved,
			done,
			done_timestamp
		FROM transactions`

	filter := ` WHERE account_id=$1 AND deleted_at IS NULL`
	params = append(params, req.AccountID)

	order := ` ORDER BY created_at`
	limit := ` LIMIT 10`
	offset := ` OFFSET 0`

	if util.IsValidTimeStamp(req.From) {
		filter += " AND created_at >= $3"
		params = append(params, req.From)
	}

	if util.IsValidTimeStamp(req.To) {
		filter += " AND created_at <= $4"
		params = append(params, req.To)
	}

	if util.IsValidUUID(req.RecipientID) {
		filter += " AND recipient_id = $5"
		params = append(params, req.RecipientID)
	}

	switch req.Desc {
	case false:
		order += " ASC"
	default:
		order += " DESC"
	}

	if req.Limit != 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Offset != 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	query = query + filter + order + limit + offset
	rows, err := r.db.QueryContext(
		ctx,
		query, params...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query rows")
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID,
			&t.AccountID,
			&t.Amount,
			&t.Type,
			&t.RecipientID,
			&t.CreatedAt,
			&count,
			&t.Approved,
			&t.Done,
			&doneTimestamp,
		)
		if err != nil {
			return nil, &customerrors.InternalServerError{Message: err.Error()}
		}
		t.DoneTimestamp = doneTimestamp.String
		transactions = append(transactions, &t)
	}
	if err != nil {
		return nil, &customerrors.InternalServerError{Message: err.Error()}
	}

	return &models.GetTransactionsByAccountIDResponse{
		Transactions: transactions,
		Count:        count,
	}, nil
}

func (r *txRepo) GetTransactionByID(ctx context.Context, req *models.GetTransactionByIDRequest) (*models.Transaction, error) {
	var (
		createdAt     sql.NullString
		doneTimestamp sql.NullString
	)
	t := &models.Transaction{}
	query :=
		`SELECT 
			guid, 
			account_id, 
			transaction_amount,
			transaction_type, 
			recipient_id, 
			created_at,
			approved,
			done,
			done_timestamp
		FROM transactions
		WHERE guid=$1 AND account_id=$2 AND deleted_at IS NULL`

	err := r.db.QueryRowContext(
		ctx,
		query,
		req.ID,
		req.AccountID,
	).Scan(
		&t.ID,
		&t.AccountID,
		&t.Amount,
		&t.Type,
		&t.RecipientID,
		&createdAt,
		&t.Approved,
		&t.Done,
		&doneTimestamp,
	)
	if err != nil && err == sql.ErrNoRows {
		return nil, &customerrors.TransactionNotFoundError{Guid: req.ID}
	} else if err != nil {
		return nil, &customerrors.InternalServerError{Message: err.Error()}
	}
	t.CreatedAt = createdAt.String
	t.DoneTimestamp = doneTimestamp.String

	return t, nil
}

func (r *txRepo) GetTransactionsByIDS(ctx context.Context, req *models.GetTransactionsByIDSRequest) (resp *models.GetTransactionsByIDSResponse, err error) {
	var (
		createdAt     sql.NullString
		doneTimestamp sql.NullString
	)
	transactions := make([]*models.Transaction, 0)
	resp = &models.GetTransactionsByIDSResponse{
		Transactions: transactions,
	}

	query :=
		`SELECT 
			guid, 
			account_id, 
			transaction_amount,
			transaction_type, 
			recipient_id, 
			created_at,
			approved,
			done,
			done_timestamp
		FROM transactions
		WHERE guid=ANY($1) AND account_id=$2 AND deleted_at IS NULL`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		pq.Array(req.IDS),
		req.AccountID,
	)
	if err != nil {
		return nil, &customerrors.InternalServerError{Message: err.Error()}
	}
	defer rows.Close()

	for rows.Next() {
		t := &models.Transaction{}
		err := rows.Scan(
			&t.ID,
			&t.AccountID,
			&t.Amount,
			&t.Type,
			&t.RecipientID,
			&createdAt,
			&t.Approved,
			&t.Done,
			&doneTimestamp,
		)
		if err != nil {
			return nil, &customerrors.InternalServerError{Message: err.Error()}
		}
		t.CreatedAt = createdAt.String
		t.DoneTimestamp = doneTimestamp.String
		transactions = append(transactions, t)
	}
	if err = rows.Err(); err != nil {
		return nil, &customerrors.InternalServerError{Message: err.Error()}
	}
	resp.Transactions = transactions

	return resp, nil
}

func (r *txRepo) ApproveTransactions(ctx context.Context, tx *sql.Tx, req *models.ApproveTransactionsRequest) (err error) {

	query :=
		`UPDATE transactions SET 
			approved=true, 
			done=true, 
			done_timestamp=CURRENT_TIMESTAMP 
		WHERE 
			guid=ANY($1) AND account_id=$2 AND deleted_at IS NULL AND approved=false AND done=false`

	result, err := tx.ExecContext(
		ctx,
		query,
		pq.Array(req.TransactionIDS),
		req.AccountID,
	)
	if err != nil {
		return &customerrors.InternalServerError{Message: err.Error()}
	}
	if cn, err := result.RowsAffected(); err != nil || cn == 0 {
		return &customerrors.TransactionNotFoundError{Guid: req.TransactionIDS[0]}
	}

	return nil
}
