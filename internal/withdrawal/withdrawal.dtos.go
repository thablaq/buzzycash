package withdrawal




type RetrieveAccountDetailsRequest struct {
	BankCode    string `json:"bank_code" validate:"required"`
	AccountNumber string `json:"account_number" binding:"required" validate:"required,len=10,numeric"`
}


type InitiateWithdrawalRequest struct {
    Amount        int64  `json:"amount" binding:"required" validate:"required"`
    AccountName   string `json:"account_name" binding:"required" validate:"required"`
    BankCode      string `json:"bank_code" binding:"required" validate:"required"`
    AccountNumber string `json:"account_number" binding:"required" validate:"required,len=10,numeric"`
    Currency      string `json:"currency" binding:"required" validate:"required,len=3"`
}