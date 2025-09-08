package gaming

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/dblaq/buzzycash/internal/config"
)

const (
	RefreshWindow = 300 // refresh 5m before expiry
	SafetyBuffer  = 60  // seconds before expiry
)

func NewGamingAuthService() *GamingAuthService {
	s := &GamingAuthService{}
	log.Println("INFO: Initializing GamingAuthService. Fetching initial token...")

	// Block until first token is fetched
	var err error
	for {
		_, err = s.fetchToken()
		if err != nil {
			log.Printf("ERROR: Gaming Initial token fetch failed: %v. Retrying in 3 seconds...", err)
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("INFO: Gaming Initial token fetched successfully.")
		break
	}

	// Start background refresh loop
	go s.startTokenRefreshLoop()

	return s
}

func (s *GamingAuthService) startTokenRefreshLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.RLock()
		t := s.token
		s.mu.RUnlock()

		if t == nil {
			log.Println("WARN: No Gaming token available, fetching new one...")
			if _, err := s.fetchToken(); err != nil {
				log.Printf("ERROR: Gaming Token fetch failed: %v", err)
			}
			continue
		}

		// Calculate expiration time based on when token was fetched + duration
		if s.tokenFetchTime.IsZero() {
			log.Println("WARN: No Gaming token fetch time recorded, fetching new token...")
			if _, err := s.fetchToken(); err != nil {
				log.Printf("ERROR: Gaming Token refresh failed: %v", err)
			}
			continue
		}

		tokenDuration, err := parseISODuration(t.ExpiresAt)
		if err != nil {
			log.Printf("WARN: Failed to parse gaming token duration, fetching new token: %v", err)
			if _, err := s.fetchToken(); err != nil {
				log.Printf("ERROR: Gaming Token refresh failed: %v", err)
			}
			continue
		}

		expiryTime := s.tokenFetchTime.Add(tokenDuration)

		// Check if token is already expired
		if time.Now().After(expiryTime) {
			log.Println("INFO: Gaming Token has expired, fetching new one...")
			if _, err := s.fetchToken(); err != nil {
				log.Printf("ERROR: Gaming Token refresh failed: %v", err)
			}
			continue
		}

		// Check if it's time to refresh (within 5 minutes of expiry)
		timeUntilExpiry := time.Until(expiryTime)
		if timeUntilExpiry <= time.Duration(RefreshWindow)*time.Second {
			log.Printf("INFO: Gaming Token expires in %v, refreshing...", timeUntilExpiry.Truncate(time.Second))
			if _, err := s.fetchToken(); err != nil {
				log.Printf("ERROR: Gaming Token refresh failed: %v", err)
			}
		}
	}
}

func (s *GamingAuthService) GetToken() (*NBTokenResponse, error) {
	s.mu.RLock()
	t := s.token
	tokenFetchTime := s.tokenFetchTime
	s.mu.RUnlock()

	// Check if we have a valid token
	if t != nil && !tokenFetchTime.IsZero() {
		tokenDuration, err := parseISODuration(t.ExpiresAt)
		if err != nil {
			log.Printf("WARN: Failed to parse gaming token duration, fetching new token: %v", err)
		} else {
			expiryTime := tokenFetchTime.Add(tokenDuration)
			// Check if token is still valid with safety buffer
			if time.Now().Before(expiryTime.Add(-time.Duration(SafetyBuffer) * time.Second)) {
				return t, nil
			}
		}
	}

	// Fetch new token
	log.Println("INFO: Fetching new gaming token for API request...")
	return s.fetchToken()
}

func (gs *GamingAuthService) fetchToken() (*NBTokenResponse, error) {
	loginData := map[string]string{
		"username":   config.AppConfig.BuzzyCashUsername,
		"password":   config.AppConfig.BuzzyCashPassword,
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}
	payload, _ := json.Marshal(loginData)

	url := config.AppConfig.MaekandexGamingUrl + "login/"
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gaming login failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenResp NBTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	// Parse the duration to calculate expiration
	tokenDuration, err := parseISODuration(tokenResp.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gaming token duration: %w", err)
	}

	expiryTime := time.Now().Add(tokenDuration)

	// Store token and fetch time
	gs.mu.Lock()
	gs.token = &tokenResp
	gs.tokenFetchTime = time.Now()
	gs.mu.Unlock()

	log.Printf("INFO: Gaming Token fetched successfully. Expires at %v (in %v)",
		expiryTime.Format("2006-01-02 15:04:05"),
		tokenDuration.Truncate(time.Second))

	return gs.token, nil
}

// parseISODuration parses ISO 8601 duration format like "P0DT01H00M00S"
func parseISODuration(durationStr string) (time.Duration, error) {
	// Regex to parse ISO 8601 duration format
	re := regexp.MustCompile(`P(?P<days>\d+D)?T(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`)
	matches := re.FindStringSubmatch(durationStr)

	if matches == nil {
		return 0, fmt.Errorf("invalid ISO duration format: %s", durationStr)
	}

	var total time.Duration

	// Helper function to parse duration components
	parsePart := func(value, unit string) time.Duration {
		if value == "" {
			return 0
		}
		// Remove the unit character and parse as integer
		numStr := value[:len(value)-1]
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return 0
		}
		switch unit {
		case "D":
			return time.Duration(num) * 24 * time.Hour
		case "H":
			return time.Duration(num) * time.Hour
		case "M":
			return time.Duration(num) * time.Minute
		case "S":
			return time.Duration(num) * time.Second
		default:
			return 0
		}
	}

	// Extract and parse each component
	total += parsePart(matches[1], "D") // days
	total += parsePart(matches[2], "H") // hours
	total += parsePart(matches[3], "M") // minutes
	total += parsePart(matches[4], "S") // seconds

	return total, nil
}
