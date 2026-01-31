package domain

import "time"

type AdminInvitation struct {
	ID        string
	Token     string
	Email     string
	UserID    string
	CreatedBy string
	ExpiresAt time.Time
	UsedAt    *time.Time
	UsedBy    string
	CreatedAt time.Time
}
