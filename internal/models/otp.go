package models

import (
	"time"
)

type UserOtpSecurity struct {
	ID                                   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID                               string    `gorm:"type:uuid;uniqueIndex"`
	
	VerificationCode                     string    `gorm:"size:255"`
	VerificationCodeCreatedAt            *time.Time
	VerificationCodeExpiresAt            *time.Time
	OtpRetryCount                        int       `gorm:"default:0"`
	IsOtpVerifiedForPasswordReset        bool      `gorm:"default:false"`
	OtpLockedUntil                       *time.Time
	
	PasswordResetVerificationCode        string    `gorm:"size:255"`
	PasswordResetVerificationCodeCreatedAt *time.Time
	PasswordResetVerificationCodeExpiresAt *time.Time
	ForgotPasswordOtpLockedUntil         *time.Time
	PasswordResetSentTo                  string    `gorm:"size:255"`
	ResetOtpRetryCount                   int       `gorm:"default:0"`

	User *User `gorm:"constraint:OnDelete:CASCADE;"`
}