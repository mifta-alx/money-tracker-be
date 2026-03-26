package repository

import (
	"context"
	"database/sql"
	"money-tracker/internal/models"

	"github.com/google/uuid"
)

type TransferRepository struct {
	DB *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{DB: db}
}

func (r *TransferRepository) CreateTransferTx(ctx context.Context, tx *sql.Tx, t *models.Transfer) error {
	query := `INSERT INTO transfers (user_id, from_account_id, to_account_id, amount, notes) 
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at`

	return tx.QueryRowContext(ctx, query,
		t.UserID, t.FromAccountID, t.ToAccountID, t.Amount, t.Notes).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TransferRepository) UpdateTransferTx(ctx context.Context, tx *sql.Tx, t *models.Transfer) error {
	query := `UPDATE transfers SET from_account_id = $1, to_account_id = $2, amount = $3, notes = $4, updated_at = NOW() WHERE id = $5 AND user_id = $6 RETURNING id, created_at, updated_at`
	return tx.QueryRowContext(ctx, query, t.FromAccountID, t.ToAccountID, t.Amount, t.Notes, t.ID, t.UserID).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TransferRepository) DeleteTransferTx(ctx context.Context, tx *sql.Tx, id, userID uuid.UUID) error {
	query := `DELETE FROM transfers WHERE id = $1 AND user_id = $2`
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

func (r *TransferRepository) GetTransferWithTx(ctx context.Context, tx *sql.Tx, id, usrID uuid.UUID) (*models.Transfer, error) {
	query := `SELECT id, from_account_id, to_account_id, amount, notes, created_at, updated_at FROM transfers WHERE id = $1 AND user_id = $2`
	var transfer models.Transfer
	err := tx.QueryRowContext(ctx, query, id, usrID).Scan(
		&transfer.ID, &transfer.FromAccountID, &transfer.ToAccountID, &transfer.Amount, &transfer.Notes, &transfer.CreatedAt, &transfer.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

func (r *TransferRepository) GetTransfers(ctx context.Context, userID uuid.UUID) ([]*models.Transfer, error) {
	query := `SELECT id, from_account_id, to_account_id, amount, notes, created_at, updated_at FROM transfers WHERE user_id = $1 ORDER BY created_at`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	var transfers []*models.Transfer
	for rows.Next() {
		var tf models.Transfer
		err := rows.Scan(&tf.ID, &tf.FromAccountID, &tf.ToAccountID, &tf.Amount, &tf.Notes, &tf.CreatedAt, &tf.UpdatedAt)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, &tf)
	}

	if transfers == nil {
		transfers = []*models.Transfer{}
	}

	return transfers, nil
}

func (r *TransferRepository) GetTransfer(ctx context.Context, id, userID uuid.UUID) (*models.Transfer, error) {
	query := `SELECT id, from_account_id, to_account_id, amount, notes, created_at, updated_at FROM transfers WHERE id = $1 AND user_id = $2`
	var transfer models.Transfer
	err := r.DB.QueryRowContext(ctx, query, id, userID).Scan(&transfer.ID, &transfer.FromAccountID, &transfer.ToAccountID, &transfer.Amount, &transfer.Notes, &transfer.CreatedAt, &transfer.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}
