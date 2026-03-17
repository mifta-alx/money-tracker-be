package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthProvider struct {
	ID             uuid.UUID `json:"id" db:"id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	Provider       string    `json:"provider" db:"provider"`
	ProviderUserID string    `json:"provider_user_id" db:"provider_user_id"`
	PasswordHash   *string   `json:"-" db:"password_hash"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
