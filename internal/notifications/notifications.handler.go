package notifications

import (
	"fmt"
	"net/http"
	"strconv"
     "log"
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
	"github.com/gin-gonic/gin"
)



func GetNotificationsHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var notifications []models.Notification
	result := config.DB.
		Where("user_id = ?", currentUser.ID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications)

	if result.Error != nil {
		log.Printf("Error fetching notifications for user %d: %v", currentUser.ID, result.Error)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}

	log.Printf("Notifications retrieved successfully for user %d", currentUser.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Notifications retrieved successfully",
		"notifications": notifications,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": result.RowsAffected,
		},
	})
}

func GetUnreadNotificationsCountHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	var count int64
	if err := config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", currentUser.ID, false).
		Count(&count).Error; err != nil {
		log.Printf("Error fetching unread notifications count for user %d: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch unread notifications")
		return
	}

	log.Printf("Unread notifications count retrieved successfully for user %d", currentUser.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Unread notifications count retrieved successfully",
		"unreadCount": count,
	})
}

func MarkAsReadHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	notificationID := ctx.Param("notificationId")
	if notificationID == "" {
		log.Printf("Invalid notification ID provided by user %d", currentUser.ID)
		utils.Error(ctx, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	var notification models.Notification
	if err := config.DB.Where("id = ? AND user_id = ?", notificationID, currentUser.ID).First(&notification).Error; err != nil {
		log.Printf("Notification %s not found for user %d: %v", notificationID, currentUser.ID, err)
		utils.Error(ctx, http.StatusNotFound, "Notification not found")
		return
	}

	notification.IsRead = true
	if err := config.DB.Save(&notification).Error; err != nil {
		log.Printf("Error marking notification %s as read for user %d: %v", notificationID, currentUser.ID, err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to mark notification as read")
		return
	}

	log.Printf("Notification %s marked as read for user %d", notificationID, currentUser.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"message":      fmt.Sprintf("Notification %s marked as read", notificationID),
		"notification": notification,
	})
}

func MarkAllAsReadHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	result := config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", currentUser.ID, false).
		Updates(map[string]interface{}{"is_read": true})

	if result.Error != nil {
		log.Printf("Error marking all notifications as read for user %d: %v", currentUser.ID, result.Error)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to mark notifications as read")
		return
	}

	log.Printf("All notifications marked as read for user %d, updated count: %d", currentUser.ID, result.RowsAffected)
	ctx.JSON(http.StatusOK, gin.H{
		"message":      "All notifications marked as read",
		"updatedCount": result.RowsAffected,
	})
}

