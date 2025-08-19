package tickets




import (
	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/middlewares"
)


func TicketRoutes(rg *gin.RouterGroup){
	ticketRoutes := rg.Group("/ticket")
	{
		ticketRoutes.POST("/purchase-ticket",middlewares.AuthMiddleware,BuyGameTicketHandler)
		ticketRoutes.GET("/get-tickets",middlewares.AuthMiddleware, GetUserGameTicketsHandler)
		
	}
}