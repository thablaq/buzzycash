package transaction

import (
	"github.com/dblaq/buzzycash/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func TransactionRoutes(rg *gin.RouterGroup) {
	transactionRoutes := rg.Group("/transactions")
	{
		transactionRoutes.GET("/history", middlewares.AuthMiddleware, GetTransactionHistoryHandler)
		transactionRoutes.GET("/search", middlewares.AuthMiddleware, SearchTransactionHistoryHandler)
		transactionRoutes.GET("/:id", middlewares.AuthMiddleware, GetTransactionByID)

	}
}
