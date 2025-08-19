package services

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	// "io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/dblaq/buzzycash/internal/config"
	// "yourproject/db"
	"github.com/dblaq/buzzycash/internal/models"
)



type EmailService struct{}

// GetEmailTemplate reads and processes an email template
func (es *EmailService) GetEmailTemplate(templateName string, context map[string]string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %v", err)
	}

	templatePath := filepath.Join(cwd, "src/templates", templateName+".html")
	fmt.Println("Using template path:", templatePath)

	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %v", err)
	}

	template := string(templateBytes)
	for key, value := range context {
		template = strings.ReplaceAll(template, "$"+key, value)
	}

	return template, nil
}

// GenerateOtp generates a 6-digit OTP
func (es *EmailService) GenerateOtp() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

// clearOtp clears OTP from database
func (es *EmailService) clearOtp(userID string) error {
	return config.DB.Model(&models.UserOtpSecurity{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"verification_code":                          nil,
		"verification_code_created_at":               nil,
		"verification_code_expires_at":               nil,
		"password_reset_verification_code":           nil,
		"password_reset_verification_code_created_at": nil,
		"password_reset_verification_code_expires_at": nil,
		"password_reset_sent_to":                      nil,
	}).Error
}

// updateOrCreateOtp updates or creates an OTP record
func (es *EmailService) updateOrCreateOtp(userID, otp string, expiresAt time.Time, isPasswordReset bool) error {
	if isPasswordReset {
		return config.DB.Model(&models.UserOtpSecurity{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
			"password_reset_verification_code":           otp,
			"password_reset_verification_code_created_at": time.Now(),
			"password_reset_verification_code_expires_at": expiresAt,
			"password_reset_sent_to":                      "phone",
		}).Error
	}

	now := time.Now()
	expires := time.Now().Add(5 * time.Minute)
   otpSecurity := models.UserOtpSecurity{
    UserID:                    userID,
    VerificationCode:          otp,
    VerificationCodeCreatedAt: &now,
    VerificationCodeExpiresAt: &expires, 
}

	var existing models.UserOtpSecurity
	err := config.DB.Where("user_id = ?", userID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return config.DB.Create(&otpSecurity).Error
	} else if err != nil {
		return err
	}
	return config.DB.Model(&models.UserOtpSecurity{}).Where("user_id = ?", userID).Updates(&otpSecurity).Error
}

