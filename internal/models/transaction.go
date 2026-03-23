package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"-" db:"user_id"`
	AccountID  uuid.UUID  `json:"account_id" db:"account_id"`
	CategoryID *uuid.UUID `json:"category_id" db:"category_id"`
	Title      string     `json:"title" db:"title"`
	Type       string     `json:"type" db:"type"`
	Amount     int64      `json:"amount" db:"amount"`
	Date       time.Time  `json:"date" db:"date"`
	Notes      *string    `json:"notes" db:"notes"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}
