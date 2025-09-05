package mailers

import (
	"fmt"
	"time"
     "log"
	"github.com/dblaq/buzzycash/internal/models"
)




// // SendOtp sends OTP to Nigerian phone numbers
// func (es *EmailService) SendNaijaOtp(phoneNumber, userID string) (interface{}, error) {
// 	if phoneNumber == "" {
// 		log.Println("Recipient phone number is missing")
// 		return nil, fmt.Errorf("recipient phone number is required")
// 	}

// 	otp := es.GenerateOtp()
// 	otpExpiresAt := time.Now().Add(5 * time.Minute)
// 	log.Println("Sending OTP to", phoneNumber[:3]+"****")

// 	if err := es.updateOrCreateOtp(userID, otp, otpExpiresAt, models.OtpActionVerifyAccount,"phone"); err != nil {
// 		log.Println("Failed to update OTP record for user:", userID, "Error:", err)
// 		return nil, fmt.Errorf("failed to update OTP record: %v", err)
// 	}

// 	formattedNumber := es.formatPhoneNumber(phoneNumber, "234")
// 	message := fmt.Sprintf("Your Otp verification code is %s. Valid for 5 minutes.", otp)

// 	result, err := es.sendSmsViaLenhub(formattedNumber, message)
// 	if err != nil {
// 		log.Println("Failed to send OTP to:", formattedNumber, "Error:", err)
// 		if clearErr := es.clearOtp(userID,models.OtpActionVerifyAccount); clearErr != nil {
// 			log.Println("Failed to clear OTP after send failure for user:", userID, "Error:", clearErr)
// 		}
// 		return nil, fmt.Errorf("failed to send OTP: %v", err)
// 	}

// 	log.Println("OTP sent successfully to:", formattedNumber)
// 	return result, nil
// }

// // SendGhanaOtp sends OTP to Ghanaian phone numbers
// func (es *EmailService) SendGhanaOtp(phoneNumber, userID string) (interface{}, error) {
// 	if phoneNumber == "" {
// 		log.Println("Recipient phone number is missing")
// 		return nil, fmt.Errorf("recipient phone number is required")
// 	}

// 	otp := es.GenerateOtp()
// 	otpExpiresAt := time.Now().Add(5 * time.Minute)
// 	log.Println("Sending OTP to", phoneNumber[:3]+"****")

// 	if err := es.updateOrCreateOtp(userID, otp, otpExpiresAt, models.OtpActionVerifyAccount,"phone"); err != nil {
// 		log.Println("Failed to update OTP record for user:", userID, "Error:", err)
// 		return nil, fmt.Errorf("failed to update OTP record: %v", err)
// 	}

// 	formattedNumber := es.formatPhoneNumber(phoneNumber, "233")
// 	message := fmt.Sprintf("Your OTP verification code is %s. Valid for 5 minutes.", otp)

// 	result, err := es.sendSmsViaHubtel(formattedNumber, message)
// 	if err != nil {
// 		log.Println("Failed to send OTP to:", formattedNumber, "Error:", err)
// 		if clearErr := es.clearOtp(userID,models.OtpActionVerifyAccount); clearErr != nil {
// 			log.Println("Failed to clear OTP after send failure for user:", userID, "Error:", clearErr)
// 		}
// 		return nil, fmt.Errorf("failed to send OTP: %v", err)
// 	}

// 	log.Println("OTP sent successfully to:", formattedNumber)
// 	return result, nil
// }

// // SendChangePasswordSuccess sends password change confirmation email
// func (es *EmailService) SendChangePasswordSuccess(recipient, firstName, lastName string) (interface{}, error) {
// 	if recipient == "" {
// 		log.Println("Recipient email is missing")
// 		return nil, fmt.Errorf("recipient email is required")
// 	}

// 	log.Println("Preparing to send change password success email to:", recipient)
// 	emailContent, err := es.GetEmailTemplate("change-password-success", map[string]string{
// 		"title":   "Change Password Success",
// 		"message": "Your password has been changed successfully.",
// 		"name":    firstName + " " + lastName,
// 		"year":    fmt.Sprintf("%d", time.Now().Year()),
// 	})
// 	if err != nil {
// 		log.Println("Failed to get email template for recipient:", recipient, "Error:", err)
// 		return nil, fmt.Errorf("failed to get email template: %v", err)
// 	}

