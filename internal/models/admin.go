package models

import (
	"time"
)

type Admin struct {
	ID             string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	
	Name           string    `gorm:"size:255"`
	Email          string    `gorm:"size:255;uniqueIndex"`
	Password       string    `gorm:"size:255"`
	CreatedAt      time.Time `gorm:"default:current_timestamp"`
	ProfilePicture *string   `gorm:"size:255"`
	RoleID         *string   `gorm:"type:uuid"`
	
	Role *Role `gorm:"constraint:OnDelete:SET NULL;"`
}