package domain

import "time"

type AuditLog struct {
	ID        string
	UserID    string
	Action    string
	Entity    string
	EntityID  string
	Payload   map[string]any
	IP        string
	UserAgent string
	CreatedAt time.Time
}
