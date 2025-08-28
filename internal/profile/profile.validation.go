package profile

import (
	"errors"
	"regexp"
	"strings"
)

// Validation errors
var (
	ErrFullNameRequired      = errors.New("full name is required")
	ErrFullNameTooShort      = errors.New("full name must be at least 2 characters")
	ErrGenderRequired        = errors.New("gender is required")
	ErrInvalidGender         = errors.New("gender must be MALE, FEMALE or OTHERS")
	ErrEmailRequired         = errors.New("email is required")
	ErrInvalidEmail          = errors.New("invalid email address")
	ErrUsernameRequired      = errors.New("username is required")
	ErrUsernameTooShort      = errors.New("username must be at least 3 characters")
	ErrUsernameTooLong       = errors.New("username must be at most 30 characters")
	ErrUsernameInvalidStart  = errors.New("username must start with a letter")
	ErrUsernameInvalidChars  = errors.New("username can only contain letters, numbers, or underscores")
	ErrUsernameDoubleUnderscore = errors.New("username cannot contain consecutive underscores")
	ErrUsernameEndsWithUnderscore = errors.New("username cannot end with an underscore")
	ErrUsernameContainsAt    = errors.New("username cannot contain '@'")
	ErrUsernameContainsDotCom = errors.New("username cannot contain '.com'")
	ErrUsernameContainsDotMail = errors.New("username cannot contain '.mail'")
	ErrUsernameIsEmail       = errors.New("username cannot be an email address")
	ErrNoFieldsToUpdate      = errors.New("at least one field must be provided to update")
)

// ValidateCreateProfile validates CreateProfileRequest
func (r *CreateProfileRequest) Validate() error {
	// Validate full name
	if strings.TrimSpace(r.FullName) == "" {
		return ErrFullNameRequired
	}
	if len(r.FullName) < 2 {
		return ErrFullNameTooShort
	}

	// Validate gender
	if strings.TrimSpace(r.Gender) == "" {
		return ErrGenderRequired
	}
	if r.Gender != "MALE" && r.Gender != "FEMALE" && r.Gender != "OTHERS" {
		return ErrInvalidGender
	}

	// Validate email
	if strings.TrimSpace(r.Email) == "" {
		return ErrEmailRequired
	}
	if !validateEmail(r.Email) {
		return ErrInvalidEmail
	}

	// Validate username
	if err := validateUsername(r.UserName); err != nil {
		return err
	}

	return nil
}

// ValidateUpdateProfile validates ProfileUpdateRequest
func (r *ProfileUpdateRequest) Validate() error {
	// Check at least one field is provided
	if strings.TrimSpace(r.FullName) == "" && 
	   strings.TrimSpace(r.Gender) == "" && 
	   strings.TrimSpace(r.DateOfBirth) == "" {
		return ErrNoFieldsToUpdate
	}

	// Validate full name if provided
	if r.FullName != "" && len(r.FullName) < 2 {
		return ErrFullNameTooShort
	}

	// Validate gender if provided
	if r.Gender != "" && r.Gender != "MALE" && r.Gender != "FEMALE" && r.Gender != "OTHERS" {
		return ErrInvalidGender
	}

	return nil
}

// validateUsername implements all username validation rules
func validateUsername(username string) error {
	if strings.TrimSpace(username) == "" {
		return ErrUsernameRequired
	}
	if len(username) < 3 {
		return ErrUsernameTooShort
	}
	if len(username) > 30 {
		return ErrUsernameTooLong
	}

	// Must start with letter
	if !regexp.MustCompile(`^[a-zA-Z]`).MatchString(username) {
		return ErrUsernameInvalidStart
	}

	// Only letters, numbers, underscores
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return ErrUsernameInvalidChars
	}

	// No consecutive underscores
	if strings.Contains(username, "__") {
		return ErrUsernameDoubleUnderscore
	}

	// Doesn't end with underscore
	if strings.HasSuffix(username, "_") {
		return ErrUsernameEndsWithUnderscore
	}

	// Doesn't contain @
	if strings.Contains(username, "@") {
		return ErrUsernameContainsAt
	}

	// Doesn't contain .com
	if strings.Contains(strings.ToLower(username), ".com") {
		return ErrUsernameContainsDotCom
	}

	// Doesn't contain .mail
	if strings.Contains(strings.ToLower(username), ".mail") {
		return ErrUsernameContainsDotMail
	}

	// Not an email address
	if regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`).MatchString(username) {
		return ErrUsernameIsEmail
	}

	return nil
}

// validateEmail checks if email is valid
func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}


