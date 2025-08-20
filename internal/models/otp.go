package models

import (
	"time"
)

type OtpAction string

const (
	OtpActionVerifyAccount OtpAction = "verify_account"
	OtpActionPasswordReset OtpAction = "password_reset"
	OtpActionVerifyEmail   OtpAction = "verify_email"
)

type UserOtpSecurity struct {
	ID     string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID string `gorm:"type:uuid;uniqueIndex"`

	Code                          string `gorm:"size:255;not null"`
	CreatedAt                     *time.Time
	ExpiresAt                     *time.Time
	RetryCount                    int `gorm:"default:0"`
	LockedUntil                   *time.Time
	IsOtpVerifiedForPasswordReset bool      `gorm:"default:false"`
	SentTo                        string    `gorm:"size:255"`
	Action                        OtpAction `gorm:"size:50;not null"`

	User *User `gorm:"constraint:OnDelete:CASCADE;"`
}
