package main

import (
	"fmt"
 ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
   "github.com/dblaq/buzzycash/docs"


	"github.com/gin-gonic/gin"
     "github.com/dblaq/buzzycash/internal/middlewares"
	"github.com/dblaq/buzzycash/internal/auth"
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/notifications"
	"github.com/dblaq/buzzycash/internal/profile"
	"github.com/dblaq/buzzycash/internal/referrals"
	"github.com/dblaq/buzzycash/internal/results"
	"github.com/dblaq/buzzycash/internal/upload-images"
	"github.com/dblaq/buzzycash/internal/virtual"
	"github.com/dblaq/buzzycash/internal/tickets"
	"github.com/dblaq/buzzycash/internal/wallets"
	"github.com/dblaq/buzzycash/internal/payments"
	"github.com/dblaq/buzzycash/internal/withdrawal"
	"github.com/dblaq/buzzycash/internal/transaction"
)


// @title BuzzyCash API
// @version 1.0
// @description REST API for BuzzyCash platform
// @host localhost:5005
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	 config.LoadConfig()

	config.InitDB()

	r := gin.Default()
	
	r.Static("/uploads/profile-pictures", "./uploads/profile-pictures")

	r.Use(gin.Logger())
	r.Use(middlewares.RecoveryAndErrorMiddleware())
	r.Use(gin.Recovery())
	
	
 // ✅ Set Swagger info
    docs.SwaggerInfo.BasePath = "/api/v1"
    docs.SwaggerInfo.Host = "localhost:5005"
	


	api := r.Group("/api/v1")

	api.GET("/welcome", func(ctx *gin.Context) {
		ctx.String(200, "Welcome to BuzzyCash API")
	})
	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,url))


	auth.AuthRoutes(api)
	notifications.NotificationRoutes(api)
	profile.ProfileRoutes(api)
	referral.ReferralRoutes(api)
	results.ResultRoutes(api)
	uploadimages.UploadRoutes(api)
	virtual.VirtualRoutes(api)
	tickets.TicketRoutes(api)
	wallets.WalletRoutes(api)
	payments.PaymentRoutes(api)
	withdrawal.WithdrawalRoutes(api)
	transaction.TransactionRoutes(api)

	fmt.Println("🚀 Server started on :" + config.AppConfig.Port)
	r.Run(":" + config.AppConfig.Port)
}
 