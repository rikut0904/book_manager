package domain

import "time"

type AdminUser struct {
	ID        string
	UserID    string
	CreatedBy string
	CreatedAt time.Time
}
