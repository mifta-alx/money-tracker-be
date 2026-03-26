package models

import (
	"time"

	"github.com/google/uuid"
)

type Transfer struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"-" db:"user_id"`
	FromAccountID uuid.UUID `json:"from_account_id" db:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id" db:"to_account_id"`
	Amount        int64     `json:"amount" db:"amount"`
	Notes         *string   `json:"notes" db:"notes"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
