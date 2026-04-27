package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"-" db:"user_id"`
	AllocationID *uuid.UUID `json:"allocation_id" db:"allocation_id"`
	Name         string     `json:"name" db:"name"`
	Type         string     `json:"type" db:"type"`
	Color        string     `json:"color" db:"color"`
	Icon         string     `json:"icon" db:"icon"`
	TargetAmount int64      `json:"target_amount" db:"target_amount"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}
