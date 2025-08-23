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







func FlutterwaveWebhookHandler(ctx *gin.Context) {
    secret := config.AppConfig.FlutterwaveHashKey
    sent := ctx.GetHeader("verif-hash")

    // Verify hash
    if secret == "" || subtle.ConstantTimeCompare([]byte(secret), []byte(sent)) != 1 {
        log.Printf("[FW Webhook] Invalid signature. expected=%s got=%s",sent)
        utils.Error(ctx,http.StatusUnauthorized,"invalid signature")
        return
    }

    // Read raw body
    body, _ := io.ReadAll(ctx.Request.Body)
    log.Printf("[FW Webhook] Raw body: %s", string(body))

 
    var evt FlutterwaveWebhook
    if err := json.Unmarshal(body, &evt); err != nil {
        log.Printf("[FW Webhook] JSON unmarshal error: %v", err)
        utils.Error(ctx,http.StatusBadRequest, "bad payload")
        return
    }

    log.Printf("[FW Webhook] Parsed event: %+v", evt)

    
    if evt.EventType == "charge.completed" && evt.Status == "successful" {
        handleSuccessfulPayment(evt)
    } else if evt.EventType == "BANK_TRANSFER_TRANSACTION" && evt.Status == "successful" {
        handleSuccessfulPayment(evt)
    } else {
        log.Printf("[FW Webhook] Ignored event=%s status=%s", evt.EventType, evt.Status)
    }

    ctx.Status(http.StatusOK)
}

func handleSuccessfulPayment(evt FlutterwaveWebhook) {
    txRef := evt.TxRef
    amount := evt.Amount
    currency := evt.Currency

    log.Printf("[FW Webhook] SUCCESSFUL payment. tx_ref=%s id=%d amount=%.2f currency=%s",
        txRef, evt.ID, amount, currency)

    // Update DB transaction
    if err := config.DB.Model(&models.TransactionHistory{}).
        Where("transaction_reference = ?", txRef).
        Updates(map[string]interface{}{
            "payment_status": models.Successful,
            "paid_at":        time.Now(),
            "amount":         amount,
        }).Error; err != nil {
        log.Printf("[FW Webhook] DB update failed: %+v\n", err)
        return
    }

    // Credit wallet
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
}



