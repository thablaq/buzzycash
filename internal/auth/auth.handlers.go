package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/helpers"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/services"
	"github.com/dblaq/buzzycash/internal/utils"

	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	OTP_RESEND_COOLDOWN                 = 60 // seconds
	MAX_OTP_RETRIES                     = 5
	OTP_LOCKOUT_DURATION                = 15 * time.Minute
	VERIFY_OTP_LOCKED_DURATION          = 2 * time.Minute
	FORGOT_PASSWORD_OTP_LOCKED_DURATION = 3 * time.Minute
)


func SignUpHandler(ctx *gin.Context) {
	var req SignUpRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.Error(ctx, http.StatusBadRequest, "Missing or invalid required field(s)")
			return
		}

		if err := utils.Validate.Struct(req); err != nil {
			utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
			return
		}

	// Check if user already exists
	var existingUser models.User
	if err := config.DB.Where("phone_number = ?", req.PhoneNumber).First(&existingUser).Error; err == nil {
		if !existingUser.IsVerified {
			utils.Error(ctx, http.StatusConflict,
				"Account exists but not verified. Please login to complete verification")
			return
		}
		utils.Error(ctx, http.StatusBadRequest,
			"Account already exists")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Generate referral code
	referralCode := helpers.GenerateReferralCode()
	var referrer *models.User
	if req.ReferralCode != "" {
		if err := config.DB.Preload("ReferralWallet").
			Where("referral_code = ?", req.ReferralCode).
			First(&referrer).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Println("Error checking referral code:", err)
				utils.Error(ctx,http.StatusBadRequest," invalid referral code")
				return
			}
			referrer = nil
		}
	}

	// Start transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user
	user := models.User{
		PhoneNumber:        req.PhoneNumber,
		Password:           hashedPassword,
		CountryOfResidence: req.CountryOfResidence,
		IsActive:           true,
		IsVerified:         false,
		ReferralCode:       referralCode,
		ReferredByID:       nil,
	}
	if referrer != nil {
		user.ReferredByID = &referrer.ID
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		 log.Println("❌ Failed to create user:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Create referral wallet
	referralWallet := models.ReferralWallet{
		UserID:          user.ID,
		ReferralBalance: 0,
		PointsUsed:      0,
		PointsExpired:   0,
	}
	if err := tx.Create(&referralWallet).Error; err != nil {
		tx.Rollback()
		utils.Error(ctx, http.StatusInternalServerError, "Failed to create referral wallet")
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.Error(ctx, http.StatusInternalServerError, "Transaction failed")
		return
	}

	// Send OTP (still outside transaction)
	countryPrefix := req.PhoneNumber[:3]
	emailService := services.EmailService{}
	
	

	switch countryPrefix {
	case "233":
		if _, err := emailService.SendGhanaOtp(user.PhoneNumber, user.ID); err != nil {
			utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Ghana OTP: %v", err))
			return
		}
	case "234":
		if _, err := emailService.SendNaijaOtp(user.PhoneNumber, user.ID); err != nil {
			utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Nigeria OTP: %v", err))
			return
		}
	default:
		utils.Error(ctx, http.StatusBadRequest, "Unsupported country code")
		return
	}

	// Process referral if applicable
	if referrer != nil {
		referralPoints := 100.0
		err = config.DB.Transaction(func(tx *gorm.DB) error {
			expiresAt := time.Now().AddDate(1, 0, 0) // 1 year
			referral := models.Referral{
				ReferrerID:     referrer.ID,
				ReferredUserID: user.ID,
				PointsEarned:   float32(referralPoints),
				ExpiresAt:      expiresAt,
			}
			if err := tx.Create(&referral).Error; err != nil {
				return err
			}
			return tx.Model(&models.ReferralWallet{}).
				Where("user_id = ?", referrer.ID).
				Update("referral_balance", gorm.Expr("referral_balance + ?", referralPoints)).Error
		})
		if err != nil {
			log.Println("Failed to process referral:", err)
		}
	}

	// Success response
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"user": gin.H{
				"id":                 user.ID,
				"phoneNumber":        user.PhoneNumber,
				"countryOfResidence": user.CountryOfResidence,
				"isActive":           user.IsActive,
				"isVerified":         user.IsVerified,
				"referralCode":       user.ReferralCode,
				"createdAt":          user.CreatedAt,
			},
			"wallets": gin.H{
				"referral": gin.H{
					"id":              referralWallet.ID,
					"referralBalance": referralWallet.ReferralBalance,
					"pointsUsed":      referralWallet.PointsUsed,
					"pointsExpired":   referralWallet.PointsExpired,
					"createdAt":       referralWallet.CreatedAt,
				},
			},
		},
		"message": "User registered successfully. Please verify your account with the OTP sent to your number.",
	})
}

