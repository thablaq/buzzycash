package gaming

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
     "log"
     "time"
	"sync"
	 "github.com/dblaq/buzzycash/internal/config"
)

var (
	gmService     *GMService
	gmServiceOnce sync.Once
	
	gmAuth       *GamingAuthService
	gmAuthOnce   sync.Once
)


func GmAuthInstance() *GamingAuthService {
	gmAuthOnce.Do(func() {
		log.Println("INFO: Initializing singleton GamingAuthService.")
		gmAuth = NewGamingAuthService()
	})
	return gmAuth
}

func GMInstance() *GMService {
	gmServiceOnce.Do(func() {
		log.Println("INFO: Initializing singleton GMService.")
		gmService = &GMService{
			client: &http.Client{
				Timeout: 30 * time.Second,
			},
			auth: GmAuthInstance(),
		}
	})
	return gmService
}



type GMService struct {
	client *http.Client
	auth *GamingAuthService
}


 





// RegisterUser registers a new user in the gaming system
func (gs *GMService) RegisterUser(username, email, firstName, lastName string) (map[string]interface{}, error) {
	log.Println("Starting RegisterUser")
     token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ",err)
		return nil, err
	}
	req := RegisterUserRequest{
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}
	log.Println("Request data prepared for RegisterUser")
    
	body, err := json.Marshal(req)
	if err != nil {
		log.Printf("ERROR: Failed to marshal payment request: %v\n", err)
		return nil, fmt.Errorf("failed to marshal payment request: %w", err)
	}
	log.Printf("DEBUG: Marshaled request body for Nomba: %s\n", string(body))
	
	url := config.AppConfig.MaekandexGamingUrl + "register/user/"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("ERROR: Failed to register use with engine: %v\n", err)
		return nil,fmt.Errorf("failed to create register request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for RegisterUser")
	
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to register user: ", err)
		return nil, fmt.Errorf("failed to register user: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for RegisterUser")
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw Nomba response: %s\n", string(b))

	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Nomba returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: Nomba API returned an error: %s\n", errorMessage)
		return nil,fmt.Errorf("nomba error %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode RegisterUser response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("RegisterUser response decoded successfully")
	return result, nil
}

// StartGame starts a game
func (gs *GMService) StartGame(gameID string) (map[string]interface{}, error) {
	log.Println("Starting StartGame")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare request payload
	reqData := GameRequest{
		GameID:    gameID,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}
	body, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("ERROR: Failed to marshal StartGame request: %v\n", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	log.Printf("DEBUG: Marshaled StartGame request body: %s\n", string(body))

	// Build HTTP request
	url := config.AppConfig.MaekandexGamingUrl + "games/start/"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for StartGame: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for StartGame")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to start game: ", err)
		return nil, fmt.Errorf("failed to start game: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for StartGame")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw StartGame response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode StartGame response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Game started successfully: ", result)

	return result, nil
}

// StopGame stops a game
func (gs *GMService) StopGame(gameID string) (map[string]interface{}, error) {
	log.Println("Starting StopGame")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare request payload
	reqData := GameRequest{
		GameID:    gameID,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}
	body, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("ERROR: Failed to marshal StopGame request: %v\n", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	log.Printf("DEBUG: Marshaled StopGame request body: %s\n", string(body))

	// Build HTTP request
	url := config.AppConfig.MaekandexGamingUrl + "games/stop/"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for StopGame: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for StopGame")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to stop game: ", err)
		return nil, fmt.Errorf("failed to stop game: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for StopGame")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw StopGame response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode StopGame response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Game stopped successfully: ", result)

	return result, nil
}

// GetDraws retrieves draws for a game
func (gs *GMService) GetDraws(gameID string) (map[string]interface{}, error) {
	log.Println("Starting GetDraws")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare request payload
	reqData := GameRequest{
		GameID:    gameID,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}
	body, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("ERROR: Failed to marshal GetDraws request: %v\n", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	log.Printf("DEBUG: Marshaled GetDraws request body: %s\n", string(body))

	// Build HTTP request
	url := config.AppConfig.MaekandexGamingUrl + "games/draw/"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for GetDraws: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for GetDraws")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to get draws: ", err)
		return nil, fmt.Errorf("failed to get draws: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetDraws")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw GetDraws response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode GetDraws response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Draws retrieved successfully: ", result)

	return result, nil
}

