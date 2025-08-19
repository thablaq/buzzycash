package virtual

import (
	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/middlewares"
)

func VirtualRoutes(rg *gin.RouterGroup) {
	virtualRoutes := rg.Group("/virtual")
	{
		virtualRoutes.POST("/start-game", middlewares.AuthMiddleware,StartVirtualGameHandler)
		virtualRoutes.GET("/get-games", middlewares.AuthMiddleware,GetVirtualGamesHandler)
	}
}