func VerifyAccountHandler(ctx *gin.Context) {
	var req VerifyAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	// 🔍 Fetch user with OTP security
	var user models.User
	if err := config.DB.Preload("OtpSecurity").
		Where("phone_number = ?", req.PhoneNumber).
		First(&user).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	// Already verified?
	if user.IsVerified {
		utils.Error(ctx, http.StatusConflict, "User already verified")
		return
	}

	otp := user.OtpSecurity
	if otp.VerificationCode == "" {
		utils.Error(ctx, http.StatusConflict, "Invalid verification code")
		return
	}

	if otp.VerificationCode != req.VerificationCode {
		utils.Error(ctx, http.StatusConflict, "Invalid verification code")
		return
	}

	if otp.VerificationCodeCreatedAt.IsZero() || otp.VerificationCodeExpiresAt == nil || otp.VerificationCodeExpiresAt.IsZero() {
		utils.Error(ctx, http.StatusBadRequest, "OTP metadata is incomplete or missing")
		return
	}

	if time.Now().After(*otp.VerificationCodeExpiresAt) {
		utils.Error(ctx, http.StatusBadRequest, "OTP has expired")
		return
	}

	// ✅ Transaction to update user + clear OTP
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).Update("is_verified", true).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.UserOtpSecurity{}).
			Where("user_id = ?", user.ID).
			Updates(map[string]interface{}{
				"verification_code":            nil,
				"verification_code_expires_at": nil,
				"verification_code_created_at": nil,
				"verify_otp_locked_until":      nil,
				"otp_retry_count":              0,
			}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to verify account")
		return
	}

	// 🎟️ Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	// Save refresh token
	expireAt := time.Now().AddDate(0, 0, config.AppConfig.RefreshTokenExpiresDays)
	rt := models.RefreshToken{
		UserID:   user.ID,
		Token:    refreshToken,
		ExpireAt: &expireAt,
	}
	if err := config.DB.Create(&rt).Error; err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User verified account successfully.",
		"user": gin.H{
			"id":                 user.ID,
			"fullName":           user.FullName,
			"phoneNumber":        user.PhoneNumber,
			"email":              user.Email,
			"isActive":           user.IsActive,
			"isVerified":         true,
			"isProfileCreated":   user.IsProfileCreated,
			"countryOfResidence": user.CountryOfResidence,
			"gender":             user.Gender,
			"dateOfBirth":        user.DateOfBirth,
			"profilePicture":     user.ProfilePicture,
			"accessToken":        accessToken,
			"refreshToken":       refreshToken,
		},
	})
}


