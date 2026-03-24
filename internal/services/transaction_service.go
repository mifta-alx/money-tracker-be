package services

import (
	"context"
	"database/sql"
	"errors"
	"money-tracker/internal/models"
	"money-tracker/internal/repository"

	"github.com/google/uuid"
)

type TransactionService struct {
	repo *repository.TransactionRepository
	db   *sql.DB
}

func NewTransactionService(r *repository.TransactionRepository, db *sql.DB) *TransactionService {
	return &TransactionService{
		repo: r,
		db:   db,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req *models.Transaction) (*models.Transaction, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrInternal
	}

	defer func() { _ = tx.Rollback() }()

	if err := s.repo.CreateTransactionTx(ctx, tx, req); err != nil {
		return nil, ErrInternal
	}

	adjustment := req.Amount
	if req.Type == "expense" {
		adjustment = -req.Amount
	}

	query := `UPDATE accounts SET balance = balance + $1, updated_at = NOW() 
              WHERE id = $2 AND user_id = $3`

	result, err := tx.ExecContext(ctx, query, adjustment, req.AccountID, req.UserID)
	if err != nil {
		return nil, ErrInternal
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, ErrAccountNotFound
	}
	if err := tx.Commit(); err != nil {
		return nil, ErrInternal
	}

	return req, nil
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, req *models.Transaction) (*models.Transaction, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrInternal
	}

	defer func() { _ = tx.Rollback() }()

	oldData, err := s.repo.GetTransactionWithTx(ctx, tx, req.ID, req.UserID)
	if err != nil {
		return nil, ErrTransactionNotFound
	}

	var oldAdjustment int64
	if oldData.Type == "expense" {
		oldAdjustment = oldData.Amount
	} else {
		oldAdjustment = -oldData.Amount
	}

	resOld, err := tx.ExecContext(ctx,
		`UPDATE accounts SET balance = balance + $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`, oldAdjustment, oldData.AccountID, req.UserID)
	if err != nil {
		return nil, ErrInternal
	}

	if rows, _ := resOld.RowsAffected(); rows == 0 {
		return nil, ErrInternal
	}

	if err := s.repo.UpdateTransactionTx(ctx, tx, req); err != nil {
		return nil, ErrInternal
	}

	var newAdjustment int64
	if req.Type == "expense" {
		newAdjustment = -req.Amount
	} else {
		newAdjustment = req.Amount
	}

	resNew, err := tx.ExecContext(ctx,
		`UPDATE accounts SET balance = balance + $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`,
		newAdjustment, req.AccountID, req.UserID)

	if err != nil {
		return nil, ErrInternal
	}

	if rows, _ := resNew.RowsAffected(); rows == 0 {
		return nil, ErrAccountNotFound
	}

	if err := tx.Commit(); err != nil {
		return nil, ErrInternal
	}

	return req, nil
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, transactionID, userID uuid.UUID) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return ErrInternal
	}

	defer func() { _ = tx.Rollback() }()

	oldData, err := s.repo.GetTransactionWithTx(ctx, tx, transactionID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTransactionNotFound
		}
		return ErrInternal
	}

	oldAdjustment := oldData.Amount
	if oldData.Type == "income" {
		oldAdjustment = -oldData.Amount
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE accounts SET balance = balance + $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`, oldAdjustment, oldData.AccountID, oldData.UserID)

	if err != nil {
		return ErrInternal
	}

	if err := s.repo.DeleteTransactionTx(ctx, tx, transactionID, userID); err != nil {
		return ErrInternal
	}
	return tx.Commit()
}

func (s *TransactionService) GetTransactions(ctx context.Context, userID uuid.UUID) ([]*models.Transaction, error) {
	transactions, err := s.repo.GetTransactions(ctx, userID)
	if err != nil {
		return nil, ErrInternal
	}
	return transactions, nil
}

func (s *TransactionService) GetTransaction(ctx context.Context, transactionID, userID uuid.UUID) (*models.Transaction, error) {
	transaction, err := s.repo.GetTransaction(ctx, transactionID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
	}
	return transaction, nil
}
