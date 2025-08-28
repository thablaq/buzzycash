package withdrawal

import (
	"log"
	"net/http"

	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/helpers"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
	"github.com/dblaq/buzzycash/pkg/gateway"
	"github.com/gin-gonic/gin"
)

// nomba routes payout
// do kyc above 250k

func ListBanksHandler(ctx *gin.Context) {

	currentUser := ctx.MustGet("currentUser").(models.User)
	var user models.User
	if err := config.DB.First(&user, "id = ?", currentUser.ID).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	ps := gateway.NBInstance()
	banks, err := ps.ListNBBanks()
	if err != nil {
		log.Printf("Flutterwave error: %v", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch banks")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Fetched all banks successfully",
		"data":    banks,
	})
}

func RetrieveAccountDetailsHandler(ctx *gin.Context) {
	var req RetrieveAccountDetailsRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	log.Println("Attempting to validate request struct")
	if err := req.Validate(); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	var user models.User
	if err := config.DB.First(&user, "id = ?", currentUser.ID).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	nbReq := gateway.NBRetrieveAccountDetails{
		AccountNumber: req.AccountNumber,
		BankCode:      req.BankCode,
	}

	ps := gateway.NBInstance()
	accountDetails, err := ps.FetchAccountDetails(nbReq)
	if err != nil {
		log.Printf("Flutterwave error: %v", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch account details")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Fetched user account successfully",
		"data":    accountDetails,
	})
}

func InitiateWithdrawalHandler(ctx *gin.Context) {
	var req InitiateWithdrawalRequest

	log.Println("Attempting to bind JSON request for Withdrawal")
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	log.Println("Attempting to validate request struct")
	if err := req.Validate(); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	userID := currentUser.ID
	email := currentUser.Email

	log.Printf("[CreditWallet] Initiating credit wallet request for userID: %s, email: %s\n", userID, email)

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	if !currentUser.IsEmailVerified {
		utils.Error(ctx, http.StatusForbidden, "Please verify your email to proceed")
		return
	}

	transactionRef := helpers.GenerateTransactionReference()
	log.Printf("[transactionRef] Generated transaction reference: %s\n", transactionRef)

	reference := helpers.GenerateFWRef()
	log.Printf("[reference] Generated  reference: %s\n", reference)

	nbReq := gateway.NBWithdrawalRequest{
		MerchantTxRef: reference,
		Amount:        req.Amount,
		BankCode:      req.BankCode,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
		Narration:     "Buzzycash withdrawal",
		SenderName:    "BuzzyCash",
	}
	ps := gateway.NBInstance()
	_, err := ps.InitiateWithdrawal(nbReq)
	if err != nil {
		log.Printf("Nomba Withdrawal error: %v", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate payment")
		return
	}

	// Record pending transaction
	history := models.TransactionHistory{
		Amount:               req.Amount,
		CustomerEmail:        email,
		UserID:               currentUser.ID,
		PaymentStatus:        models.Pending,
		PaymentMethod:        models.Nomba,
		TransactionReference: transactionRef,
		Reference:            reference,
		TransactionType:      models.Withdrawal,
		Category:             models.WithdrawRequest,
		PaymentType:         models.Payout,
		Currency:             "NGN",
	}
	log.Printf("DEBUG TransactionHistory UserID: '%s'", history.UserID)
	if err := config.DB.Create(&history).Error; err != nil {
		log.Printf("DB history creation failed: %+v\n", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to record transaction")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":              "Generated payment link successfully",
		"amountPaid":           req.Amount,
		"customerEmail":        email,
		"userID":               userID,
		"paymentStatus":        models.Pending,
		"paymentMethod":        models.Nomba,
		"transactionReference": transactionRef,
		"reference":            reference,
		"transactionType":      models.Withdrawal,
		"category":             models.WithdrawRequest,
		"paymentType":          models.Payout,
		"currency":             "NGN",
	})
}
