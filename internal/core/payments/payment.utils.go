// package payments

// import (
// 	"log"
// 	"time"

// 	"github.com/dblaq/buzzycash/internal/config"
// 	"github.com/dblaq/buzzycash/internal/models"
// 	// "github.com/dblaq/buzzycash/internal/notifications"
// 	"errors"
// 	"fmt"
// 	"github.com/dblaq/buzzycash/external/gaming"
// 	"gorm.io/gorm"
// 	"gorm.io/gorm/clause"
// )

// func handleSuccessfulPaymentFW(evt FlutterwaveWebhook) error {
// 	reference := evt.TxRef
// 	amount := evt.Amount

// 	db := config.DB
// 	var history models.Transaction
// 	log.Printf("[FW Webhook] Processing payment - Reference: %s, Amount: %v (type: %T)", reference, amount, amount)

// 	// First: update history + credit wallet atomically
// 	if err := db.Transaction(func(tx *gorm.DB) error {
// 		// 1) Lock + load the history row by reference (FOR UPDATE) and preload user
// 		if err := tx.Preload("User").
// 			Clauses(clause.Locking{Strength: "UPDATE"}).
// 			Where("reference = ?", reference).
// 			First(&history).Error; err != nil {

// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				log.Printf("[FW Webhook] No TransactionHistory for reference=%s", reference)
// 				return nil
// 			}
// 			return fmt.Errorf("load history failed: %w", err)
// 		}

// 		// 2) Idempotency check
// 		if history.PaymentStatus == models.Successful {
// 			log.Printf("[FW Webhook] Reference=%s already processed; skipping", reference)
// 			return nil
// 		}

// 		// 3) Update status
// 		if err := tx.Model(&history).Updates(map[string]interface{}{
// 			"payment_status": models.Successful,
// 			"paid_at":        time.Now(),
// 		}).Error; err != nil {
// 			return fmt.Errorf("update history failed: %w", err)
// 		}

// 		// 4) Credit wallet
// 		gs := gaming.GMInstance()
// 		if _, err := gs.CreditUserWallet(history.User.PhoneNumber, amount); err != nil {
// 			return fmt.Errorf("wallet credit failed: %w", err)
// 		}

// 		return nil
// 	}); err != nil {
// 		return err
// 	}

// 	// ✅ Outside transaction: Create notification
// 	if history.ID != "" {
// 		// title, subtitle := notifications.BuildTxNotifContent(history)
// 		amountInt := int64(amount)
// 		log.Printf("[FW Webhook] Creating notification - UserID: %s, Amount: %d", history.UserID, amountInt)
// 		notif := models.Notification{
// 			UserID:   history.UserID,
// 			Type:     models.Transactions,
// 			Title:    "Deposit Successful",
// 			Subtitle: "You have successfully deposited into your wallet.",
// 			Amount:   amountInt,
// 			Currency: string(history.Currency),
// 			Status:   "successful",
// 		}

// 		if err := config.DB.Create(&notif).Error; err != nil {
// 			// don’t rollback payment, just log
// 			log.Printf("[FW Webhook] WARNING: could not create notification for ref=%s: %v", reference, err, err.Error())
// 		} else {
// 			log.Printf("[FW Webhook] SUCCESS ref=%s | wallet credited & notification created", reference)
// 		}
// 	}

// 	return nil
// }



// func handleSuccessfulPaymentNB(evt FlutterwaveWebhook) error {
// 	reference := evt.TxRef
// 	amount := evt.Amount

// 	db := config.DB
// 	var history models.Transaction
// 	log.Printf("[FW Webhook] Processing payment - Reference: %s, Amount: %v (type: %T)", reference, amount, amount)

// 	// First: update history + credit wallet atomically
// 	if err := db.Transaction(func(tx *gorm.DB) error {
// 		// 1) Lock + load the history row by reference (FOR UPDATE) and preload user
// 		if err := tx.Preload("User").
// 			Clauses(clause.Locking{Strength: "UPDATE"}).
// 			Where("reference = ?", reference).
// 			First(&history).Error; err != nil {

// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				log.Printf("[FW Webhook] No TransactionHistory for reference=%s", reference)
// 				return nil
// 			}
// 			return fmt.Errorf("load history failed: %w", err)
// 		}

// 		// 2) Idempotency check
// 		if history.PaymentStatus == models.Successful {
// 			log.Printf("[FW Webhook] Reference=%s already processed; skipping", reference)
// 			return nil
// 		}

// 		// 3) Update status
// 		if err := tx.Model(&history).Updates(map[string]interface{}{
// 			"payment_status": models.Successful,
// 			"paid_at":        time.Now(),
// 		}).Error; err != nil {
// 			return fmt.Errorf("update history failed: %w", err)
// 		}

// 		// 4) Credit wallet
// 		gs := gaming.GMInstance()
// 		if _, err := gs.CreditUserWallet(history.User.PhoneNumber, amount); err != nil {
// 			return fmt.Errorf("wallet credit failed: %w", err)
// 		}

// 		return nil
// 	}); err != nil {
// 		return err
// 	}

