package payments


import (
	// "fmt"
	"net/http"
	"encoding/json"
	"log"
	"io"
	"time"
	"crypto/subtle"

	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/pkg/externals"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
)



func VerifyPaymentHandler(ctx *gin.Context) {
	log.Printf("INFO: VerifyPaymentHandler started.\n")

	paymentID := ctx.Query("paymentId")
	if paymentID == "" {
		log.Printf("WARN: VerifyPaymentHandler: Missing payment ID in request.\n")
		utils.Error(ctx, http.StatusBadRequest, "payment id is required")
		return
	}
	log.Printf("INFO: VerifyPaymentHandler: Received request for paymentID: %s\n", paymentID)

	// Find transaction by paymentID
	var transaction models.TransactionHistory
	if err := config.DB.Where("payment_id = ?", paymentID).First(&transaction).Error; err != nil {
		log.Printf("ERROR: VerifyPaymentHandler: Transaction not found for paymentID %s: %v\n", paymentID, err)
		utils.Error(ctx, http.StatusNotFound, "Transaction not found")
		return
	}
	log.Printf("INFO: VerifyPaymentHandler: Found transaction (ID: %s, UserID: %s, Status: %s) for paymentID: %s.\n", transaction.ID, transaction.UserID, transaction.PaymentStatus, paymentID)

	// Check for duplicate successful transaction
	if transaction.PaymentStatus == models.Successful {
		log.Printf("WARN: VerifyPaymentHandler: Duplicate successful transaction detected for paymentID %s.\n", paymentID)
		utils.Error(ctx, http.StatusBadRequest, "duplicate reference detected")
		return
	}
	log.Printf("INFO: VerifyPaymentHandler: Transaction %s is not yet successful. Proceeding with external verification.\n", paymentID)

	// Call external gaming service to verify payment
	gs := externals.NewGamingService()
	log.Printf("INFO: VerifyPaymentHandler: Calling GamingService to verify payment for paymentID: %s.\n", paymentID)
	result, err := gs.VerifyPayment(paymentID)
	if err != nil {
		log.Printf("ERROR: VerifyPaymentHandler: Failed to verify payment with GamingService for paymentID %s: %v\n", paymentID, err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to verify payment")
		return
	}
	log.Printf("INFO: VerifyPaymentHandler: GamingService verification successful for paymentID %s. Raw result: %+v\n", paymentID, result)

	// Extract message from result safely
	msg, ok := result["message"].(string)
	if !ok {
		log.Printf("ERROR: VerifyPaymentHandler: Invalid response from GamingService (missing 'message' field or not a string) for paymentID %s. Result: %+v\n", paymentID, result)
		utils.Error(ctx, http.StatusInternalServerError, "Invalid response from GamingService")
		return
	}
	log.Printf("INFO: VerifyPaymentHandler: Extracted message from GamingService response: '%s' for paymentID: %s.\n", msg, paymentID)

	userID := transaction.UserID

	// Prepare notification based on payment result
	var notification models.Notification
	if msg == "PAYMENT SUCCESSFUL" && transaction.PaymentStatus == models.Pending {
		log.Printf("INFO: VerifyPaymentHandler: Payment %s is successful and transaction was pending. Updating status to Successful.\n", paymentID)
		transaction.PaymentStatus = models.Successful
		if err := config.DB.Save(&transaction).Error; err != nil {
			log.Printf("ERROR: VerifyPaymentHandler: Failed to update transaction status to Successful for paymentID %s: %v\n", paymentID, err)
			utils.Error(ctx, http.StatusInternalServerError, "Failed to update transaction status")
			return
		}
		log.Printf("INFO: VerifyPaymentHandler: Transaction %s status successfully updated to Successful.\n", paymentID)

		// notification = models.Notification{
		// 	UserID:  userID,
		// 	Title:   "Wallet top-up",
		// 	Message: fmt.Sprintf("Your payment of ₦%.2f was successful.", *transaction.AmountPaid),
		// 	Type:    models.Wallet,
		// 	IsRead:  false,
		// }
		// log.Printf("INFO: VerifyPaymentHandler: Prepared successful wallet top-up notification for UserID: %s, Amount: %.2f.\n", userID, *transaction.AmountPaid)
	} else {
		log.Printf("INFO: VerifyPaymentHandler: Payment %s either not successful or not pending. Preparing failed notification.\n", paymentID)
		// notification = models.Notification{
		// 	UserID:  userID,
		// 	Title:   "Wallet top-up failed",
		// 	Message: fmt.Sprintf("Your payment of ₦%.2f could not be verified.", *transaction.AmountPaid),
		// 	Type:    models.Wallet,
		// 	IsRead:  false,
		// }
		// log.Printf("INFO: VerifyPaymentHandler: Prepared failed wallet top-up notification for UserID: %s, Amount: %.2f.\n", userID, *transaction.AmountPaid)
	}

	// Save notification
	if err := config.DB.Create(&notification).Error; err != nil {
		log.Printf("ERROR: VerifyPaymentHandler: Failed to create notification for UserID %s, paymentID %s: %v\n", userID, paymentID, err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to create notification")
		return
	}
	log.Printf("INFO: VerifyPaymentHandler: Notification successfully created for UserID: %s.\n", userID)

	// Return response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    result,
		"message": "Payment verification processed successfully",
	})
	log.Printf("INFO: VerifyPaymentHandler: Payment verification processed successfully and response sent for paymentID: %s.\n", paymentID)
}




