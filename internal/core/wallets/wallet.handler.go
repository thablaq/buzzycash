package wallets

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"strings"
     "gorm.io/gorm"
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/external/gaming"
	"github.com/dblaq/buzzycash/external/gateway"
	"github.com/dblaq/buzzycash/internal/helpers"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
)

type WalletHandler struct {
	db *gorm.DB
}

// NewProfileHandler creates a new ProfileHandler with dependencies
func NewWalletHandler(db *gorm.DB) *WalletHandler {
	return &WalletHandler{
		db: db,
	}
}


func (h *WalletHandler)GetUserBalanceHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	userID := currentUser.ID
	username := currentUser.PhoneNumber
	log.Printf("[GetUserBalance] Initiating get user balance request for userID: %s, username: %s\n", userID, username)

	// Fetch user from DB
	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "User not found")
		return
	}

	if !user.IsVerified {
		utils.Error(ctx, http.StatusBadRequest, "User not verified")
		return
	}

	gs := gaming.GMInstance()
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
		log.Printf("[FundWallet] JSON binding error: %v\n", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}
	
 if err := req.Validate(); err != nil {
        utils.Error(ctx, http.StatusBadRequest, err.Error())
        return
    }

	currentUser := ctx.MustGet("currentUser").(models.User)
	fullName := currentUser.FullName
	email := currentUser.Email

	log.Printf("[CreditWallet] Initiating credit wallet request for userID: %s, email: %s, method: %s\n",
		currentUser.ID, email, req.PaymentMethod)

	transactionRef := helpers.GenerateTransactionReference()
	log.Printf("[FundWallet] Generated transaction reference: %s\n", transactionRef)
	reference := helpers.GenerateFWRef()
	log.Printf("[FundWallet] Generated payment gateway reference (for Flutterwave/Nomba): %s\n", reference)

	var checkoutLink, orderRef string
	var err error

	switch strings.ToLower(req.PaymentMethod) {
	case "flutterwave":
		fwReq := gateway.FWPaymentRequest{
			Reference:   reference,
			Amount:      req.Amount,
			Currency:    "NGN",
			RedirectURL: "Buzzycash://Home",
			Customer: gateway.FWCustomer{
				Email:    email,
				FullName: fullName,
			},
		}
		ps := gateway.FWInstance()
		checkoutLink, err = ps.CreateCheckout(fwReq)
		orderRef = reference 
	case "nomba":
	nbReq := gateway.NBPaymentRequest{
		Order: gateway.NBOrder{
			CallbackURL:   "Buzzycash://Home",
			CustomerEmail: email,
			Amount:        req.Amount,
			Currency:      "NGN",
			CustomerID:    currentUser.ID,
		},
		TokenizeCard: true,
}

		ps := gateway.NBInstance()
		checkoutLink, orderRef, err = ps.CreateNBCheckout(nbReq)
	default:
		log.Printf("[FundWallet] Invalid payment method requested: %s for userID: %s\n", req.PaymentMethod, currentUser.ID)
		utils.Error(ctx, http.StatusBadRequest, "Invalid payment method")
		return
	}

	if err != nil {
		log.Printf("Payment provider error (%s): %v", req.PaymentMethod, err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate payment")
		return
	}

	// Save transaction history
	history := models.Transaction{
		Amount:              req.Amount,
		CustomerEmail:       email,
		UserID:              currentUser.ID,
		PaymentStatus:       models.Pending,
		PaymentMethod:       models.EPaymentMethod(req.PaymentMethod),
		TransactionReference: transactionRef,
		Reference:            orderRef,
		TransactionType:     models.Credit,
		Category:            models.Deposit,
		PaymentType:         models.Topup,
		Currency:            "NGN",
	}
	if err := config.DB.Create(&history).Error; err != nil {
		log.Printf("DB history creation failed: %+v\n", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to record transaction")
		return
	}

	// Respond
	ctx.JSON(http.StatusOK, gin.H{
		"message":              "Generated payment link successfully",
		"checkoutLink":         checkoutLink,
		"amountPaid":           req.Amount,
		"customerEmail":        email,
		"userID":               currentUser.ID,
		"paymentStatus":        models.Pending,
		"paymentMethod":        req.PaymentMethod,
		 "paymentType":          models.Topup,
		"transactionReference": transactionRef,
		"reference":            orderRef,
		"transactionType":      models.Credit,
		"category":             models.Deposit,
		"currency":             "NGN",
	})
}



