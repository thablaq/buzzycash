package payments

type FlutterwaveWebhook struct {
	ID            int     `json:"id"`
	TxRef         string  `json:"txRef"`
	FlwRef        string  `json:"flwRef"`
	OrderRef      string  `json:"orderRef"`
	Amount        float64 `json:"amount"`
	ChargedAmount float64 `json:"charged_amount"`
	Status        string  `json:"status"`
	Currency      string  `json:"currency"`
	EventType     string  `json:"event.type"`
	Customer      struct {
		ID    int    `json:"id"`
		Phone string `json:"phone"`
		Email string `json:"email"`
		Name  string `json:"fullName"`
	} `json:"customer"`
}





type NombaWebhook struct {
	EventType string `json:"event_type"`
	RequestID string `json:"requestId"`
	Data      struct {
		Customer struct {
			BillerID  string `json:"billerId"`
			ProductID string `json:"productId"`
		} `json:"customer"`

		Merchant struct {
			UserID        string  `json:"userId"`
			WalletBalance float64 `json:"walletBalance"`
			WalletID      string  `json:"walletId"`
		} `json:"merchant"`

		Order struct {
			AccountID              string  `json:"accountId"`
			Amount                 float64 `json:"amount"`
			CallbackURL            string  `json:"callbackUrl"`
			CardCurrency           string  `json:"cardCurrency"`
			CardLast4Digits        string  `json:"cardLast4Digits"`
			CardType               string  `json:"cardType"`
			Currency               string  `json:"currency"`
			CustomerEmail          string  `json:"customerEmail"`
			CustomerID             string  `json:"customerId"`
			IsTokenizedCardPayment string  `json:"isTokenizedCardPayment"`
			OrderID                string  `json:"orderId"`
			OrderReference         string  `json:"orderReference"`
			PaymentMethod          string  `json:"paymentMethod"`
		} `json:"order"`

		Terminal struct{} `json:"terminal"`

		TokenizedCardData struct {
			CardPan         string `json:"cardPan"`
			CardType        string `json:"cardType"`
			TokenExpiryMonth string `json:"tokenExpiryMonth"`
			TokenExpiryYear  string `json:"tokenExpiryYear"`
			TokenKey        string `json:"tokenKey"`
		} `json:"tokenizedCardData"`

		Transaction struct {
			Fee              float64 `json:"fee"`
			MerchantTxRef    string  `json:"merchantTxRef"`
			OriginatingFrom  string  `json:"originatingFrom"`
			ResponseCode     string  `json:"responseCode"`
			Time             string  `json:"time"`
			TransactionAmount float64 `json:"transactionAmount"`
			TransactionID    string  `json:"transactionId"`
			Type             string  `json:"type"`
		} `json:"transaction"`
	} `json:"data"`
}
