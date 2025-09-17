package profile




import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	
	// "github.com/dblaq/buzzycash/internal/handlers"
	"github.com/dblaq/buzzycash/internal/middlewares"
)

func ProfileRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// Initialize the profile handler with database dependency
	profileHandler := NewProfileHandler(db)
	
	profileRoutes := rg.Group("/profile")
	{
		profileRoutes.POST("/create-profile", middlewares.AuthMiddleware, profileHandler.CreateProfileHandler)
		profileRoutes.GET("/get-profile", middlewares.AuthMiddleware, profileHandler.GetUserProfileHandler)
		profileRoutes.PATCH("/update-profile", middlewares.AuthMiddleware, profileHandler.UpdateUserProfileHandler)
		profileRoutes.POST("/request-verification", middlewares.AuthMiddleware, profileHandler.RequestEmailVerificationHandler)
		profileRoutes.POST("/verify-email", middlewares.AuthMiddleware, profileHandler.VerifyAccountEmailHandler)
		profileRoutes.GET("/", profileHandler.ChooseUsernameHandler)
	}
}