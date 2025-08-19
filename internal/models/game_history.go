package models

import (
	"time"
)

type GameStatus string

const (
	Ongoing GameStatus = "ONGOING"
	Won     GameStatus = "WON"
	Lost    GameStatus = "LOST"
)

type GameHistory struct {
	ID                   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketTypeID         string    `gorm:"type:uuid"`
	UserID               string    `gorm:"type:uuid"`
	TransactionHistoryID *string   `gorm:"type:uuid"`
	
	Prize        float64
	Status       GameStatus
	WinningBalls *int
	PlayedAt     time.Time
	CreatedAt    time.Time `gorm:"default:current_timestamp"`
	
	User               User                `gorm:"constraint:OnDelete:CASCADE;"`
	TransactionHistory *TransactionHistory `gorm:"constraint:OnDelete:SET NULL;"`
}