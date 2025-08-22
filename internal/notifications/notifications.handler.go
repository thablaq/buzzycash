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

	// Optional type filter (games / transactions)
	notifType := ctx.Query("type")

	var notifications []models.Notification
	query := config.DB.Where("user_id = ?", currentUser.ID)

	if notifType != "" {
		query = query.Where("type = ?", notifType)
	}

	result := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit + 1). 
		Find(&notifications)

	if result.Error != nil {
		log.Printf("Error fetching notifications for user %d: %v", currentUser.ID, result.Error)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}

	// Determine if there are more
	hasMore := false
	if len(notifications) > limit {
		hasMore = true
		notifications = notifications[:limit]
	}

	log.Printf("Notifications retrieved successfully for user %d", currentUser.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Notifications retrieved successfully",
		"notifications": notifications,
			"page":    page,
			"limit":   limit,
			"hasMore": hasMore,
		
	})
}


func GetUnreadNotificationsCountHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	notifType := ctx.Query("type")

	query := config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", currentUser.ID, false)

	if notifType != "" {
		query = query.Where("type = ?", notifType)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		log.Printf("Error fetching unread notifications count for user %d: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch unread notifications")
		return
	}

	log.Printf("Unread notifications count retrieved successfully for user %d", currentUser.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Unread notifications count retrieved successfully",
		"unreadCount": count,
		"type":        notifType, 
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

	// Optional filter by type
	notificationType := ctx.Query("type")

	db := config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", currentUser.ID, false)

	if notificationType != "" {
		db = db.Where("type = ?", notificationType)
	}

	result := db.Updates(map[string]interface{}{"is_read": true})

	if result.Error != nil {
		log.Printf("Error marking notifications as read for user %d (type: %s): %v",
			currentUser.ID, notificationType, result.Error)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to mark notifications as read")
		return
	}

	log.Printf("Notifications marked as read for user %d, type: %s, updated count: %d",
		currentUser.ID, notificationType, result.RowsAffected)

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Notifications marked as read",
		"type":         notificationType,
		"updatedCount": result.RowsAffected,
	})
}


