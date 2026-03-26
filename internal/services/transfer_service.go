package services

import (
	"context"
	"database/sql"
	"errors"
	"money-tracker/internal/models"
	"money-tracker/internal/repository"

	"github.com/google/uuid"
)

type TransferService struct {
	repo *repository.TransferRepository
	db   *sql.DB
}

func NewTransferService(r *repository.TransferRepository, db *sql.DB) *TransferService {
	return &TransferService{
		repo: r,
		db:   db,
	}
}

func (s *TransferService) CreateTransfer(ctx context.Context, req *models.Transfer) (*models.Transfer, error) {
	if req.FromAccountID == req.ToAccountID {
		return nil, ErrTransferSameAccount
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrInternal
	}

	defer func() { _ = tx.Rollback() }()

	if err := s.repo.CreateTransferTx(ctx, tx, req); err != nil {
		return nil, ErrInternal
	}

	resFrom, err := tx.ExecContext(ctx,
		`UPDATE accounts SET balance = balance - $1, updated_at = NOW() 
         WHERE id = $2 AND user_id = $3`,
		req.Amount, req.FromAccountID, req.UserID)
	if err != nil {
		return nil, ErrInternal
	}

	rowsFrom, _ := resFrom.RowsAffected()
	if rowsFrom == 0 {
		return nil, ErrAccountNotFound
	}

	resTo, err := tx.ExecContext(ctx,
		`UPDATE accounts SET balance = balance + $1, updated_at = NOW() 
         WHERE id = $2 AND user_id = $3`,
		req.Amount, req.ToAccountID, req.UserID)

	if err != nil {
		return nil, ErrInternal
	}
	rowsTo, _ := resTo.RowsAffected()
	if rowsTo == 0 {
		return nil, ErrAccountNotFound
	}

	if err := tx.Commit(); err != nil {
		return nil, ErrInternal
	}

	return req, nil
}

func (s *TransferService) UpdateTransfer(ctx context.Context, req *models.Transfer) (*models.Transfer, error) {
	if req.FromAccountID == req.ToAccountID {
		return nil, ErrTransferSameAccount
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrInternal
	}

	defer func() { _ = tx.Rollback() }()

	oldTransfer, err := s.repo.GetTransferWithTx(ctx, tx, req.ID, req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransferNotFound
		}
		return nil, ErrInternal
	}

	_, err = tx.ExecContext(ctx, `UPDATE accounts SET balance = balance + $1 WHERE id = $2 AND user_id = $3`,
		oldTransfer.Amount, oldTransfer.FromAccountID, req.UserID)
	if err != nil {
		return nil, ErrInternal
	}

	_, err = tx.ExecContext(ctx, `UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND user_id = $3`,
		oldTransfer.Amount, oldTransfer.ToAccountID, req.UserID)
	if err != nil {
		return nil, ErrInternal
	}

	if err := s.repo.UpdateTransferTx(ctx, tx, req); err != nil {
		return nil, ErrInternal
	}

	resFrom, err := tx.ExecContext(ctx, `UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND user_id = $3`,
		req.Amount, req.FromAccountID, req.UserID)
	if err != nil {
		return nil, ErrInternal
	}

	if r, _ := resFrom.RowsAffected(); r == 0 {
		return nil, ErrAccountNotFound
	}

	_, err = tx.ExecContext(ctx, `UPDATE accounts SET balance = balance + $1 WHERE id = $2 AND user_id = $3`,
		req.Amount, req.ToAccountID, req.UserID)
	if err != nil {
		return nil, ErrInternal
	}

	if err := tx.Commit(); err != nil {
		return nil, ErrInternal
	}
	return req, nil
}

func (s *TransferService) DeleteTransfer(ctx context.Context, id, userID uuid.UUID) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return ErrInternal
	}
	defer func() { _ = tx.Rollback() }()

	transfer, err := s.repo.GetTransferWithTx(ctx, tx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTransferNotFound
		}
		return ErrInternal
	}

	_, err = tx.ExecContext(ctx, `UPDATE accounts SET balance = balance + $1 WHERE id = $2 AND user_id = $3`,
		transfer.Amount, transfer.FromAccountID, userID)
	if err != nil {
		return ErrInternal
	}

	_, err = tx.ExecContext(ctx, `UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND user_id = $3`,
		transfer.Amount, transfer.ToAccountID, userID)
	if err != nil {
		return ErrInternal
	}

	if err := s.repo.DeleteTransferTx(ctx, tx, id, userID); err != nil {
		return ErrInternal
	}

	return tx.Commit()
}

func (s *TransferService) GetTransfers(ctx context.Context, userID uuid.UUID) ([]*models.Transfer, error) {
	transfers, err := s.repo.GetTransfers(ctx, userID)
	if err != nil {
		return nil, ErrInternal
	}
	return transfers, nil
}

func (s *TransferService) GetTransfer(ctx context.Context, transferID, userID uuid.UUID) (*models.Transfer, error) {
	transfer, err := s.repo.GetTransfer(ctx, transferID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransferNotFound
		}
	}
	return transfer, nil
}
