package repository

import (
	"context"
	"database/sql"
	"money-tracker/internal/models"

	"github.com/google/uuid"
)

type AccountRepository struct {
	DB *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		DB: db,
	}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, a *models.Account) error {
	query := `INSERT INTO accounts (user_id, name, type, balance, icon, color) 
              VALUES ($1, $2, $3, $4, $5, $6) 
              RETURNING id, created_at, updated_at`
	return r.DB.QueryRowContext(ctx, query,
		a.UserID,
		a.Name,
		a.Type,
		a.Balance,
		a.Icon,
		a.Color,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *AccountRepository) UpdateAccount(ctx context.Context, a *models.Account) error {
	query := `UPDATE accounts SET name = $1, type = $2, balance = $3, icon = $4, color = $5, updated_at = NOW() WHERE id = $6 AND user_id = $7 RETURNING id, created_at, updated_at`
	return r.DB.QueryRowContext(ctx, query, a.Name, a.Type, a.Balance, a.Icon, a.Color, a.ID, a.UserID).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *AccountRepository) DeleteAccount(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM accounts WHERE id = $1 AND user_id = $2`
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

func (r *AccountRepository) GetAccounts(ctx context.Context, userID uuid.UUID) ([]*models.Account, float64, error) {
	query := `SELECT id, name, type, balance, icon, color, created_at, updated_at, SUM(balance) OVER() as total_balance FROM accounts WHERE user_id = $1 ORDER BY created_at`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, 0, err
	}

	defer func() { _ = rows.Close() }()

	var accounts []*models.Account
	var totalBalance float64

	for rows.Next() {
		var acc models.Account
		err := rows.Scan(&acc.ID, &acc.Name, &acc.Type, &acc.Balance, &acc.Icon, &acc.Color, &acc.CreatedAt, &acc.UpdatedAt, &totalBalance)
		if err != nil {
			return nil, 0, err
		}
		accounts = append(accounts, &acc)
	}

	if accounts == nil {
		accounts = []*models.Account{}
	}

	return accounts, totalBalance, nil
}

func (r *AccountRepository) GetAccount(ctx context.Context, id, userID uuid.UUID) (*models.Account, error) {
	query := `SELECT id, name, type, balance, icon, color, created_at, updated_at FROM accounts WHERE id = $1 AND user_id = $2`
	var acc models.Account
	err := r.DB.QueryRowContext(ctx, query, id, userID).Scan(&acc.ID, &acc.Name, &acc.Type, &acc.Balance, &acc.Icon, &acc.Color, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}