// ResendOtpHandler handles resending OTP to users
func ResendOtpHandler(ctx *gin.Context) {
	var req ResendOtpRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	
	if err := utils.Validate.Struct(req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
}
	


	
	var user models.User
	if err := config.DB.Preload("OtpSecurity").
		Where("phone_number = ?", req.PhoneNumber).
		First(&user).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	if user.IsVerified {
		utils.Error(ctx, http.StatusConflict, "User already verified")
		return
	}

	otp := user.OtpSecurity
	currentTime := time.Now()

	// Check lockout
	if otp != nil && otp.OtpLockedUntil != nil && currentTime.Before(*otp.OtpLockedUntil) {
		remainingTime := int((*otp.OtpLockedUntil).Sub(currentTime).Minutes())
		utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Please wait %d minute(s) before requesting a new OTP.", remainingTime))
		return
	}

	// Check cooldown
	if otp != nil && otp.VerificationCodeCreatedAt != nil {
		timeSinceLastOtp := currentTime.Sub(*otp.VerificationCodeCreatedAt)
		if timeSinceLastOtp < time.Duration(OTP_RESEND_COOLDOWN)*time.Second {
			remainingCooldown := OTP_RESEND_COOLDOWN - int(timeSinceLastOtp.Seconds())
			utils.Error(ctx, http.StatusTooManyRequests, fmt.Sprintf("Please wait %d seconds before requesting a new OTP.", remainingCooldown))
			return
		}
	}

	// Check retry limit
	if otp != nil && otp.OtpRetryCount >= MAX_OTP_RETRIES {
		config.DB.Model(&models.UserOtpSecurity{}).
			Where("user_id = ?", user.ID).
			Updates(map[string]interface{}{
				"verify_otp_locked_until": time.Now().Add(OTP_LOCKOUT_DURATION),
				"otp_retry_count":         0,
			})
		utils.Error(ctx, http.StatusTooManyRequests, "Too many OTP attempts. Please try again later.")
		return
	}

	// Determine country by phone prefix
	countryPrefix := req.PhoneNumber[:3]
	emailService := services.EmailService{}

	switch countryPrefix {
	case "233":
		if _, err := emailService.SendGhanaOtp(user.PhoneNumber, user.ID); err != nil {
			utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Ghana OTP: %v", err))
			return
		}
	case "234":
		if _, err := emailService.SendNaijaOtp(user.PhoneNumber, user.ID); err != nil {
			utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Nigeria OTP: %v", err))
			return
		}
	default:
		utils.Error(ctx, http.StatusBadRequest, "Unsupported country code")
		return
	}

	// Update retry count + lockout timestamp
	config.DB.Model(&models.UserOtpSecurity{}).
		Where("user_id = ?", user.ID).
		Updates(map[string]interface{}{
			"otp_retry_count":         otp.OtpRetryCount + 1,
			"verify_otp_locked_until": time.Now().Add(VERIFY_OTP_LOCKED_DURATION),
		})

	ctx.JSON(http.StatusOK, gin.H{
		"message": "New OTP sent successfully. Please check your SMS.",
		"user": gin.H{
			"id":                 user.ID,
			"fullName":           user.FullName,
			"phoneNumber":        user.PhoneNumber,
			"email":              user.Email,
			"isActive":           user.IsActive,
			"isVerified":         user.IsVerified,
			"isProfileCreated":   user.IsProfileCreated,
			"countryOfResidence": user.CountryOfResidence,
			"gender":             user.Gender,
			"dateOfBirth":        user.DateOfBirth,
		},
	})
}

func LoginHandler(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
}

	var user models.User
	if err := config.DB.Where("email = ? OR phone_number = ?", req.Email, req.PhoneNumber).
		Preload("OtpSecurity").
		First(&user).Error;

	err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}
	


	if !user.IsActive {
		utils.Error(ctx, http.StatusBadRequest, "Account is blocked, please contact support")
		return
	}

	if !user.IsVerified {
		countryPrefix := req.PhoneNumber[:3]
		emailService := services.EmailService{}

		switch countryPrefix {
		case "233":
			if _, err := emailService.SendGhanaOtp(user.PhoneNumber, user.ID); err != nil {
				utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Ghana OTP: %v", err))
				return
			}
		case "234":
			if _, err := emailService.SendNaijaOtp(user.PhoneNumber, user.ID); err != nil {
				utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Nigeria OTP: %v", err))
				return
			}
		default:
			utils.Error(ctx, http.StatusBadRequest, "Unsupported country code")
			return
		}
	}

	if !utils.ComparePassword(user.Password,req.Password) {
		utils.Error(ctx, http.StatusForbidden, "Invalid credentials")
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	config.DB.Model(&user).Update("last_login", time.Now())

	// Upsert refresh token
	expiresAt := time.Now().AddDate(0, 0, config.AppConfig.RefreshTokenExpiresDays)
	rt := models.RefreshToken{
		UserID:   user.ID,
		Token:    refreshToken,
		ExpireAt: &expiresAt,
	}
	config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"token", "expire_at"}),
	}).Create(&rt)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
		"user": gin.H{
			"id":                 user.ID,
			"fullName":           user.FullName,
			"username":           user.Username,
			"phoneNumber":        user.PhoneNumber,
			"email":              user.Email,
			"isActive":           user.IsActive,
			"gender":             user.Gender,
			"isProfileCreated":   user.IsProfileCreated,
			"countryOfResidence": user.CountryOfResidence,
			"dateOfBirth":        user.DateOfBirth,
			"isVerified":         user.IsVerified,
			"lastLogin":          user.LastLogin,
			"profilePicture":     user.ProfilePicture,
			"accessToken":        accessToken,
			"refreshToken":       refreshToken,
		},
	})
}

