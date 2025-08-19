package models

import (
	"time"
)

type NotificationType string

const (
	Withdrawal      NotificationType = "WITHDRAWAL"
	Cashout         NotificationType = "CASHOUT"
	Ticket  NotificationType = "TICKET_PURCHASE"
	Wallet          NotificationType = "WALLET"
	PasswordChange  NotificationType = "PASSWORD_CHANGE"
)

type Notification struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string    `gorm:"type:uuid"`
	
	Title     string    `gorm:"size:255"`
	Message   string
	Type      NotificationType
	IsRead    bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	
	User User `gorm:"constraint:OnDelete:CASCADE;"`
}