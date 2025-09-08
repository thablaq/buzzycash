package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"io"
	"net/http"
	"github.com/dblaq/buzzycash/internal/config"
	"time"
)


const (
	SafetyBuffer  = 60  // seconds before expiry
	RefreshWindow = 300 // refresh 5m before expiry
	CountdownInterval = 1 * time.Minute  
)




func NewNombaAuthService() *NombaAuthService {
	s := &NombaAuthService{}
	log.Println("INFO: Initializing NombaAuthService. Fetching initial token...")

	// Block until first token is fetched
	var err error
	for {
		_, err = s.fetchToken()
		if err != nil {
			log.Printf("ERROR: Initial NB token fetch failed: %v. Retrying in 3 seconds...", err)
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("INFO: Initial NB token fetched successfully.")
		break
	}

	// Start background refresh loop
	go s.startTokenRefreshLoop()

	return s
}


func (s *NombaAuthService) startTokenRefreshLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.RLock()
		t := s.token
		s.mu.RUnlock()

		if t == nil {
			log.Println("WARN: No NB token available, fetching new one...")
			if _, err := s.fetchToken(); err != nil {
				log.Printf("ERROR: NB Token fetch failed: %v", err)
			}
			continue
		}

		// Parse the expiration time
		expiryTime, err := time.Parse(time.RFC3339, t.ExpiresAt)
		if err != nil {
			log.Printf("WARN: Failed to parse NB token expiration, fetching new token: %v", err)
			if _, err := s.fetchToken(); err != nil {
				log.Printf("ERROR: NB Token refresh failed: %v", err)
			}
			continue
		}

		// Check if token is already expired
		if time.Now().After(expiryTime) {
			log.Println("INFO: NB Token has expired, fetching new one...")
			if _, err := s.fetchToken(); err != nil {
				log.Printf("ERROR: NB Token refresh failed: %v", err)
			}
			continue
		}

		// Check if it's time to refresh (within 5 minutes of expiry)
		timeUntilExpiry := time.Until(expiryTime)
		if timeUntilExpiry <= time.Duration(RefreshWindow)*time.Second {
			log.Printf("INFO: NB Token expires in %v, refreshing...", timeUntilExpiry.Truncate(time.Second))
			if _, err := s.fetchToken(); err != nil {
				log.Printf("ERROR: NB Token refresh failed: %v", err)
			}
		} else {
			// Log how long until refresh for debugging
			if timeUntilExpiry < 10*time.Minute {
				log.Printf("DEBUG: NB Token refresh in %v", timeUntilExpiry.Truncate(time.Second))
			}
		}
	}
}


func (s *NombaAuthService) GetToken() (*TokenResponse, error) {
	s.mu.RLock()
	t := s.token
	s.mu.RUnlock()

	// Check if we have a valid token
	if t != nil {
		// Parse the expiration time from ISO string
		expiryTime, err := time.Parse(time.RFC3339, t.ExpiresAt)
		if err != nil {
			log.Printf("WARN: Failed to parse NB token expiration, fetching new token: %v", err)
		} else {
			// Check if token is still valid with safety buffer (60 seconds)
			if time.Now().Before(expiryTime.Add(-time.Duration(SafetyBuffer) * time.Second)) {
				return t, nil
			}
		}
	}

	// Fetch new token
	log.Println("INFO: Fetching new NB token for API request...")
	return s.fetchToken()
}



func (s *NombaAuthService) fetchToken() (*TokenResponse, error) {
	log.Println("INFO: Attempting to fetch a new access token from Nomba API.")
	
	body := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     config.AppConfig.NombaClientID,
		"client_secret": config.AppConfig.NombaApiKey,
	}
	payload, _ := json.Marshal(body)

	url := config.AppConfig.NombaApiBase + "auth/token/issue"
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("accountId", config.AppConfig.NombaAccountID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("NB token issue API error: status %d, body: %s",
			resp.StatusCode, string(bodyBytes))
	}

	// Parse the response into a temporary struct that matches the API
	var apiResp struct {
		Code        string `json:"code"`
		Description string `json:"description"`
		Data        struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresAt    string `json:"expiresAt"`
			BusinessID   string `json:"businessId"`
		} `json:"data"`
		Status bool `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode issue error: %w", err)
	}

	// Check if the API response was successful
	if apiResp.Code != "00" {
		return nil, fmt.Errorf("API returned error: %s", apiResp.Description)
	}

	t := &TokenResponse{
		AccessToken:  apiResp.Data.AccessToken,
		RefreshToken: apiResp.Data.RefreshToken,
		ExpiresAt:    apiResp.Data.ExpiresAt,
		BusinessID:   apiResp.Data.BusinessID,
	}

	// Store token
	s.mu.Lock()
	s.token = t
	s.mu.Unlock()

	// Parse the expiration time for logging
	expiryTime, err := time.Parse(time.RFC3339, t.ExpiresAt)
	if err != nil {
		log.Printf("WARN: Failed to parse NB expiration time: %v", err)
	} else {
		timeUntilExpiry := time.Until(expiryTime)
		log.Printf("INFO: NB Token fetched successfully. Expires at %v (in %v)", 
			expiryTime.Format("2006-01-02 15:04:05"), 
			timeUntilExpiry.Truncate(time.Second))
	}
	
	return t, nil
}



