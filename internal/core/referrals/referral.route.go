package referral 


import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/dblaq/buzzycash/internal/middlewares"
)


func ReferralRoutes(rg *gin.RouterGroup, db *gorm.DB){
	referralHandler := NewReferralHandler(db)
	referralRoutes := rg.Group("/referrals")
	{
		referralRoutes.GET("/referral-details",middlewares.AuthMiddleware,referralHandler.GetReferralDetailsHandler)
		
	}
}