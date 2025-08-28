package payments

import (
	"log"
	"time"

	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/notifications"
	"github.com/dblaq/buzzycash/pkg/gaming"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"errors"
	"fmt"
)



func handleSuccessfulPayment(evt FlutterwaveWebhook) error {
	reference := evt.TxRef
	amount := evt.Amount

	db := config.DB
	var history models.TransactionHistory

	// First: update history + credit wallet atomically
	if err := db.Transaction(func(tx *gorm.DB) error {
		// 1) Lock + load the history row by reference (FOR UPDATE) and preload user
		if err := tx.Preload("User").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("reference = ?", reference).
			First(&history).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("[FW Webhook] No TransactionHistory for reference=%s", reference)
				return nil
			}
			return fmt.Errorf("load history failed: %w", err)
		}

		// 2) Idempotency check
		if history.PaymentStatus == models.Successful {
			log.Printf("[FW Webhook] Reference=%s already processed; skipping", reference)
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
		title, subtitle := notifications.BuildTxNotifContent(history)
		notif := models.Notification{
			UserID:   history.UserID,
			Type:     "transaction",
			Title:    title,
			Subtitle: subtitle,
			Amount:   amount,
			Currency: string(history.Currency),
			Status:   "successful",
		}

		if err := db.Create(&notif).Error; err != nil {
			// don’t rollback payment, just log
			log.Printf("[FW Webhook] WARNING: could not create notification for ref=%s: %v", reference, err)
		} else {
			log.Printf("[FW Webhook] SUCCESS ref=%s | wallet credited & notification created", reference)
		}
	}

	return nil
}


