package tickets




import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/dblaq/buzzycash/internal/middlewares"
)
func TicketRoutes(rg *gin.RouterGroup,db *gorm.DB){
	ticketHandler := NewTicketHandler(db)
	ticketRoutes := rg.Group("/ticket")
	{
		ticketRoutes.POST("/purchase-ticket",middlewares.AuthMiddleware,ticketHandler.BuyGameTicketHandler)
		ticketRoutes.GET("/get-tickets",middlewares.AuthMiddleware, GetUserGameTicketsHandler)
		ticketRoutes.GET("/gaming",middlewares.AuthMiddleware, GetAllGamesHandler)
		ticketRoutes.POST("/create-game",middlewares.AuthMiddleware, CreateGameHandler)

	}
}