// sendSmsViaLenhub sends SMS using LENHUB API
func (es *EmailService) sendSmsViaLenhub(phoneNumber, message string) (interface{}, error) {
	payload := map[string]interface{}{
		"client_id":       config.AppConfig.LenhubClientID,
		"receiver_number": phoneNumber,
		"message":         message,
		"sender_id":       config.AppConfig.BuzzyCashSenderID,
		"types":           "2",
	}

	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", config.AppConfig.LenhubApiBase+"sendsms/api", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.AppConfig.LenhubApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// bodyBytes, _ := io.ReadAll(resp.Body)
	// bodyString := string(bodyBytes)
	// fmt.Println("📤 Lenhub Request Payload:", string(jsonPayload))
	// 	fmt.Println("📥 Lenhub Raw Response:", bodyString)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to send SMS, status code: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// sendSmsViaHubtel sends SMS using Hubtel API
func (es *EmailService) sendSmsViaHubtel(phoneNumber, message string) (interface{}, error) {
	payload := map[string]interface{}{
		"From":    config.AppConfig.HubtelSenderID,
		"To":      phoneNumber,
		"Content": message,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", config.AppConfig.HubtelApiBase+"/messages/send", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(config.AppConfig.HubtelClientID + ":" + config.AppConfig.HubtelClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to send SMS, status code: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// sendEmailViaLenhub sends email using LENHUB API
func (es *EmailService) sendEmailViaLenhub(recipient, subject, message, greetings string) (interface{}, error) {
	payload := map[string]interface{}{
		"client_id":  config.AppConfig.LenhubClientID,
		"subject":    subject,
		"message":    message,
		"Sender_id":  config.AppConfig.LenhubSenderID,
		"recipient":  recipient,
		"greetings":  greetings,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", config.AppConfig.LenhubApiBase+"send/email/api", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.AppConfig.LenhubApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to send email, status code: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// formatPhoneNumber formats phone number based on country code
func (es *EmailService) formatPhoneNumber(phoneNumber, countryCode string) string {
    if strings.HasPrefix(phoneNumber, "0") {
        return countryCode + phoneNumber[1:]
    }
    if strings.HasPrefix(phoneNumber, countryCode) {
        return phoneNumber
    }
    return countryCode + phoneNumber
}



// SendOtp sends OTP to Nigerian phone numbers
func (es *EmailService) SendNaijaOtp(phoneNumber, userID string) (interface{}, error) {
	if phoneNumber == "" {
		return nil, fmt.Errorf("recipient phone number is required")
	}

	otp := es.GenerateOtp()
	otpExpiresAt := time.Now().Add(5 * time.Minute)
	fmt.Println("Sending OTP to", phoneNumber[:3])
	fmt.Println("📌 Saving OTP:", otp, "for user:", userID)

	if err := es.updateOrCreateOtp(userID, otp, otpExpiresAt, false); err != nil {
		return nil, fmt.Errorf("failed to update OTP record: %v", err)
	}

	formattedNumber := es.formatPhoneNumber(phoneNumber, "234")
	message := fmt.Sprintf("Your Otp verification code is %s. Valid for 5 minutes.", otp)

	result, err := es.sendSmsViaLenhub(formattedNumber, message)
	if err != nil {
		if clearErr := es.clearOtp(userID); clearErr != nil {
			fmt.Println("Failed to clear OTP after send failure:", clearErr)
		}
		return nil, fmt.Errorf("failed to send OTP: %v", err)
	}

	fmt.Println("OTP sent successfully to", formattedNumber)
	return result, nil
}

// SendGhanaOtp sends OTP to Ghanaian phone numbers
func (es *EmailService) SendGhanaOtp(phoneNumber, userID string) (interface{}, error) {
	if phoneNumber == "" {
		return nil, fmt.Errorf("recipient phone number is required")
	}

	otp := es.GenerateOtp()
	otpExpiresAt := time.Now().Add(5 * time.Minute)
	fmt.Println("Sending OTP to", phoneNumber[:3]+"****")

	if err := es.updateOrCreateOtp(userID, otp, otpExpiresAt, false); err != nil {
		return nil, fmt.Errorf("failed to update OTP record: %v", err)
	}

	formattedNumber := es.formatPhoneNumber(phoneNumber, "233")
	message := fmt.Sprintf("Your OTP verification code is %s. Valid for 5 minutes.", otp)

	result, err := es.sendSmsViaHubtel(formattedNumber, message)
	if err != nil {
		if clearErr := es.clearOtp(userID); clearErr != nil {
			fmt.Println("Failed to clear OTP after send failure:", clearErr)
		}
		return nil, fmt.Errorf("failed to send OTP: %v", err)
	}

	fmt.Println("OTP sent successfully to", formattedNumber)
	return result, nil
}

// SendChangePasswordSuccess sends password change confirmation email
func (es *EmailService) SendChangePasswordSuccess(recipient, firstName, lastName string) (interface{}, error) {
	if recipient == "" {
		return nil, fmt.Errorf("recipient email is required")
	}

	emailContent, err := es.GetEmailTemplate("change-password-success", map[string]string{
		"title":   "Change Password Success",
		"message": "Your password has been changed successfully.",
		"name":    firstName + " " + lastName,
		"year":    fmt.Sprintf("%d", time.Now().Year()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get email template: %v", err)
	}

	result, err := es.sendEmailViaLenhub(
		recipient,
		"Change Password Success",
		emailContent,
		"Dear "+firstName+" "+lastName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send email: %v", err)
	}

	return result, nil
}

// SendForgotPasswordNGNOtp sends forgot password OTP to Nigerian numbers
func (es *EmailService) SendForgotPasswordNGNOtp(phoneNumber, userID string) (interface{}, error) {
	if phoneNumber == "" {
		return nil, fmt.Errorf("recipient phone number is required")
	}

	otp := es.GenerateOtp()
	otpExpiresAt := time.Now().Add(5 * time.Minute)
	fmt.Println("Sending forgot password OTP to", phoneNumber[:3]+"****")

	if err := es.updateOrCreateOtp(userID, otp, otpExpiresAt, true); err != nil {
		return nil, fmt.Errorf("failed to update OTP record: %v", err)
	}

	formattedNumber := es.formatPhoneNumber(phoneNumber, "234")
	message := fmt.Sprintf("Your Otp verification code is %s. Valid for 5 minutes.", otp)

	result, err := es.sendSmsViaLenhub(formattedNumber, message)
	if err != nil {
		if clearErr := es.clearOtp(userID); clearErr != nil {
			fmt.Println("Failed to clear OTP after send failure:", clearErr)
		}
		return nil, fmt.Errorf("failed to send OTP: %v", err)
	}

	fmt.Println("Forgot password OTP sent successfully to", formattedNumber)
	return result, nil
}

// SendForgotPasswordGHCOtp sends forgot password OTP to Ghanaian numbers
func (es *EmailService) SendForgotPasswordGHCOtp(phoneNumber, userID string) (interface{}, error) {
	if phoneNumber == "" {
		return nil, fmt.Errorf("recipient phone number is required")
	}

	otp := es.GenerateOtp()
	otpExpiresAt := time.Now().Add(5 * time.Minute)
	fmt.Println("Sending forgot password OTP to", phoneNumber[:3]+"****")

	if err := es.updateOrCreateOtp(userID, otp, otpExpiresAt, true); err != nil {
		return nil, fmt.Errorf("failed to update OTP record: %v", err)
	}

	formattedNumber := es.formatPhoneNumber(phoneNumber, "233")
	message := fmt.Sprintf("Your OTP verification code is %s. Valid for 5 minutes.", otp)

	result, err := es.sendSmsViaHubtel(formattedNumber, message)
	if err != nil {
		if clearErr := es.clearOtp(userID); clearErr != nil {
			fmt.Println("Failed to clear OTP after send failure:", clearErr)
		}
		return nil, fmt.Errorf("failed to send OTP: %v", err)
	}

	fmt.Println("Forgot password OTP sent successfully to", formattedNumber)
	return result, nil
}

// SendForgotPasswordOtp sends forgot password OTP via email
func (es *EmailService) SendForgotPasswordEmailOtp(recipient, fullName, userID string) (interface{}, error) {
	if recipient == "" {
		return nil, fmt.Errorf("recipient email is required")
	}

	otp := es.GenerateOtp()
	otpExpiresAt := time.Now().Add(5 * time.Minute)
	fmt.Println("Sending password reset OTP to", recipient[:3]+"****@***.com")

	if err := config.DB.Model(&models.UserOtpSecurity{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"password_reset_verification_code":           otp,
		"password_reset_sent_to":                     "email",
		"password_reset_verification_code_created_at": time.Now(),
		"password_reset_verification_code_expires_at": otpExpiresAt,
	}).Error; err != nil {
		return nil, fmt.Errorf("failed to update OTP record: %v", err)
	}

	emailContent, err := es.GetEmailTemplate("forgot-password", map[string]string{
		"title":   "🔑 Reset Your Password",
		"message": "Use the OTP below to reset your password.",
		"name":    fullName,
		"otp":     otp,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get email template: %v", err)
	}

	result, err := es.sendEmailViaLenhub(
		recipient,
		"Your OTP for Password Reset",
		emailContent,
		"Dear "+fullName,
	)
	if err != nil {
		if clearErr := es.clearOtp(userID); clearErr != nil {
			fmt.Println("Failed to clear OTP after send failure:", clearErr)
		}
		return nil, fmt.Errorf("failed to send password reset OTP: %v", err)
	}

	return result, nil
}

// NombaPaymentSuccess sends payment success notification
func (es *EmailService) NombaPaymentSuccess(recipient, firstName, lastName, transactionReference string, amount float64) (map[string]interface{}, error) {
	if recipient == "" {
		return nil, fmt.Errorf("recipient email is required")
	}

	emailContent, err := es.GetEmailTemplate("nomba-payment-confirmation", map[string]string{
		"title":              "Payment Success",
		"message":            "Payment successful.",
		"name":              firstName + " " + lastName,
		"transactionReference": transactionReference,
		"amount":            fmt.Sprintf("%.2f", amount),
		"year":              fmt.Sprintf("%d", time.Now().Year()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get email template: %v", err)
	}

	result, err := es.sendEmailViaLenhub(
		recipient,
		"Payment Successful",
		emailContent,
		"Dear "+firstName+" "+lastName,
	)
	if err != nil {
		return map[string]interface{}{
			"success":    false,
			"statusCode": 500,
			"message":    "Failed to send message",
			"error":      err.Error(),
		}, nil
	}

	return map[string]interface{}{
		"success":    true,
		"statusCode": 200,
		"message":    "Payment Success",
		"token":      result,
	}, nil
}