package profile



import (
	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/middlewares"
)


func ProfileRoutes(rg *gin.RouterGroup){
	profileRoutes := rg.Group("/profile")
	{
		profileRoutes.POST("/create-profile",middlewares.AuthMiddleware,CreateProfileHandler)
		profileRoutes.GET("/get-profile",middlewares.AuthMiddleware, GetUserProfileHandler)
		profileRoutes.PATCH("/update-profile",middlewares.AuthMiddleware,UpdateUserProfileHandler)
		
	}
}