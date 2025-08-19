package wallets





type creditWalletRequest struct {
    Amount float64 `json:"amount" validate:"gt=0"`
}