// ChangePasswordHandler changes user password
func ChangePasswordHandler(ctx *gin.Context) {
	var req PasswordChangeRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
}

	currentUser := ctx.MustGet("currentUser").(models.User)

	if !currentUser.IsVerified {
		utils.Error(ctx, http.StatusBadRequest, "Only verified account can change password")
	}
	if !utils.ComparePassword(currentUser.Password, req.CurrentPassword) {
		utils.Error(ctx, http.StatusBadRequest, "Current password is incorrect")
		return
	}

	if req.CurrentPassword == req.NewPassword {
		utils.Error(ctx, http.StatusBadRequest, "New password can not be the same as current password")
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to process password")
		return
	}

	if err := config.DB.Model(&currentUser).Update("password", hashedPassword).Error; err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to update password")
		return
	}

	accessToken, err := utils.GenerateAccessToken(currentUser.ID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(currentUser.ID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	// Delete old refresh tokens
	if err := config.DB.Where("user_id = ?", currentUser.ID).Delete(&models.RefreshToken{}).Error; err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to delete old refresh tokens")
		return
	}

	// Create new refresh token
	expireAt := time.Now().Add(time.Hour * 24 * time.Duration(config.AppConfig.RefreshTokenExpiresDays))
	newRefresh := models.RefreshToken{
		UserID:   currentUser.ID,
		Token:    refreshToken,
		ExpireAt: &expireAt,
	}
	if err := config.DB.Create(&newRefresh).Error; err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	// Add notification
	notification := models.Notification{
		UserID:  currentUser.ID,
		Title:   "Password Change Successful",
		Message: fmt.Sprintf("Your password change for %s has been successful.", currentUser.Email),
		Type:    models.PasswordChange,
		IsRead:  false,
	}
	if err := config.DB.Create(&notification).Error; err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to save notification")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
		"user": gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	})
}

