package models

import (
	"time"
)

type ReferralWallet struct {
	ID              string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID          string `gorm:"type:uuid;not null;uniqueIndex"`
	ReferralBalance int64  `gorm:"default:0"`
	PointsUsed      int64  `gorm:"default:0"`
	PointsExpired   int64  `gorm:"default:0"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Earnings        []ReferralEarning `gorm:"foreignKey:WalletID"`
}

type ReferralEarning struct {
	ID         string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WalletID   string `gorm:"type:uuid;not null;index"`
	ReferrerID string `gorm:"type:uuid;not null"`
	ReferredID string `gorm:"type:uuid;not null"`
	Points     int64  `gorm:"not null"`
	CreatedAt  time.Time
	ExpiresAt  time.Time
	Used       bool `gorm:"default:false"`
}
