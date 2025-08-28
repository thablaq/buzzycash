package withdrawal




import (
	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/middlewares"
)

func WithdrawalRoutes(rg *gin.RouterGroup) {
	withdrawalRoutes := rg.Group("/withdrawal")
	{
		withdrawalRoutes.GET("/list-banks", middlewares.AuthMiddleware,ListBanksHandler)
		withdrawalRoutes.POST("/account-details", middlewares.AuthMiddleware,RetrieveAccountDetailsHandler)
		withdrawalRoutes.POST("/initiate-withdrawal", middlewares.AuthMiddleware,InitiateWithdrawalHandler)
	}
}