package withdrawal



import (
	"errors"
	"regexp"
	"strings"
)

// Validation errors
var (
	ErrAmountTooShort          = errors.New("amount must be at least 1000 naira")
	ErrAccountNumberLength      = errors.New("account number must be exactly 10 digits")

)


func (r *InitiateWithdrawalRequest) Validate() error {
	if err := validateAccountNumber(r.AccountNumber); err != nil {
		return err
	}
	if err := validateAmount(r.Amount); err != nil {
		return err
	}
	return nil
}

func (r *RetrieveAccountDetailsRequest) Validate() error {
	if err := validateAccountNumber(r.AccountNumber); err != nil {
		return err
	}
	return nil
}

func validateAccountNumber(accountNumber string) error {
	accountNumber = strings.TrimSpace(accountNumber)
	if len(accountNumber) != 10 {
		return ErrAccountNumberLength
	}
	matched, err := regexp.MatchString(`^\d{10}$`, accountNumber)
	if err != nil {
		return err
	}
	if !matched {
		return ErrAccountNumberLength
	}
	return nil
}
func validateAmount(amount int64) error {
	if amount < 1000 {
		return ErrAmountTooShort
	}
	return nil
}