// 	log.Println("Sending email to:", recipient)
// 	result, err := es.sendEmailViaLenhub(
// 		recipient,
// 		"Change Password Success",
// 		emailContent,
// 		"Dear "+firstName+" "+lastName,
// 	)
// 	if err != nil {
// 		log.Println("Failed to send change password success email to:", recipient, "Error:", err)
// 		return nil, fmt.Errorf("failed to send email: %v", err)
// 	}

// 	log.Println("Change password success email sent successfully to:", recipient)
// 	return result, nil
// }

// // SendForgotPasswordNGNOtp sends forgot password OTP to Nigerian numbers
// func (es *EmailService) SendForgotPasswordNGNOtp(phoneNumber, userID string) (interface{}, error) {
// 	if phoneNumber == "" {
// 		log.Println("Recipient phone number is missing")
// 		return nil, fmt.Errorf("recipient phone number is required")
// 	}

// 	otp := es.GenerateOtp()
// 	otpExpiresAt := time.Now().Add(5 * time.Minute)
// 	log.Println("Generated OTP:", otp, "for user:", userID)
// 	log.Println("Sending forgot password OTP to", phoneNumber[:3]+"****")

// 	if err := es.updateOrCreateOtp(userID, otp, otpExpiresAt,models.OtpActionPasswordReset,"phone"); err != nil {
// 		log.Println("Failed to update OTP record for user:", userID, "Error:", err)
// 		return nil, fmt.Errorf("failed to update OTP record: %v", err)
// 	}

// 	formattedNumber := es.formatPhoneNumber(phoneNumber, "234")
// 	message := fmt.Sprintf("Your Otp verification code is %s. Valid for 5 minutes.", otp)

// 	result, err := es.sendSmsViaLenhub(formattedNumber, message)
// 	if err != nil {
// 		log.Println("Failed to send OTP to:", formattedNumber, "Error:", err)
// 		if clearErr := es.clearOtp(userID,models.OtpActionPasswordReset); clearErr != nil {
// 			log.Println("Failed to clear OTP after send failure for user:", userID, "Error:", clearErr)
// 		}
// 		return nil, fmt.Errorf("failed to send OTP: %v", err)
// 	}

// 	log.Println("Forgot password OTP sent successfully to:", formattedNumber)
// 	return result, nil
// }


// // SendForgotPasswordGHCOtp sends forgot password OTP to Ghanaian numbers
// func (es *EmailService) SendForgotPasswordGHCOtp(phoneNumber, userID string) (interface{}, error) {
// 	if phoneNumber == "" {
// 		log.Println("Recipient phone number is missing")
// 		return nil, fmt.Errorf("recipient phone number is required")
// 	}

// 	otp := es.GenerateOtp()
// 	otpExpiresAt := time.Now().Add(5 * time.Minute)
// 	log.Println("Generated OTP:", otp, "for user:", userID)
// 	log.Println("Sending forgot password OTP to", phoneNumber[:3]+"****")

// 	if err := es.updateOrCreateOtp(userID, otp, otpExpiresAt, models.OtpActionPasswordReset,"phone"); err != nil {
// 		log.Println("Failed to update OTP record for user:", userID, "Error:", err)
// 		return nil, fmt.Errorf("failed to update OTP record: %v", err)
// 	}

// 	formattedNumber := es.formatPhoneNumber(phoneNumber, "233")
// 	message := fmt.Sprintf("Your OTP verification code is %s. Valid for 5 minutes.", otp)

// 	result, err := es.sendSmsViaHubtel(formattedNumber, message)
// 	if err != nil {
// 		log.Println("Failed to send OTP to:", formattedNumber, "Error:", err)
// 		if clearErr := es.clearOtp(userID,models.OtpActionPasswordReset); clearErr != nil {
// 			log.Println("Failed to clear OTP after send failure for user:", userID, "Error:", clearErr)
// 		}
// 		return nil, fmt.Errorf("failed to send OTP: %v", err)
// 	}

// 	log.Println("Forgot password OTP sent successfully to:", formattedNumber)
// 	return result, nil
// }