// ForgotPasswordUser handles initiating password reset
func ForgotPasswordHandler(ctx *gin.Context) {
	var req ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	
	if err := utils.Validate.Struct(req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
}

	var user models.User
	if err := config.DB.Preload("OtpSecurity").
		Where("email = ? OR phone_number = ?", req.Email, req.PhoneNumber).
		First(&user).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	now := time.Now()

	// Lockout check
	if user.OtpSecurity != nil && user.OtpSecurity.ForgotPasswordOtpLockedUntil != nil &&
		now.Before(*user.OtpSecurity.ForgotPasswordOtpLockedUntil) {
		remaining := int(user.OtpSecurity.ForgotPasswordOtpLockedUntil.Sub(now).Minutes())
		utils.Error(ctx, http.StatusBadRequest,
			fmt.Sprintf("Please wait %d minute(s) before requesting a new OTP.", remaining))
		return
	}

	// Active OTP check
	if user.OtpSecurity != nil && user.OtpSecurity.PasswordResetVerificationCodeExpiresAt != nil &&
		now.Before(*user.OtpSecurity.PasswordResetVerificationCodeExpiresAt) {
		remaining := int(user.OtpSecurity.PasswordResetVerificationCodeExpiresAt.Sub(now).Minutes())
		utils.Error(ctx, http.StatusBadRequest,
			fmt.Sprintf("An active OTP already exists. Please wait %d minute(s) before requesting a new OTP.", remaining))
		return
	}

	// Retry count
	if user.OtpSecurity != nil && user.OtpSecurity.ResetOtpRetryCount >= MAX_OTP_RETRIES {
		lockUntil := now.Add(OTP_LOCKOUT_DURATION)
		config.DB.Model(&models.UserOtpSecurity{}).
			Where("user_id = ?", user.ID).
			Updates(map[string]interface{}{
				"forgot_password_otp_locked_until": lockUntil,
				"reset_otp_retry_count":            0,
			})
		utils.Error(ctx, http.StatusTooManyRequests,
			"You have exceeded the maximum OTP attempts. Please wait before trying again.")
		return
	}

	// Send OTP (phone or email)
	emailService := services.EmailService{}

	if req.PhoneNumber != "" {
		cleanPhone := strings.ReplaceAll(req.PhoneNumber, " ", "")

		var err error
		switch {
		case strings.HasPrefix(cleanPhone, "234"):
			if _, err = emailService.SendForgotPasswordNGNOtp(cleanPhone, user.ID); err != nil {
				utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Nigeria OTP: %v", err))
				return
			}

		case strings.HasPrefix(cleanPhone, "233"):
			if _, err = emailService.SendForgotPasswordGHCOtp(cleanPhone, user.ID); err != nil {
				utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Ghana OTP: %v", err))
				return
			}

		default:
			utils.Error(ctx, http.StatusBadRequest, "Unsupported phone country code")
			return
		}

		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "Failed to send OTP")
			return
		}

	} else if req.Email != "" {
		if _, err := emailService.SendForgotPasswordEmailOtp(req.Email, user.FullName, user.ID); err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "Failed to send OTP email")
			return
		}
	} else {
		utils.Error(ctx, http.StatusBadRequest, "Either phone number or email is required")
		return
	}

	// Update retry count + lock until
	lockUntil := now.Add(FORGOT_PASSWORD_OTP_LOCKED_DURATION)
	config.DB.Model(&models.UserOtpSecurity{}).
		Where("user_id = ?", user.ID).
		Updates(map[string]interface{}{
			"reset_otp_retry_count":            gorm.Expr("reset_otp_retry_count + ?", 1),
			"forgot_password_otp_locked_until": lockUntil,
		})

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully. Please check your email or phone.",
		"user": gin.H{
			"id":          user.ID,
			"fullName":    user.FullName,
			"phoneNumber": user.PhoneNumber,
			"email":       user.Email,
		},
	})
}

// VerifyPasswordForgotOtp handles OTP verification
func VerifyPasswordForgotOtpHandler(ctx *gin.Context) {
	var req VerifyPasswordForgotOtpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	
	if err := utils.Validate.Struct(req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
}

	var user models.User
	if err := config.DB.Preload("OtpSecurity").
		Where("email = ? OR phone_number = ?", req.Email, req.PhoneNumber).
		First(&user).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	if !user.IsVerified {
		utils.Error(ctx, http.StatusForbidden, "Please verify your account before resetting password")
		return
	}

	otpSec := user.OtpSecurity
	if otpSec == nil {
		utils.Error(ctx, http.StatusBadRequest, "OTP metadata missing")
		return
	}

	// Check OTP destination
	if otpSec.PasswordResetSentTo == "email" && req.Email == "" {
		utils.Error(ctx, http.StatusBadRequest, "OTP was sent to email, please provide email")
		return
	}
	if otpSec.PasswordResetSentTo == "phone" && req.PhoneNumber == "" {
		utils.Error(ctx, http.StatusBadRequest, "OTP was sent to phone, please provide phone number")
		return
	}

	// Verify code
	if otpSec.PasswordResetVerificationCode != req.VerificationCode {
		utils.Error(ctx, http.StatusConflict, "Invalid verification code")
		return
	}

	if time.Now().After(*otpSec.PasswordResetVerificationCodeExpiresAt) {
		utils.Error(ctx, http.StatusBadRequest, "OTP has expired")
		return
	}

	// Mark verified
	config.DB.Model(&models.UserOtpSecurity{}).
		Where("user_id = ?", user.ID).
		Update("is_otp_verified_for_password_reset", true)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP verification successful",
		"userId":  user.ID,
		"email":   user.Email,
	})
}

