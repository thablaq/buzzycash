package uploadimages

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	// "strings"

	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/models"
)

// Allowed MIME types
var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"application/pdf": true,
}

const maxFileSize = 10 * 1024 * 1024 // 10 MB
const uploadDir = "uploads/profile-pictures"

// helper to generate random UUID
func randomUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ensure upload directory exists
func init() {
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}
}

// UploadProfileHandler handles profile picture uploads for both users and admins
func UploadProfileHandler(ctx *gin.Context) {
	// Get current user/admin from context
	currentUser, userExists := ctx.Get("currentUser")
	// currentAdmin, adminExists := ctx.Get("currentAdmin")

	if !userExists{
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	
	file, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	
	if file.Size > maxFileSize {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	
	mime := file.Header.Get("Content-Type")
	if !allowedMimeTypes[mime] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}


	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", randomUUID(), ext)
	filePath := filepath.Join(uploadDir, fileName)

	// Save uploaded file
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	avatarURL := fmt.Sprintf("/%s/%s", uploadDir, fileName)

	// Update DB
	if userExists {
		user := currentUser.(models.User)

		// Delete old file if exists
		if user.ProfilePicture != "" {
			oldFile := filepath.Join(uploadDir, filepath.Base(user.ProfilePicture))
			os.Remove(oldFile)
		}

		if err := config.DB.Model(&user).Update("profile_picture", avatarURL).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"message":   "User profile picture updated",
			"avatarUrl": avatarURL,
		})
		return
	}

	// if adminExists {
	// 	admin := currentAdmin.(models.Admin)

	// 	// Delete old file if exists
	// 	if admin.ProfilePicture != "" {
	// 		oldFile := filepath.Join(uploadDir, filepath.Base(admin.ProfilePicture))
	// 		os.Remove(oldFile)
	// 	}

	// 	if err := config.DB.Model(&admin).Update("profile_picture", avatarURL).Error; err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
	// 		return
	// 	}

	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"status":    "success",
	// 		"message":   "Admin profile picture updated",
	// 		"avatarUrl": avatarURL,
	// 	})
	// 	return
	// }
}
