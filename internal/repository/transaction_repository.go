package repository

import (
	"context"
	"database/sql"
	"money-tracker/internal/models"

	"github.com/google/uuid"
)

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		DB: db,
	}
}

func (r *TransactionRepository) CreateTransactionTx(ctx context.Context, tx *sql.Tx, t *models.Transaction) error {
	query := `INSERT INTO transactions (user_id, account_id, category_id, title, type, amount, date, notes) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
              RETURNING id, created_at, updated_at`

	return tx.QueryRowContext(ctx, query,
		t.UserID, t.AccountID, t.CategoryID, t.Title, t.Type, t.Amount, t.Date, t.Notes,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TransactionRepository) UpdateTransactionTx(ctx context.Context, tx *sql.Tx, t *models.Transaction) error {
	query := `UPDATE transactions SET account_id = $1, category_id = $2, title = $3, type = $4, amount= $5, date=$6, notes = $7, updated_at = NOW() WHERE id = $8 AND user_id = $9 RETURNING id, created_at, updated_at`
	return tx.QueryRowContext(ctx, query,
		t.AccountID, t.CategoryID, t.Title, t.Type, t.Amount, t.Date, t.Notes, t.ID, t.UserID).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TransactionRepository) DeleteTransactionTx(ctx context.Context, tx *sql.Tx, id, userID uuid.UUID) error {
	query := `DELETE FROM transactions WHERE id = $1 AND user_id = $2`
	result, err := tx.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (r *TransactionRepository) GetTransactionWithTx(ctx context.Context, tx *sql.Tx, id, usrID uuid.UUID) (*models.Transaction, error) {
	query := `SELECT id, user_id, account_id, category_id, title, type, amount, date, notes, created_at, updated_at 
              FROM transactions WHERE id = $1 AND user_id = $2`
	var t models.Transaction
	err := tx.QueryRowContext(ctx, query, id, usrID).Scan(
		&t.ID, &t.UserID, &t.AccountID, &t.CategoryID, &t.Title, &t.Type, &t.Amount, &t.Date, &t.Notes, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TransactionRepository) GetTransactions(ctx context.Context, userID uuid.UUID) ([]*models.Transaction, error) {
	query := `SELECT id, account_id, category_id, title, type, amount, date, notes, created_at, updated_at FROM transactions WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	var transactions []*models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.AccountID, &transaction.CategoryID, &transaction.Title, &transaction.Type, &transaction.Amount, &transaction.Date, &transaction.Notes, &transaction.CreatedAt, &transaction.UpdatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	if transactions == nil {
		transactions = []*models.Transaction{}
	}
	return transactions, nil
}

func (r *TransactionRepository) GetTransaction(ctx context.Context, id, usrID uuid.UUID) (*models.Transaction, error) {
	query := `SELECT id, account_id, category_id, title, type, amount, date, notes, created_at, updated_at FROM transactions WHERE id = $1 AND user_id = $2`
	var transaction models.Transaction
	err := r.DB.QueryRowContext(ctx, query, id, usrID).Scan(
		&transaction.ID,
		&transaction.AccountID,
		&transaction.CategoryID,
		&transaction.Title,
		&transaction.Type,
		&transaction.Amount,
		&transaction.Date,
		&transaction.Notes,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
