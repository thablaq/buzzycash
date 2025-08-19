package uploadimages


import (
	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/middlewares"
)


func UploadRoutes(rg *gin.RouterGroup){
	uploadRoutes := rg.Group("/upload")
	{
		uploadRoutes.POST("/user",middlewares.AuthMiddleware,UploadProfileHandler)
		// profileRoutes.GET("/get-profile",middlewares.AuthMiddleware, GetUserProfileHandler)
		// profileRoutes.PATCH("/update-profile",middlewares.AuthMiddleware,UpdateUserProfileHandler)
		
	}
}