package repository

import (
	"sync"
	"time"

	"book_manager/backend/internal/domain"
)

type MemoryAuditLogRepository struct {
	mu    sync.RWMutex
	items []domain.AuditLog
}

func NewMemoryAuditLogRepository() *MemoryAuditLogRepository {
	return &MemoryAuditLogRepository{
		items: make([]domain.AuditLog, 0),
	}
}

func (r *MemoryAuditLogRepository) Create(log domain.AuditLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items = append(r.items, log)
	return nil
}

func (r *MemoryAuditLogRepository) DeleteBefore(t time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filtered := r.items[:0]
	for _, item := range r.items {
		if item.CreatedAt.After(t) {
			filtered = append(filtered, item)
		}
	}
	r.items = filtered
	return nil
}
