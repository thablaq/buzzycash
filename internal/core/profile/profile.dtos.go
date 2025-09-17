package profile




import "time"

type CreateProfileRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100" validate:"required,min=2,max=100"`
	Gender   string `json:"gender" binding:"required,oneof=MALE FEMALE OTHERS" validate:"required,oneof=MALE FEMALE OTHERS"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	UserName string `json:"user_name" binding:"required,min=3,max=30" validate:"required,min=3,max=30"`
}

type ProfileUpdateRequest struct {
	FullName    string `json:"full_name,omitempty" binding:"omitempty,min=2,max=100" validate:"omitempty,min=2,max=100"`
	Gender      string `json:"gender,omitempty" binding:"omitempty,oneof=MALE FEMALE OTHERS" validate:"omitempty,oneof=MALE FEMALE OTHERS"`
	DateOfBirth string `json:"date_of_birth,omitempty" binding:"omitempty,datetime=2006-01-02" validate:"omitempty,datetime=2006-01-02"`
}


type VerifyEmailProfileRequest struct {
    Email    string `json:"email" binding:"required,email" validate:"required,email"`
    VerificationCode string `json:"verification_code" binding:"required" validate:"len=6"`
}


type ChooseUsernameRequest struct {
	UserName string `json:"user_name" binding:"required,min=3,max=30" validate:"required,min=3,max=30"`
}


type ProfileResponse struct {
	ID               string    `json:"id"`
	PhoneNumber      string    `json:"phoneNumber"`
	FullName         string    `json:"fullName"`
	Gender           string    `json:"gender"`
	DateOfBirth      string    `json:"dateOfBirth"`
	Email            string    `json:"email"`
	IsProfileCreated bool      `json:"isProfileCreated"`
	IsVerified       bool      `json:"isVerified"`
	IsActive         bool      `json:"isActive"`
	Username         string    `json:"username"`
	LastLogin        time.Time `json:"lastLogin,omitempty"`
	CountryOfResidence string  `json:"countryOfResidence,omitempty"`
}