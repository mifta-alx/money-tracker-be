package services

import (
	"context"
	"database/sql"
	"errors"
	"money-tracker/internal/models"
	"money-tracker/internal/repository"

	"github.com/google/uuid"
)

type BudgetAllocationService struct {
	repo *repository.BudgetAllocationRepository
}

func NewBudgetAllocationService(r *repository.BudgetAllocationRepository) *BudgetAllocationService {
	return &BudgetAllocationService{repo: r}
}

func (s *BudgetAllocationService) CreateBudgetAllocation(ctx context.Context, req *models.BudgetAllocation) (*models.BudgetAllocation, error) {
	if req.Name == "" {
		return nil, ErrMissingRequiredFields
	}
	if req.TargetAmount < 0 {
		return nil, ErrBalanceCannotBeNegative
	}

	err := s.repo.CreateBudgetAllocation(ctx, req)
	if err != nil {
		return nil, ErrInternal
	}

	return req, nil
}

func (s *BudgetAllocationService) UpdateBudgetAllocation(ctx context.Context, req *models.BudgetAllocation) (*models.BudgetAllocation, error) {
	err := s.repo.UpdateBudgetAllocation(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBudgetNotFound
		}
		return nil, ErrInternal
	}
	return req, nil
}

func (s *BudgetAllocationService) DeleteBudgetAllocation(ctx context.Context, budgetID, userID uuid.UUID) error {
	err := s.repo.DeleteBudgetAllocation(ctx, budgetID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrBudgetNotFound
		}
		return ErrInternal
	}
	return nil
}

func (s *BudgetAllocationService) GetBudgetAllocations(ctx context.Context, userID uuid.UUID) ([]*models.BudgetAllocation, error) {
	budgets, err := s.repo.GetBudgetAllocations(ctx, userID)
	if err != nil {
		return nil, ErrInternal
	}
	return budgets, nil
}

func (s *BudgetAllocationService) GetBudgetAllocation(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.BudgetAllocation, error) {
	budget, err := s.repo.GetBudgetAllocation(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBudgetNotFound
		}
		return nil, ErrInternal
	}
	return budget, nil
}
