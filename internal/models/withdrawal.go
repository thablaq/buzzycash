package models

import (
	"time"

)

type WithdrawalRequest struct {
	ID               string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID           string    `gorm:"type:uuid"`
	
	Amount           float64
	Currency         string    `gorm:"size:3;default:'NGN'"`
	PaymentStatus    *EPaymentStatus
	Reason           *string `gorm:"size:255"`
	PaymentReference *string `gorm:"size:255"`
	RequestedAt      time.Time `gorm:"default:current_timestamp"`
	ProcessedAt      *time.Time
	CreatedAt        time.Time `gorm:"default:current_timestamp"`
	UpdatedAt        time.Time
	
	Transaction *TransactionHistory `gorm:"foreignKey:WithdrawalsID"`
	User        User                `gorm:"constraint:OnDelete:CASCADE;"`
}