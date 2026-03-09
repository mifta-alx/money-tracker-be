package repository

import (
	"context"
	"database/sql"
	"errors"
	"money-tracker/internal/models"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) GetUserWithPassword(ctx context.Context, email string) (*models.User, string, error) {
	query := `SELECT u.id, u.email, u.name, ap.password_hash 
		FROM users u JOIN auth_providers ap ON u.id = ap.user_id WHERE u.email = $1 AND ap.provider = 'local'`
	var u models.User
	var passwordHash string
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Name, &passwordHash)

	if err != nil {
		return nil, "", errors.New("invalid_credentials")
	}

	return &u, passwordHash, nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, name, created_at FROM users WHERE email = $1`
	var u models.User
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, u *models.User, passwordHash string) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	userQuery := `INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRowContext(ctx, userQuery, u.Email, u.Name).Scan(&u.ID)

	if err != nil {
		tx.Rollback()
		return err
	}

	authQuery := `INSERT INTO auth_providers (user_id, provider, password_hash) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, authQuery, u.ID, "local", passwordHash)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *AuthRepository) CreateUserOAuth(ctx context.Context, u *models.User, ap *models.AuthProvider) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	userQuery := `INSERT INTO users (email, name, avatar_url) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRowContext(ctx, userQuery, u.Email, u.Name, u.AvatarURL).Scan(&u.ID)

	if err != nil {
		tx.Rollback()
		return err
	}

	ap.UserID = u.ID
	authQuery := `INSERT INTO auth_providers (user_id, provider, provider_user_id) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, authQuery, ap.UserID, ap.Provider, ap.ProviderUserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
