package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID  `json:"id" db:"id"`
	Email               string     `json:"email" db:"email"`
	Name                string     `json:"name" db:"name"`
	AvatarURL           *string    `json:"avatar_url" db:"avatar_url"`
	OnboardingCompleted bool       `json:"onboarding_completed" db:"onboarding_completed"`
	CreatedAt           time.Time  `json:"-" db:"created_at"`
	UpdatedAt           time.Time  `json:"-" db:"updated_at"`
	LastSignInAt        *time.Time `json:"-" db:"last_sign_in_at"`
}
