package sms




import (
	"fmt"
	"time"
     "log"
	"github.com/dblaq/buzzycash/internal/models"
)




// SendOtp sends OTP to Nigerian phone numbers
func (es *SmsService) SendNaijaOtp(phoneNumber, userID string) (interface{}, error) {
	if phoneNumber == "" {
		log.Println("Recipient phone number is missing")
		return nil, fmt.Errorf("recipient phone number is required")
	}

	otp := es.GenerateOtp()
	otpExpiresAt := time.Now().Add(5 * time.Minute)
	log.Println("Sending OTP to", phoneNumber[:3]+"****")

	if err := es.UpdateOrCreateOtp(userID, otp, otpExpiresAt, models.OtpActionVerifyAccount,"phone"); err != nil {
		log.Println("Failed to update OTP record for user:", userID, "Error:", err)
		return nil, fmt.Errorf("failed to update OTP record: %v", err)
	}

	formattedNumber := es.formatPhoneNumber(phoneNumber, "234")
	message := fmt.Sprintf("Your Otp verification code is %s. Valid for 5 minutes.", otp)

	result, err := es.sendSmsViaLenhub(formattedNumber, message)
	if err != nil {
		log.Println("Failed to send OTP to:", formattedNumber, "Error:", err)
		if clearErr := es.ClearOtp(userID,models.OtpActionVerifyAccount); clearErr != nil {
			log.Println("Failed to clear OTP after send failure for user:", userID, "Error:", clearErr)
		}
		return nil, fmt.Errorf("failed to send OTP: %v", err)
	}

	log.Println("OTP sent successfully to:", formattedNumber)
	return result, nil
}

// SendGhanaOtp sends OTP to Ghanaian phone numbers
func (es *SmsService) SendGhanaOtp(phoneNumber, userID string) (interface{}, error) {
	if phoneNumber == "" {
		log.Println("Recipient phone number is missing")
		return nil, fmt.Errorf("recipient phone number is required")
	}

	otp := es.GenerateOtp()
	otpExpiresAt := time.Now().Add(5 * time.Minute)
	log.Println("Sending OTP to", phoneNumber[:3]+"****")

	if err := es.UpdateOrCreateOtp(userID, otp, otpExpiresAt, models.OtpActionVerifyAccount,"phone"); err != nil {
		log.Println("Failed to update OTP record for user:", userID, "Error:", err)
		return nil, fmt.Errorf("failed to update OTP record: %v", err)
	}

	formattedNumber := es.formatPhoneNumber(phoneNumber, "233")
	message := fmt.Sprintf("Your OTP verification code is %s. Valid for 5 minutes.", otp)

	result, err := es.sendSmsViaHubtel(formattedNumber, message)
	if err != nil {
		log.Println("Failed to send OTP to:", formattedNumber, "Error:", err)
		if clearErr := es.ClearOtp(userID,models.OtpActionVerifyAccount); clearErr != nil {
			log.Println("Failed to clear OTP after send failure for user:", userID, "Error:", clearErr)
		}
		return nil, fmt.Errorf("failed to send OTP: %v", err)
	}

	log.Println("OTP sent successfully to:", formattedNumber)
	return result, nil
}


// SendForgotPasswordNGNOtp sends forgot password OTP to Nigerian numbers
func (es *SmsService) SendForgotPasswordNGNOtp(phoneNumber, userID string) (interface{}, error) {
	if phoneNumber == "" {
		log.Println("Recipient phone number is missing")
		return nil, fmt.Errorf("recipient phone number is required")
	}

	otp := es.GenerateOtp()
	otpExpiresAt := time.Now().Add(5 * time.Minute)
	log.Println("Generated OTP:", otp, "for user:", userID)
	log.Println("Sending forgot password OTP to", phoneNumber[:3]+"****")

	if err := es.UpdateOrCreateOtp(userID, otp, otpExpiresAt,models.OtpActionPasswordReset,"phone"); err != nil {
		log.Println("Failed to update OTP record for user:", userID, "Error:", err)
		return nil, fmt.Errorf("failed to update OTP record: %v", err)
	}

	formattedNumber := es.formatPhoneNumber(phoneNumber, "234")
	message := fmt.Sprintf("Your Otp verification code is %s. Valid for 5 minutes.", otp)

	result, err := es.sendSmsViaLenhub(formattedNumber, message)
	if err != nil {
		log.Println("Failed to send OTP to:", formattedNumber, "Error:", err)
		if clearErr := es.ClearOtp(userID,models.OtpActionPasswordReset); clearErr != nil {
			log.Println("Failed to clear OTP after send failure for user:", userID, "Error:", clearErr)
		}
		return nil, fmt.Errorf("failed to send OTP: %v", err)
	}

	log.Println("Forgot password OTP sent successfully to:", formattedNumber)
	return result, nil
}


// SendForgotPasswordGHCOtp sends forgot password OTP to Ghanaian numbers
func (es *SmsService) SendForgotPasswordGHCOtp(phoneNumber, userID string) (interface{}, error) {
	if phoneNumber == "" {
		log.Println("Recipient phone number is missing")
		return nil, fmt.Errorf("recipient phone number is required")
	}

	otp := es.GenerateOtp()
	otpExpiresAt := time.Now().Add(5 * time.Minute)
	log.Println("Generated OTP:", otp, "for user:", userID)
	log.Println("Sending forgot password OTP to", phoneNumber[:3]+"****")

	if err := es.UpdateOrCreateOtp(userID, otp, otpExpiresAt, models.OtpActionPasswordReset,"phone"); err != nil {
		log.Println("Failed to update OTP record for user:", userID, "Error:", err)
		return nil, fmt.Errorf("failed to update OTP record: %v", err)
	}

	formattedNumber := es.formatPhoneNumber(phoneNumber, "233")
	message := fmt.Sprintf("Your OTP verification code is %s. Valid for 5 minutes.", otp)

	result, err := es.sendSmsViaHubtel(formattedNumber, message)
	if err != nil {
		log.Println("Failed to send OTP to:", formattedNumber, "Error:", err)
		if clearErr := es.ClearOtp(userID,models.OtpActionPasswordReset); clearErr != nil {
			log.Println("Failed to clear OTP after send failure for user:", userID, "Error:", clearErr)
		}
		return nil, fmt.Errorf("failed to send OTP: %v", err)
	}

	log.Println("Forgot password OTP sent successfully to:", formattedNumber)
	return result, nil
}




