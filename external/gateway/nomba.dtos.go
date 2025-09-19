package gateway

import (
	"sync"

)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string  `json:"expiresAt"`
	BusinessID   string `json:"businessId"`
}



type NombaAuthService struct {
	mu           sync.RWMutex
	token        *TokenResponse
}



type NBCustomer struct {
	Email    string `json:"email"`
	FullName string `json:"full_name,omitempty"`
}

type NBOrder struct {
	CallbackURL   string `json:"callbackUrl"`
	CustomerEmail string `json:"customerEmail"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	CustomerID    string `json:"customerId"`
}
type NBPaymentRequest struct {
	Order        NBOrder     `json:"order"`
	Metadata     interface{} `json:"metadata,omitempty"`
	TokenizeCard bool        `json:"tokenizeCard"`
}

type NBCheckoutResponse struct {
	Status bool `json:"status"`
	Data   struct {
		CheckoutLink   string `json:"checkoutLink"`
		OrderReference string `json:"orderReference"`
	} `json:"data"`
}

type NBBankResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    []Bank `json:"data"`
}

type Bank struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Logo string   `json:"logo"`
}

type NBRetrieveAccountDetails struct {
	AccountNumber string `json:"accountNumber"`
	BankCode      string `json:"bankCode"`
}

type NBAccountDetails struct {
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
}

type NBAccountDetailsResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AccountNumber string `json:"accountNumber"`
		AccountName   string `json:"accountName"`
	} `json:"data"`
}

type NBWithdrawalResp struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
	} `json:"data"`
}

type NBWithdrawalRequest struct {
	Amount        int64          `json:"amount"`
	BankCode      string         `json:"bankCode"`
	AccountNumber string         `json:"accountNumber"`
	AccountName   string         `json:"accountName"`
	SenderName    string         `json:"senderName,omitempty"`
	MerchantTxRef string         `json:"merchantTxRef"`
	Narration     string         `json:"narration,omitempty"`
	Meta          map[string]any `json:"meta,omitempty"`
}