// SendForgotPasswordOtp sends forgot password OTP via email
func (es *EmailService) SendForgotPasswordEmailOtp(recipient, fullName, userID string) (interface{}, error) {
	if recipient == "" {
		log.Println("Recipient email is missing")
		return nil, fmt.Errorf("recipient email is required")
	}

	otp := es.GenerateOtp()
	otpExpiresAt := time.Now().Add(5 * time.Minute)
	log.Println("Generated OTP:", otp, "for user:", userID)
	log.Println("Sending password reset OTP to", recipient[:3]+"****@***.com")

	
	if err := es.updateOrCreateOtp(userID, otp, otpExpiresAt, models.OtpActionPasswordReset, "email"); err != nil {
		log.Println("Failed to update OTP record for user:", userID, "Error:", err)
		return nil, fmt.Errorf("failed to update OTP record: %v", err)
	}

	emailContent, err := es.GetEmailTemplate("forgot-password", map[string]string{
		"title":   "🔑 Reset Your Password",
		"message": "Use the OTP below to reset your password.",
		"name":    fullName,
		"otp":     otp,
	})
	if err != nil {
		log.Println("Failed to get email template for user:", userID, "Error:", err)
		return nil, fmt.Errorf("failed to get email template: %v", err)
	}

	result, err := es.sendEmailViaLenhub(
		recipient,
		"Your OTP for Password Reset",
		emailContent,
		"Dear "+fullName,
	)
	if err != nil {
		log.Println("Failed to send password reset OTP to:", recipient, "Error:", err)
		
		if clearErr := es.clearOtp(userID,models.OtpActionPasswordReset); clearErr != nil {
			log.Println("Failed to clear OTP after send failure for user:", userID, "Error:", clearErr)
		}
		return nil, fmt.Errorf("failed to send password reset OTP: %v", err)
	}

	log.Println("Password reset OTP sent successfully to:", recipient)
	return result, nil
}



func (es *EmailService) SendEmailVerificationOtp(recipient, fullName, userID string) (interface{}, error) {
	if recipient == "" {
		log.Println("Recipient email is missing")
		return nil, fmt.Errorf("recipient email is required")
	}

	otp := es.GenerateOtp()
	otpExpiresAt := time.Now().Add(5 * time.Minute)
	log.Println("Generated OTP:", otp, "for user:", userID)
	log.Println("Sending email verification OTP to", recipient[:3]+"****@***.com")

	if err := es.updateOrCreateOtp(userID, otp, otpExpiresAt, models.OtpActionVerifyEmail, "email"); err != nil {
		log.Println("Failed to update OTP record for user:", userID, "Error:", err)
		return nil, fmt.Errorf("failed to update OTP record: %v", err)
	}

	emailContent, err := es.GetEmailTemplate("email-verification", map[string]string{
		"title":   "🔑 Verify Your Email",
		"message": "Use the OTP below to verify your email.",
		"name":    fullName,
		"otp":     otp,
	})
	if err != nil {
		log.Println("Failed to get email template for user:", userID, "Error:", err)
		return nil, fmt.Errorf("failed to get email template: %v", err)
	}

	result, err := es.sendEmailViaLenhub(
		recipient,
		"Your OTP for Email Verification",
		emailContent,
		"Dear "+fullName,
	)
	if err != nil {
		log.Println("Failed to send email verification OTP to:", recipient, "Error:", err)
		if clearErr := es.clearOtp(userID,models.OtpActionVerifyEmail); clearErr != nil {
			log.Println("Failed to clear OTP after send failure for user:", userID, "Error:", clearErr)
		}
		return nil, fmt.Errorf("failed to send email verification OTP: %v", err)
	}

	log.Println("Email verification  OTP sent successfully to:", recipient)
	return result, nil
}




// // NombaPaymentSuccess sends payment success notification
// func (es *EmailService) NombaPaymentSuccess(recipient, firstName, lastName, transactionReference string, amount float64) (map[string]interface{}, error) {
// 	if recipient == "" {
// 		return nil, fmt.Errorf("recipient email is required")
// 	}

// 	emailContent, err := es.GetEmailTemplate("nomba-payment-confirmation", map[string]string{
// 		"title":              "Payment Success",
// 		"message":            "Payment successful.",
// 		"name":              firstName + " " + lastName,
// 		"transactionReference": transactionReference,
// 		"amount":            fmt.Sprintf("%.2f", amount),
// 		"year":              fmt.Sprintf("%d", time.Now().Year()),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get email template: %v", err)
// 	}

// 	result, err := es.sendEmailViaLenhub(
// 		recipient,
// 		"Payment Successful",
// 		emailContent,
// 		"Dear "+firstName+" "+lastName,
// 	)
// 	if err != nil {
// 		return map[string]interface{}{
// 			"success":    false,
// 			"statusCode": 500,
// 			"message":    "Failed to send message",
// 			"error":      err.Error(),
// 		}, nil
// 	}

// 	return map[string]interface{}{
// 		"success":    true,
// 		"statusCode": 200,
// 		"message":    "Payment Success",
// 		"token":      result,
// 	}, nil
// }