// 	// ✅ Outside transaction: Create notification
// 	if history.ID != "" {
// 		// title, subtitle := notifications.BuildTxNotifContent(history)
// 		amountInt := int64(amount)
// 		log.Printf("[FW Webhook] Creating notification - UserID: %s, Amount: %d", history.UserID, amountInt)
// 		notif := models.Notification{
// 			UserID:   history.UserID,
// 			Type:     models.Transactions,
// 			Title:    "Deposit Successful",
// 			Subtitle: "You have successfully deposited into your wallet.",
// 			Amount:   amountInt,
// 			Currency: string(history.Currency),
// 			Status:   "successful",
// 		}

// 		if err := config.DB.Create(&notif).Error; err != nil {
// 			// don’t rollback payment, just log
// 			log.Printf("[FW Webhook] WARNING: could not create notification for ref=%s: %v", reference, err, err.Error())
// 		} else {
// 			log.Printf("[FW Webhook] SUCCESS ref=%s | wallet credited & notification created", reference)
// 		}
// 	}

// 	return nil
// }
// 


package payments

import (
	"log"
	"time"
	"errors"
	"fmt"

	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/external/gaming"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PaymentService handles payment-related operations
type PaymentService struct {
	db *gorm.DB
}

// NewPaymentService creates a new PaymentService with dependencies
func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{
		db: db,
	}
}

// HandleSuccessfulPaymentFW handles successful Flutterwave payments
func (p *PaymentService) HandleSuccessfulPaymentFW(evt FlutterwaveWebhook) error {
	return p.handleSuccessfulPayment(evt, "FW")
}

// // HandleSuccessfulPaymentNB handles successful payment for another provider
// func (p *PaymentService) HandleSuccessfulPaymentNB(evt FlutterwaveWebhook) error {
// 	return p.handleSuccessfulPayment(evt, "NB")
// }

// handleSuccessfulPayment is the common implementation
func (p *PaymentService) handleSuccessfulPayment(evt FlutterwaveWebhook, provider string) error {
	reference := evt.TxRef
	amount := evt.Amount
	db := p.db // Create local variable for shorter syntax

	var history models.Transaction
	log.Printf("[%s Webhook] Processing payment - Reference: %s, Amount: %v (type: %T)", provider, reference, amount, amount)

	// First: update history + credit wallet atomically
	if err := db.Transaction(func(tx *gorm.DB) error {
		// 1) Lock + load the history row by reference (FOR UPDATE) and preload user
		if err := tx.Preload("User").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("reference = ?", reference).
			First(&history).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("[%s Webhook] No TransactionHistory for reference=%s", provider, reference)
				return nil
			}
			return fmt.Errorf("load history failed: %w", err)
		}

		// 2) Idempotency check
		if history.PaymentStatus == models.Successful {
			log.Printf("[%s Webhook] Reference=%s already processed; skipping", provider, reference)
			return nil
		}

		// 3) Update status
		if err := tx.Model(&history).Updates(map[string]interface{}{
			"payment_status": models.Successful,
			"paid_at":        time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("update history failed: %w", err)
		}

		// 4) Credit wallet
		gs := gaming.GMInstance()
		if _, err := gs.CreditUserWallet(history.User.PhoneNumber, amount); err != nil {
			return fmt.Errorf("wallet credit failed: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	// ✅ Outside transaction: Create notification
	if history.ID != "" {
		amountInt := int64(amount)
		log.Printf("[%s Webhook] Creating notification - UserID: %s, Amount: %d", provider, history.UserID, amountInt)
		notif := models.Notification{
			UserID:   history.UserID,
			Type:     models.Transactions,
			Title:    "Deposit Successful",
			Subtitle: "You have successfully deposited into your wallet.",
			Amount:   amountInt,
			Currency: string(history.Currency),
			Status:   "successful",
		}

		if err := db.Create(&notif).Error; err != nil {
			// don't rollback payment, just log
			log.Printf("[%s Webhook] WARNING: could not create notification for ref=%s: %v", provider, reference, err)
		} else {
			log.Printf("[%s Webhook] SUCCESS ref=%s | wallet credited & notification created", provider, reference)
		}
	}

	return nil
}

// If you have webhook handlers that need to be Gin handlers
// type WebhookHandler struct {
// 	paymentService *PaymentService
// }

// func NewWebhookHandler(db *gorm.DB) *WebhookHandler {
// 	return &WebhookHandler{
// 		paymentService: NewPaymentService(db),
// 	}
// }

// func (w *WebhookHandler) FlutterwaveWebhookHandler(ctx *gin.Context) {
// 	var evt FlutterwaveWebhook
// 	if err := ctx.ShouldBindJSON(&evt); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook data"})
// 		return
// 	}

// 	// Verify webhook signature here if needed
	
// 	if err := w.paymentService.HandleSuccessfulPaymentFW(evt); err != nil {
// 		log.Printf("Webhook processing failed: %v", err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Processing failed"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
// }

// // Alternative: If you prefer to keep the original function signatures 
// // but still use dependency injection, you can create wrapper functions
// func HandleSuccessfulPaymentFW(db *gorm.DB, evt FlutterwaveWebhook) error {
// 	service := NewPaymentService(db)
// 	return service.HandleSuccessfulPaymentFW(evt)
// }

// func HandleSuccessfulPaymentNB(db *gorm.DB, evt FlutterwaveWebhook) error {
// 	service := NewPaymentService(db)
// 	return service.HandleSuccessfulPaymentNB(evt)
// }