package models

import (
	"time"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type EPaymentStatus string
type EPaymentType string
type ECurrency string
type ETransactionType string
type EPaymentMethod string
type TransactionCategory string

const (
	Pending    EPaymentStatus = "PENDING"
	Successful EPaymentStatus = "SUCCESSFUL"
	Reversed   EPaymentStatus = "REVERSED"
	Failed     EPaymentStatus = "FAILED"
	Rejected   EPaymentStatus = "REJECTED"
)

const (
	Topup  EPaymentType = "TOPUP"
	Bonus  EPaymentType = "BONUS"
	Profit EPaymentType = "PROFIT"
	Payout EPaymentType = "PAYOUT"
)

const (
	NGN ECurrency = "NGN"
	CED ECurrency = "CED"
)

const (
	Credit ETransactionType = "CREDIT"
	Debit  ETransactionType = "DEBIT"
	Withdrawal ETransactionType = "WITHDRAWAL"
)

const (
	Nomba  EPaymentMethod = "NOMBA"
	Wallet EPaymentMethod = "WALLET"
	Flutterwave EPaymentMethod = "FLUTTERWAVE"
)

const (
	PrizeMoney        TransactionCategory = "PRIZE_MONEY"
	Ticket    TransactionCategory = "TICKET"
	Deposit           TransactionCategory = "DEPOSIT"
	Cashout           TransactionCategory = "CASHOUT"
	WithdrawRequest  TransactionCategory = "WITHDRAWAL"
	WithdrawReversed  TransactionCategory = "WITHDRAW_REVERSED"
)

type TransactionHistory struct {
	ID                   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID               string    `gorm:"type:uuid"`
	TicketPurchaseID string `gorm:"type:uuid;default:null"`
	
	Amount               int64 
	TransactionReference string  `gorm:"size:255;uniqueIndex"`
	Reference             string  `gorm:"size:255;uniqueIndex"`
	Metadata             JSONB
	CustomerEmail        string `gorm:"size:255"`
	PaymentStatus        EPaymentStatus
	PaymentType          EPaymentType
	Currency             ECurrency `gorm:"default:NGN"`
	PaidAt               time.Time
	DeletedAt            *time.Time
	TransactionType      ETransactionType
	CreatedAt            time.Time `gorm:"default:current_timestamp"`
	UpdatedAt            time.Time
	PaymentMethod        EPaymentMethod
	Category             TransactionCategory
	
	User            User               `gorm:"constraint:OnDelete:CASCADE;"`
	TicketPurchase  []TicketPurchase  `gorm:"foreignKey:TransactionHistoryID"`
	GameHistories  []GameHistory     `gorm:"foreignKey:TransactionHistoryID"`
}

type JSONB map[string]interface{}

func (j JSONB) GormDataType() string {
	return "jsonb"
}

// For saving into DB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

// For reading from DB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}

	bytes, ok := value.([]uint8)
	if !ok {
		return fmt.Errorf("failed to scan JSONB: %v", value)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(bytes, &m); err != nil {
		return err
	}
	*j = m
	return nil
}