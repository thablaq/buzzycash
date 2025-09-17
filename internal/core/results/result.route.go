package results

import (
	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/middlewares"
)


func ResultRoutes(rg *gin.RouterGroup){
	resultRoutes := rg.Group("/result")
	{
		resultRoutes.GET("/winners",middlewares.AuthMiddleware,GetWinnerLogsHandler)
		resultRoutes.GET("/leaderboard",middlewares.AuthMiddleware, GetLeaderBoardHandler)
		resultRoutes.GET("/user-results",middlewares.AuthMiddleware,GetUserResultsHandler)
		
	}
}