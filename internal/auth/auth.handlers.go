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
	log.Println("SignUpHandler invoked")
	var req SignUpRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to bind JSON:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		log.Println("Validation error:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := config.DB.Where("phone_number = ?", req.PhoneNumber).First(&existingUser).Error; err == nil {
		if !existingUser.IsVerified {
			log.Println("Account exists but not verified for phone number:", req.PhoneNumber)
			utils.Error(ctx, http.StatusConflict,
				"Account exists but not verified. Please login to complete verification")
			return
		}
		log.Println("Account already exists for phone number:", req.PhoneNumber)
		utils.Error(ctx, http.StatusBadRequest,
			"Account already exists")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Println("Failed to hash password:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Generate referral code
	referralCode := helpers.GenerateReferralCode()
	log.Println("Generated referral code:", referralCode)
	var referrer *models.User
	if req.ReferralCode != "" {
		log.Println("Checking referral code:", req.ReferralCode)
		if err := config.DB.Preload("ReferralWallet").
			Where("referral_code = ?", req.ReferralCode).
			First(&referrer).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Println("Error checking referral code:", err)
				utils.Error(ctx, http.StatusBadRequest, " invalid referral code")
				return
			}
			log.Println("Referral code not found:", req.ReferralCode)
			referrer = nil
		} else {
			log.Println("Referral code valid, referrer ID:", referrer.ID)
		}
	}

	// Start transaction
	log.Println("Starting transaction for user creation")
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Println("Transaction panicked, rolling back:", r)
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
	log.Println("User created successfully with ID:", user.ID)

	// Create referral wallet
	referralWallet := models.ReferralWallet{
		UserID:          user.ID,
		ReferralBalance: 0,
		PointsUsed:      0,
		PointsExpired:   0,
	}
	if err := tx.Create(&referralWallet).Error; err != nil {
		tx.Rollback()
		log.Println("Failed to create referral wallet:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to create referral wallet")
		return
	}
	log.Println("Referral wallet created successfully for user ID:", user.ID)

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Println("Transaction commit failed:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Transaction failed")
		return
	}
	log.Println("Transaction committed successfully")

	// Send OTP (still outside transaction)
	countryPrefix := req.PhoneNumber[:3]
	emailService := services.EmailService{}
	log.Println("Sending OTP to phone number with prefix:", countryPrefix)

	switch countryPrefix {
	case "233":
		if _, err := emailService.SendGhanaOtp(user.PhoneNumber, user.ID); err != nil {
			log.Println("Failed to send Ghana OTP:", err)
			utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Ghana OTP: %v", err))
			return
		}
		log.Println("Ghana OTP sent successfully to:", user.PhoneNumber)
	case "234":
		if _, err := emailService.SendNaijaOtp(user.PhoneNumber, user.ID); err != nil {
			log.Println("Failed to send Nigeria OTP:", err)
			utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Nigeria OTP: %v", err))
			return
		}
		log.Println("Nigeria OTP sent successfully to:", user.PhoneNumber)
	default:
		log.Println("Unsupported country code for phone number:", req.PhoneNumber)
		utils.Error(ctx, http.StatusBadRequest, "Unsupported country code")
		return
	}

	// Process referral if applicable
	if referrer != nil {
		referralPoints := 100.0
		log.Println("Processing referral for referrer ID:", referrer.ID)
		err = config.DB.Transaction(func(tx *gorm.DB) error {
			expiresAt := time.Now().AddDate(1, 0, 0) // 1 year
			referral := models.Referral{
				ReferrerID:     referrer.ID,
				ReferredUserID: user.ID,
				PointsEarned:   float32(referralPoints),
				ExpiresAt:      expiresAt,
			}
			if err := tx.Create(&referral).Error; err != nil {
				log.Println("Failed to create referral record:", err)
				return err
			}
			log.Println("Referral record created successfully for referrer ID:", referrer.ID)
			return tx.Model(&models.ReferralWallet{}).
				Where("user_id = ?", referrer.ID).
				Update("referral_balance", gorm.Expr("referral_balance + ?", referralPoints)).Error
		})
		if err != nil {
			log.Println("Failed to process referral:", err)
		} else {
			log.Println("Referral processed successfully for referrer ID:", referrer.ID)
		}
	}

	// Success response
	log.Println("SignUpHandler completed successfully for phone number:", req.PhoneNumber)
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
	log.Println("VerifyAccountHandler invoked")
	var req VerifyAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to bind JSON:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		log.Println("Validation error:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	// 🔍 Fetch user with OTP security
	log.Println("Fetching user with phone number:", req.PhoneNumber)
	var user models.User
	if err := config.DB.Preload("OtpSecurity").
		Where("phone_number = ?", req.PhoneNumber).
		First(&user).Error; err != nil {
		log.Println("User not found for phone number:", req.PhoneNumber)
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	// Already verified?
	if user.IsVerified {
		log.Println("User already verified for phone number:", req.PhoneNumber)
		utils.Error(ctx, http.StatusConflict, "User already verified")
		return
	}

	
	
	var otp models.UserOtpSecurity
		if err := config.DB.Where("user_id = ? AND action = ?", user.ID, models.OtpActionVerifyAccount).
			First(&otp).Error; err != nil {
			log.Println("OTP not found for account verification:", err)
			utils.Error(ctx, http.StatusNotFound, "OTP not found for account verification")
			return
		}

	if otp.Code != req.VerificationCode {
		log.Println("Verification code mismatch for user ID:", user.ID)
		utils.Error(ctx, http.StatusConflict, "Invalid verification code")
		return
	}

	if otp.CreatedAt.IsZero() || otp.ExpiresAt == nil || otp.ExpiresAt.IsZero() {
		log.Println("OTP metadata is incomplete or missing for user ID:", user.ID)
		utils.Error(ctx, http.StatusBadRequest, "OTP metadata is incomplete or missing")
		return
	}

	if time.Now().After(*otp.ExpiresAt) {
		log.Println("OTP has expired for user ID:", user.ID)
		utils.Error(ctx, http.StatusBadRequest, "OTP has expired")
		return
	}

	// ✅ Transaction to update user + clear OTP
	log.Println("Starting transaction to verify user ID:", user.ID)
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).Update("is_verified", true).Error; err != nil {
			log.Println("Failed to update user verification status for user ID:", user.ID, "Error:", err)
			return err
		}

		if err := tx.Model(&models.UserOtpSecurity{}).
			Where("user_id = ?", user.ID).
			Updates(map[string]interface{}{
				"code":         "",
				"expires_at":   nil,
				"created_at":   nil,
				"locked_until": nil,
				"action":      "",
				"retry_count":  0,
			}).Error; err != nil {
			log.Println("Failed to clear OTP fields for user ID:", user.ID, "Error:", err)
			return err
		}

		log.Println("Transaction completed successfully for user ID:", user.ID)
		return nil
	})

	if err != nil {
		log.Println("VerifyAccount transaction failed for user ID:", user.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to verify account")
		return
	}

	// 🎟️ Generate tokens
	log.Println("Generating access token for user ID:", user.ID)
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		log.Println("Failed to generate access token for user ID:", user.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	log.Println("Generating refresh token for user ID:", user.ID)
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Println("Failed to generate refresh token for user ID:", user.ID, "Error:", err)
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
	log.Println("Saving refresh token for user ID:", user.ID)
	if err := config.DB.Create(&rt).Error; err != nil {
		log.Println("Failed to save refresh token for user ID:", user.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	log.Println("VerifyAccountHandler completed successfully for user ID:", user.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User verified account successfully.",
		"user": gin.H{
			"id":                 user.ID,
			"fullName":           user.FullName,
			"phoneNumber":        user.PhoneNumber,
			"email":              user.Email,
			"isActive":           user.IsActive,
			"isVerified":          user.IsVerified,
			"isEmailVerified":   user.IsEmailVerified,
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
	log.Println("ResendOtpHandler invoked")
	var req ResendOtpRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to bind JSON:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		log.Println("Validation error:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	var user models.User
	log.Println("Fetching user with phone number:", req.PhoneNumber)
	if err := config.DB.Preload("OtpSecurity").
		Where("phone_number = ?", req.PhoneNumber).
		First(&user).Error; err != nil {
		log.Println("User not found for phone number:", req.PhoneNumber)
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	if user.IsVerified {
		log.Println("User already verified for phone number:", req.PhoneNumber)
		utils.Error(ctx, http.StatusConflict, "User already verified")
		return
	}

	otp := user.OtpSecurity
	currentTime := time.Now()

	// Check lockout
	if otp != nil && otp.LockedUntil != nil && currentTime.Before(*otp.LockedUntil) {
		remainingTime := int((*otp.LockedUntil).Sub(currentTime).Minutes())
		log.Printf("User ID %d is locked out. Remaining time: %d minute(s)\n", user.ID, remainingTime)
		utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Please wait %d minute(s) before requesting a new OTP.", remainingTime))
		return
	}

	// Check cooldown
	if otp != nil && otp.CreatedAt != nil {
		timeSinceLastOtp := currentTime.Sub(*otp.CreatedAt)
		if timeSinceLastOtp < time.Duration(OTP_RESEND_COOLDOWN)*time.Second {
			remainingCooldown := OTP_RESEND_COOLDOWN - int(timeSinceLastOtp.Seconds())
			log.Printf("User ID %d is in cooldown period. Remaining cooldown: %d seconds\n", user.ID, remainingCooldown)
			utils.Error(ctx, http.StatusTooManyRequests, fmt.Sprintf("Please wait %d seconds before requesting a new OTP.", remainingCooldown))
			return
		}
	}

	// Check retry limit
	if otp != nil && otp.RetryCount >= MAX_OTP_RETRIES {
		log.Printf("User ID %d has exceeded maximum OTP retries. Locking out for %d minutes\n", user.ID, OTP_LOCKOUT_DURATION.Minutes())
		config.DB.Model(&models.UserOtpSecurity{}).
			Where("user_id = ?", user.ID).
			Updates(map[string]interface{}{
				"locked_until": time.Now().Add(OTP_LOCKOUT_DURATION),
				"retry_count":  0,
			})
		utils.Error(ctx, http.StatusTooManyRequests, "Too many OTP attempts. Please try again later.")
		return
	}

	// Determine country by phone prefix
	countryPrefix := req.PhoneNumber[:3]
	emailService := services.EmailService{}
	log.Println("Determining country by phone prefix:", countryPrefix)

	switch countryPrefix {
	case "233":
		log.Println("Sending Ghana OTP to phone number:", user.PhoneNumber)
		if _, err := emailService.SendGhanaOtp(user.PhoneNumber, user.ID); err != nil {
			log.Println("Failed to send Ghana OTP:", err)
			utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Ghana OTP: %v", err))
			return
		}
		log.Println("Ghana OTP sent successfully to:", user.PhoneNumber)
	case "234":
		log.Println("Sending Nigeria OTP to phone number:", user.PhoneNumber)
		if _, err := emailService.SendNaijaOtp(user.PhoneNumber, user.ID); err != nil {
			log.Println("Failed to send Nigeria OTP:", err)
			utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Nigeria OTP: %v", err))
			return
		}
		log.Println("Nigeria OTP sent successfully to:", user.PhoneNumber)
	default:
		log.Println("Unsupported country code for phone number:", req.PhoneNumber)
		utils.Error(ctx, http.StatusBadRequest, "Unsupported country code")
		return
	}

	// Update retry count + lockout timestamp
	log.Printf("Updating OTP retry count and lockout timestamp for user ID %d\n", user.ID)
	config.DB.Model(&models.UserOtpSecurity{}).
		Where("user_id = ?", user.ID).
		Updates(map[string]interface{}{
			"retry_count":  otp.RetryCount + 1,
			"locked_until": time.Now().Add(VERIFY_OTP_LOCKED_DURATION),
		})

	log.Println("ResendOtpHandler completed successfully for user ID:", user.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "New OTP sent successfully. Please check your SMS.",
		"user": gin.H{
			"id":                 user.ID,
			"fullName":           user.FullName,
			"phoneNumber":        user.PhoneNumber,
			"email":              user.Email,
			"isActive":           user.IsActive,
			"isVerified":         user.IsVerified,
			"isEmailVerified":   user.IsEmailVerified,
			"isProfileCreated":   user.IsProfileCreated,
			"countryOfResidence": user.CountryOfResidence,
			"gender":             user.Gender,
			"dateOfBirth":        user.DateOfBirth,
		},
	})
}

func LoginHandler(ctx *gin.Context) {
	log.Println("LoginHandler invoked")
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to bind JSON:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		log.Println("Validation error:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := req.Validate(); err != nil {
		log.Println("Request validation failed:", err)
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	log.Println("Fetching user with email or phone number:", req.Email, req.PhoneNumber)
	if err := config.DB.Where("email = ? OR phone_number = ?", req.Email, req.PhoneNumber).
		Preload("OtpSecurity").
		First(&user).Error; err != nil {
		log.Println("Raw query failed:", err)
		log.Println("User not found for email or phone number:", req.Email, req.PhoneNumber)
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}
	

	if !user.IsActive {
		log.Println("Account is blocked for user ID:", user.ID)
		utils.Error(ctx, http.StatusBadRequest, "Account is blocked, please contact support")
		return
	}
   
	if req.Email != "" && !user.IsEmailVerified {
    log.Println("Login attempt with unverified email for user ID:", user.ID)
    utils.Error(ctx, http.StatusBadRequest, "Your email is not verified. Please visit your profile to complete verification.")
    return
}


	if !user.IsVerified {
		log.Println("User not verified for user ID:", user.ID)
		countryPrefix := req.PhoneNumber[:3]
		emailService := services.EmailService{}

		switch countryPrefix {
		case "233":
			log.Println("Sending Ghana OTP to phone number:", user.PhoneNumber)
			if _, err := emailService.SendGhanaOtp(user.PhoneNumber, user.ID); err != nil {
				log.Println("Failed to send Ghana OTP:", err)
				utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Ghana OTP: %v", err))
				return
			}
		case "234":
			log.Println("Sending Nigeria OTP to phone number:", user.PhoneNumber)
			if _, err := emailService.SendNaijaOtp(user.PhoneNumber, user.ID); err != nil {
				log.Println("Failed to send Nigeria OTP:", err)
				utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Nigeria OTP: %v", err))
				return
			}
		default:
			log.Println("Unsupported country code for phone number:", req.PhoneNumber)
			utils.Error(ctx, http.StatusBadRequest, "Unsupported country code")
			return
		}
		utils.Error(ctx,http.StatusForbidden,"Verification OTP sent. Please verify your account to continue.")
		return
	}

	if !utils.ComparePassword(user.Password, req.Password) {
		log.Println("Invalid credentials for user ID:", user.ID)
		utils.Error(ctx, http.StatusForbidden, "Invalid credentials")
		return
	}

	log.Println("Generating access token for user ID:", user.ID)
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		log.Println("Failed to generate access token for user ID:", user.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	log.Println("Generating refresh token for user ID:", user.ID)
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Println("Failed to generate refresh token for user ID:", user.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	log.Println("Updating last login timestamp for user ID:", user.ID)
	config.DB.Model(&user).Update("last_login", time.Now())

	// Upsert refresh token
	expiresAt := time.Now().AddDate(0, 0, config.AppConfig.RefreshTokenExpiresDays)
	rt := models.RefreshToken{
		UserID:   user.ID,
		Token:    refreshToken,
		ExpireAt: &expiresAt,
	}
	log.Println("Upserting refresh token for user ID:", user.ID)
	config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"token", "expire_at"}),
	}).Create(&rt)

	log.Println("LoginHandler completed successfully for user ID:", user.ID)
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
			"isEmailVerified":   user.IsEmailVerified,
			"lastLogin":          user.LastLogin,
			"profilePicture":     user.ProfilePicture,
			"accessToken":        accessToken,
			"refreshToken":       refreshToken,
		},
	})
}

// ChangePasswordHandler changes user password
func ChangePasswordHandler(ctx *gin.Context) {
	log.Println("ChangePasswordHandler invoked")
	var req PasswordChangeRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to bind JSON:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		log.Println("Validation error:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := req.Validate(); err != nil {
		log.Println("Request validation failed:", err)
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	log.Println("Current user fetched with ID:", currentUser.ID)

	if !currentUser.IsVerified {
		log.Println("User is not verified. User ID:", currentUser.ID)
		utils.Error(ctx, http.StatusBadRequest, "Only verified account can change password")
		return
	}

	if !utils.ComparePassword(currentUser.Password, req.CurrentPassword) {
		log.Println("Current password is incorrect for user ID:", currentUser.ID)
		utils.Error(ctx, http.StatusBadRequest, "Current password is incorrect")
		return
	}

	if req.CurrentPassword == req.NewPassword {
		log.Println("New password is the same as current password for user ID:", currentUser.ID)
		utils.Error(ctx, http.StatusBadRequest, "New password can not be the same as current password")
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Println("Failed to hash new password for user ID:", currentUser.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to process password")
		return
	}

	log.Println("Updating password for user ID:", currentUser.ID)
	if err := config.DB.Model(&currentUser).Update("password", hashedPassword).Error; err != nil {
		log.Println("Failed to update password for user ID:", currentUser.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to update password")
		return
	}

	log.Println("Generating access token for user ID:", currentUser.ID)
	accessToken, err := utils.GenerateAccessToken(currentUser.ID)
	if err != nil {
		log.Println("Failed to generate access token for user ID:", currentUser.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	log.Println("Generating refresh token for user ID:", currentUser.ID)
	refreshToken, err := utils.GenerateRefreshToken(currentUser.ID)
	if err != nil {
		log.Println("Failed to generate refresh token for user ID:", currentUser.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	log.Println("Deleting old refresh tokens for user ID:", currentUser.ID)
	if err := config.DB.Where("user_id = ?", currentUser.ID).Delete(&models.RefreshToken{}).Error; err != nil {
		log.Println("Failed to delete old refresh tokens for user ID:", currentUser.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to delete old refresh tokens")
		return
	}

	expireAt := time.Now().Add(time.Hour * 24 * time.Duration(config.AppConfig.RefreshTokenExpiresDays))
	newRefresh := models.RefreshToken{
		UserID:   currentUser.ID,
		Token:    refreshToken,
		ExpireAt: &expireAt,
	}
	log.Println("Creating new refresh token for user ID:", currentUser.ID)
	if err := config.DB.Create(&newRefresh).Error; err != nil {
		log.Println("Failed to save new refresh token for user ID:", currentUser.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	notification := models.Notification{
		UserID:  currentUser.ID,
		Title:   "Password Change Successful",
		Message: fmt.Sprintf("Your password change for %s has been successful.", currentUser.Email),
		Type:    models.PasswordChange,
		IsRead:  false,
	}
	log.Println("Creating notification for password change for user ID:", currentUser.ID)
	if err := config.DB.Create(&notification).Error; err != nil {
		log.Println("Failed to save notification for user ID:", currentUser.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to save notification")
		return
	}

	log.Println("Password change successful for user ID:", currentUser.ID)
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
	log.Println("ForgotPasswordHandler invoked")
	var req ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to bind JSON:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		log.Println("Validation error:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := req.Validate(); err != nil {
		log.Println("Request validation failed:", err)
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	query := config.DB.Preload("OtpSecurity")
	if req.Email != "" {
		log.Println("Searching user by email:", req.Email)
		query = query.Where("email = ?", req.Email)
	} else {
		log.Println("Searching user by phone number:", req.PhoneNumber)
		query = query.Where("phone_number = ?", req.PhoneNumber)
	}

	if err := query.First(&user).Error; err != nil {
		log.Println("User not found for provided email or phone number")
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	now := time.Now()

	// Lockout check
	if user.OtpSecurity != nil && user.OtpSecurity.LockedUntil != nil &&
		now.Before(*user.OtpSecurity.LockedUntil) {
		remaining := int(user.OtpSecurity.LockedUntil.Sub(now).Minutes())
		log.Printf("User ID %d is locked out. Remaining time: %d minute(s)\n", user.ID, remaining)
		utils.Error(ctx, http.StatusBadRequest,
			fmt.Sprintf("Please wait %d minute(s) before requesting a new OTP.", remaining))
		return
	}

	// Active OTP check
	if user.OtpSecurity != nil && user.OtpSecurity.ExpiresAt != nil &&
		now.Before(*user.OtpSecurity.ExpiresAt) {
		remaining := int(user.OtpSecurity.ExpiresAt.Sub(now).Minutes())
		log.Printf("Active OTP exists for user ID %d. Remaining time: %d minute(s)\n", user.ID, remaining)
		utils.Error(ctx, http.StatusBadRequest,
			fmt.Sprintf("An active OTP already exists. Please wait %d minute(s) before requesting a new OTP.", remaining))
		return
	}

	// Retry count
	if user.OtpSecurity != nil && user.OtpSecurity.RetryCount >= MAX_OTP_RETRIES {
		lockUntil := now.Add(OTP_LOCKOUT_DURATION)
		log.Printf("User ID %d has exceeded maximum OTP retries. Locking out until: %v\n", user.ID, lockUntil)
		config.DB.Model(&models.UserOtpSecurity{}).
			Where("user_id = ?", user.ID).
			Updates(map[string]interface{}{
				"locked_until": lockUntil,
				"retry_count":  0,
			})
		utils.Error(ctx, http.StatusTooManyRequests,
			"You have exceeded the maximum OTP attempts. Please wait before trying again.")
		return
	}

	// Send OTP (phone or email)
	emailService := services.EmailService{}

	if req.PhoneNumber != "" {
		cleanPhone := strings.ReplaceAll(req.PhoneNumber, " ", "")
		log.Println("Sending OTP to phone number:", cleanPhone)

		var err error
		switch {
		case strings.HasPrefix(cleanPhone, "234"):
			log.Println("Detected Nigeria phone number. Sending OTP...")
			if _, err = emailService.SendForgotPasswordNGNOtp(cleanPhone, user.ID); err != nil {
				log.Println("Failed to send Nigeria OTP:", err)
				utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Nigeria OTP: %v", err))
				return
			}

		case strings.HasPrefix(cleanPhone, "233"):
			log.Println("Detected Ghana phone number. Sending OTP...")
			if _, err = emailService.SendForgotPasswordGHCOtp(cleanPhone, user.ID); err != nil {
				log.Println("Failed to send Ghana OTP:", err)
				utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to send Ghana OTP: %v", err))
				return
			}

		default:
			log.Println("Unsupported phone country code:", cleanPhone)
			utils.Error(ctx, http.StatusBadRequest, "Unsupported phone country code")
			return
		}

		if err != nil {
			log.Println("Failed to send OTP:", err)
			utils.Error(ctx, http.StatusInternalServerError, "Failed to send OTP")
			return
		}

	} else if req.Email != "" {
		log.Println("Sending OTP to email:", req.Email)
		if _, err := emailService.SendForgotPasswordEmailOtp(req.Email, user.FullName, user.ID); err != nil {
			log.Println("Failed to send OTP email:", err)
			utils.Error(ctx, http.StatusInternalServerError, "Failed to send OTP email")
			return
		}
	} else {
		log.Println("Neither phone number nor email provided")
		utils.Error(ctx, http.StatusBadRequest, "Either phone number or email is required")
		return
	}

	// Update retry count + lock until
	lockUntil := now.Add(FORGOT_PASSWORD_OTP_LOCKED_DURATION)
	log.Printf("Updating OTP retry count and lockout timestamp for user ID %d\n", user.ID)
	config.DB.Model(&models.UserOtpSecurity{}).
		Where("user_id = ?", user.ID).
		Updates(map[string]interface{}{
			"retry_count":  gorm.Expr("retry_count + ?", 1),
			"locked_until": lockUntil,
		})

	log.Println("ForgotPasswordHandler completed successfully for user ID:", user.ID)
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
	log.Println("VerifyPasswordForgotOtpHandler invoked")
	var req VerifyPasswordForgotOtpRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to bind JSON:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		log.Println("Validation error:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	var user models.User
	log.Println("Fetching user with email or phone number:", req.Email, req.PhoneNumber)
	if err := config.DB.Preload("OtpSecurity").
		Where("email = ? OR phone_number = ?", req.Email, req.PhoneNumber).
		First(&user).Error; err != nil {
		log.Println("User not found for provided email or phone number")
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	if !user.IsVerified {
		log.Println("User is not verified. User ID:", user.ID)
		utils.Error(ctx, http.StatusForbidden, "Please verify your account before resetting password")
		return
	}

	var otpSec models.UserOtpSecurity
		if err := config.DB.Where("user_id = ? AND action = ?", user.ID, models.OtpActionPasswordReset).
			First(&otpSec).Error; err != nil {
			log.Println("Password reset OTP not found for user ID:", user.ID, "Error:", err)
			utils.Error(ctx, http.StatusNotFound, "Password reset OTP not found")
			return
		}

	// Check OTP destination
	if otpSec.SentTo == "email" && req.Email == "" {
		log.Println("OTP was sent to email but email not provided for user ID:", user.ID)
		utils.Error(ctx, http.StatusBadRequest, "OTP was sent to email, please provide email")
		return
	}
	if otpSec.SentTo == "phone" && req.PhoneNumber == "" {
		log.Println("OTP was sent to phone but phone number not provided for user ID:", user.ID)
		utils.Error(ctx, http.StatusBadRequest, "OTP was sent to phone, please provide phone number")
		return
	}

	// Verify code
	if otpSec.Code != req.VerificationCode {
		log.Println("Invalid verification code for user ID:", user.ID)
		utils.Error(ctx, http.StatusConflict, "Invalid verification code")
		return
	}

	if time.Now().After(*otpSec.ExpiresAt) {
		log.Println("OTP has expired for user ID:", user.ID)
		utils.Error(ctx, http.StatusBadRequest, "OTP has expired")
		return
	}

	// Mark verified
	log.Println("Marking OTP as verified for password reset for user ID:", user.ID)
	config.DB.Model(&models.UserOtpSecurity{}).
		Where("user_id = ?", user.ID).
		Update("is_otp_verified_for_password_reset", true)

	log.Println("VerifyPasswordForgotOtpHandler completed successfully for user ID:", user.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP verification successful",
		"userId":  user.ID,
		"email":   user.Email,
	})
}

// ResetPasswordUser handles actual password reset
func ResetPasswordHandler(ctx *gin.Context) {
	log.Println("ResetPasswordHandler invoked")
	var req ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to bind JSON:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		log.Println("Validation error:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	var user models.User
	log.Println("Fetching user with ID:", req.UserId)
	if err := config.DB.Preload("OtpSecurity").
		Where("id = ?", req.UserId).
		First(&user).Error; err != nil {
		log.Println("User not found for ID:", req.UserId)
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	if user.OtpSecurity == nil || !user.OtpSecurity.IsOtpVerifiedForPasswordReset {
		log.Println("OTP verification required for user ID:", user.ID)
		utils.Error(ctx, http.StatusForbidden, "OTP verification required")
		return
	}

	// Hash new password
	log.Println("Hashing new password for user ID:", user.ID)
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Println("Failed to hash new password for user ID:", user.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to process password")
		return
	}

	log.Println("Updating password for user ID:", user.ID)
	if err := config.DB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("password", hashedPassword).Error; err != nil {
		log.Println("Failed to update password for user ID:", user.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to update password")
		return
	}

	// Clear OTP fields
	log.Println("Clearing OTP fields for user ID:", user.ID)
	if err := config.DB.Model(&models.UserOtpSecurity{}).
		Where("user_id = ?", user.ID).
		Updates(map[string]interface{}{
			"code":         "",
			"is_otp_verified_for_password_reset":  false,
			"sent_to":      nil,
			"action":       "",
			"locked_until": nil,
			"created_at":   nil,
			"expires_at":   nil,
			"retry_count":  0,
		}).Error; err != nil {
		log.Println("Failed to clear OTP fields for user ID:", user.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to clear OTP fields")
		return
	}

	log.Println("Password reset successful for user ID:", user.ID)
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
	log.Println("LogoutHandler invoked")
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		log.Println("Authorization header missing")
		utils.Error(ctx, http.StatusUnauthorized, "Authorization header missing")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		log.Println("Invalid token format")
		utils.Error(ctx, http.StatusUnauthorized, "Invalid token format")
		return
	}

	// Decode token to get expiration
	log.Println("Decoding token")
	claims, err := utils.DecodeToken(tokenString)
	if err != nil {
		log.Println("Invalid token:", err)
		utils.Error(ctx, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Determine expiration time
	var expireAt time.Time
	if exp, ok := claims["exp"].(float64); ok {
		expireAt = time.Unix(int64(exp), 0)
		log.Println("Token expiration time determined:", expireAt)
	} else {
		// fallback if token has no exp: blacklist for 1 hour
		expireAt = time.Now().Add(1 * time.Hour)
		log.Println("Token has no expiration time, defaulting to 1 hour from now:", expireAt)
	}

	// Blacklist the token
	log.Println("Blacklisting token")
	if err := utils.BlacklistToken(tokenString, expireAt); err != nil {
		log.Println("Failed to blacklist token:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to blacklist token")
		return
	}

	log.Println("LogoutHandler completed successfully")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User logged out successfully",
	})
}

func RefreshTokenHandler(ctx *gin.Context) {
	log.Println("RefreshTokenHandler invoked")
	var req RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to bind JSON:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		log.Println("Validation error:", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	var tokenEntry models.RefreshToken
	log.Println("Fetching refresh token from database")
	err := config.DB.
		Where("token = ?", req.RefreshToken).
		Order("created_at DESC").
		First(&tokenEntry).Error
	if err != nil {
		log.Println("Refresh token not found or session expired:", err)
		utils.Error(ctx, http.StatusUnauthorized, "Session expired")
		return
	}

	isExpired := tokenEntry.ExpireAt != nil && time.Now().After(*tokenEntry.ExpireAt)
	isValid := false
	var userId string

	log.Println("Verifying refresh token")
	userId, err = utils.VerifyJWTRefreshToken(tokenEntry.Token)
	if err == nil && !isExpired {
		isValid = true
	}

	log.Println("Deleting refresh token from database (one-time use)")
	config.DB.Delete(&tokenEntry)

	if !isValid || userId == "" {
		log.Println("Refresh token is invalid or expired")
		utils.Error(ctx, http.StatusUnauthorized, "Session expired")
		return
	}

	var user models.User
	log.Println("Fetching user with ID:", userId)
	if err := config.DB.First(&user, "id = ?", userId).Error; err != nil {
		log.Println("User not found for ID:", userId, "Error:", err)
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	log.Println("Generating new access token for user ID:", user.ID)
	accessToken, err := utils.GenerateAccessToken(user.ID)

	if err != nil {
		log.Println("Failed to generate access token for user ID:", user.ID, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	log.Println("RefreshTokenHandler completed successfully for user ID:", user.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
	})
}
