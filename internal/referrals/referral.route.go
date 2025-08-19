package referral 


import (
	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/middlewares"
)


func ReferralRoutes(rg *gin.RouterGroup){
	referralRoutes := rg.Group("/referrals")
	{
		referralRoutes.GET("/referral-details",middlewares.AuthMiddleware,GetReferralDetailsHandler)
		
	}
}