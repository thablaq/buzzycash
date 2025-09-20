package payments

import (
	"log"
	"time"
	"errors"
	"fmt"
    "strings"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/external/gaming"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)


type PaymentService struct {
	db *gorm.DB
}


func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{
		db: db,
	}
}


func (p *PaymentService) HandleSuccessfulPaymentFW(evt FlutterwaveWebhook) error {
	return p.handleFWSuccessfulPayment(evt, "FW")
}


func (p *PaymentService) HandleSuccessfulPaymentNB(evt NombaWebhook) error {
	return p.handleNBSuccessfulPayment(evt, "NB")
}


func (p *PaymentService) handleFWSuccessfulPayment(evt FlutterwaveWebhook, provider string) error {
	reference := evt.TxRef
	amount := evt.Amount
	db := p.db

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

	// Create notification
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

func (p *PaymentService) handleNBSuccessfulPayment(evt NombaWebhook, provider string) error {
	reference := evt.Data.Order.OrderID
	event_type := evt.EventType
	amount := evt.Data.Order.Amount - evt.Data.Transaction.Fee
	db := p.db 

	var history models.Transaction
	log.Printf("[%s Webhook] Processing payment - Reference: %s, Amount: %v (type: %T)", provider, reference, amount, amount)
	
	
	if strings.ToLower(event_type) != "payment_success" {
			log.Printf("[%s Webhook] deposit status not successful (status=%s). Skipping update.", provider, event_type)
			return nil
		}

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

	// Create notification
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
			log.Printf("[%s Webhook] WARNING: could not create notification for ref=%s: %v", provider, reference, err)
		} else {
			log.Printf("[%s Webhook] SUCCESS ref=%s | wallet credited & notification created", provider, reference)
		}
	}

	return nil
}




func (p *PaymentService) handleNBSuccessfulWithdrawal(evt NombaWithdrawalResponse) error {
	reference := evt.Data.Meta.MerchantTxRef
	amount := evt.Data.Amount
	status := evt.Data.Status
	db := p.db

	var history models.Transaction
	log.Printf("[%s Webhook] Processing withdrawal - Reference: %s, Amount: %v, Status: %s", 
		reference, amount, status)
	
	if strings.ToUpper(status) != "SUCCESS" {
			log.Printf("[%s Webhook] Withdrawal status not successful (status=%s). Skipping update.",status)
			return nil
		}

	// Update history atomically
	if err := db.Transaction(func(tx *gorm.DB) error {
		// 1) Lock + load the history row
		if err := tx.Preload("User").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("reference = ?", reference).
			First(&history).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("[%s Webhook] No TransactionHistory for reference=%s", reference)
				return nil
			}
			return fmt.Errorf("load history failed: %w", err)
		}

		// 2) Idempotency check
		if history.PaymentStatus == models.Successful {
			log.Printf("[%s Webhook] Reference=%s already processed; skipping",reference)
			return nil
		}

		// 3) Update status
		if err := tx.Model(&history).Updates(map[string]interface{}{
			"payment_status": models.Successful,
			"paid_at":        time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("update history failed: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	// Create notification outside the transaction
	if history.ID != "" {
		amountInt := int64(amount)
		log.Printf("[%s Webhook] Creating withdrawal notification - UserID: %s, Amount: %d", history.UserID, amountInt)
		notif := models.Notification{
			UserID:   history.UserID,
			Type:     models.Transactions,
			Title:    "Withdrawal Successful",
			Subtitle: "Your withdrawal was processed successfully.",
			Amount:   amountInt,
			Currency: string(history.Currency),
			Status:   "successful",
		}

		if err := db.Create(&notif).Error; err != nil {
			log.Printf("[%s Webhook] WARNING: could not create notification for ref=%s: %v", reference, err)
		} else {
			log.Printf("[%s Webhook] SUCCESS ref=%s | withdrawal processed & notification created",reference)
		}
	}

	return nil
}
