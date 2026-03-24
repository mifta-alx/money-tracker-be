package models

import (
	"time"

	"github.com/google/uuid"
)

type BudgetAllocation struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"-" db:"user_id"`
	Name         string    `json:"name" db:"name"`
	Percentage   int64     `json:"percentage" db:"percentage"`
	TargetAmount int64     `json:"target_amount" db:"target_amount"`
	Period       string    `json:"period" db:"period"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