func FlutterwaveWebhookHandler(ctx *gin.Context) {
    secret := config.AppConfig.FlutterwaveHashKey
    sent := ctx.GetHeader("verif-hash")

    // Verify hash
    if secret == "" || subtle.ConstantTimeCompare([]byte(secret), []byte(sent)) != 1 {
        log.Printf("[FW Webhook] Invalid signature. expected=%s got=%s", secret, sent)
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
        return
    }

    // Read raw body
    body, _ := io.ReadAll(ctx.Request.Body)
    log.Printf("[FW Webhook] Raw body: %s", string(body)) // 🔥 full webhook payload

    // Try to unmarshal
    var evt FlutterwaveWebhook
    if err := json.Unmarshal(body, &evt); err != nil {
        log.Printf("[FW Webhook] JSON unmarshal error: %v", err)
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad payload"})
        return
    }

    log.Printf("[FW Webhook] Parsed event: %+v", evt) // 🔥 full struct after unmarshal

    // Only care about successful charges
    if (evt.EventType == "charge.completed"|| evt.EventType == "BANK_TRANSFER_TRANSACTION") && evt.Status == "successful" {
        txRef := evt.TxRef
        amount := evt.Amount
        currency := evt.Currency

        log.Printf("[FW Webhook] SUCCESSFUL payment. tx_ref=%s id=%d amount=%.2f currency=%s",
            txRef, evt.ID, amount, currency)

        // Update the pending transaction to successful
        if err := config.DB.Model(&models.TransactionHistory{}).
            Where("transaction_reference = ?", txRef).
            Updates(map[string]interface{}{
                "payment_status": models.Successful,
                "paid_at":        time.Now(),
                "amount":         amount,
            }).Error; err != nil {
            log.Printf("[FW Webhook] DB update failed: %+v\n", err)
        }

        // Now credit the user’s wallet
        var history models.TransactionHistory
        if err := config.DB.Preload("User").
            Where("transaction_reference = ?", txRef).First(&history).Error; err == nil {
            gs := externals.NewGamingService()
            if _, err := gs.CreditUserWallet(history.User.PhoneNumber, amount); err != nil {
                log.Printf("[FW Webhook] Wallet credit failed: %v", err)
            } else {
                log.Printf("[FW Webhook] Wallet credited for user=%s amount=%.2f", history.User.PhoneNumber, amount)
            }
        } else {
            log.Printf("[FW Webhook] Could not find transaction history for tx_ref=%s", txRef)
        }
    } else {
        log.Printf("[FW Webhook] Ignored event=%s status=%s", evt, evt.Status)
    }

    ctx.Status(http.StatusOK)
}

