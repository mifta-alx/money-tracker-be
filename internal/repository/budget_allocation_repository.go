package repository

import (
	"context"
	"database/sql"
	"money-tracker/internal/models"

	"github.com/google/uuid"
)

type BudgetAllocationRepository struct {
	DB *sql.DB
}

func NewBudgetAllocationRepository(db *sql.DB) *BudgetAllocationRepository {
	return &BudgetAllocationRepository{
		DB: db,
	}
}

func (r *BudgetAllocationRepository) CreateBudgetAllocation(ctx context.Context, b *models.BudgetAllocation) error {
	query := `INSERT INTO budget_allocations (user_id, name, percentage, target_amount, period) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	return r.DB.QueryRowContext(ctx, query,
		b.UserID,
		b.Name,
		b.Percentage,
		b.TargetAmount,
		b.Period,
	).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *BudgetAllocationRepository) UpdateBudgetAllocation(ctx context.Context, b *models.BudgetAllocation) error {
	query := `UPDATE budget_allocations SET name = $1, percentage = $2, target_amount = $3, period = $4, updated_at = NOW() WHERE id = $5 AND user_id = $6 RETURNING id, created_at, updated_at`
	return r.DB.QueryRowContext(ctx, query,
		b.Name,
		b.Percentage,
		b.TargetAmount,
		b.Period,
		b.ID,
		b.UserID,
	).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *BudgetAllocationRepository) DeleteBudgetAllocation(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM budget_allocations WHERE id = $1 AND user_id = $2`
	result, err := r.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (r *BudgetAllocationRepository) GetBudgetAllocations(ctx context.Context, userID uuid.UUID) ([]*models.BudgetAllocation, error) {
	query := `SELECT id, name, percentage, target_amount, period, created_at, updated_at FROM budget_allocations WHERE user_id = $1`
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()
	var budgets []*models.BudgetAllocation
	for rows.Next() {
		var budget models.BudgetAllocation
		err := rows.Scan(&budget.ID, &budget.Name, &budget.Percentage, &budget.TargetAmount, &budget.Period, &budget.CreatedAt, &budget.UpdatedAt)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, &budget)
	}

	if budgets == nil {
		budgets = []*models.BudgetAllocation{}
	}

	return budgets, nil
}

func (r *BudgetAllocationRepository) GetBudgetAllocation(ctx context.Context, id, userID uuid.UUID) (*models.BudgetAllocation, error) {
	query := `SELECT id, name, percentage, target_amount, period, created_at, updated_at FROM budget_allocations WHERE id = $1 AND user_id = $2`
	var budget models.BudgetAllocation
	err := r.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&budget.ID,
		&budget.Name,
		&budget.Percentage,
		&budget.TargetAmount,
		&budget.Period,
		&budget.CreatedAt,
		&budget.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &budget, nil
}
