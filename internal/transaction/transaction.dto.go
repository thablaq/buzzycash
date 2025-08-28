package transaction

import "time"

type TransactionHistoryResponse struct {
	ID                   string    `json:"id"`
	TicketPurchaseID     string    `json:"ticket_purchase_id,omitempty"`
	Amount               float64     `json:"amount"`
	TransactionReference string    `json:"transaction_reference"`
	Reference            string    `json:"reference"`
	CustomerEmail        string    `json:"customer_email,omitempty"`
	PaymentStatus        string    `json:"payment_status"`
	PaymentType          string    `json:"payment_type"`
	Currency             string    `json:"currency"`
	PaidAt               time.Time `json:"paid_at"`
	TransactionType      string    `json:"transaction_type"`
	PaymentMethod        string    `json:"payment_method"`
	Category             string    `json:"category"`
}