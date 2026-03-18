package repository

import (
	"context"
	"database/sql"
	"money-tracker/internal/models"
	"time"

	"github.com/google/uuid"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) GetUserWithPassword(ctx context.Context, email string) (*models.User, string, error) {
	query := `SELECT u.id, u.email, u.name, u.created_at, u.updated_at, ap.password_hash 
		FROM users u JOIN auth_providers ap ON u.id = ap.user_id WHERE u.email = $1 AND ap.provider = 'email'`

	var u models.User
	var passwordHash string
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt, &passwordHash)

	if err != nil {
		return nil, "", err
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

func (r *AuthRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, email, name, avatar_url, created_at, updated_at FROM users WHERE id = $1`
	var u models.User
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
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

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	userQuery := `INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err = tx.QueryRowContext(ctx, userQuery, u.Email, u.Name).Scan(&u.ID, &u.CreatedAt,
		&u.UpdatedAt)

	if err != nil {
		return err
	}

	authQuery := `INSERT INTO auth_providers (user_id, provider, password_hash) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, authQuery, u.ID, "email", passwordHash)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *AuthRepository) CheckProviderExists(ctx context.Context, userId uuid.UUID, provider string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM auth_providers WHERE user_id = $1 AND provider = $2);`
	var exist bool
	err := r.DB.QueryRowContext(ctx, query, userId, provider).Scan(&exist)

	if err != nil {
		return false, err
	}

	return exist, nil
}

func (r *AuthRepository) AddAuthProvider(ctx context.Context, userId uuid.UUID, provider, providerUserId string) error {
	query := `INSERT INTO auth_providers (user_id, provider, provider_user_id) VALUES ($1, $2, $3);`
	_, err := r.DB.ExecContext(ctx, query, userId, provider, providerUserId)
	return err
}

func (r *AuthRepository) CreateUserOAuth(ctx context.Context, u *models.User, ap *models.AuthProvider) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	userQuery := `INSERT INTO users (email, name, avatar_url) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRowContext(ctx, userQuery, u.Email, u.Name, u.AvatarURL).Scan(&u.ID)
	if err != nil {
		return err
	}

	ap.UserID = u.ID
	authQuery := `INSERT INTO auth_providers (user_id, provider, provider_user_id) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, authQuery, ap.UserID, ap.Provider, ap.ProviderUserID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *AuthRepository) UpdateLastSign(ctx context.Context, userId uuid.UUID) error {
	query := `UPDATE users SET last_sign_in_at = $1, updated_at = $1 WHERE id = $2;`
	_, err := r.DB.ExecContext(ctx, query, time.Now(), userId)
	return err
}
