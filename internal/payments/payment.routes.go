package payments




import (
	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/middlewares"
)

func PaymentRoutes(rg *gin.RouterGroup) {
	paymentRoutes := rg.Group("/payments")
	{
		paymentRoutes.GET("/verify", middlewares.AuthMiddleware,VerifyPaymentHandler)
		paymentRoutes.POST("/webhook", FlutterwaveWebhookHandler) 
	}
}