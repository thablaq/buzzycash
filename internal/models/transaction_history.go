package models

import (
	"time"
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
)

const (
	Nomba  EPaymentMethod = "NOMBA"
	WalletTX EPaymentMethod = "WALLET"
)

const (
	PrizeMoney        TransactionCategory = "PRIZE_MONEY"
	Purchase    TransactionCategory = "TICKET_PURCHASE"
	Deposit           TransactionCategory = "DEPOSIT"
	CashoutTX           TransactionCategory = "CASHOUT"
	WithdrawReversed  TransactionCategory = "WITHDRAW_REVERSED"
)

type TransactionHistory struct {
	ID                   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID               string    `gorm:"type:uuid"`
	TicketPurchaseID     string    `gorm:"type:uuid;index"`
	WithdrawalsID        string   `gorm:"type:uuid;uniqueIndex"`
	
	DebitAmount          float64 `gorm:"type:decimal(10,2);default:0.0"`
	AmountPaid           float64 `gorm:"type:decimal(10,2);default:0.0"`
	PaymentID            string  `gorm:"size:255;uniqueIndex"`
	TransactionReference string  `gorm:"size:255;uniqueIndex"`
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
	Withdrawals     []WithdrawalRequest `gorm:"foreignKey:TransactionHistoryID"`
}

type JSONB map[string]interface{}

func (j JSONB) GormDataType() string {
	return "jsonb"
}