package repository

import (
	"time"

	"book_manager/backend/internal/domain"
)

type AuditLogRepository interface {
	Create(log domain.AuditLog) error
	DeleteBefore(t time.Time) error
}
