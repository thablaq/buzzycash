package wallets



	import (
		"errors"
)




// Validation errors
var (
	ErrAmountTooShort          = errors.New("topup amount must be at least 100 naira")

)


func (r *CreditWalletRequest) Validate() error {
	if err := validateAmount(r.Amount); err != nil {
		return err
	}
	return nil
}

func validateAmount(amount int64) error {
	if amount < 100 {
		return ErrAmountTooShort
	}
	return nil
}


