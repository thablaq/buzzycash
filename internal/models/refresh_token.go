package models

import (
	"time"
)

type RefreshToken struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string    `gorm:"type:uuid;uniqueIndex"`
	
	Token     string    `gorm:"size:255"`
	ExpireAt  *time.Time
	CreatedAt *time.Time `gorm:"default:current_timestamp"`
	UpdatedAt *time.Time
	
	User User `gorm:"constraint:OnDelete:CASCADE;"`
}

type BlacklistedToken struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Token     string    `gorm:"not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}