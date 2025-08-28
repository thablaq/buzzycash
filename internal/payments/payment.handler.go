package payments

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/utils"
	"github.com/gin-gonic/gin"
)

// func FlutterwaveWebhookHandler(ctx *gin.Context) {
// 	secret := config.AppConfig.FlutterwaveHashKey
// 	sent := ctx.GetHeader("verif-hash")

// 	// Verify hash
// 	if secret == "" || subtle.ConstantTimeCompare([]byte(secret), []byte(sent)) != 1 {
// 		log.Printf("[FW Webhook] Invalid signature. expected=%s got=%s", sent)
// 		utils.Error(ctx, http.StatusUnauthorized, "invalid signature")
// 		return
// 	}

// 	// Read raw body
// 	body, _ := io.ReadAll(ctx.Request.Body)
// 	log.Printf("[FW Webhook] Raw body: %s", string(body))

// 	var evt FlutterwaveWebhook
// 	if err := json.Unmarshal(body, &evt); err != nil {
// 		log.Printf("[FW Webhook] JSON unmarshal error: %v", err)
// 		utils.Error(ctx, http.StatusBadRequest, "bad payload")
// 		return
// 	}

// 	log.Printf("[FW Webhook] Parsed event: %+v", evt)

// 	if evt.EventType == "charge.completed" && evt.Status == "successful" {
// 		handleSuccessfulPayment(evt)
// 	} else if evt.EventType == "BANK_TRANSFER_TRANSACTION" && evt.Status == "successful" {
// 		handleSuccessfulPayment(evt)
// 	} else {
// 		log.Printf("[FW Webhook] Ignored event=%s status=%s", evt.EventType, evt.Status)
// 	}

// 	ctx.Status(http.StatusOK)
// }

// func handleSuccessfulPayment(evt FlutterwaveWebhook) {
// 	reference := evt.TxRef
// 	amount := evt.Amount
// 	currency := evt.Currency

// 	log.Printf("[FW Webhook] SUCCESSFUL payment. reference=%s id=%d amount=%.2f currency=%s",
// 		reference, evt.ID, amount, currency)

// 	// Update DB transaction
// 	if err := config.DB.Model(&models.TransactionHistory{}).
// 		Where("reference = ?", reference).
// 		Updates(map[string]interface{}{
// 			"payment_status": models.Successful,
// 			"paid_at":        time.Now(),
// 			"amount":         amount,
// 		}).Error; err != nil {
// 		log.Printf("[FW Webhook] DB update failed: %+v\n", err)
// 		return
// 	}

// 	// Credit wallet
// 	var history models.TransactionHistory
// 	if err := config.DB.Preload("User").
// 		Where("reference = ?", reference).First(&history).Error; err == nil {
// 		gs := gaming.GMInstance()
// 		if _, err := gs.CreditUserWallet(history.User.PhoneNumber, amount); err != nil {
// 			log.Printf("[FW Webhook] Wallet credit failed: %v", err)
// 		} else {
// 			log.Printf("[FW Webhook] Wallet credited for user=%s amount=%.2f", history.User.PhoneNumber, amount)
// 			notif := models.Notification{
// 				UserID:    history.UserID,
// 				Type:      "transaction",
// 				Title:     "Deposit Successful",
// 				Subtitle:  "Your wallet has been credited.",
// 				Amount:    amount,
// 				Currency:  currency,
// 				Status:    "successful",
// 				CreatedAt: time.Now(),
// 			}
// 			if err := config.DB.Create(&notif).Error; err != nil {
// 				log.Printf("[FW Webhook] Failed to save notification: %v", err)
// 			} else {
// 				log.Printf("[FW Webhook] Notification created for user=%d", history.UserID)
// 			}
// 		}
// 	} else {
// 		log.Printf("[FW Webhook] Could not find transaction history for reference=%s", reference)
// 	}
// }

// FlutterwaveWebhookHandler handles incoming FW webhooks
func FlutterwaveWebhookHandler(ctx *gin.Context) {
	secret := config.AppConfig.FlutterwaveHashKey
	sent := ctx.GetHeader("verif-hash")

	// Verify hash (fixed log placeholders)
	if secret == "" || subtle.ConstantTimeCompare([]byte(secret), []byte(sent)) != 1 {
		log.Printf("[FW Webhook] Invalid signature. expected=%s got=%s", secret, sent)
		utils.Error(ctx, http.StatusUnauthorized, "invalid signature")
		return
	}

	// Read raw body
	body, _ := io.ReadAll(ctx.Request.Body)
	log.Printf("[FW Webhook] Raw body: %s", string(body))

	// Parse
	var evt FlutterwaveWebhook
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("[FW Webhook] JSON unmarshal error: %v", err)
		utils.Error(ctx, http.StatusBadRequest, "bad payload")
		return
	}
	log.Printf("[FW Webhook] Parsed event: %+v", evt)

	// Normalize checks
	event := strings.ToUpper(evt.EventType)
	status := strings.ToLower(evt.Status)

	// Accept only successful money-in events (adjust as needed)
	if (event == "CHARGE.COMPLETED" || event == "BANK_TRANSFER_TRANSACTION") && status == "successful" {
		if err := handleSuccessfulPayment(evt); err != nil {
			// We return 200 so FW doesn't keep retrying; we log the failure for investigation.
			log.Printf("[FW Webhook] Processing error for ref=%s: %v", evt.TxRef, err)
		}
	} else {
		log.Printf("[FW Webhook] Ignored event=%s status=%s", evt.EventType, evt.Status)
	}

	// Always 200 OK to acknowledge receipt
	ctx.Status(http.StatusOK)
}
