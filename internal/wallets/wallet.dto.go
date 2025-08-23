package wallets





type CreditWalletRequest struct {
    Amount float64 `json:"amount" validate:"gt=0"`
}


// type flutterwaveWebhook struct {
// 	Event string `json:"event"`
// 	Data  struct {
// 		ID       int64  `json:"id"`
// 		TxRef    string `json:"tx_ref"`
// 		Status   string `json:"status"` // "successful"
// 		Amount   float64 `json:"amount"`
// 		Currency string  `json:"currency"`
// 		FlwRef   string  `json:"flw_ref"`
// 	} `json:"data"`
// }