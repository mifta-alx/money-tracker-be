package services

import (
	"context"
	"database/sql"
	"errors"
	"money-tracker/internal/models"
	"money-tracker/internal/repository"

	"github.com/google/uuid"
)

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(r *repository.AccountRepository) *AccountService {
	return &AccountService{r}
}

func (s *AccountService) CreateAccount(ctx context.Context, req *models.Account) (*models.Account, error) {
	if req.Name == "" || req.UserID == uuid.Nil || req.Icon == "" {
		return nil, ErrMissingRequiredFields
	}
	if req.Balance < 0 {
		return nil, ErrBalanceCannotBeNegative
	}

	err := s.repo.CreateAccount(ctx, req)
	if err != nil {
		return nil, ErrInternal
	}

	return req, nil
}

func (s *AccountService) UpdateAccount(ctx context.Context, req *models.Account) (*models.Account, error) {
	err := s.repo.UpdateAccount(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, ErrInternal
	}
	return req, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, accountID, userID uuid.UUID) error {
	err := s.repo.DeleteAccount(ctx, accountID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAccountNotFound
		}
		return ErrInternal
	}
	return nil
}

func (s *AccountService) GetAccounts(ctx context.Context, userID uuid.UUID) ([]*models.Account, error) {
	accounts, err := s.repo.GetAccounts(ctx, userID)
	if err != nil {
		return nil, ErrInternal
	}
	return accounts, nil
}

func (s *AccountService) GetAccount(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Account, error) {
	acc, err := s.repo.GetAccount(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, ErrInternal
	}
	return acc, nil
}
