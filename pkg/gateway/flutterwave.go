package gateway



import (
	"bytes"
	"encoding/json"
	"log"
	"fmt"
	"io"
	"net/http"
	// "os"
	"time"
	"github.com/dblaq/buzzycash/internal/config"
)

type PaymentService struct{
	client *http.Client
}

func NewPaymentService() *PaymentService{
	return &PaymentService{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}


func (s *PaymentService)CreateCheckout(req FWPaymentRequest) (string, error) {
	log.Printf("INFO: CreateCheckout initiated for payment request: %+v\n", req)

	body, err := json.Marshal(req)
	if err != nil {
		log.Printf("ERROR: Failed to marshal payment request: %v\n", err)
		return "", fmt.Errorf("failed to marshal payment request: %w", err)
	}
	log.Printf("DEBUG: Marshaled request body for Flutterwave: %s\n", string(body))

	url := config.AppConfig.FlutterwaveApiBase + "payments"
	log.Printf("INFO: Preparing to send payment request to URL: %s\n", url)

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		log.Printf("ERROR: Failed to create new HTTP request: %v\n", err)
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+config.AppConfig.FlutterwaveSecretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	log.Printf("INFO: Executing HTTP POST request to Flutterwave API...\n")
	resp, err := s.client.Do(httpReq)
	if err != nil {
		log.Printf("ERROR: HTTP request to Flutterwave failed: %v\n", err)
		return "", fmt.Errorf("flutterwave HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	log.Printf("INFO: Received response from Flutterwave with status: %s\n", resp.Status)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read response body from Flutterwave: %v\n", err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Flutterwave raw response body: %s\n", string(b))

	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("flutterwave error %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: Flutterwave API returned an error status: %s\n", errorMessage)
		return "", fmt.Errorf("flutterwave error %d", resp.StatusCode)
	}

	var fr fwCreateResp
	if err := json.Unmarshal(b, &fr); err != nil {
		log.Printf("ERROR: Failed to unmarshal Flutterwave response: %v, raw body: %s\n", err, string(b))
		return "", fmt.Errorf("failed to unmarshal flutterwave response: %w", err)
	}
	log.Printf("DEBUG: Unmarshaled Flutterwave response: %+v\n", fr)

	if fr.Status != "success" || fr.Data.Link == "" {
		errorMessage := fmt.Sprintf("flutterwave create payment failed: status='%s', message='%s', link='%s'", fr.Status, fr.Message, fr.Data.Link)
		log.Printf("ERROR: Flutterwave create payment failed: %s\n", errorMessage)
		return "", fmt.Errorf("flutterwave create payment failed: status='%s', message='%s', link='%s'", fr.Status, fr.Message, fr.Data.Link)
	}

	log.Printf("INFO: Successfully created Flutterwave checkout link: %s\n", fr.Data.Link)
	return fr.Data.Link, nil
}