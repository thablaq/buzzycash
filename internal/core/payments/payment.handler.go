package payments

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
"gorm.io/gorm"
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/utils"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	paymentService *PaymentService
}

func NewWebhookHandler(db *gorm.DB) *WebhookHandler {
	return &WebhookHandler{
		paymentService: NewPaymentService(db),
	}
}

// FlutterwaveWebhookHandler handles incoming FW webhooks
func (w *WebhookHandler)FlutterwaveWebhookHandler(ctx *gin.Context) {
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
		if err := w.paymentService.HandleSuccessfulPaymentFW(evt); err != nil {
			// We return 200 so FW doesn't keep retrying; we log the failure for investigation.
			log.Printf("[FW Webhook] Processing error for ref=%s: %v", evt.TxRef, err)
		}
	} else {
		log.Printf("[FW Webhook] Ignored event=%s status=%s", evt.EventType, evt.Status)
	}

	// Always 200 OK to acknowledge receipt
	ctx.Status(http.StatusOK)
}
