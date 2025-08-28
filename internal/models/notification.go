package models

import (
	"time"
)

type NotificationType string

const (
	Transaction       NotificationType = "TRANSACTION"
	Games          NotificationType = "GAMES"
)


type Notification struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"type:uuid"`
	
	Title     string    `gorm:"size:255"`
	Message   string
	Type      NotificationType
	Subtitle  string    `gorm:"size:500"`
	Amount    float64   
	Currency  string    `gorm:"size:10"`
	Status    string    `gorm:"size:50"`
	IsRead    bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	
	User User `gorm:"constraint:OnDelete:CASCADE;"`
}