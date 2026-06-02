package models

import "time"

type Organisation struct {
	ID                  string     `json:"id"`
	Name                string     `json:"name"`
	OnboardingCompleted bool       `json:"onboarding_completed"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}
