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

type NombaWebhookWrapper struct {
	EventType string `json:"eventType,omitempty"` // present in payment webhooks
	Code      string `json:"code,omitempty"`      // present in withdrawal response
	Status    string `json:"status,omitempty"`    // for withdrawal
}


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

	
	if (event == "CHARGE.COMPLETED" || event == "BANK_TRANSFER_TRANSACTION") && status == "successful" {
		if err := w.paymentService.HandleSuccessfulPaymentFW(evt); err != nil {
			log.Printf("[FW Webhook] Processing error for ref=%s: %v", evt.TxRef, err)
		}
	} else {
		log.Printf("[FW Webhook] Ignored event=%s status=%s", evt.EventType, evt.Status)
	}
	ctx.Status(http.StatusOK)
}


// func (w *WebhookHandler)NombaWebhookHandler(ctx *gin.Context) {
// 	body, _ := io.ReadAll(ctx.Request.Body)
// 	log.Printf("[NB Webhook] Raw body: %s", string(body))

// 	// Parse
// 	var evt NombaWebhook
// 	if err := json.Unmarshal(body, &evt); err != nil {
// 		log.Printf("[NB Webhook] JSON unmarshal error: %v", err)
// 		utils.Error(ctx, http.StatusBadRequest, "bad payload")
// 		return
// 	}
// 	log.Printf("[NB Webhook] Parsed event: %+v", evt)

// 	// Normalize checks
// 	event := strings.ToLower(evt.EventType)

// 	if (event == "payment_success"){
// 		if err := w.paymentService.HandleSuccessfulPaymentNB(evt); err != nil {
// 			log.Printf("[NB Webhook] Processing error for ref=%s: %v", evt.Data.Order.OrderID, err)
// 		}
// 	} else {
// 		log.Printf("[NB Webhook] Ignored event=%s status=%s", evt.EventType)
// 	}
// 	ctx.Status(http.StatusOK)
// }

// func (w *WebhookHandler)NombaWebhookWithdrawHandler(ctx *gin.Context) {
// 	body, _ := io.ReadAll(ctx.Request.Body)
// 	log.Printf("[NB Webhook] Raw body: %s", string(body))

// 	// Parse
// 	var evt NombaWithdrawalResponse
// 	if err := json.Unmarshal(body, &evt); err != nil {
// 		log.Printf("[NB Webhook] JSON unmarshal error: %v", err)
// 		utils.Error(ctx, http.StatusBadRequest, "bad payload")
// 		return
// 	}
// 	log.Printf("[NB Webhook] Parsed event: %+v", evt)

// 	// Normalize checks
// 	status := strings.ToUpper(evt.Data.Status)

// 	if (status == "SUCCESS"){
// 		if err := w.paymentService.handleNBSuccessfulWithdrawal(evt); err != nil {
// 			log.Printf("[NB Webhook] Processing error for ref=%s: %v", evt.Data.Status, err)
// 		}
// 	} else {
// 		log.Printf("[NB Webhook] Ignored event=%s status=%s", evt.Data.Status)
// 	}
// 	ctx.Status(http.StatusOK)
// }



func (w *WebhookHandler) NombaWebhookHandler(ctx *gin.Context) {
	body, _ := io.ReadAll(ctx.Request.Body)
	log.Printf("[NB Webhook] Raw body: %s", string(body))

	// Step 1: Peek into the payload
	var wrapper NombaWebhookWrapper
	if err := json.Unmarshal(body, &wrapper); err != nil {
		log.Printf("[NB Webhook] JSON peek unmarshal error: %v", err)
		utils.Error(ctx, http.StatusBadRequest, "bad payload")
		return
	}

	// Step 2: Branch logic
	switch {
	// ----------- Deposit flow ------------
	case strings.ToLower(wrapper.EventType) == "payment_success":
		var evt NombaWebhook
		if err := json.Unmarshal(body, &evt); err != nil {
			log.Printf("[NB Webhook] Deposit unmarshal error: %v", err)
			utils.Error(ctx, http.StatusBadRequest, "bad deposit payload")
			return
		}

		if err := w.paymentService.handleNBSuccessfulPayment(evt, "nomba"); err != nil {
			log.Printf("[NB Webhook] Deposit processing error for ref=%s: %v", evt.Data.Order.OrderID, err)
		}

	// ----------- Withdrawal flow ----------
	case strings.ToUpper(wrapper.Status) == "SUCCESS":
		var evt NombaWithdrawalResponse
		if err := json.Unmarshal(body, &evt); err != nil {
			log.Printf("[NB Webhook] Withdrawal unmarshal error: %v", err)
			utils.Error(ctx, http.StatusBadRequest, "bad withdrawal payload")
			return
		}

		if err := w.paymentService.handleNBSuccessfulWithdrawal(evt); err != nil {
			log.Printf("[NB Webhook] Withdrawal processing error for ref=%s: %v", evt.Data.Meta.MerchantTxRef, err)
		}

	default:
		log.Printf("[NB Webhook] Ignored eventType=%s status=%s", wrapper.EventType, wrapper.Status)
	}

	ctx.Status(http.StatusOK)
}
