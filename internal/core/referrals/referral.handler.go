package referral

import (
	"net/http"
     "log"
	"github.com/gin-gonic/gin"
		"gorm.io/gorm"
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
)


type ReferralHandler struct {
	db *gorm.DB
}


func NewReferralHandler(db *gorm.DB) *ReferralHandler {
	return &ReferralHandler{
		db: db,
	}
}


func (h *ReferralHandler)GetReferralDetailsHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)


	var referralWallet models.ReferralWallet
	if err := h.db.
		Preload("Earnings").
		Where("user_id = ?", currentUser.ID).
		First(&referralWallet).Error; err != nil {
		log.Printf("Error fetching referral wallet for user ID %s: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusNotFound, "Referral wallet not found")
		return
	}


	var inviteesCount int64
	if err := config.DB.Model(&models.ReferralEarning{}).
		Where("referrer_id = ?", currentUser.ID).
		Count(&inviteesCount).Error; err != nil {
		log.Printf("Error counting referrals for user ID %s: %v", currentUser.ID, err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch referrals")
		return
	}

	
	totalEarned := referralWallet.ReferralBalance


	ctx.JSON(http.StatusOK, gin.H{
		"referralCode": currentUser.ReferralCode,
		"totalEarned":  totalEarned,
		"invitees":     inviteesCount,
		"message":      "Referral details retrieved successfully",
	})
}
