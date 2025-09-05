package auth

// @Summary Register new user
// @Description Create a user account
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body SignUpRequest true "Registration data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /register [post]
func _() {}

// @Summary User login
// @Description Authenticate user
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /login [post]
func _() {}

// @Summary Verify account
// @Description Verify user account with OTP
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body VerifyAccountRequest true "Verification data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /verify-account [post]
func _() {}

// @Summary Resend OTP
// @Description Resend verification OTP to user
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body ResendOtpRequest true "Phone number"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /resend-otp [post]
func _() {}

// @Summary Change password
// @Description Change password for authenticated user
// @Tags authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body PasswordChangeRequest true "Password change data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /change-password [patch]
func _() {}

// @Summary Forgot password
// @Description Request password reset using email or phone number
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Forgot password data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /forgot-password [post]
func _() {}

// @Summary Verify reset password OTP
// @Description Verify OTP for password reset flow
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body VerifyPasswordForgotOtpRequest true "Verification data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /verify-reset-password-otp [post]
func _() {}

// @Summary Reset password
// @Description Reset password after OTP verification
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset password data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /reset-password [put]
func _() {}

// @Summary Logout user
// @Description Logout the authenticated user and invalidate tokens
// @Tags authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /logout [post]
func _() {}

// @Summary Refresh token
// @Description Refresh access token using refresh token
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /refresh-token [post]
func _() {}
