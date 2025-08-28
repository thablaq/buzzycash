package wallets





type CreditWalletRequest struct {
	Amount        int64  `json:"amount" validate:"required,gt=0"`
	PaymentMethod string `json:"payment_method" validate:"required,oneof=flutterwave nomba"`
}

