package auth




// Validation rules
const (
    MinPasswordLength    = 8
    MaxPhoneNumberLength = 18
    MinPhoneNumberLength = 7
    OtpLength            = 6
)

// DTOs for authentication-related requests

type SignUpRequest struct {
    PhoneNumber        string `json:"phoneNumber" binding:"required" validate:"min=7,max=18"`
    Password           string `json:"password" binding:"required" validate:"min=8"`
    ConfirmPassword    string `json:"confirmPassword" binding:"required" validate:"eqfield=Password"`
    CountryOfResidence string `json:"countryOfResidence" binding:"required"`
    ReferralCode       string `json:"referralCode,omitempty"`
}

type LoginRequest struct {
    Email       string `json:"email,omitempty" validate:"omitempty,email"`
    PhoneNumber string `json:"phoneNumber,omitempty" validate:"omitempty,min=7,max=18"`
    Password    string `json:"password" binding:"required"`
}

type PasswordChangeRequest struct {
    CurrentPassword    string `json:"currentPassword" binding:"required" validate:"min=8"`
    NewPassword        string `json:"newPassword" binding:"required" validate:"min=8"`
    ConfirmNewPassword string `json:"confirmNewPassword" binding:"required" validate:"eqfield=NewPassword"`
}

type VerifyAccountRequest struct {
    PhoneNumber      string `json:"phoneNumber" binding:"required" validate:"min=7,max=18"`
    VerificationCode string `json:"verificationCode" binding:"required" validate:"len=6"`
}

type ResendOtpRequest struct {
    PhoneNumber string `json:"phoneNumber" binding:"required" validate:"min=7,max=18"`
}

type ForgotPasswordRequest struct {
    Email       string `json:"email,omitempty" validate:"omitempty,email"`
    PhoneNumber string `json:"phoneNumber,omitempty" validate:"omitempty,min=7,max=18"`
}

type ResetPasswordRequest struct {
    UserId             string `json:"userId" binding:"required"`
    NewPassword        string `json:"newPassword" binding:"required" validate:"min=8"`
    ConfirmNewPassword string `json:"confirmNewPassword" binding:"required" validate:"eqfield=NewPassword"`
}

type VerifyPasswordForgotOtpRequest struct {
    Email            string `json:"email,omitempty" validate:"omitempty,email"`
    PhoneNumber      string `json:"phoneNumber,omitempty" validate:"omitempty,min=7,max=18"`
    VerificationCode string `json:"verificationCode" binding:"required" validate:"len=6"`
}

type RefreshTokenRequest struct {
    RefreshToken string `json:"refreshToken" binding:"required"`
}
