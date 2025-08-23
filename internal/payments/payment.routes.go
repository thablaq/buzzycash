package payments




import (
	"github.com/gin-gonic/gin"
)

func PaymentRoutes(rg *gin.RouterGroup) {
	paymentRoutes := rg.Group("/payments")
	{
		paymentRoutes.POST("/webhook", FlutterwaveWebhookHandler) 
	}
}