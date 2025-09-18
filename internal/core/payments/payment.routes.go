package payments

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PaymentRoutes(rg *gin.RouterGroup,db *gorm.DB) {
	webhookHandler := NewWebhookHandler(db)
	paymentRoutes := rg.Group("/webhook")
	{
		paymentRoutes.POST("/fw", webhookHandler.FlutterwaveWebhookHandler)
	}
}
