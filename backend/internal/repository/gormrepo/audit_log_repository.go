package gormrepo

import (
	"encoding/json"
	"time"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(log domain.AuditLog) error {
	payload, err := json.Marshal(log.Payload)
	if err != nil {
		payload = []byte("{}")
	}
	model := AuditLog{
		ID:        log.ID,
		UserID:    log.UserID,
		Action:    log.Action,
		Entity:    log.Entity,
		EntityID:  log.EntityID,
		Payload:   datatypes.JSON(payload),
		IP:        log.IP,
		UserAgent: log.UserAgent,
		CreatedAt: log.CreatedAt,
	}
	return r.db.Create(&model).Error
}

func (r *AuditLogRepository) DeleteBefore(t time.Time) error {
	return r.db.Delete(&AuditLog{}, "created_at < ?", t).Error
}

var _ repository.AuditLogRepository = (*AuditLogRepository)(nil)
