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

	return db.Transaction(func(tx *gorm.DB) error {
		// 1) Lock + load the history row by reference (FOR UPDATE) and preload user
		var history models.TransactionHistory
		if err := tx.Preload("User").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("reference = ?", reference).
			First(&history).Error; err != nil {

			// Not found -> nothing to do (ack webhook)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("[FW Webhook] No TransactionHistory for reference=%s", reference)
				return nil
			}
			return fmt.Errorf("load history failed: %w", err)
		}

		// 2) Idempotency: if already successful, skip (prevents double-credit & double-notif)
		if history.PaymentStatus == models.Successful {
			log.Printf("[FW Webhook] Reference=%s already processed; skipping", reference)
			return nil
		}

		// 3) Update status + paid_at (keep amount if you already stored it at init; otherwise set here)
		if err := tx.Model(&history).Updates(map[string]interface{}{
			"payment_status": models.Successful,
			"paid_at":        time.Now(),
			// If your column is INT (kobo/naira as int), cast appropriately:
			// "amount": int64(math.Round(amount)), // only if you *want* to overwrite
		}).Error; err != nil {
			return fmt.Errorf("update history failed: %w", err)
		}

		// 4) Credit wallet (external op). If it fails, rollback the whole transaction.
		gs := gaming.GMInstance()
		if _, err := gs.CreditUserWallet(history.User.PhoneNumber, amount); err != nil {
			return fmt.Errorf("wallet credit failed: %w", err)
		}

		// 5) Create notification ONLY NOW (after successful credit)
		title, subtitle := notifications.BuildTxNotifContent(history)
		notif := models.Notification{
			UserID:   history.UserID, // ensure type matches your models (UUID string vs uint)
			Type:     "transaction",
			Title:    title,
			Subtitle: subtitle,
			Amount:   amount,          // align with your Notification.Amount type
			Currency: string(history.Currency),
			Status:   "successful",
			// CreatedAt is auto if you have gorm.Model; otherwise set time.Now()
		}
		if err := tx.Create(&notif).Error; err != nil {
			return fmt.Errorf("create notification failed: %w", err)
		}

		log.Printf("[FW Webhook] SUCCESS ref=%s | wallet credited & notification created", reference)
		return nil
	})
}


