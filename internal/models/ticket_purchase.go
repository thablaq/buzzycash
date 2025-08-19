package models

import (
	"time"
)

type TicketPurchase struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      string    `gorm:"type:uuid"`
	
	TotalAmount float64
	UnitPrice   float64
	Quantity    int
	PurchasedAt time.Time `gorm:"default:current_timestamp"`
	Currency    string    `gorm:"size:3;default:'NGN'"`
	
	User         User                `gorm:"constraint:OnDelete:CASCADE;"`
	Transaction  *TransactionHistory `gorm:"foreignKey:TicketPurchaseID"`
}