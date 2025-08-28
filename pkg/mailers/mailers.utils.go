package mailers


import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
    "text/template"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
     "log"
	// "gorm.io/gorm"

	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/models"
)



type EmailService struct{}

// GetEmailTemplate reads and processes an email template
func (es *EmailService) GetEmailTemplate(templateName string, context map[string]string) (string, error) {
    cwd, _ := os.Getwd()
    templatePath := filepath.Join(cwd, "internal/templates", templateName+".html")

    tmpl, err := template.ParseFiles(templatePath)
    if err != nil {
        return "", fmt.Errorf("failed to parse template file: %v", err)
    }
    
    context["year"] = fmt.Sprintf("%d", time.Now().Year())

    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, context); err != nil {
        return "", fmt.Errorf("failed to execute template: %v", err)
    }

    return buf.String(), nil
}

// GenerateOtp generates a 6-digit OTP
func (es *EmailService) GenerateOtp() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

// clearOtp clears OTP from database
func (es *EmailService) clearOtp(userID string, action models.OtpAction) error {
	log.Println("Clearing OTP for user:", userID, "action:", action)

	err := config.DB.Model(&models.UserOtpSecurity{}).
		Where("user_id = ?" , userID).
		Updates(map[string]interface{}{
			"code":       "",
			"created_at": nil,
			"expires_at": nil,
			"sent_to":    nil,
			"action":     "",
		}).Error

	if err != nil {
		log.Println("Failed to clear OTP for user:", userID, "Error:", err)
	} else {
		log.Println("OTP cleared successfully for user:", userID)
	}
	return err
}


func (es *EmailService) updateOrCreateOtp(userID, otp string, expiresAt time.Time, action models.OtpAction, sentTo string) error {
	now := time.Now()
	
	log.Println("updateOrCreateOtp called with action:", string(action))

	// Try to update existing record for this user
	updateResult := config.DB.Model(&models.UserOtpSecurity{}).
		Where("user_id = ?", userID). // Look for any existing OTP for this user
		Updates(map[string]interface{}{
			"code":                               otp,
			"created_at":                         now,
			"expires_at":                         expiresAt,
			"sent_to":                            sentTo,
			"action":                             action,
		})
	
	if updateResult.Error != nil {
		log.Println("Error updating OTP:", updateResult.Error)
		return updateResult.Error
	}
	
	// If no record was updated (user has no existing OTP), create a new one
	if updateResult.RowsAffected == 0 {
		otpSecurity := models.UserOtpSecurity{
			UserID:    userID,
			Code:      otp,
			CreatedAt: now,
			ExpiresAt: expiresAt,
			IsOtpVerifiedForPasswordReset: false,
			SentTo:    sentTo,
			Action:    action,
		}
		
		log.Println("Creating new OTP record with action:", string(action))
		return config.DB.Create(&otpSecurity).Error
	}
	
	log.Println("Updated existing OTP record with new action:", string(action))
	return nil
}

// sendSmsViaLenhub sends SMS using LENHUB API
func (es *EmailService) sendSmsViaLenhub(phoneNumber, message string) (interface{}, error) {
	log.Println("Preparing to send SMS via Lenhub")
	log.Printf("Phone Number: %s, Message: %s", phoneNumber, message)

	payload := map[string]interface{}{
		"client_id":       config.AppConfig.LenhubClientID,
		"receiver_number": phoneNumber,
		"message":         message,
		"sender_id":       config.AppConfig.BuzzyCashSenderID,
		"types":           "2",
	}

	jsonPayload, _ := json.Marshal(payload)
	log.Println("Payload prepared for Lenhub SMS API:")

	req, err := http.NewRequest("POST", config.AppConfig.LenhubApiBase+"sendsms/api", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("Failed to create HTTP request:", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.AppConfig.LenhubApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	log.Println("Sending SMS request to Lenhub API")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send SMS request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Received response from Lenhub API with status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to send SMS, status code:", resp.StatusCode)
		return nil, fmt.Errorf("failed to send SMS, status code: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode response body:", err)
		return nil, err
	}

	log.Println("SMS sent successfully via Lenhub")
	return result, nil
}

// sendSmsViaHubtel sends SMS using Hubtel API
func (es *EmailService) sendSmsViaHubtel(phoneNumber, message string) (interface{}, error) {
	log.Println("Preparing to send SMS via Hubtel")
	log.Printf("Phone Number: %s, Message: %s", phoneNumber, message)

	payload := map[string]interface{}{
		"From":    config.AppConfig.HubtelSenderID,
		"To":      phoneNumber,
		"Content": message,
	}

	jsonPayload, _ := json.Marshal(payload)
	log.Println("Payload prepared for Hubtel SMS API:")

	req, err := http.NewRequest("POST", config.AppConfig.HubtelApiBase+"/messages/send", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("Failed to create HTTP request:", err)
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(config.AppConfig.HubtelClientID + ":" + config.AppConfig.HubtelClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	log.Println("Sending SMS request to Hubtel API")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send SMS request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Received response from Hubtel API with status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to send SMS, status code:", resp.StatusCode)
		return nil, fmt.Errorf("failed to send SMS, status code: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode response body:", err)
		return nil, err
	}

	log.Println("SMS sent successfully via Hubtel")
	return result, nil
}

// sendEmailViaLenhub sends email using LENHUB API
func (es *EmailService) sendEmailViaLenhub(recipient, subject, message, greetings string) (interface{}, error) {
	log.Println("Preparing to send email via Lenhub")
	log.Printf("Recipient: %s, Subject: %s, Greetings: %s", recipient, subject, greetings)

	payload := map[string]interface{}{
		"client_id":  config.AppConfig.LenhubClientID,
		"subject":    subject,
		"message":    message,
		"Sender_id":  config.AppConfig.LenhubSenderID,
		"recipient":  recipient,
		"greetings":  greetings,
	}

	jsonPayload, _ := json.Marshal(payload)
	log.Println("Payload prepared for Lenhub email API")

	req, err := http.NewRequest("POST", config.AppConfig.LenhubApiBase+"send/email/api", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("Failed to create HTTP request:", err)
		return nil, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.AppConfig.LenhubApiKey)

	client := &http.Client{}
	log.Println("Sending email request to Lenhub API")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send email request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Received response from Lenhub API with status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to send email, status code:", resp.StatusCode)
		return nil, fmt.Errorf("failed to send email, status code: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode response body:", err)
		return nil, err
	}

	log.Println("Email sent successfully via Lenhub")
	return result, nil
}

// formatPhoneNumber formats phone number based on country code
func (es *EmailService) formatPhoneNumber(phoneNumber, countryCode string) string {
    log.Println("Formatting phone number:", phoneNumber, "with country code:", countryCode)
    if strings.HasPrefix(phoneNumber, "0") {
        formattedNumber := countryCode + phoneNumber[1:]
        log.Println("Formatted phone number (removed leading 0):", formattedNumber)
        return formattedNumber
    }
    if strings.HasPrefix(phoneNumber, countryCode) {
        log.Println("Phone number already formatted:", phoneNumber)
        return phoneNumber
    }
    formattedNumber := countryCode + phoneNumber
    log.Println("Formatted phone number (added country code):", formattedNumber)
    return formattedNumber
}
