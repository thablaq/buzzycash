package referral

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
)

func GetReferralDetailsHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	// Fetch referral wallet
	var referralWallet models.ReferralWallet
		if err := config.DB.Where("user_id = ?", currentUser.ID).First(&referralWallet).Error; err != nil {
			utils.Error(ctx,http.StatusNotFound,"Referral wallet not found")
			return
		}

		// Fetch referrals made by the user
		var referrals []models.Referral
		if err := config.DB.Where("referrer_id = ?", currentUser.ID).Find(&referrals).Error; err != nil {
			utils.Error(ctx,http.StatusInternalServerError, "Failed to fetch referrals")
			return
		}

	invitees := len(referrals)
	totalEarned := referralWallet.ReferralBalance
	

	// Return response
	ctx.JSON(http.StatusOK, gin.H{
			"referralCode":  currentUser.ReferralCode,
			"totalEarned":  totalEarned,
			"invitees":     invitees,
			"message":      "Referral details retrieved successfully",
		})
}
