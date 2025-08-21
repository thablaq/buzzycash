package wallets



import (
	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/middlewares"
)

func WalletRoutes(rg *gin.RouterGroup) {
	walletRoutes := rg.Group("/wallet")
	{
		walletRoutes.POST("/fund-wallet", middlewares.AuthMiddleware,FundWalletHandler)
		walletRoutes.GET("/get-wallet", middlewares.AuthMiddleware,GetUserBalanceHandler)
	}
}