package notifications



import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/dblaq/buzzycash/internal/middlewares"
)


func NotificationRoutes(rg *gin.RouterGroup, db *gorm.DB){
	notifyHandler := NewNotifyHandler(db)
	notificationRoutes := rg.Group("/notification")
	{
		notificationRoutes.GET("/",middlewares.AuthMiddleware,notifyHandler.GetNotificationsHandler)
		notificationRoutes.GET("/unread",middlewares.AuthMiddleware, notifyHandler.GetUnreadNotificationsCountHandler)
		notificationRoutes.PATCH("/:notificationId/read",middlewares.AuthMiddleware,notifyHandler.MarkAsReadHandler)
		notificationRoutes.PATCH("/read-all",middlewares.AuthMiddleware,notifyHandler.MarkAllAsReadHandler)
	}
}