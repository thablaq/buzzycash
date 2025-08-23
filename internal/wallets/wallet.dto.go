package wallets





type CreditWalletRequest struct {
    Amount float64 `json:"amount" validate:"gt=0"`
}

