package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
      "time"
	"github.com/dblaq/buzzycash/internal/config"
)



var (
	nbService     *NBService
	nbServiceOnce sync.Once
	
	nombaAuth     *NombaAuthService
	nombaAuthOnce sync.Once
)

// GetNombaAuthService returns the shared auth service singleton
func GetNombaAuthService() *NombaAuthService {
	nombaAuthOnce.Do(func() {
		nombaAuth = NewNombaAuthService()
	})
	return nombaAuth
}

// GetNBService returns the shared NB service singleton with proper auth dependency
func NBInstance() *NBService {
	nbServiceOnce.Do(func() {
		nbService = &NBService{
			client: &http.Client{
				Timeout: 30 * time.Second,
				// Add other client configurations if needed
			},
			auth: GetNombaAuthService(),
		}
	})
	return nbService
}

type NBService struct {
	client *http.Client
	auth   *NombaAuthService
}


func (s *NBService) CreateNBCheckout(req NBPaymentRequest) (string, string, error) {

	token, err := s.auth.GetToken()
	if err != nil {
    return "", "", fmt.Errorf("failed to retrieve access token: %v", err)
}

	body, err := json.Marshal(req)
	if err != nil {
		log.Printf("ERROR: Failed to marshal payment request: %v\n", err)
		return "", "", fmt.Errorf("failed to marshal payment request: %w", err)
	}
	log.Printf("DEBUG: Marshaled request body for Nomba: %s\n", string(body))
	url := config.AppConfig.NombaApiBase + "checkout/order"
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		log.Printf("ERROR: Failed to create new HTTP request: %v\n", err)
		return "", "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("accountId", config.AppConfig.NombaAccountID)
	httpReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", "", fmt.Errorf("failed to call nomba API: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}
	// log.Printf("DEBUG: Raw Nomba response: %s\n", string(b))

	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Nomba returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: Nomba API returned an error: %s\n", errorMessage)
		return "", "", fmt.Errorf("nomba error %d", resp.StatusCode)
	}

	// 5. Decode response
	var nb NBCheckoutResponse
    if err := json.Unmarshal(b, &nb); err != nil {
    return "", "", fmt.Errorf("invalid response decode: %w", err)
    }

	if nb.Data.CheckoutLink == "" || nb.Data.OrderReference == "" {
		return "", "", fmt.Errorf("invalid Nomba API response: %+v", nb)
	}

	log.Printf("INFO: Successfully created Nomba checkout link: %s\n", nb.Data.CheckoutLink)
	return nb.Data.CheckoutLink, nb.Data.OrderReference, nil
}

func (s *NBService) ListNBBanks() ([]Bank, error) {

	// 1. Get token
	token, err := s.auth.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access token: %w", err)
	}

	url := config.AppConfig.NombaApiBase + "transfers/banks"
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create new HTTP request: %v\n", err)
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("accountId", config.AppConfig.NombaAccountID)
	httpReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call nomba API: %w", err)
	}
	defer resp.Body.Close()
	
	// bodyBytes, err := io.ReadAll(resp.Body)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to read Nomba API response body: %w", err)
	// 	}
	// 	log.Printf("[NOMBA] Raw response body: %s", string(bodyBytes))

	// 5. Decode response
	var nb NBBankResponse
	if err := json.NewDecoder(resp.Body).Decode(&nb); err != nil {
		return nil, fmt.Errorf("invalid response decode: %w", err)
	}

	if nb.Data == nil {
		return nil, fmt.Errorf("invalid Nomba API response: %+v", nb)
	}

	return nb.Data, nil
}

func (s *NBService) FetchAccountDetails(req NBRetrieveAccountDetails) (*NBAccountDetails, error) {
	token, err := s.auth.GetToken()
	if err != nil  {
		return nil, fmt.Errorf("failed to retrieve access token: %w", err)
	}

	body, err := json.Marshal(req)
	if err != nil {
		log.Printf("ERROR: Failed to marshal account details request: %v\n", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	log.Printf("DEBUG: Marshaled request body for Nomba account resolution: %s\n", string(body))

	url := config.AppConfig.NombaApiBase + "transfers/bank/lookup"
	log.Printf("INFO: Sending account resolution request to URL: %s\n", url)

	// Create HTTP request
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for Nomba: %v\n", err)
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("accountId", config.AppConfig.NombaAccountID)
	httpReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	log.Printf("INFO: Executing HTTP POST request to Nomba API for account resolution...\n")
	resp, err := s.client.Do(httpReq)
	if err != nil {
		log.Printf("ERROR: HTTP request to Nomba failed: %v\n", err)
		return nil, fmt.Errorf("nomba HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	log.Printf("INFO: Received response from Nomba with status: %s\n", resp.Status)

	var nb NBAccountDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&nb); err != nil {
		log.Printf("ERROR: Failed to decode Nomba response: %v\n", err)
		return nil, fmt.Errorf("invalid response decode: %w", err)
	}
	// Validate response data
	if nb.Data.AccountName == "" || nb.Data.AccountNumber == "" {
		errorMessage := fmt.Sprintf(
			"Nomba account resolution failed: message='%s', accountName='%s'",
			nb.Message, nb.Data.AccountName,
		)
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("nomba account resolution failed: message='%s', accountName='%s'", nb.Message, nb.Data.AccountName)
	}

	return &NBAccountDetails{
		AccountNumber: nb.Data.AccountNumber,
		AccountName:   nb.Data.AccountName,
	}, nil
}

func (s *NBService) InitiateWithdrawal(req NBWithdrawalRequest) (string, error) {
	token, err := s.auth.GetToken()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve access token: %w", err)
	}

	body, err := json.Marshal(req)
	if err != nil {
		log.Printf("ERROR: Failed to marshal withdrawal request: %v\n", err)
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	log.Printf("DEBUG: Marshaled request body for Nomba withdrawal: %s\n", string(body))

	url := config.AppConfig.NombaApiBase + "transfers/bank"
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("accountId", config.AppConfig.NombaAccountID)
	httpReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("nomba HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw Nomba withdrawal response: %s\n", string(b))

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("nomba error %d: %s", resp.StatusCode, string(b))
	}

	var nb NBWithdrawalResp
	if err := json.Unmarshal(b, &nb); err != nil {
		return "", fmt.Errorf("failed to unmarshal Nomba response: %w", err)
	}
	log.Printf("DEBUG: Unmarshaled Nomba withdrawal response: %+v\n", nb)

	if !nb.Status {
		return "", fmt.Errorf("nomba withdrawal failed: message='%s'", nb.Message)
	}

	return nb.Message, nil
}
