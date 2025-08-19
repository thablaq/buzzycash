package models

import (
	"time"
)


type ReferralWallet struct {
	ID              string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID          string  `gorm:"unique"`
	ReferralBalance float32 `gorm:"type:decimal(10,2);default:0.0"`
	PointsUsed      float32 `gorm:"type:decimal(10,2);default:0.0"`
	PointsExpired   float32 `gorm:"type:decimal(10,2);default:0.0"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Referral struct {
	ID                   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ReferrerID           string
	ReferredUserID       string
	PointsEarned         float32   `gorm:"type:decimal(10,2);default:0.0"`
	PointsUsed           float32   `gorm:"type:decimal(10,2);default:0.0"`
	PointsExpired        float32   `gorm:"type:decimal(10,2);default:0.0"`
	SignupDate           time.Time `gorm:"default:now()"`
	FirstTransactionDate *time.Time
	TransactionCount     int `gorm:"default:0"`
	CreatedAt            time.Time
	ExpiresAt            time.Time

	ReferredUser User `gorm:"foreignKey:ReferredUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Referrer     User `gorm:"foreignKey:ReferrerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
