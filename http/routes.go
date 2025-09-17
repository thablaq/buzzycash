package http




import (
	"fmt"
	"github.com/dblaq/buzzycash/internal/core/auth"
	"github.com/dblaq/buzzycash/internal/core/notifications"
	"github.com/dblaq/buzzycash/internal/core/payments"
	"github.com/dblaq/buzzycash/internal/core/profile"
	"github.com/dblaq/buzzycash/internal/core/referrals"
	"github.com/dblaq/buzzycash/internal/core/results"
	"github.com/dblaq/buzzycash/internal/core/tickets"
	"github.com/dblaq/buzzycash/internal/core/transaction"
	"github.com/dblaq/buzzycash/internal/core/upload-images"
	"github.com/dblaq/buzzycash/internal/core/virtual"
	"github.com/dblaq/buzzycash/internal/core/wallets"
	"github.com/dblaq/buzzycash/internal/core/withdrawal"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api/v1")

	// Webhook
	api.POST("/webhook/nomba", func(ctx *gin.Context) {
		var payload map[string]interface{}
		if err := ctx.BindJSON(&payload); err != nil {
			ctx.JSON(400, gin.H{"error": "invalid payload"})
			return
		}
		fmt.Printf("🔔 Nomba Webhook Received: %+v\n", payload)
		ctx.JSON(202, gin.H{"status": "success"})
	})

	// Feature routes
	auth.AuthRoutes(api,db)
	notifications.NotificationRoutes(api,db)
	profile.ProfileRoutes(api,db)
	referral.ReferralRoutes(api,db)
	results.ResultRoutes(api)
	uploadimages.UploadRoutes(api)
	virtual.VirtualRoutes(api)
	tickets.TicketRoutes(api)
	wallets.WalletRoutes(api,db)
	payments.PaymentRoutes(api,db)
	withdrawal.WithdrawalRoutes(api,db)
	transaction.TransactionRoutes(api,db)
}