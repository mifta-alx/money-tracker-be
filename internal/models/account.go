package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"-" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Type      string    `json:"type" db:"type"`
	Balance   int64     `json:"balance" db:"balance"`
	Icon      string    `json:"icon" db:"icon"`
	Color     string    `json:"color" db:"color"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
