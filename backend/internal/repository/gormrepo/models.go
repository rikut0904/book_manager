package gormrepo

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID           string `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex"`
	Username     string
	PasswordHash string
}

type ProfileSettings struct {
	UserID     string `gorm:"primaryKey"`
	Visibility string
}

type Book struct {
	ID            string  `gorm:"primaryKey"`
	ISBN13        *string `gorm:"uniqueIndex"`
	Title         string
	Authors       datatypes.JSON `gorm:"type:jsonb"`
	Publisher     string
	PublishedDate string
	ThumbnailURL  string
	Source        string
}

type UserBook struct {
	ID           string `gorm:"primaryKey"`
	UserID       string `gorm:"uniqueIndex:idx_user_book"`
	BookID       string `gorm:"uniqueIndex:idx_user_book"`
	Note         string
	AcquiredAt   string
	SeriesID     string
	VolumeNumber *int
}

type Favorite struct {
	ID       string `gorm:"primaryKey"`
	UserID   string `gorm:"uniqueIndex:idx_favorite_book;uniqueIndex:idx_favorite_series"`
	Type     string
	BookID   *string `gorm:"uniqueIndex:idx_favorite_book"`
	SeriesID *string `gorm:"uniqueIndex:idx_favorite_series"`
}

type NextToBuyManual struct {
	ID           string `gorm:"primaryKey"`
	UserID       string `gorm:"index"`
	Title        string
	SeriesName   string
	VolumeNumber *int
	Note         string
}

type Tag struct {
	ID          string `gorm:"primaryKey"`
	OwnerUserID string `gorm:"uniqueIndex:idx_tag_owner_name"`
	Name        string `gorm:"uniqueIndex:idx_tag_owner_name"`
}

type BookTag struct {
	UserID string `gorm:"primaryKey"`
	BookID string `gorm:"primaryKey"`
	TagID  string `gorm:"primaryKey"`
}

type Recommendation struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	BookID    string
	Comment   string
	CreatedAt time.Time
}
