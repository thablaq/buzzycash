package models

import (
	"time"
)

type TicketPurchase struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      string    `gorm:"type:uuid"`
	TransactionHistoryID string   `gorm:"type:uuid"`
	TotalAmount int64
	UnitPrice   int64
	Quantity    int
	PurchasedAt time.Time `gorm:"default:current_timestamp"`
	Currency    string    `gorm:"size:3;default:'NGN'"`
	
	User         User                `gorm:"constraint:OnDelete:CASCADE;"`
	TransactionHistory TransactionHistory `gorm:"foreignKey:TransactionHistoryID;references:ID;constraint:OnDelete:SET NULL"`
}