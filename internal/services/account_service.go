package services

import (
	"context"
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
	if req.Name == "" || req.UserID == uuid.Nil {
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

func (s *AccountService) GetAccounts(ctx context.Context, userID uuid.UUID) ([]*models.Account, error) {
	accounts, err := s.repo.GetAccounts(ctx, userID)
	if err != nil {
		return nil, ErrInternal
	}
	return accounts, nil
}
