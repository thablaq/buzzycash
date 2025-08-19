package profile

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
	"github.com/dblaq/buzzycash/pkg/externals"
	"github.com/gin-gonic/gin"
)

// CreateProfile creates a new profile for the logged-in user
func CreateProfileHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	// Check if profile already exists
	var existingUser models.User
	if err := config.DB.First(&existingUser, "id = ?", currentUser.ID).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "User not found")
		return
	}
	if existingUser.IsProfileCreated {
		utils.Error(ctx, http.StatusBadRequest, "Profile has already been created")
		return
	}

	// Validate request body
	var req CreateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
}

	// Check for existing email or username
	var emailTaken, usernameTaken models.User
	config.DB.Where("email = ? AND id <> ?", req.Email, currentUser.ID).First(&emailTaken)
	config.DB.Where("username = ? AND id <> ?", req.UserName, currentUser.ID).First(&usernameTaken)
	if emailTaken.ID != "" {
		utils.Error(ctx, http.StatusBadRequest, "Email already taken")
		return
	}
	if usernameTaken.ID != "" {
		utils.Error(ctx, http.StatusBadRequest, "Username already taken")
		return
	}

	// Update user profile
	updateData := map[string]interface{}{
		"full_name":          req.FullName,
		"gender":             req.Gender,
		"email":              strings.ToLower(req.Email),
		"username":           req.UserName,
		"is_profile_created": true,
	}
	if err := config.DB.Model(&existingUser).Updates(updateData).Error; err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to create profile")
		return
	}

	if err := config.DB.First(&existingUser, "id = ?", currentUser.ID).Error; err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to load updated user")
		return
	}

	// Split full name for GamingService
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
		utils.Error(ctx, http.StatusBadRequest, fmt.Sprintf("Failed to register with gaming service: %v", err))
		return
	}

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

// GetUserProfile fetches the logged-in user's profile
func GetUserProfileHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	var profile models.User
	if err := config.DB.Select(
		"full_name", "email", "is_active", "is_verified", "last_login",
		"phone_number", "country_of_residence", "username", "gender", "is_profile_created",
	).First(&profile, "id = ?", currentUser.ID).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

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
			"countryOfResidence":currentUser.CountryOfResidence,
			"lastLogin": currentUser.LastLogin,
		},
	})
}

// UpdateProfile updates specific fields in the user's profile
func UpdateUserProfileHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	var req ProfileUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	
	if err := utils.Validate.Struct(req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
}

	// Check if user exists
	var existingUser models.User
	if err := config.DB.First(&existingUser, "id = ?", currentUser.ID).Error; err != nil {
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
		utils.Error(ctx, http.StatusInternalServerError, "Failed to update profile")
		return
	}

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
