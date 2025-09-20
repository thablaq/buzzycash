package notifications

import (
	"fmt"
	"net/http"
	"strconv"
     "log"
    "sort"
	// "github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
	"github.com/gin-gonic/gin"
		"gorm.io/gorm"
)


type NotifyHandler struct {
	db *gorm.DB
}


func NewNotifyHandler(db *gorm.DB) *NotifyHandler {
	return &NotifyHandler{
		db: db,
	}
}



func (h *NotifyHandler)GetNotificationsHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	notifType := ctx.Query("type") // "games", "transactions", or empty

	var responses []NotificationResponse

	// Fetch Transactions
	if notifType == "" || notifType == "transactions" {
		var notifs []models.Notification
		if err := h.db.Where("user_id = ?", currentUser.ID).
			Order("created_at desc").Offset(offset).Limit(limit).
			Find(&notifs).Error; err == nil {
			for _, n := range notifs {
				responses = append(responses, mapNotificationToResponse(n))
			}
		}
}


	// Fetch Games
	if notifType == "" || notifType == "games" {
		var notifs []models.Notification
		if err := h.db.Where("user_id = ?", currentUser.ID).
			Order("created_at desc").Offset(offset).Limit(limit).
			Find(&notifs).Error; err == nil {
			for _, g := range notifs {
				responses = append(responses, mapNotificationToResponse(g))
			}
		}
	}

	// Sort all responses by CreatedAt DESC
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].CreatedAt.After(responses[j].CreatedAt)
	})

	hasMore := false
	if len(responses) > limit {
		hasMore = true
		responses = responses[:limit]
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Notifications retrieved successfully",
		"notifications": responses,
		"page":          page,
		"limit":         limit,
		"hasMore":       hasMore,
	})
}




func (h *NotifyHandler)GetUnreadNotificationsCountHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	notifType := ctx.Query("type")

	query := h.db.Model(&models.Notification{}).
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


func(h *NotifyHandler) MarkAsReadHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	notificationID := ctx.Param("notificationId")
	if notificationID == "" {
		log.Printf("Invalid notification ID provided by user %d", currentUser.ID)
		utils.Error(ctx, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	var notification models.Notification
	if err := h.db.Where("id = ? AND user_id = ?", notificationID, currentUser.ID).First(&notification).Error; err != nil {
		log.Printf("Notification %s not found for user %d: %v", notificationID, currentUser.ID, err)
		utils.Error(ctx, http.StatusNotFound, "Notification not found")
		return
	}

	notification.IsRead = true
	if err := h.db.Save(&notification).Error; err != nil {
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


func (h *NotifyHandler)MarkAllAsReadHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	// Optional filter by type
	notificationType := ctx.Query("type")

	n := h.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", currentUser.ID, false)

	if notificationType != "" {
		n = n.Where("type = ?", notificationType)
	}

	result := n.Updates(map[string]interface{}{"is_read": true})

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


