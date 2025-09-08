package auth

type SignUpRequest struct {
	PhoneNumber        string `json:"phone_number" binding:"required" validate:"min=7,max=18"`
	Password           string `json:"password" binding:"required" validate:"min=8"`
	ConfirmPassword    string `json:"confirm_password" binding:"required" validate:"eqfield=Password"`
	CountryOfResidence string `json:"country_of_residence" binding:"required"`
	ReferralCode       string `json:"referral_code,omitempty"`
}

type LoginRequest struct {
	Email       string `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber string `json:"phone_number,omitempty" validate:"omitempty,min=7,max=18"`
	Password    string `json:"password" binding:"required"`
}

type PasswordChangeRequest struct {
	CurrentPassword    string `json:"current_password" binding:"required" validate:"min=8"`
	NewPassword        string `json:"new_password" binding:"required" validate:"min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required" validate:"eqfield=NewPassword"`
}

type VerifyAccountRequest struct {
	PhoneNumber      string `json:"phone_number" binding:"required" validate:"min=7,max=18"`
	VerificationCode string `json:"verification_code" binding:"required" validate:"len=6"`
}

type ResendOtpRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" validate:"min=7,max=18"`
}

type ForgotPasswordRequest struct {
	Email       string `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber string `json:"phone_number,omitempty" validate:"omitempty,min=7,max=18"`
}

type ResetPasswordRequest struct {
	UserId             string `json:"user_id" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required" validate:"min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required" validate:"eqfield=NewPassword"`
}

type VerifyPasswordForgotOtpRequest struct {
	Email            string `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber      string `json:"phone_number,omitempty" validate:"omitempty,min=7,max=18"`
	VerificationCode string `json:"verification_code" binding:"required" validate:"len=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