// ResetPasswordUser handles actual password reset
func ResetPasswordHandler(ctx *gin.Context) {
	var req ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	
	if err := utils.Validate.Struct(req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
}

	var user models.User
	if err := config.DB.Preload("OtpSecurity").
		Where("id = ?", req.UserId).
		First(&user).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	if user.OtpSecurity == nil || !user.OtpSecurity.IsOtpVerifiedForPasswordReset {
		utils.Error(ctx, http.StatusForbidden, "OTP verification required")
		return
	}

	// Check if same password
	// if utils.ComparePassword(user.Password,req.NewPassword,) {
	// 	utils.Error(ctx, http.StatusBadRequest, "New password must be different from old password")
	// 	return
	// }

	// Hash new password
	hashedPassword, _ := utils.HashPassword(req.NewPassword)
	config.DB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("password", hashedPassword)

	// Clear OTP fields
	config.DB.Model(&models.UserOtpSecurity{}).
		Where("user_id = ?", user.ID).
		Updates(map[string]interface{}{
			"password_reset_verification_code":            nil,
			"is_otp_verified_for_password_reset":          false,
			"password_reset_sent_to":                      nil,
			"forgot_password_otp_locked_until":            nil,
			"password_reset_verification_code_created_at": nil,
			"password_reset_verification_code_expires_at": nil,
			"reset_otp_retry_count":                       0,
		})

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password reset successful",
		"user": gin.H{
			"id":          user.ID,
			"fullName":    user.FullName,
			"phoneNumber": user.PhoneNumber,
			"email":       user.Email,
		},
	})
}




func LogoutHandler(ctx *gin.Context) {
    authHeader := ctx.GetHeader("Authorization")
    if authHeader == "" {
        utils.Error(ctx, http.StatusUnauthorized, "Authorization header missing")
        return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    if tokenString == authHeader {
        utils.Error(ctx, http.StatusUnauthorized, "Invalid token format")
        return
    }

    // Decode token to get expiration
    claims, err := utils.DecodeToken(tokenString)
    if err != nil {
        utils.Error(ctx, http.StatusUnauthorized, "Invalid token")
        return
    }

    // Determine expiration time
    var expireAt time.Time
    if exp, ok := claims["exp"].(float64); ok {
        expireAt = time.Unix(int64(exp), 0)
    } else {
        // fallback if token has no exp: blacklist for 1 hour
        expireAt = time.Now().Add(1 * time.Hour)
    }

    // Blacklist the token
    if err := utils.BlacklistToken(tokenString, expireAt); err != nil {
        utils.Error(ctx, http.StatusInternalServerError, "Failed to blacklist token")
        return
    }

    // Return success
    ctx.JSON(http.StatusOK, gin.H{
        "message": "User logged out successfully",
    })
}






func RefreshTokenHandler(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Session expired")
		return
	}
	
	if err := utils.Validate.Struct(req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
}

	var tokenEntry models.RefreshToken
	err := config.DB.
		Where("token = ?", req.RefreshToken).
		Order("created_at DESC").
		First(&tokenEntry).Error
	if err != nil {
		utils.Error(ctx, http.StatusUnauthorized, "Session expired")
		return
	}

	isExpired := tokenEntry.ExpireAt != nil && time.Now().After(*tokenEntry.ExpireAt)
	isValid := false
	var userId string

	userId, err = utils.VerifyJWTRefreshToken(tokenEntry.Token)
if err == nil && !isExpired {
		isValid = true
}


	// Always delete the token (one-time use)
	config.DB.Delete(&tokenEntry)

	if !isValid || userId == "" {
		utils.Error(ctx, http.StatusUnauthorized, "Session expired")
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", userId).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
	})
}