// GetAllTickets retrieves all tickets
func (gs *GMService) GetAllTickets() (map[string]interface{}, error) {
	log.Println("Starting GetAllTickets")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare query parameters
	params := map[string]string{
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	urlStr := config.AppConfig.MaekandexGamingUrl + "admin/get_tickets/?" + query.Encode()

	// Build HTTP request
	httpReq, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for GetAllTickets: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for GetAllTickets")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to get tickets: ", err)
		return nil, fmt.Errorf("failed to get tickets: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetAllTickets")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw GetAllTickets response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode GetAllTickets response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Tickets retrieved successfully: ", result)

	return result, nil
}

// GetWalletBalances retrieves all wallet balances
func (gs *GMService) GetWalletBalances() (map[string]interface{}, error) {
	log.Println("Starting GetWalletBalances")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare query parameters
	params := map[string]string{
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	urlStr := config.AppConfig.MaekandexGamingUrl + "all/wallet?" + query.Encode()

	// Build HTTP request
	httpReq, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for GetWalletBalances: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for GetWalletBalances")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to get user balances: ", err)
		return nil, fmt.Errorf("failed to get user balances: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetWalletBalances")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw GetWalletBalances response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode GetWalletBalances response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("User balances retrieved successfully: ", result)

	return result, nil
}

// GetUserWallet retrieves wallet for a specific user
func (gs *GMService) GetUserWallet(username string) (map[string]interface{}, error) {
	log.Println("Starting GetUserWallet")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare query parameters
	params := map[string]string{
		"user_id": username,
	}
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	urlStr := config.AppConfig.MaekandexGamingUrl + "check/wallet?" + query.Encode()

	// Build HTTP request
	httpReq, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for GetUserWallet: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for GetUserWallet")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to get wallet for username: ", err)
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetUserWallet")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw GetUserWallet response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode GetUserWallet response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("User wallet retrieved successfully for username")

	return result, nil
}

// BuyTicket purchases a ticket for a game
func (gs *GMService) BuyTicket(gameID, username string, quantity int, amountPaid int64) (*BuyTicketResponse, error) {
	log.Println("Starting BuyTicket")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare request payload
	reqData := BuyTicketRequest{
		GameID:     gameID,
		Username:   username,
		Quantity:   quantity,
		AmountPaid: amountPaid,
		CompanyID:  config.AppConfig.BuzzyCashCompanyID,
	}
	body, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("ERROR: Failed to marshal BuyTicket request: %v\n", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	log.Printf("DEBUG: Marshaled BuyTicket request body: %s\n", string(body))

	// Build HTTP request
	url := config.AppConfig.MaekandexGamingUrl + "games/buy/"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for BuyTicket: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for BuyTicket")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to purchase ticket: ", err)
		return nil, fmt.Errorf("failed to purchase ticket: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for BuyTicket")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw BuyTicket response: %s\n", string(b))

	// error handling
	if resp.StatusCode >= 300 {
		var apiResp map[string]interface{}
		msg := "unknown error"

		if err := json.Unmarshal(b, &apiResp); err == nil {
			if m, ok := apiResp["message"].(string); ok {
				msg = m
			} else if m, ok := apiResp["error"].(string); ok {
				msg = m
			}
		} else {
			// fallback if response isn't JSON
			msg = string(b)
		}

		log.Printf("ERROR: Gaming API returned %d: %s", resp.StatusCode, msg)
		return nil, &APIError{StatusCode: resp.StatusCode, Message: msg}
	}

	// Decode response
	var result BuyTicketResponse
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode BuyTicket response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Ticket purchased successfully: ", result)

	return &result, nil
}


// GetUserTickets retrieves tickets for a specific user
func (gs *GMService) GetUserTickets(username string) (map[string]interface{}, error) {
	log.Println("Starting GetUserTickets")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare query parameters
	params := map[string]string{
		"user_id":    username,
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	urlStr := config.AppConfig.MaekandexGamingUrl + "users/tickets?" + query.Encode()

	// Build HTTP request
	httpReq, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for GetUserTickets: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for GetUserTickets")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to get tickets for username: ", err)
		return nil, fmt.Errorf("failed to get tickets: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetUserTickets")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw GetUserTickets response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode GetUserTickets response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("User tickets retrieved successfully for username")

	return result, nil
}

// GetUserResults retrieves results for a specific user
func (gs *GMService) GetUserResults(username string) (map[string]interface{}, error) {
	log.Println("Starting GetUserResults")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare query parameters
	params := map[string]string{
		"username":   username,
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	urlStr := config.AppConfig.MaekandexGamingUrl + "games/results?" + query.Encode()

	// Build HTTP request
	httpReq, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for GetUserResults: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for GetUserResults")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to get results for username: ", err)
		return nil, fmt.Errorf("failed to get results: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetUserResults")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw GetUserResults response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode GetUserResults response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("User results retrieved successfully for username")

	return result, nil
}

