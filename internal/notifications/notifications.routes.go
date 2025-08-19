package notifications



import (
	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/middlewares"
)


func NotificationRoutes(rg *gin.RouterGroup){
	notificationRoutes := rg.Group("/notification")
	{
		notificationRoutes.GET("/",middlewares.AuthMiddleware,GetNotificationsHandler)
		notificationRoutes.GET("/unread",middlewares.AuthMiddleware, GetUnreadNotificationsCountHandler)
		notificationRoutes.PATCH("/:notificationId/read",middlewares.AuthMiddleware,MarkAsReadHandler)
		notificationRoutes.PATCH("/read-all",middlewares.AuthMiddleware,MarkAllAsReadHandler)
	}
}