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

func (r *AccountRepository) CreateAccount(ctx context.Context, u *models.Account) error {
	query := `INSERT INTO accounts (user_id, name, type, balance, icon, color) 
              VALUES ($1, $2, $3, $4, $5, $6) 
              RETURNING id, created_at, updated_at`
	err := r.DB.QueryRowContext(ctx, query,
		u.UserID,
		u.Name,
		u.Type,
		u.Balance,
		u.Icon,
		u.Color,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	return err
}

func (r *AccountRepository) GetAccounts(ctx context.Context, userID uuid.UUID) ([]*models.Account, error) {
	query := `SELECT id, name, type, balance, color, icon, created_at FROM accounts WHERE user_id = $1`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if err == nil {
			err = closeErr
		}
	}()

	var accounts []*models.Account
	for rows.Next() {
		var acc models.Account
		err := rows.Scan(&acc.ID, &acc.Name, &acc.Type, &acc.Balance, &acc.Icon, &acc.Color, &acc.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &acc)
	}

	if accounts == nil {
		accounts = []*models.Account{}
	}

	return accounts, nil
}

func (r *AccountRepository) GetAccount(ctx context.Context, id, userID uuid.UUID) (*models.Account, error) {
	query := `SELECT id, name, type, balance, icon, color, created_at FROM accounts WHERE id = $1 AND user_id = $2`
	var acc models.Account
	err := r.DB.QueryRowContext(ctx, query, id, userID).Scan(&acc.ID, &acc.Name, &acc.Type, &acc.Balance, &acc.Icon, &acc.Color, &acc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}
