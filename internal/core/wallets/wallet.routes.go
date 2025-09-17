package wallets



import (
	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/middlewares"
	"gorm.io/gorm"
)

func WalletRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	walletHandler := NewWalletHandler(db)
	walletRoutes := rg.Group("/wallet")
	{
		walletRoutes.POST("/fund-wallet", middlewares.AuthMiddleware,FundWalletHandler)
		walletRoutes.GET("/get-wallet", middlewares.AuthMiddleware,walletHandler.GetUserBalanceHandler)
	}
}