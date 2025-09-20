package models

import (
	"time"
)

type GenderType string

const (
	Male   GenderType = "MALE"
	Female GenderType = "FEMALE"
	Others GenderType = "OTHERS"
)

type User struct {
	ID                 string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FullName           string     `gorm:"size:255"`
	PhoneNumber        string     `gorm:"size:255;uniqueIndex"`
	Email              string     `gorm:"size:255;uniqueIndex"`
	Username           string     `gorm:"size:255;uniqueIndex"`
	DateOfBirth        string     `gorm:"size:255"`
	Password           string     `gorm:"size:255"`
	ProfilePicture     string     `gorm:"size:255"`
	IsProfileCreated   bool       `gorm:"default:false"`
	ReferralCode       string     `gorm:"size:255;uniqueIndex"`
	IsActive           bool       `gorm:"default:true"`
	IsEmailVerified    bool       `gorm:"default:false"`
	IsKycVerified      bool        `gorm:"default:false"`
	IsVerified         bool       `gorm:"default:false"`
	CreatedAt          time.Time  `gorm:"default:current_timestamp"`
	LastLogin          time.Time  `gorm:"default:current_timestamp"`
	ReferredByID       *string    `gorm:"type:uuid"` 
	Gender             GenderType
	CountryOfResidence string     `gorm:"size:255"`

	// Relations
	Notifications      []Notification        `gorm:"foreignKey:UserID"`
	ReferralWallet     ReferralWallet        `gorm:"foreignKey:UserID"`    
	ReferralEarnings   []ReferralEarning     `gorm:"foreignKey:ReferrerID"` 
	Transaction        []Transaction  `gorm:"foreignKey:UserID"`
	RefreshTokens      []RefreshToken        `gorm:"foreignKey:UserID"`
	TicketPurchases    []TicketPurchase      `gorm:"foreignKey:UserID"`
	GameHistories      []GameHistory         `gorm:"foreignKey:UserID"`
	OtpSecurity        *UserOtpSecurity      `gorm:"foreignKey:UserID"`
}
