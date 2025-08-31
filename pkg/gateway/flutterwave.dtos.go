package gateway




type FWCustomer struct {
	Email string `json:"email"`
	FullName string `json:"full_name,omitempty"`
}


type FWPaymentRequest struct {
	Reference       string         `json:"tx_ref"`
	Amount      int64         `json:"amount"`   
	Currency    string         `json:"currency"` 
	RedirectURL string         `json:"redirect_url"`
	Customer    FWCustomer     `json:"customer"`

	// Optional
	PaymentOptions string          `json:"payment_options,omitempty"` // "card,banktransfer,ussd"
	Meta           map[string]any  `json:"meta,omitempty"`
}

type fwCreateResp struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Data struct {
		Link string `json:"link"`
	} `json:"data"`
}



type FWVerifyResp struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Data struct {
		ID            int64   `json:"id"`
		TxRef         string  `json:"tx_ref"`
		FlwRef        string  `json:"flw_ref"`
		Amount        float64 `json:"amount"`
		Currency      string  `json:"currency"`
		ChargedAmount float64 `json:"charged_amount"`
		AppFee        float64 `json:"app_fee"`
		MerchantFee   float64 `json:"merchant_fee"`
		ProcessorResp string  `json:"processor_response"`
		AuthModel     string  `json:"auth_model"`
		Status        string  `json:"status"`
		PaymentType   string  `json:"payment_type"`
		CreatedAt     string  `json:"created_at"`
		AccountID     int64   `json:"account_id"`
		Customer      struct {
			ID    int64  `json:"id"`
			Name  string `json:"name"`
			Phone string `json:"phone_number"`
			Email string `json:"email"`
		} `json:"customer"`
	} `json:"data"`
}



