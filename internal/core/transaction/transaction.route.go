package transaction

import (
	"github.com/dblaq/buzzycash/internal/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TransactionRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	transactionHandler := NewTransactionHandler(db)
	transactionRoutes := rg.Group("/transactions")
	{
		transactionRoutes.GET("/history", middlewares.AuthMiddleware, GetTransactionHistoryHandler)
		transactionRoutes.GET("/search", middlewares.AuthMiddleware, SearchTransactionHistoryHandler)
		transactionRoutes.GET("/:id", middlewares.AuthMiddleware, transactionHandler.GetTransactionByID)

	}
}