// GetWinnerLogs retrieves winner logs
func (gs *GMService) GetWinnerLogs() (map[string]interface{}, error) {
	log.Println("Starting GetWinnerLogs")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil || token.Accesstoken == "" {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Build HTTP request (no query params)
	url := config.AppConfig.MaekandexGamingUrl + "winners/logs/"
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for GetWinnerLogs: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for GetWinnerLogs")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to get winner logs: ", err)
		return nil, fmt.Errorf("failed to get winner logs: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetWinnerLogs")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw GetWinnerLogs response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode winner logs response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Winner logs retrieved successfully: ", result)

	return result, nil
}

// GetLeaderBoard retrieves the leaderboard
func (gs *GMService) GetLeaderBoard() (map[string]interface{}, error) {
	log.Println("Starting GetLeaderBoard")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Build HTTP request
	url := config.AppConfig.MaekandexGamingUrl + "leaderboard"
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for GetLeaderBoard: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for GetLeaderBoard")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to get leaderboard: ", err)
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetLeaderBoard")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw GetLeaderBoard response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode leaderboard response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Leaderboard retrieved successfully: ", result)

	return result, nil
}

// GetVirtualGames retrieves available virtual games
func (gs *GMService) GetVirtualGames() ([]interface{}, error) {
	log.Println("Starting GetVirtualGames")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Build HTTP request
	url := config.AppConfig.MaekandexGamingUrl + "list/virtual/games"
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for GetVirtualGames: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for GetVirtualGames")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to get virtual games: ", err)
		return nil, fmt.Errorf("failed to get virtual games: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetVirtualGames")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw GetVirtualGames response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode virtual games response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract data
	data, ok := result["data"]
	if !ok {
		log.Println("Invalid games data received from Meakindex")
		return nil, fmt.Errorf("invalid games data received from Meakindex")
	}

	gamesList, ok := data.([]interface{})
	if !ok {
		log.Println("Invalid games data format received from Meakindex")
		return nil, fmt.Errorf("invalid games data format received from Meakindex")
	}

	log.Println("Virtual games retrieved successfully: ", gamesList)
	return gamesList, nil
}

// StartVirtualGame starts a virtual game
func (gs *GMService) StartVirtualGame(gameType, username string) (map[string]interface{}, error) {
	log.Println("Starting StartVirtualGame")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare query parameters
	params := map[string]string{
		"gameType": gameType,
		"username": username,
	}
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	urlStr := config.AppConfig.MaekandexGamingUrl + "virtualstart/game?" + query.Encode()

	// Build HTTP request
	httpReq, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for StartVirtualGame: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for StartVirtualGame")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to start virtual game for gameType: ", err)
		return nil, fmt.Errorf("failed to start virtual game: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for StartVirtualGame")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw StartVirtualGame response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode virtual game start response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Virtual game started successfully for gameType: ", result)

	return result, nil
}

// CreateGames creates a new game (admin function)
func (gs *GMService) CreateGames(gameName string, amount int64, drawInterval int, winningPercentage float64, maxWinners int, date string, weightedDistribution bool) (map[string]interface{}, error) {
	log.Println("Starting CreateGames")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare request payload
	reqData := CreateGameRequest{
		CompanyID:            config.AppConfig.BuzzyCashCompanyID,
		GameName:             gameName,
		Amount:               amount,
		DrawInterval:         drawInterval,
		WinningPercentage:    winningPercentage,
		MaxWinners:           maxWinners,
		Date:                 date,
		WeightedDistribution: weightedDistribution,
	}
	body, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("ERROR: Failed to marshal CreateGames request: %v\n", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	log.Printf("DEBUG: Marshaled CreateGames request body: %s\n", string(body))

	// Build HTTP request
	url := config.AppConfig.MaekandexGamingUrl + "create/games/"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for CreateGames: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for CreateGames")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Error creating game: ", err)
		return nil, fmt.Errorf("failed to create game with Meakindex: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for CreateGames")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw CreateGames response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode create game response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Game created successfully: ", result)

	return result, nil
}

// GetGames retrieves all games
func (gs *GMService) GetGames() (map[string]interface{}, error) {
	log.Println("Starting GetGames")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare query parameters
	params := map[string]string{
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	urlStr := config.AppConfig.MaekandexGamingUrl + "games/?" + query.Encode()

	// Build HTTP request
	httpReq, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for GetGames: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for GetGames")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to get games: ", err)
		return nil, fmt.Errorf("failed to get games: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetGames")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw GetGames response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode games response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Games retrieved successfully: ", result)

	return result, nil
}

