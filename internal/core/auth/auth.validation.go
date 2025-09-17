package auth

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// Validation rules
const (
	MinPasswordLength    = 8
	MaxPhoneNumberLength = 18
	MinPhoneNumberLength = 7
	OtpLength            = 6
)

// Validation errors
var (
	ErrPasswordTooShort         = errors.New("password must be at least 8 characters long")
	ErrPasswordNoUppercase      = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoNumber         = errors.New("password must contain at least one number")
	ErrPasswordNoSpecialChar    = errors.New("password must contain at least one special character")
	ErrPhoneNumberLength        = errors.New("phone number must be between 7 to 18 digits")
	ErrPasswordsDontMatch       = errors.New("passwords do not match")
	ErrEmailOrPhoneRequired     = errors.New("either email or phone number is required")
	ErrInvalidEmail             = errors.New("invalid email address")
	ErrOtpLength                = errors.New("OTP must be exactly 6 characters")
	ErrCountryRequired          = errors.New("country is required")
	ErrCurrentPasswordRequired  = errors.New("current password is required")
	ErrNewPasswordRequired      = errors.New("new password is required")
	ErrConfirmPasswordRequired  = errors.New("confirm password is required")
	ErrUserIdRequired           = errors.New("user ID is required")
	ErrVerificationCodeRequired = errors.New("verification code is required")
	ErrProvidedEmailAndPhone    = errors.New("provide either email and password or phone number and password, not both")
)

// Validation methods

func (r *SignUpRequest) Validate() error {
	if err := validatePhoneNumber(r.PhoneNumber); err != nil {
		return err
	}
	if r.Password != r.ConfirmPassword {
		return ErrPasswordsDontMatch
	}

	if err := validatePassword(r.Password); err != nil {
		return err
	}
	if strings.TrimSpace(r.CountryOfResidence) == "" {
		return ErrCountryRequired
	}
	return nil
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" && r.PhoneNumber == "" {
		return ErrEmailOrPhoneRequired
	}

	// Cannot provide both email and phone
	if r.Email != "" && r.PhoneNumber != "" {
		return ErrProvidedEmailAndPhone
	}

	if r.Email != "" {
		if err := validateEmail(r.Email); err != nil {
			return err
		}
	}

	if r.PhoneNumber != "" {
		if err := validatePhoneNumber(r.PhoneNumber); err != nil {
			return err
		}
	}

	// Validate password length
	if len(r.Password) < MinPasswordLength {
		return ErrPasswordTooShort
	}

	return nil
}

func (r *PasswordChangeRequest) Validate() error {
	if strings.TrimSpace(r.CurrentPassword) == "" {
		return ErrCurrentPasswordRequired
	}
	if r.NewPassword != r.ConfirmNewPassword {
		return ErrPasswordsDontMatch
	}
	if err := validatePassword(r.NewPassword); err != nil {
		return err
	}
	return nil
}

func (r *VerifyAccountRequest) Validate() error {
	if err := validatePhoneNumber(r.PhoneNumber); err != nil {
		return err
	}
	if len(r.VerificationCode) != OtpLength {
		return ErrOtpLength
	}
	return nil
}

func (r *ResendOtpRequest) Validate() error {
	return validatePhoneNumber(r.PhoneNumber)
}

func (r *ForgotPasswordRequest) Validate() error {
	if r.Email == "" && r.PhoneNumber == "" {
		return errors.New("provide either email or phone number")
	}
	if r.Email != "" && r.PhoneNumber != "" {
		return errors.New("cannot provide both email and phone number")
	}
	if r.Email != "" {
		if err := validateEmail(r.Email); err != nil {
			return err
		}
	}
	if r.PhoneNumber != "" {
		if err := validatePhoneNumber(r.PhoneNumber); err != nil {
			return err
		}
	}
	return nil
}

func (r *ResetPasswordRequest) Validate() error {
	if strings.TrimSpace(r.UserId) == "" {
		return ErrUserIdRequired
	}
	if r.NewPassword != r.ConfirmNewPassword {
		return ErrPasswordsDontMatch
	}
	if err := validatePassword(r.NewPassword); err != nil {
		return err
	}
	return nil
}

func (r *VerifyPasswordForgotOtpRequest) Validate() error {
	if r.Email == "" && r.PhoneNumber == "" {
		return ErrEmailOrPhoneRequired
	}
	if r.Email != "" {
		if err := validateEmail(r.Email); err != nil {
			return err
		}
	}
	if r.PhoneNumber != "" {
		if err := validatePhoneNumber(r.PhoneNumber); err != nil {
			return err
		}
	}
	if len(r.VerificationCode) != OtpLength {
		return ErrOtpLength
	}
	return nil
}

func (r *RefreshTokenRequest) Validate() error {
	if strings.TrimSpace(r.RefreshToken) == "" {
		return errors.New("refresh token is required")
	}
	return nil
}

// Helper validation functions
func validatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}

	var (
		hasUpper   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsNumber(c):
			hasNumber = true
		case !unicode.IsLetter(c) && !unicode.IsNumber(c):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUppercase
	}
	if !hasNumber {
		return ErrPasswordNoNumber
	}
	if !hasSpecial {
		return ErrPasswordNoSpecialChar
	}

	return nil
}

func validatePhoneNumber(phone string) error {
	// Remove all non-digit characters
	re := regexp.MustCompile(`\D`)
	digits := re.ReplaceAllString(phone, "")

	if len(digits) < MinPhoneNumberLength || len(digits) > MaxPhoneNumberLength {
		return ErrPhoneNumberLength
	}
	return nil
}

func validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}
