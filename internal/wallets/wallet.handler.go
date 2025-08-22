package wallets

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/pkg/externals"
	"github.com/dblaq/buzzycash/internal/helpers"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
)



func GetUserBalanceHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	userID := currentUser.ID
	username := currentUser.PhoneNumber
	log.Printf("[GetUserBalance] Initiating get user balance request for userID: %s, username: %s\n", userID, username)

	// Fetch user from DB
	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	if !user.IsVerified {
		utils.Error(ctx, http.StatusBadRequest, "User not verified")
		return
	}

	gs := externals.NewGamingService()
	result, err := gs.GetUserWallet(username)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch user wallet")
		return
	}

	ctx.JSON(http.StatusOK,
		gin.H{
			"message": "User wallet retrieved successfully",
			"result":  result,
		},
	)
}


func FundWalletHandler(ctx *gin.Context) {
	var req CreditWalletRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	userID := currentUser.ID
	username := currentUser.PhoneNumber
	email := currentUser.Email

	log.Printf("[CreditWallet] Initiating credit wallet request for userID: %s, username: %s, email: %s\n", userID, username, email)

	// Ensure user exists
	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	transactionRef := helpers.GenerateTransactionReference()
	log.Printf("[transactionRef] Generated transaction reference: %s\n", transactionRef)

	gs := externals.NewGamingService()
	result, err := gs.CreditUserWallet(username, req.Amount)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate payment")
		return
	}

	// Record pending transaction
	history := models.TransactionHistory{
		AmountPaid:           req.Amount,
		CustomerEmail:        email,
		UserID:               userID,
		PaymentStatus:        models.Pending,
		PaymentMethod:        models.Nomba,
		TransactionReference: transactionRef,
		TransactionType:      models.Credit,
		Category:             models.Deposit,
		Currency:             "NGN",
	
	}
	if err := config.DB.Create(&history).Error; err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to create transaction record")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":             "Generated payment link successfully",
		"checkoutLink":        &result.Message,
		"amountPaid":          req.Amount,
		"customerEmail":       email,
		"userID":              userID,
		"paymentStatus":       models.Pending,
		"paymentMethod":       models.Nomba,
		"transactionReference": transactionRef,
		"transactionType":     models.Credit,
		"category":            models.Deposit,
		"currency":            "NGN",
	})
}

