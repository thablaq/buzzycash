package models

import (
	"time"
)

type NotificationType string

const (
	Trasaction       NotificationType = "TRANSACTION"
	Games          NotificationType = "GAMES"
)

type Notification struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string    `gorm:"type:uuid"`
	
	Title     string    `gorm:"size:255"`
	Message   string
	Type      NotificationType
    Subtitle  string    `json:"subtitle"`
    Amount    float64   `json:"amount"`
    Currency  string    `json:"currency"`
    Status    string    `json:"status"`
	IsRead    bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	
	User User `gorm:"constraint:OnDelete:CASCADE;"`
}