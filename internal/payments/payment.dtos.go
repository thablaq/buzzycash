package payments




type FlutterwaveWebhook struct {
    ID            int     `json:"id"`
   	TxRef  string `json:"tx_ref"`
    FlwRef string `json:"flw_ref"`
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