// DebitUserWallet debits amount from user's wallet
func (gs *GMService) DebitUserWallet(username string, amount float64) (map[string]interface{}, error) {
	log.Println("Starting DebitUserWallet")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare request payload
	reqData := DebitWalletRequest{
		UserID:    username,
		Amount:    amount,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}
	body, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("ERROR: Failed to marshal DebitUserWallet request: %v\n", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	log.Printf("DEBUG: Marshaled DebitUserWallet request body: %s\n", string(body))

	// Build HTTP request
	url := config.AppConfig.MaekandexGamingUrl + "debit/wallet/"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for DebitUserWallet: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for DebitUserWallet")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Failed to debit wallet: ", err)
		return nil, fmt.Errorf("failed to debit wallet: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for DebitUserWallet")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw DebitUserWallet response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode debit wallet response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Wallet debited successfully: ", result)

	return result, nil
}

// CreditUserWallet credits amount to user's wallet
func (gs *GMService) CreditUserWallet(username string, amount float64) (*PaymentResponse, error) {
	log.Println("Starting CreditUserWallet")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare request payload
	reqData := CreditWalletRequest{
		UserID:    username,
		Amount:    amount,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}
	body, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("ERROR: Failed to marshal CreditUserWallet request: %v\n", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	log.Printf("DEBUG: Marshaled CreditUserWallet request body: %s\n", string(body))

	// Build HTTP request
	url := config.AppConfig.MaekandexGamingUrl + "credit/wallet/"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for CreditUserWallet: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for CreditUserWallet")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Error getting credit wallet response: ", err)
		return nil, fmt.Errorf("failed to get wallet response: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for CreditUserWallet")

	// Read and log raw response
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body: ", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw CreditUserWallet response: %s\n", string(rawBody))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(rawBody))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result PaymentResponse
	if err := json.Unmarshal(rawBody, &result); err != nil {
		log.Println("Failed to decode credit wallet response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Decoded credit wallet response:", result)

	return &result, nil
}

// ListPayouts lists all payouts
func (gs *GMService) ListPayouts() (map[string]interface{}, error) {
	log.Println("Starting ListPayouts")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Build HTTP request
	url := config.AppConfig.MaekandexGamingUrl + "payouts/"
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for ListPayouts: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for ListPayouts")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Error listing payouts: ", err)
		return nil, fmt.Errorf("failed to list payouts: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for ListPayouts")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw ListPayouts response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode payouts response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Payouts listed successfully: ", result)

	return result, nil
}

// ListUserPayout lists payouts for a specific user
func (gs *GMService) ListUserPayout(username string) (map[string]interface{}, error) {
	log.Println("Starting ListUserPayout")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare query parameters
	params := map[string]string{
		"username": username,
	}
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	urlStr := config.AppConfig.MaekandexGamingUrl + "payouts/?" + query.Encode()

	// Build HTTP request
	httpReq, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for ListUserPayout: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for ListUserPayout")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Error listing user payouts: ", err)
		return nil, fmt.Errorf("failed to list user payouts: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for ListUserPayout")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw ListUserPayout response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode user payouts response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("User payouts listed successfully: ", result)

	return result, nil
}

// PayoutByAdmin processes payout by admin
func (gs *GMService) PayoutByAdmin(payoutID string) (map[string]interface{}, error) {
	log.Println("Starting PayoutByAdmin")

	// Get access token
	token, err := gs.auth.GetToken()
	if err != nil  {
		log.Println("Failed to get access token: ", err)
		return nil, err
	}

	// Prepare query parameters
	params := map[string]string{
		"payout_id": payoutID,
	}
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	urlStr := config.AppConfig.MaekandexGamingUrl + "payouts/?" + query.Encode()

	// Build HTTP request
	httpReq, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for PayoutByAdmin: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.Accesstoken)
	log.Println("Headers set successfully for PayoutByAdmin")

	// Send request
	resp, err := gs.client.Do(httpReq)
	if err != nil {
		log.Println("Error processing payout: ", err)
		return nil, fmt.Errorf("failed to process payout: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for PayoutByAdmin")

	// Read and log raw response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("DEBUG: Raw PayoutByAdmin response: %s\n", string(b))

	// Check HTTP status
	if resp.StatusCode >= 300 {
		errorMessage := fmt.Sprintf("Gaming API returned error status %d: %s", resp.StatusCode, string(b))
		log.Printf("ERROR: %s\n", errorMessage)
		return nil, fmt.Errorf("gaming API error %d", resp.StatusCode)
	}

	// Decode response
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to decode payout processing response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Payout processed successfully: ", result)

	return result, nil
}