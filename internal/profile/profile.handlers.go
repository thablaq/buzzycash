package profile

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"gorm.io/gorm"
	"time"
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/services"
	"github.com/dblaq/buzzycash/internal/utils"
	"github.com/dblaq/buzzycash/pkg/externals"
	"github.com/gin-gonic/gin"
)

func CreateProfileHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	var existingUser models.User
	if err := config.DB.First(&existingUser, "id = ?", currentUser.ID).Error; err != nil {
		log.Printf("Error finding user: %v", err)
		utils.Error(ctx, http.StatusBadRequest, "User not found")
		return
	}
	if existingUser.IsProfileCreated {
		log.Printf("Profile creation attempt for user %s, but profile already exists", currentUser.ID)
		utils.Error(ctx, http.StatusBadRequest, "Profile has already been created")
		return
	}

	var req CreateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON for user %s: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		log.Printf("Validation error for user %s: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	var emailTaken, usernameTaken models.User
	config.DB.Where("email = ? AND id <> ?", req.Email, currentUser.ID).First(&emailTaken)
	config.DB.Where("username = ? AND id <> ?", req.UserName, currentUser.ID).First(&usernameTaken)
	if emailTaken.ID != "" {
		log.Printf("Email %s is already taken for user %s", req.Email, currentUser.ID)
		utils.Error(ctx, http.StatusBadRequest, "Email already taken")
		return
	}
	if usernameTaken.ID != "" {
		log.Printf("Username %s is already taken for user %s", req.UserName, currentUser.ID)
		utils.Error(ctx, http.StatusBadRequest, "Username already taken")
		return
	}

	updateData := map[string]interface{}{
		"full_name":          req.FullName,
		"gender":             req.Gender,
		"email":              strings.ToLower(req.Email),
		"username":           req.UserName,
		"is_profile_created": true,
	}
	if err := config.DB.Model(&existingUser).Updates(updateData).Error; err != nil {
		log.Printf("Error updating profile for user %s: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to create profile")
		return
	}

	if err := config.DB.First(&existingUser, "id = ?", currentUser.ID).Error; err != nil {
		log.Printf("Error reloading updated user %s: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to load updated user")
		return
	}

	nameParts := strings.Fields(currentUser.FullName)
	firstName := "-"
	lastName := "-"
	if len(nameParts) > 0 {
		firstName = nameParts[0]
	}
	if len(nameParts) > 1 {
		lastName = strings.Join(nameParts[1:], " ")
	}

	gs := externals.NewGamingService()

	_, err := gs.RegisterUser(currentUser.PhoneNumber, req.Email, firstName, lastName)
	if err != nil {
		log.Printf("Error registering user %s with gaming service: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Failed to register with gaming service: %v", err))
		return
	}

	log.Printf("Profile created successfully for user %s", currentUser.ID)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User profile created successfully",
		"user": gin.H{
			"id":               existingUser.ID,
			"phoneNumber":      existingUser.PhoneNumber,
			"fullName":         existingUser.FullName,
			"gender":           existingUser.Gender,
			"dateOfBirth":      existingUser.DateOfBirth,
			"email":            existingUser.Email,
			"isProfileCreated": existingUser.IsProfileCreated,
			"isVerified":       existingUser.IsVerified,
			"isActive":         existingUser.IsActive,
			"username":         existingUser.Username,
		},
	})
}

func GetUserProfileHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	var profile models.User
	if err := config.DB.Select(
		"full_name", "email", "is_active", "is_verified", "last_login",
		"phone_number", "country_of_residence", "username", "gender", "is_profile_created",
	).First(&profile, "id = ?", currentUser.ID).Error; err != nil {
		log.Printf("Error retrieving profile for user %s: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	log.Printf("Profile retrieved successfully for user %s", currentUser.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User profile retrieved successfully",
		"user": gin.H{
			"id":               currentUser.ID,
			"phoneNumber":      currentUser.PhoneNumber,
			"fullName":         currentUser.FullName,
			"gender":           currentUser.Gender,
			"dateOfBirth":      currentUser.DateOfBirth,
			"email":            currentUser.Email,
			"isProfileCreated": currentUser.IsProfileCreated,
			"isVerified":       currentUser.IsVerified,
			"isActive":         currentUser.IsActive,
			"username":         currentUser.Username,
			"countryOfResidence": currentUser.CountryOfResidence,
			"lastLogin": currentUser.LastLogin,
		},
	})
}

func UpdateUserProfileHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	var req ProfileUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON for user %s: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		log.Printf("Validation error for user %s: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	var existingUser models.User
	if err := config.DB.First(&existingUser, "id = ?", currentUser.ID).Error; err != nil {
		log.Printf("Error finding user %s: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusBadRequest, "User not found")
		return
	}

	updateData := map[string]interface{}{}
	if req.FullName != "" {
		updateData["full_name"] = req.FullName
	}
	if req.Gender != "" {
		updateData["gender"] = req.Gender
	}
	if req.DateOfBirth != "" {
		updateData["date_of_birth"] = req.DateOfBirth
	}

	var updatedUser models.User
	if err := config.DB.Model(&existingUser).Updates(updateData).Scan(&updatedUser).Error; err != nil {
		log.Printf("Error updating profile for user %s: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	log.Printf("Profile updated successfully for user %s", currentUser.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user": gin.H{
			"id":            updatedUser.ID,
			"full_name":     updatedUser.FullName,
			"gender":        updatedUser.Gender,
			"date_of_birth": updatedUser.DateOfBirth,
		},
	})
}



func RequestEmailVerificationHandler(ctx *gin.Context) {
	log.Println("EmailVerificationHandler invoked")

	currentUser := ctx.MustGet("currentUser").(models.User)
	log.Printf("Current user retrieved: ID=%s, Email=%s", currentUser.ID, currentUser.Email)

	var existingUser models.User
	if err := config.DB.Where("email = ?", currentUser.Email).First(&existingUser).Error; err != nil {
		log.Printf("No email associated with this account for user ID=%s: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusBadRequest, "No email associated with this account")
		return
	}

	if existingUser.IsEmailVerified {
		log.Printf("Email already verified for user ID=%s", currentUser.ID)
		utils.Error(ctx, http.StatusForbidden, "Email already verified")
		return
	}

	emailService := services.EmailService{}
	log.Printf("EmailService initialized for user ID=%s", currentUser.ID)

	if currentUser.Email != nil && *currentUser.Email != "" {
		log.Printf("Sending email verification OTP to user ID=%s, Email=%s", currentUser.ID, *currentUser.Email)
		_, err := emailService.SendEmailVerificationOtp(*currentUser.Email, currentUser.FullName, currentUser.ID)
		if err != nil {
			log.Printf("Failed to send verification email to user ID=%s: %v", currentUser.ID, err)
			utils.Error(ctx, http.StatusInternalServerError, "Failed to send verification email")
			return
		}

		log.Printf("Verification email sent successfully to user ID=%s, Email=%s", currentUser.ID, *currentUser.Email)
		ctx.JSON(http.StatusOK, gin.H{
			"message":          "Verification email sent successfully",
			"id":               currentUser.ID,
			"email":            currentUser.Email,
			"isEmailVerified":  currentUser.IsEmailVerified,
		})
		return
	}

	log.Printf("No valid email found for user ID=%s", currentUser.ID)
	utils.Error(ctx, http.StatusBadRequest, "No valid email found")
}



func VerifyAccountEmailHandler(ctx *gin.Context) {
	log.Println("VerifyAccountHandler invoked")
	var req VerifyEmailProfileRequest
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
	log.Println("Fetching user with phone number:", req.Email)
	var user models.User
	if err := config.DB.Preload("OtpSecurity").
		Where("email = ?", req.Email).
		First(&user).Error; err != nil {
		log.Println("User not found for email:", req.Email)
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	// Already verified?
	if user.IsEmailVerified {
		log.Println("User already verified for email:", req.Email)
		utils.Error(ctx, http.StatusConflict, "User already verified")
		return
	}

	
	
	var otp models.UserOtpSecurity
		if err := config.DB.Where("user_id = ? AND action = ?", user.ID, models.OtpActionVerifyEmail).
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
		if err := tx.Model(&user).Update("is_email_verified", true).Error; err != nil {
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


	log.Println("VerifyEmailHandler completed successfully for user ID:", user.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User verified email successfully.",
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
		},
	})
}