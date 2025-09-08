package tickets



	import (
		"errors"
)




// Validation errors
var (
	ErrAmountTooShort          = errors.New("amount must be at least 100 naira")
	ErrQuantityTooShort      = errors.New("quantity must be at least 1")

)


func (r *BuyTicketRequest) Validate() error {
	if err := validateAmount(r.AmountPaid); err != nil {
		return err
	}
	if err := validateQuantity(r.Quantity); err != nil {
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

func validateQuantity(quantity int) error {
	if quantity < 1 {
		return ErrQuantityTooShort
	}
	return nil
}
