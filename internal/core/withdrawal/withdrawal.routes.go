package withdrawal




import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/dblaq/buzzycash/internal/middlewares"
)

func WithdrawalRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	withdrawHandler := NewWithdrawHandler(db)
	withdrawalRoutes := rg.Group("/withdrawal")
	{
		withdrawalRoutes.GET("/list-banks", middlewares.AuthMiddleware,withdrawHandler.ListBanksHandler)
		withdrawalRoutes.POST("/account-details", middlewares.AuthMiddleware,withdrawHandler.RetrieveAccountDetailsHandler)
		withdrawalRoutes.POST("/initiate-withdrawal", middlewares.AuthMiddleware,withdrawHandler.InitiateWithdrawalHandler)
	}
}